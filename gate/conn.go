package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira/gate/crypto"
	"github.com/lujingwei002/gira/gate/message"
	"github.com/lujingwei002/gira/gate/packet"
	"golang.org/x/sync/errgroup"
)

const (
	_ int32 = iota
	conn_status_start
	conn_status_handshake
	conn_status_working
	conn_status_closed
)

type Conn struct {
	session *Session
	conn    net.Conn
	state   int32
	lastAt  int64
	decoder *packet.Decoder
	gateway *Gateway

	chSend    chan []byte
	chMessage chan *Message

	ctx        context.Context
	cancelFunc context.CancelFunc
	errGroup   *errgroup.Group
	errCtx     context.Context

	ctrlCtx        context.Context
	ctrlCancelFunc context.CancelFunc
	lastErr        error
}
type handShake_request struct {
	Sys struct {
		Token   string `json:"token"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"sys"`
}

// Create new agent instance
func newConn(gateway *Gateway) *Conn {
	self := &Conn{
		gateway: gateway,
		state:   conn_status_start,
		lastAt:  time.Now().Unix(),
		decoder: packet.NewDecoder(),
	}
	self.session = newSession(self)
	return self
}

func (a *Conn) send(data []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = ErrBrokenPipe
		}
	}()
	a.chSend <- data
	return
}

func (self *Conn) push(route string, data []byte) error {
	if self.gateway.debug {
		log.Infow("gate conn push", "session_id", self.session.ID(), "route", route, "len", len(data))
	}
	if self.status() != conn_status_working {
		return ErrNotWorking
	}
	if self.lastErr != nil {
		return self.lastErr
	}
	var err error
	payload, err := self.serialize(data)
	if err != nil {
		log.Infow("gate conn push fail", "session_id", self.session.ID(), "route", route, "error", err)
	}
	m := &message.Message{
		Type:  message.Push,
		Data:  payload,
		Route: "",
		ID:    0,
	}
	em, err := m.Encode()
	if err != nil {
		log.Info(err.Error())
		return err
	}
	p, err := packet.Encode(packet.Data, em)
	if err != nil {
		log.Info(err)
		return err
	}

	return self.send(p)
}

func (self *Conn) response(mid uint64, data []byte) error {
	if self.gateway.debug {
		log.Infow("gate conn response", "session_id", self.session.ID(), "req_id", mid, "len", len(data))
	}
	if self.lastErr != nil {
		return self.lastErr
	}
	if self.status() != conn_status_working {
		return ErrBrokenPipe
	}
	if mid <= 0 {
		return ErrSessionOnNotify
	}
	var err error
	payload, err := self.serialize(data)
	if err != nil {
		log.Infow("gate conn response fail", "session_id", self.session.ID(), "req_id", mid, "error", err)
	}
	m := &message.Message{
		Type:  message.Response,
		Data:  payload,
		Route: "",
		ID:    mid,
	}
	em, err := m.Encode()
	if err != nil {
		return err
	}
	p, err := packet.Encode(packet.Data, em)
	if err != nil {
		return err
	}
	return self.send(p)
}

func (self *Conn) kick(reason string) error {
	data, err := packet.Encode(packet.Kick, []byte(reason))
	if err != nil {
		return err
	}
	return self.send(data)
}

func (self *Conn) sendServerSuspendPacket(reason string) error {
	data, err := packet.Encode(packet.ServerSuspend, []byte(reason))
	if err != nil {
		return err
	}
	return self.send(data)
}

func (self *Conn) sendServerResumePacket(reason string) error {
	data, err := packet.Encode(packet.ServerResume, []byte(reason))
	if err != nil {
		return err
	}
	return self.send(data)
}

func (self *Conn) sendServerDownPacket(reason string) error {
	data, err := packet.Encode(packet.ServerDown, []byte(reason))
	if err != nil {
		return err
	}
	return self.send(data)
}

func (self *Conn) sendServerMaintainPacket(reason string) error {
	data, err := packet.Encode(packet.ServerMaintain, []byte(reason))
	if err != nil {
		return err
	}
	return self.send(data)
}

