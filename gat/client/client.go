package client

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lujingwei/gira"
	"github.com/lujingwei/gira/gat/crypto"
	"github.com/lujingwei/gira/gat/message"
	"github.com/lujingwei/gira/gat/packet"
	"github.com/lujingwei/gira/gat/ws"
)

const (
	_ int32 = iota
	client_status_start
	client_status_handshake
	client_status_working
	client_status_closed
)

const (
	client_write_backlog = 16
)

// Errors that could be occurred during message handling.
var (
	ErrSessionOnNotify    = errors.New("current session working on notify mode")
	ErrCloseClosedSession = errors.New("close closed session")
	ErrInvalidRegisterReq = errors.New("invalid register request")
	// ErrBrokenPipe represents the low-level connection has broken.
	ErrBrokenPipe = errors.New("broken low-level pipe")
	// ErrBufferExceed indicates that the current session buffer is full and
	// can not receive more data.
	ErrBufferExceed       = errors.New("session send buffer exceed")
	ErrCloseClosedGroup   = errors.New("close closed group")
	ErrClosedGroup        = errors.New("group closed")
	ErrMemberNotFound     = errors.New("member not found in the group")
	ErrSessionDuplication = errors.New("session has existed in the current group")
	ErrSprotoRequestType  = errors.New("sproto request type")
	ErrSprotoResponseType = errors.New("sproto response type")
	ErrHandShake          = errors.New("handshake failed")
	ErrHeartbeatTimeout   = errors.New("heartbeat timeout")
	ErrDialTimeout        = errors.New("dial timeout")
	ErrDialInterrupt      = errors.New("dial interrupt")
)

type pendingMessage struct {
	typ     message.Type
	route   string
	mid     uint64
	payload []byte
}

type handShakeResponse struct {
	Sys struct {
		Heartbeat int    `json:"heartbeat"`
		Session   uint64 `json:"session"`
	} `json:"sys"`
	Code int `json:"code"`
}

type ClientConn struct {
	// options
	isWebsocket        bool
	tslCertificate     string
	tslKey             string
	handshakeValidator func([]byte) error
	heartbeat          time.Duration
	debug              bool
	wsPath             string
	rsaPublicKey       string

	serverAddr      string
	heartbeatPacket []byte
	sessionId       uint64
	conn            net.Conn
	decoder         *packet.Decoder
	secretKey       string
	dialTimeout     time.Duration

	responseRouter sync.Map
	lastAt         int64
	chSend         chan pendingMessage
	chWrite        chan []byte
	state          int32
	// 接收到的packet
	chRecvMessage chan *message.Message
	ctx           context.Context
	cancelCtx     context.Context
	cancelFunc    context.CancelFunc
	errCtx        context.Context
	lastErr       error
}

func newClientConn() *ClientConn {
	self := &ClientConn{
		decoder:            packet.NewDecoder(),
		heartbeat:          30 * time.Second,
		debug:              false,
		handshakeValidator: func(_ []byte) error { return nil },
		rsaPublicKey:       "",
		lastAt:             time.Now().Unix(),
		chSend:             make(chan pendingMessage, client_write_backlog),
		chWrite:            make(chan []byte, client_write_backlog),
		state:              client_status_start,

		chRecvMessage: make(chan *message.Message, client_write_backlog),
	}
	return self
}

type opt func(conn *ClientConn)

func WithHandshakeValidator(fn func([]byte) error) opt {
	return func(conn *ClientConn) {
		conn.handshakeValidator = fn
	}
}

func WithHeartbeatInterval(d time.Duration) opt {
	return func(conn *ClientConn) {
		conn.heartbeat = d
	}
}

func WithDictionary(dict map[string]uint16) opt {
	return func(conn *ClientConn) {
	}
}

func WithWSPath(path string) opt {
	return func(conn *ClientConn) {
		conn.wsPath = path
	}
}

func WithServerAdd(addr string) opt {
	return func(conn *ClientConn) {
		conn.serverAddr = addr
	}
}

