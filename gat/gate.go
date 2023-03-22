package gat

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lujingwei/gira/log"

	"github.com/gorilla/websocket"
	"github.com/lujingwei/gira"
	"github.com/lujingwei/gira/gat/packet"
	"github.com/lujingwei/gira/gat/serialize"
	"github.com/lujingwei/gira/gat/serialize/protobuf"
	"github.com/lujingwei/gira/gat/ws"
	"golang.org/x/sync/errgroup"
)

const (
	_ int32 = iota
	gate_status_start
	gate_status_maintain
	gate_status_working
	gate_status_closed
)

var heartbeat_packet []byte

func init() {
	var err error
	heartbeat_packet, err = packet.Encode(packet.Heartbeat, nil)
	if err != nil {
		panic(err)
	}
}

type Gate struct {
	ClientAddr string
	Host       string
	Port       int32

	// options
	ctx                context.Context
	cancelFunc         context.CancelFunc
	errCtx             context.Context
	errGroup           *errgroup.Group
	isWebsocket        bool
	tslCertificate     string
	tslKey             string
	handshakeValidator func([]byte) error
	serializer         serialize.Serializer
	heartbeat          time.Duration
	checkOrigin        func(*http.Request) bool
	debug              bool
	wsPath             string
	rsaPrivateKey      string
	sessionModifer     uint64
	handshakeTimeout   time.Duration
	sendBacklog        int
	recvBacklog        int
	recvBuffSize       int

	// 状态
	state         int32
	startAt       time.Time
	mu            sync.RWMutex
	sessions      map[uint64]*Session
	handler       gira.GateHandler
	middlewareArr []MiddleWareInterface
	listener      net.Listener
	server        *http.Server
	Stat          Stat
}

type Stat struct {
	ActiveSessionCount        int64 // 当前会话数量
	CumulativeSessionCount    int64 // 累计会话数量
	CumulativeConnectionCount int64 // 累计连接数量
}

func newGate() *Gate {
	gate := &Gate{
		state:              gate_status_start,
		heartbeat:          30 * time.Second,
		debug:              false,
		checkOrigin:        func(_ *http.Request) bool { return true },
		handshakeValidator: func(_ []byte) error { return nil },
		serializer:         protobuf.NewSerializer(),
		middlewareArr:      make([]MiddleWareInterface, 0),
		sessions:           map[uint64]*Session{},
		handshakeTimeout:   2 * time.Second,
		sendBacklog:        16,
		recvBacklog:        16,
		recvBuffSize:       4096,
	}
	return gate
}

func NewConfigGate(facade gira.ApplicationFacade, config gira.GateConfig) (*Gate, error) {
	var handler gira.GateHandler
	var ok bool
	if handler, ok = facade.(gira.GateHandler); !ok {
		return nil, gira.ErrGateHandlerNotImplement
	}
	opts := []Option{
		WithDebugMode(config.Debug),
		WithIsWebsocket(true),
		WithSessionModifer(uint64(facade.GetId()) << 48),
	}
	var server *Gate
	var err error
	if server, err = Listen(facade.Context(), config.Bind, opts...); err != nil {
		return nil, err
	}
	log.Info("gate start...", config.Bind)
	go server.Serve(handler)
	return server, nil
}

type Option func(gate *Gate)

func WithSessionModifer(v uint64) Option {
	return func(gate *Gate) {
		gate.sessionModifer = v
	}
}

func WithHandshakeValidator(fn func([]byte) error) Option {
	return func(gate *Gate) {
		gate.handshakeValidator = fn
	}
}

func WithSerializer(serializer serialize.Serializer) Option {
	return func(gate *Gate) {
		gate.serializer = serializer
	}
}

func WithDebugMode(v bool) Option {
	return func(gate *Gate) {
		gate.debug = v
	}
}

func WithCheckOriginFunc(fn func(*http.Request) bool) Option {
	return func(gate *Gate) {
		gate.checkOrigin = fn
	}
}

func WithHeartbeatInterval(d time.Duration) Option {
	return func(gate *Gate) {
		gate.heartbeat = d
	}
}

func WithRecvBuffSize(v int) Option {
	return func(gate *Gate) {
		gate.recvBuffSize = v
	}
}

func WithSendBacklog(v int) Option {
	return func(gate *Gate) {
		gate.sendBacklog = v
	}
}

func WithRecvBacklog(v int) Option {
	return func(gate *Gate) {
		gate.recvBacklog = v
	}
}

func WithDictionary(dict map[string]uint16) Option {
	return func(gate *Gate) {
	}
}

func WithWSPath(path string) Option {
	return func(gate *Gate) {
		gate.wsPath = path
	}
}

func WithIsWebsocket(enableWs bool) Option {
	return func(gate *Gate) {
		gate.isWebsocket = enableWs
	}
}

func WithTSLConfig(certificate, key string) Option {
	return func(gate *Gate) {
		gate.tslCertificate = certificate
		gate.tslKey = key
	}
}

func WithRSAPrivateKey(keyFile string) Option {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Info(err)
		return nil
	}
	return func(gate *Gate) {
		gate.rsaPrivateKey = string(data)
	}
}

