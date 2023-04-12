package gate

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

	"github.com/lujingwei002/gira/log"

	"github.com/gorilla/websocket"
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/gate/packet"
	"github.com/lujingwei002/gira/gate/serialize"
	"github.com/lujingwei002/gira/gate/serialize/protobuf"
	"github.com/lujingwei002/gira/gate/ws"
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

type Gateway struct {
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
	handler       gira.GatewayHandler
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

func newGateway() *Gateway {
	gate := &Gateway{
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

func NewConfigGateway(facade gira.Application, handler gira.GatewayHandler, config gira.GatewayConfig) (*Gateway, error) {
	opts := []Option{
		WithDebugMode(config.Debug),
		WithIsWebsocket(true),
		WithSessionModifer(uint64(facade.GetAppId()) << 48),
	}
	var server *Gateway
	var err error
	if server, err = Listen(facade.Context(), config.Bind, opts...); err != nil {
		return nil, err
	}
	log.Info("gate start...", config.Bind)
	go server.Serve(handler)
	return server, nil
}

type Option func(gateway *Gateway)

func WithSessionModifer(v uint64) Option {
	return func(gateway *Gateway) {
		gateway.sessionModifer = v
	}
}

func WithHandshakeValidator(fn func([]byte) error) Option {
	return func(gateway *Gateway) {
		gateway.handshakeValidator = fn
	}
}

func WithSerializer(serializer serialize.Serializer) Option {
	return func(gateway *Gateway) {
		gateway.serializer = serializer
	}
}

func WithDebugMode(v bool) Option {
	return func(gateway *Gateway) {
		gateway.debug = v
	}
}

func WithCheckOriginFunc(fn func(*http.Request) bool) Option {
	return func(gateway *Gateway) {
		gateway.checkOrigin = fn
	}
}

func WithHeartbeatInterval(d time.Duration) Option {
	return func(gateway *Gateway) {
		gateway.heartbeat = d
	}
}

func WithRecvBuffSize(v int) Option {
	return func(gateway *Gateway) {
		gateway.recvBuffSize = v
	}
}

func WithSendBacklog(v int) Option {
	return func(gateway *Gateway) {
		gateway.sendBacklog = v
	}
}

func WithRecvBacklog(v int) Option {
	return func(gateway *Gateway) {
		gateway.recvBacklog = v
	}
}

func WithDictionary(dict map[string]uint16) Option {
	return func(gateway *Gateway) {
	}
}

func WithWSPath(path string) Option {
	return func(gateway *Gateway) {
		gateway.wsPath = path
	}
}

func WithIsWebsocket(enableWs bool) Option {
	return func(gateway *Gateway) {
		gateway.isWebsocket = enableWs
	}
}

func WithTSLConfig(certificate, key string) Option {
	return func(gateway *Gateway) {
		gateway.tslCertificate = certificate
		gateway.tslKey = key
	}
}

func WithRSAPrivateKey(keyFile string) Option {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Info(err)
		return nil
	}
	return func(gateway *Gateway) {
		gateway.rsaPrivateKey = string(data)
	}
}

func Listen(ctx context.Context, addr string, opts ...Option) (*Gateway, error) {
	gateway := newGateway()
	for _, opt := range opts {
		if opt != nil {
			opt(gateway)
		}
	}
	gateway.startAt = time.Now()
	if gateway.ctx == nil {
		gateway.ctx = context.Background()
	}
	gateway.ctx, gateway.cancelFunc = context.WithCancel(ctx)
	gateway.errGroup, gateway.errCtx = errgroup.WithContext(gateway.ctx)
	addrPat := strings.SplitN(addr, ":", 2)
	if gateway.isWebsocket {
		if len(gateway.tslCertificate) != 0 {
			gateway.Host = fmt.Sprintf("wss://%s", addrPat[0])
		} else {
			gateway.Host = fmt.Sprintf("ws://%s", addrPat[0])
		}
	} else {
		gateway.Host = addrPat[0]
	}
	if port, err := strconv.Atoi(addrPat[1]); err != nil {
		return nil, ErrInvalidAddress
	} else {
		gateway.Port = int32(port)
	}
	if len(addrPat) < 2 {
		return nil, ErrInvalidAddress
	}
	gateway.ClientAddr = fmt.Sprintf(":%d", gateway.Port)
	return gateway, nil
}

func (gateway *Gateway) Serve(h gira.GatewayHandler) error {
	gateway.handler = h
	gateway.setStatus(gate_status_working)
	if gateway.isWebsocket {
		if len(gateway.tslCertificate) != 0 {
			return gateway.listenAndServeWSTLS()
		} else {
			return gateway.listenAndServeWS()
		}
	} else {
		return gateway.listenAndServe()
	}
}

func (gateway *Gateway) UseMiddleware(m MiddleWareInterface) {
	gateway.middlewareArr = append(gateway.middlewareArr, m)
}

func (gateway *Gateway) NewGroup(name string) *Group {
	return newGroup(gateway, name)
}

func (gateway *Gateway) Maintain(m bool) {
	status := gateway.status()
	if status != gate_status_maintain && status != gate_status_working {
		return
	}
	if m {
		gateway.setStatus(gate_status_maintain)
	} else {
		gateway.setStatus(gate_status_working)
	}
}

func (gateway *Gateway) Shutdown() {
	if gateway.status() == gate_status_closed {
		return
	}
	gateway.setStatus(gate_status_closed)
	if gateway.server != nil {
		timeoutCtx, timeoutFunc := context.WithTimeout(gateway.ctx, 5*time.Second)
		defer timeoutFunc()
		// shutdown 会关闭listener,拒绝新的连接， 但并不会关闭websocket连接
		gateway.server.Shutdown(timeoutCtx)
	} else if gateway.listener != nil {
		gateway.listener.Close()
	}
	gateway.Kick("shutdown")
	gateway.cancelFunc()
	var err error
	err = gateway.errGroup.Wait()
	log.Infow("gate shutdown", "error", err)
}

func (gateway *Gateway) status() int32 {
	return atomic.LoadInt32(&gateway.state)
}

func (gateway *Gateway) setStatus(state int32) {
	atomic.StoreInt32(&gateway.state, state)
}

func (gateway *Gateway) Kick(reason string) {
	//断开已有的链接
	gateway.mu.RLock()
	for _, s := range gateway.sessions {
		s.Kick(reason)
	}
	gateway.mu.RUnlock()
	now := time.Now().Unix()
	for {
		if len(gateway.sessions) <= 0 {
			break
		}
		if time.Now().Unix()-now >= 120 {
			log.Info(fmt.Sprintf("[gate] some session not close, count=%d", len(gateway.sessions)))
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (gateway *Gateway) listenAndServe() error {
	listener, err := net.Listen("tcp", gateway.ClientAddr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	gateway.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Info(err)
			return err
		}
		go gateway.handleConn(conn)
	}
}

func (gateway *Gateway) listenAndServeWS() error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     gateway.checkOrigin,
	}
	http.HandleFunc("/"+strings.TrimPrefix(gateway.wsPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		log.Infow("gate listen and server ws")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Infof("[gate] Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
			return
		}
		gateway.serveWsConn(conn)
	})
	server := &http.Server{Addr: gateway.ClientAddr}
	gateway.server = server
	if err := server.ListenAndServe(); err != nil {
		log.Info(err)
		return err
	}
	return nil
}

func (gateway *Gateway) listenAndServeWSTLS() error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     gateway.checkOrigin,
	}
	http.HandleFunc("/"+strings.TrimPrefix(gateway.wsPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Infof("[gate] Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
			return
		}
		gateway.serveWsConn(conn)
	})
	if err := http.ListenAndServeTLS(gateway.ClientAddr, gateway.tslCertificate, gateway.tslKey, nil); err != nil {
		log.Info(err)
		return err
	}
	return nil
}