func WithTSLConfig(certificate, key string) opt {
	return func(conn *ClientConn) {
		conn.tslCertificate = certificate
		conn.tslKey = key
	}
}

func WithDialTimeout(timeout time.Duration) opt {
	return func(conn *ClientConn) {
		conn.dialTimeout = timeout
	}
}

func WithContext(ctx context.Context) opt {
	return func(conn *ClientConn) {
		conn.ctx = ctx
	}
}

func WithIsWebsocket(enableWs bool) opt {
	return func(conn *ClientConn) {
		conn.isWebsocket = enableWs
	}
}
func WithDebugMode() opt {
	return func(conn *ClientConn) {
		conn.debug = true
	}
}

func WithRSAPublicKey(keyFile string) opt {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	return func(conn *ClientConn) {
		conn.rsaPublicKey = string(data)
	}
}

func Dial(addr string, opts ...opt) (gira.GateClient, error) {
	conn := newClientConn()
	for _, v := range opts {
		v(conn)
	}
	if err := conn.dial(addr); err != nil {
		return nil, err
	}
	return conn, nil
}

func (self *ClientConn) dial(addr string) error {
	var err error
	heartbeatPacket, err := packet.Encode(packet.Heartbeat, nil)
	if err != nil {
		return err
	}
	self.heartbeatPacket = heartbeatPacket
	self.serverAddr = addr
	if self.ctx == nil {
		self.ctx = context.Background()
	}
	self.cancelCtx, self.cancelFunc = context.WithCancel(self.ctx)
	var conn net.Conn
	if self.isWebsocket {
		if len(self.tslCertificate) != 0 {
			//self.dialWSTLS()
			// TODO
			return errors.New("tls certificate not implement")
		} else {
			if conn, err = self.dialWS(); err != nil {
				return err
			}
		}
	} else {
		if conn, err = self.dialTcp(); err != nil {
			return err
		}
	}
	self.conn = conn
	self.setStatus(client_status_handshake)
	err = self.sendHandShake()
	if err != nil {
		log.Println(fmt.Sprintf("[agent] client sendHandShake error: %s", err.Error()))
		conn.Close()
		return err
	}
	if err := self.recvHandShakeAck(self.ctx); err != nil {
		conn.Close()
		return err
	}
	self.setStatus(client_status_working)
	go self.readRoutine()
	go self.writeRoutine()
	return nil
}