func Listen(ctx context.Context, addr string, opts ...Option) (*Gate, error) {
	gate := newGate()
	for _, opt := range opts {
		if opt != nil {
			opt(gate)
		}
	}
	gate.startAt = time.Now()
	if gate.ctx == nil {
		gate.ctx = context.Background()
	}
	gate.ctx, gate.cancelFunc = context.WithCancel(ctx)
	gate.errGroup, gate.errCtx = errgroup.WithContext(gate.ctx)
	addrPat := strings.SplitN(addr, ":", 2)
	if gate.isWebsocket {
		if len(gate.tslCertificate) != 0 {
			gate.Host = fmt.Sprintf("wss://%s", addrPat[0])
		} else {
			gate.Host = fmt.Sprintf("ws://%s", addrPat[0])
		}
	} else {
		gate.Host = addrPat[0]
	}
	if port, err := strconv.Atoi(addrPat[1]); err != nil {
		return nil, ErrInvalidAddress
	} else {
		gate.Port = int32(port)
	}
	if len(addrPat) < 2 {
		return nil, ErrInvalidAddress
	}
	gate.ClientAddr = fmt.Sprintf(":%d", gate.Port)
	return gate, nil
}

func (gate *Gate) Serve(h gira.GateHandler) error {
	gate.handler = h
	gate.setStatus(gate_status_working)
	if gate.isWebsocket {
		if len(gate.tslCertificate) != 0 {
			return gate.listenAndServeWSTLS()
		} else {
			return gate.listenAndServeWS()
		}
	} else {
		return gate.listenAndServe()
	}
}

func (gate *Gate) UseMiddleware(m MiddleWareInterface) {
	gate.middlewareArr = append(gate.middlewareArr, m)
}

func (gate *Gate) NewGroup(name string) *Group {
	return newGroup(gate, name)
}

func (gate *Gate) Maintain(m bool) {
	status := gate.status()
	if status != gate_status_maintain && status != gate_status_working {
		return
	}
	if m {
		gate.setStatus(gate_status_maintain)
	} else {
		gate.setStatus(gate_status_working)
	}
}

func (gate *Gate) Shutdown() {
	if gate.status() == gate_status_closed {
		return
	}
	gate.setStatus(gate_status_closed)
	if gate.server != nil {
		timeoutCtx, timeoutFunc := context.WithTimeout(gate.ctx, 5*time.Second)
		defer timeoutFunc()
		// shutdown 会关闭listener,拒绝新的连接， 但并不会关闭websocket连接
		gate.server.Shutdown(timeoutCtx)
	} else if gate.listener != nil {
		gate.listener.Close()
	}
	gate.Kick("shutdown")
	gate.cancelFunc()
	var err error
	err = gate.errGroup.Wait()
	log.Infow("gate shutdown", "error", err)
}

func (gate *Gate) status() int32 {
	return atomic.LoadInt32(&gate.state)
}

func (gate *Gate) setStatus(state int32) {
	atomic.StoreInt32(&gate.state, state)
}

func (gate *Gate) Kick(reason string) {
	//断开已有的链接
	gate.mu.RLock()
	for _, s := range gate.sessions {
		s.Kick(reason)
	}
	gate.mu.RUnlock()
	now := time.Now().Unix()
	for {
		if len(gate.sessions) <= 0 {
			break
		}
		if time.Now().Unix()-now >= 120 {
			log.Info(fmt.Sprintf("[gate] some session not close, count=%d", len(gate.sessions)))
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (gate *Gate) listenAndServe() error {
	listener, err := net.Listen("tcp", gate.ClientAddr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	gate.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Info(err)
			return err
		}
		go gate.handleConn(conn)
	}
}

func (gate *Gate) listenAndServeWS() error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     gate.checkOrigin,
	}
	http.HandleFunc("/"+strings.TrimPrefix(gate.wsPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		log.Infow("gate listen and server ws")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Infof("[gate] Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
			return
		}
		gate.serveWsConn(conn)
	})
	server := &http.Server{Addr: gate.ClientAddr}
	gate.server = server
	if err := server.ListenAndServe(); err != nil {
		log.Info(err)
		return err
	}
	return nil
}

func (gate *Gate) listenAndServeWSTLS() error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     gate.checkOrigin,
	}
	http.HandleFunc("/"+strings.TrimPrefix(gate.wsPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Infof("[gate] Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
			return
		}
		gate.serveWsConn(conn)
	})
	if err := http.ListenAndServeTLS(gate.ClientAddr, gate.tslCertificate, gate.tslKey, nil); err != nil {
		log.Info(err)
		return err
	}
	return nil
}

func (gate *Gate) findSession(sid uint64) *Session {
	gate.mu.RLock()
	s := gate.sessions[sid]
	gate.mu.RUnlock()
	return s
}

func (gate *Gate) handleConn(conn net.Conn) {
	if gate.status() != gate_status_working {
		log.Info("[gate] gate is not running")
		conn.Close()
		return
	}
	atomic.AddInt64(&gate.Stat.CumulativeConnectionCount, 1)
	c := newConn(gate)
	if gate.debug {
		log.Infow("accept a client", "session_id", c.session.ID(), "remote_addr", conn.RemoteAddr())
	}
	err := c.serve(gate.ctx, conn)
	if gate.debug {
		log.Infow("gate conn serve exit", "error", err)
	}
}

func (gate *Gate) serveWsConn(conn *websocket.Conn) {
	if gate.status() != gate_status_working {
		log.Info("[gate] gate is not running")
		conn.Close()
		return
	}
	c, err := ws.NewConn(conn)
	if err != nil {
		log.Info(err)
		return
	}
	gate.errGroup.Go(func() error {
		gate.handleConn(c)
		return nil
	})
}

func (gate *Gate) storeSession(s *Session) {
	gate.mu.Lock()
	gate.sessions[s.ID()] = s
	gate.mu.Unlock()
}

func (gate *Gate) sessionClosed(s *Session) error {
	gate.mu.Lock()
	delete(gate.sessions, s.ID())
	gate.mu.Unlock()
	return nil
}