func (gateway *Gateway) findSession(sid uint64) *Session {
	gateway.mu.RLock()
	s := gateway.sessions[sid]
	gateway.mu.RUnlock()
	return s
}

func (gateway *Gateway) handleConn(conn net.Conn) {
	if gateway.status() != gate_status_working {
		log.Info("[gate] gate is not running")
		conn.Close()
		return
	}
	atomic.AddInt64(&gateway.Stat.CumulativeConnectionCount, 1)
	c := newConn(gateway)
	if gateway.debug {
		log.Infow("accept a client", "session_id", c.session.ID(), "remote_addr", conn.RemoteAddr())
	}
	err := c.serve(gateway.ctx, conn)
	if gateway.debug {
		log.Infow("gate conn serve exit", "error", err)
	}
}

func (gateway *Gateway) serveWsConn(conn *websocket.Conn) {
	if gateway.status() != gate_status_working {
		log.Info("[gate] gate is not running")
		conn.Close()
		return
	}
	c, err := ws.NewConn(conn)
	if err != nil {
		log.Info(err)
		return
	}
	gateway.errGroup.Go(func() error {
		gateway.handleConn(c)
		return nil
	})
}

func (gateway *Gateway) storeSession(s *Session) {
	gateway.mu.Lock()
	gateway.sessions[s.ID()] = s
	gateway.mu.Unlock()
}

func (gateway *Gateway) sessionClosed(s *Session) error {
	gateway.mu.Lock()
	delete(gateway.sessions, s.ID())
	gateway.mu.Unlock()
	return nil
}