func (self *ClientConn) dialTcp() (net.Conn, error) {
	netDialer := net.Dialer{}
	if self.dialTimeout != 0 {
		netDialer.Timeout = self.dialTimeout
	}
	conn, err := netDialer.DialContext(self.cancelCtx, "tcp", self.serverAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (self *ClientConn) dialWS() (net.Conn, error) {
	netDialer := &net.Dialer{}
	if self.dialTimeout != 0 {
		netDialer.Timeout = self.dialTimeout
	}
	// websocket.DefaultDialer
	d := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		NetDial: func(network, addr string) (net.Conn, error) {
			return netDialer.DialContext(self.cancelCtx, network, addr)
		},
	}
	u := url.URL{Scheme: "ws", Host: self.serverAddr, Path: self.wsPath}
	conn, _, err := d.DialContext(self.cancelCtx, u.String(), nil)
	if err != nil {
		return nil, err
	}
	wsconn, err := ws.NewConn(conn)
	if err != nil {
		return nil, err
	}
	return wsconn, nil
}

func (self *ClientConn) setSecretKey(key string) {
	self.secretKey = key
}

func (self *ClientConn) getSecretKey() string {
	return self.secretKey
}

func (self *ClientConn) sendHandShake() error {
	tokenByte := make([]byte, 8)
	_, err := rand.Read(tokenByte)
	if err != nil {
		return err
	}
	token := ""
	if self.rsaPublicKey != "" {
		token = base32.StdEncoding.EncodeToString(tokenByte)[:8]
		self.setSecretKey(token)
		token, err = crypto.RsaEncryptWithSha1Base64(token, self.rsaPublicKey)
		if err != nil {
			return err
		}
	}
	payload, err := json.Marshal(map[string]interface{}{
		"sys": map[string]interface{}{
			"type":    "go-websocket",
			"version": "0.0.1",
			"token":   token,
		},
		"user": map[string]interface{}{},
	})
	if err != nil {
		return err
	}
	data, err := packet.Encode(packet.Handshake, payload)
	if err != nil {
		return err
	}
	if _, err := self.conn.Write(data); err != nil {
		return err
	}
	if self.debug {
		gira.Infow("send handshake", "token", token)
	}
	return nil
}

func (self *ClientConn) recvHandShakeAck(ctx context.Context) error {
	buf := make([]byte, 2048)
	for {
		var n int
		var err error
		n, err = self.conn.Read(buf)
		if err != nil {
			if self.debug {
				log.Println(fmt.Sprintf("[agent] read message error: %s, session will be closed immediately, sessionid=%d", err.Error(), self.sessionId))
			}
			return err
		}
		packets, err := self.decoder.Decode(buf[:n])
		if err != nil {
			log.Println(err)
			return err
		}
		if len(packets) < 1 {
			// TODO 增加timeout功能
			continue
		}
		if len(packets) != 1 {
			return errors.New("unexpect packet count")
		}
		p := packets[0]
		if p.Type != packet.Handshake {
			return errors.New("expect handshake packet")
		}
		msg := &handShakeResponse{}
		err = json.Unmarshal(p.Data, msg)
		if err != nil {
			return err
		}
		if msg.Code != 200 {
			return ErrHandShake
		}
		payload, err := json.Marshal(map[string]interface{}{})
		if err != nil {
			return err
		}
		data, err := packet.Encode(packet.HandshakeAck, payload)
		if err != nil {
			return err
		}
		self.sessionId = msg.Sys.Session
		if self.debug {
			gira.Infow("handshake success", "session_id", self.sessionId, "remote_addr", self.conn.RemoteAddr())
		}
		if _, err := self.conn.Write(data); err != nil {
			return err
		}
		return nil
	}
}

func (self *ClientConn) writeRoutine() error {
	ticker := time.NewTicker(self.heartbeat)
	defer func() {
		ticker.Stop()
		close(self.chSend)
		close(self.chWrite)
		if self.debug {
			gira.Infow("write goroutine exit", "session_id", self.sessionId)
		}
	}()
	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Add(-2 * self.heartbeat).Unix()
			if atomic.LoadInt64(&self.lastAt) < deadline {
				log.Println(fmt.Sprintf("[agent] session heartbeat timeout, sessionid=%d, lastTime=%d, deadline=%d",
					self.sessionId, atomic.LoadInt64(&self.lastAt), deadline))
				return ErrHeartbeatTimeout
			}
		case data := <-self.chWrite:
			if _, err := self.conn.Write(data); err != nil {
				log.Println(fmt.Sprintf("[agent] conn write failed, error:%s", err.Error()))
				return err
			}
		case data := <-self.chSend:
			m := &message.Message{
				Type:  data.typ,
				Data:  data.payload,
				Route: data.route,
				ID:    data.mid,
			}
			em, err := m.Encode()
			if err != nil {
				log.Println(err.Error())
				break
			}
			p, err := packet.Encode(packet.Data, em)
			if err != nil {
				log.Println(err)
				break
			}
			self.chWrite <- p
		case <-self.cancelCtx.Done():
			return self.cancelCtx.Err()
		}
	}
}