func (self *Conn) close() error {
	if self.status() == conn_status_closed {
		return nil
	}
	self.setStatus(conn_status_closed)
	if self.gateway.debug {
		log.Infow("close gate conn", "session_id", self.session.ID(), "remote_addr", self.conn.RemoteAddr())
	}
	self.cancelFunc()
	return nil
}

func (self *Conn) remoteAddr() net.Addr {
	return self.conn.RemoteAddr()
}

func (self *Conn) String() string {
	return fmt.Sprintf("[gate conn] remote=%s, lastTime=%d", self.conn.RemoteAddr().String(), atomic.LoadInt64(&self.lastAt))
}

func (self *Conn) status() int32 {
	return atomic.LoadInt32(&self.state)
}

func (self *Conn) setStatus(state int32) {
	atomic.StoreInt32(&self.state, state)
}

func (self *Conn) serialize(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}
	data, err := self.gateway.serializer.Marshal(v)
	if err != nil {
		return nil, err
	}
	var session = self.session
	if session.getSecret() != "" {
		data, err = crypto.DesEncrypt(data, session.getSecret())
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (self *Conn) recvHandShake(ctx context.Context) error {
	buf := make([]byte, 2048)
	cancelCtx, cancelFunc := context.WithTimeout(ctx, self.gateway.handshakeTimeout)
	defer cancelFunc()

	go func() {
		select {
		case <-cancelCtx.Done():
			if cancelCtx.Err() == context.DeadlineExceeded {
				log.Info("recv handshake timeout", cancelCtx.Err())
				self.conn.Close()
			}
		}
	}()
	for {
		n, err := self.conn.Read(buf)
		if err != nil {
			log.Debugw("conn read fail", "err", err, "session_id", self.session.ID())
			return err
		}
		packets, err := self.decoder.Decode(buf[:n])
		if err != nil {
			log.Info(err.Error())
			return err
		}
		if len(packets) < 1 {
			continue
		}
		if len(packets) != 1 {
			return ErrHandShakeAck
		}
		p := packets[0]
		if p.Type != packet.Handshake {
			return ErrInvalidPacket
		}
		msg := &handShake_request{}
		err = json.Unmarshal(p.Data, msg)
		if err != nil {
			return err
		}
		if self.gateway.rsaPrivateKey != "" {
			token, err := crypto.RsaDecryptWithSha1Base64(msg.Sys.Token, self.gateway.rsaPrivateKey)
			if err != nil {
				return err
			}
			self.session.setSecret(token)
		}
		if err := self.gateway.handshakeValidator(p.Data); err != nil {
			return err
		}
		data, err := json.Marshal(map[string]interface{}{
			"code": 200,
			"sys": map[string]interface{}{
				"heartbeat": self.gateway.heartbeat.Seconds(),
				"session":   self.session.ID(),
			},
		})
		if err != nil {
			return err
		}
		handsharkResponse, err := packet.Encode(packet.Handshake, data)
		if err != nil {
			return err
		}
		if _, err := self.conn.Write(handsharkResponse); err != nil {
			return err
		}
		if self.gateway.debug {
			log.Infow("handshake success", "session_id", self.session.ID(), "remote_addr", self.conn.RemoteAddr(), "secret", self.session.getSecret())
		}
		return nil
	}
}

func (self *Conn) recvHandShakeAck(ctx context.Context) ([]*packet.Packet, error) {
	cancelCtx, cancelFunc := context.WithTimeout(ctx, self.gateway.handshakeTimeout)
	defer cancelFunc()
	go func() {
		select {
		case <-cancelCtx.Done():
			if cancelCtx.Err() == context.DeadlineExceeded {
				log.Info("recv handshake ack timeout", cancelCtx.Err())
				self.conn.Close()
			}
		}
	}()
	buf := make([]byte, 2048)
	for {
		n, err := self.conn.Read(buf)
		if err != nil {
			log.Infow("gate conn read fail", "err", err, "session_id", self.session.ID())
			return nil, err
		}
		packets, err := self.decoder.Decode(buf[:n])
		if err != nil {
			log.Info(err.Error())
			return nil, err
		}
		if len(packets) < 1 {
			continue
		}
		p := packets[0]
		if p.Type != packet.HandshakeAck {
			return nil, ErrInvalidPacket
		}
		self.setStatus(conn_status_working)
		if self.gateway.debug {
			log.Infow("recv handshake ack success", "session_id", self.session.ID(), "remote_addr", self.conn.RemoteAddr())
		}
		return packets[1:], nil
	}
}

func (self *Conn) serve(ctx context.Context, conn net.Conn) error {
	var err error
	var packets []*packet.Packet
	self.conn = conn
	sessionId := self.session.ID()
	if self.gateway.debug {
		log.Infow("session established", "session_id", sessionId)
	}
	defer func() {
		if self.gateway.debug {
			log.Infow("gate conn session goroutine exit", "session_id", sessionId)
		}
		self.setStatus(conn_status_closed)
		self.conn.Close()
	}()
	if self.gateway.handler == nil {
		return ErrInvalidHandler
	}
	self.ctx, self.cancelFunc = context.WithCancel(ctx)
	defer self.cancelFunc()
	if err := self.recvHandShake(self.ctx); err != nil {
		return err
	}
	self.setStatus(conn_status_handshake)
	if packets, err = self.recvHandShakeAck(self.ctx); err != nil {
		return err
	}
	self.setStatus(conn_status_working)

	// 握手成功，开始收发消息
	self.chMessage = make(chan *Message, self.gateway.recvBacklog)
	self.chSend = make(chan []byte, self.gateway.sendBacklog)

	errGroup, errCtx := errgroup.WithContext(self.ctx)
	self.errGroup, self.errCtx = errGroup, errCtx
	self.ctrlCtx, self.ctrlCancelFunc = context.WithCancel(self.ctx)

	//开启读协程
	errGroup.Go(func() (err error) {
		defer func() {
			if err != nil && self.lastErr == nil {
				self.lastErr = err
			}
			if self.gateway.debug {
				log.Infow("gate conn recv goroutine exit", "sessionid", sessionId)
			}
		}()
		// 处理一些多接收到的消息
		for i := range packets {
			if err = self.processPacket(packets[i]); err != nil {
				log.Info(err)
				return
			}
		}
		var n int
		// 持续读数据
		buf := make([]byte, self.gateway.recvBuffSize)
		var packets []*packet.Packet
		for {
			n, err = self.conn.Read(buf)
			if err != nil {
				if self.gateway.debug {
					log.Infow("gate conn read fail", "err", err, "session_id", self.session.ID())
				}
				return
			}
			if self.lastErr != nil {
				return self.lastErr
			}
			packets, err = self.decoder.Decode(buf[:n])
			if err != nil {
				log.Info(err)
				return err
			}
			if len(packets) < 1 {
				continue
			}
			for i := range packets {
				if err = self.processPacket(packets[i]); err != nil {
					log.Info(err)
					return
				}
			}
		}
	})
	//开启写协程
	errGroup.Go(func() (err error) {
		ticker := time.NewTicker(self.gateway.heartbeat)
		// clean func
		defer func() {
			if err != nil && self.lastErr == nil {
				self.lastErr = err
			}
			// 发送完数据后，可以关闭socket了，使recv可以解除阻塞
			self.conn.Close()
			if self.gateway.debug {
				log.Infow("gate conn send goroutine exit", "session_id", sessionId)
			}
		}()
		for {
			select {
			case <-ticker.C:
				deadline := time.Now().Add(-2 * self.gateway.heartbeat).Unix()
				if atomic.LoadInt64(&self.lastAt) < deadline {
					log.Infof("gate connection heartbeat timeout, sessionid=%d, lastTime=%d, deadline=%d\n", sessionId, atomic.LoadInt64(&self.lastAt), deadline)
					err = ErrHeartbeatTimeout
					return
				}
				self.chSend <- heartbeat_packet
			case data := <-self.chSend:
				if data == nil {
					// unexpect
					err = ErrBrokenPipe
					return
				} else {
					if _, err = self.conn.Write(data); err != nil {
						log.Infof("gate connection write failed, sessionid=%d, error:%s\n", sessionId, err.Error())
						return
					}
				}
			case <-errCtx.Done():
				// 关闭心跳
				ticker.Stop()
				close(self.chMessage)
				// 等写完数据再退出协程
				close(self.chSend)
				goto TIMEOUT_SEND
			}
		}
	TIMEOUT_SEND:
		for {
			select {
			case data := <-self.chSend:
				if data == nil {
					err = errCtx.Err()
					return
				} else {
					if _, err = self.conn.Write(data); err != nil {
						log.Infof("gate connection write failed, sessionid=%d, error:%s\n", sessionId, err.Error())
						return err
					}
				}
			}
		}
	})

	atomic.AddInt64(&self.gateway.Stat.ActiveSessionCount, 1)
	atomic.AddInt64(&self.gateway.Stat.CumulativeSessionCount, 1)

	self.gateway.storeSession(self.session)
	defer self.gateway.sessionClosed(self.session)
	// middleware
	for _, middleware := range self.gateway.middlewareArr {
		middleware.OnSessionOpen(self.session)
	}
	defer func() {
		for _, middleware := range self.gateway.middlewareArr {
			middleware.OnSessionClose(self.session)
		}
	}()

	if self.gateway.handler != nil {
		//处理消息
		self.gateway.handler.OnClientStream(self.session)
	}
	// 主动关闭
	if self.errCtx.Err() == nil {
		self.cancelFunc()
	}
	err = errGroup.Wait()
	if self.gateway.debug {
		log.Infow("gate conn wait group exit", "last_error", self.lastErr, "error", err)
	}
	atomic.AddInt64(&self.gateway.Stat.ActiveSessionCount, -1)
	return err
}

func (self *Conn) processPacket(p *packet.Packet) error {
	self.lastAt = time.Now().Unix()
	switch p.Type {
	case packet.Data:
		if self.status() < conn_status_working {
			return ErrNotHandshake
		}
		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		return self.processMessage(msg)
	case packet.Heartbeat:
		// expected
	default:
		return ErrInvalidPacket
	}
	return nil
}

func (self *Conn) processMessage(msg *message.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Info("process message panic", e)
			err = ErrBrokenPipe
		}
	}()
	var reqId uint64
	switch msg.Type {
	case message.Request:
		reqId = msg.ID
	case message.Notify:
		reqId = 0
	default:
		log.Info("[gate conn] invalid message type: " + msg.Type.String())
		err = ErrInvalidMessage
		return
	}
	var session = self.session
	// WARN: 当前data指向缓冲区，要复制出来
	var payload = make([]byte, len(msg.Data))
	copy(payload, msg.Data)
	if session.getSecret() != "" {
		payload, err = crypto.DesDecrypt(payload, session.getSecret())
		if err != nil {
			log.Infof("[gate conn] des decrypt failed, error:%s, payload:(%v)\n", err.Error(), payload)
			return
		}
	}
	if self.gateway.handler == nil {
		log.Infof("[gate conn] handler not found, route:%s\n", msg.Route)
		err = ErrInvalidHandler
	}
	r := &Message{
		session: session,
		route:   msg.Route,
		payload: payload,
		reqId:   reqId,
	}
	// 处理request
	for _, middleware := range self.gateway.middlewareArr {
		middleware.ServeMessage(r)
	}
	self.chMessage <- r
	return
}

func (self *Conn) recv(ctx context.Context) (*Message, error) {
	select {
	case r := <-self.chMessage:
		if self.lastErr != nil {
			return nil, self.lastErr
		}
		if r == nil {
			return nil, ErrBrokenPipe
		}
		return r, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}