func (self *ClientConn) readRoutine() error {
	defer func() {
		close(self.chRecvMessage)
		if self.debug {
			gira.Infow("read goroutine exit", "session_id", self.sessionId)
		}
	}()
	buf := make([]byte, 2048)
	for {
		n, err := self.conn.Read(buf)
		if err != nil {
			if self.debug {
				gira.Infow("client read fail", "err", err, "session_id", self.sessionId)
			}
			return err
		}
		packets, err := self.decoder.Decode(buf[:n])
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if len(packets) < 1 {
			continue
		}
		for _, p := range packets {
			if err := self.processPacket(p); err != nil {
				log.Println(err.Error())
				return err
			}
		}
	}
}

func (self *ClientConn) processMessage(msg *message.Message) {
	switch msg.Type {
	case message.Response:
	case message.Push:
	default:
		log.Println("Invalid message type: " + msg.Type.String())
		return
	}
	if self.getSecretKey() != "" {
		payload, err := crypto.DesDecrypt(msg.Data, self.getSecretKey())
		if err != nil {
			log.Println(fmt.Sprintf("crypto.DesDecrypt failed: %+v (%v)", err, payload))
			return
		}
		msg.Data = payload
	}
	self.chRecvMessage <- msg
}

// 处理内部消息包
func (self *ClientConn) processPacket(p *packet.Packet) error {
	if self.debug {
		log.Println("got packet", p.Type)
	}
	switch p.Type {
	case packet.Data:
		if self.status() < client_status_working {
			return fmt.Errorf("receive data on socket which not yet ACK, session will be closed immediately, sessionid=%d, remote=%s, status=%d",
				self.sessionId, self.conn.RemoteAddr().String(), self.status())
		}
		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		self.processMessage(msg)
	case packet.Heartbeat:
		self.chWrite <- self.heartbeatPacket
	}
	self.lastAt = time.Now().Unix()
	return nil
}

func (self *ClientConn) status() int32 {
	return atomic.LoadInt32(&self.state)
}

func (self *ClientConn) setStatus(state int32) {
	atomic.StoreInt32(&self.state, state)
}

func (self *ClientConn) Recv(ctx context.Context) (typ int, route string, reqId uint64, data []byte, err error) {
	select {
	case msg := <-self.chRecvMessage:
		if msg == nil {
			err = gira.ErrReadOnClosedClient
			return
		}
		typ = int(msg.Type)
		switch msg.Type {
		case message.Response:
			reqId = msg.ID
			if s, ok := self.responseRouter.Load(reqId); ok {
				route = s.(string)
				self.responseRouter.Delete(reqId)
			}
		case message.Push:
			reqId = 0
		}
		data = msg.Data
		return
	case <-ctx.Done():
		err = ctx.Err()
		return
	}
}

func (self *ClientConn) send(typ message.Type, reqId uint64, route string, data []byte) error {
	if self.status() == client_status_closed {
		return ErrBrokenPipe
	}
	if len(self.chSend) >= client_write_backlog {
		return ErrBufferExceed
	}
	if self.debug {
		gira.Infow("send", "type", typ, "sessioni_id", self.sessionId, "req_id", reqId, "route", route, "len", len(data))
	}
	if typ == message.Request {
		self.responseRouter.Store(reqId, route)
	}
	m := &message.Message{
		Type:  typ,
		Data:  data,
		Route: route,
		ID:    reqId,
	}
	em, err := m.Encode()
	if err != nil {
		return err
	}
	p, err := packet.Encode(packet.Data, em)
	if err != nil {
		return err
	}
	if _, err := self.conn.Write(p); err != nil {
		return err
	}
	return nil
}

func (self *ClientConn) Notify(route string, data []byte) error {
	return self.send(message.Notify, 0, route, data)
}

func (self *ClientConn) Request(route string, reqId uint64, data []byte) error {
	return self.send(message.Request, reqId, route, data)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (self *ClientConn) Close() error {
	if self.status() == client_status_closed {
		return ErrCloseClosedSession
	} else {
		self.setStatus(client_status_closed)
		if self.conn != nil {
			self.conn.Close()
		}
		select {
		case <-self.cancelCtx.Done():
			// expect
		default:
			self.cancelFunc()
		}
		return nil
	}
}
