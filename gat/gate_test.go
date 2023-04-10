package gat

import (
	"context"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/gat/client"
	"golang.org/x/sync/errgroup"
)

func setup() {

}
func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	err := m.Run()
	teardown()
	os.Exit(err)
}

type GateHandler struct {
}

func (self *GateHandler) OnClientStream(s gira.GateConn) {
	//var req gira.GateRequest
	var err error
	for {
		_, err = s.Recv(context.TODO())
		if err != nil {
			//log.Infow("recv", "error", err)
			break
		} else {
			//	log.Infow("recv", "data", string(req.Payload()))
		}
	}
}

// 模拟很多个客户端连接
func TestManyClient(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234",
		//WithDebugMode(),
		WithIsWebsocket(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	errGroup, errCtx := errgroup.WithContext(context.TODO())
	errGroup.Go(func() error {
		go func() {
			gate.Serve(&GateHandler{})
		}()
		select {
		case <-errCtx.Done():
			time.Sleep(1 * time.Second)
			log.Infow("gate stat", "ActiveSessionCount", gate.Stat.ActiveSessionCount)
			log.Infow("gate stat", "CumulativeSessionCount", gate.Stat.CumulativeSessionCount)
			log.Infow("gate stat", "CumulativeConnectionCount", gate.Stat.CumulativeConnectionCount)
			gate.Shutdown()
		}
		return nil
	})
	for i := 0; i < 100; i++ {
		errGroup.Go(func() error {
			time.Sleep(1 * time.Millisecond)
			switch rand.Intn(8) {
			case 0:
				// 不握手，不关连接
				netDialer := &net.Dialer{}
				netDialer.Timeout = 10 * time.Second
				// websocket.DefaultDialer
				d := &websocket.Dialer{
					Proxy:            http.ProxyFromEnvironment,
					HandshakeTimeout: 45 * time.Second,
					NetDial: func(network, addr string) (net.Conn, error) {
						return netDialer.DialContext(context.TODO(), network, addr)
					},
				}
				u := url.URL{Scheme: "ws", Host: "127.0.0.1:1234", Path: "/"}
				conn, _, err := d.DialContext(context.TODO(), u.String(), nil)
				if err != nil {
					t.Fatal(err)
				}
				t, n, err := conn.ReadMessage()
				log.Info(t, n, err)
			default:
				var conn gira.GateClient
				conn, err = client.Dial("127.0.0.1:1234",
					//client.WithDebugMode(),
					client.WithIsWebsocket(true),
				)
				if err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 2; i++ {
					err = conn.Request("", 1, []byte("hello"))
					if err != nil {
						t.Fatal(err)
					}
					time.Sleep(1 * time.Second)
				}
				conn.Close()
			}

			return io.EOF
		})
	}
	err = errGroup.Wait()
	if err != nil {
		log.Infow("errGroup", "error", err)
	}
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	t.Logf("Allocated memory: %v bytes", mem.Alloc)
	return
}

type GateHandler1 struct {
}

func (self *GateHandler1) OnClientStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	for {
		req, err = s.Recv(context.TODO())
		if err != nil {
			log.Infow("recv", "error", err)
			break
		} else {
			log.Infow("recv", "data", string(req.Payload()))
		}
	}
}

// 客户端连接上服务器后，每隔1秒发送1次消息，然后主动关闭
func TestClientClose1(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode(true), WithIsWebsocket(true))
	if err != nil {
		t.Fatal(err)
	}
	errGroup, errCtx := errgroup.WithContext(context.TODO())
	errGroup.Go(func() error {
		go func() {
			gate.Serve(&GateHandler1{})
		}()
		select {
		case <-errCtx.Done():
			time.Sleep(1 * time.Second)
			log.Info("aaaaaaaaaa", gate.Stat.ActiveSessionCount)
			log.Info("aaaaaaaaaa", gate.Stat.CumulativeSessionCount)
			gate.Shutdown()
		}
		return nil
	})
	errGroup.Go(func() error {
		time.Sleep(1 * time.Millisecond)
		var conn gira.GateClient
		conn, err = client.Dial("127.0.0.1:1234", client.WithDebugMode(), client.WithIsWebsocket(true))
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 2; i++ {
			err = conn.Request("", 1, []byte("hello"))
			if err != nil {
				t.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
		conn.Close()
		return io.EOF
	})
	err = errGroup.Wait()
	if err != nil {
		log.Infow("errGroup", "error", err)
	}
	return
}

type GateHandler2 struct {
}

func (self *GateHandler2) OnClientStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	req, err = s.Recv(context.TODO())
	if err != nil {
		log.Infow("recv", "error", err)

	} else {
		log.Infow("recv", "data", string(req.Payload()))

	}
	s.Close()
}

// 客户端连接上服务器后，每隔1秒发送1次消息，然后由服务器主动关闭链接
func TestClientClose2(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode(true), WithIsWebsocket(true))
	if err != nil {
		t.Fatal(err)
	}
	errGroup, errCtx := errgroup.WithContext(context.TODO())
	errGroup.Go(func() error {
		go func() {
			gate.Serve(&GateHandler2{})
		}()
		select {
		case <-errCtx.Done():
			gate.Shutdown()
		}
		return nil
	})
	errGroup.Go(func() error {
		time.Sleep(1 * time.Millisecond)
		var conn gira.GateClient
		conn, err = client.Dial("127.0.0.1:1234", client.WithDebugMode(), client.WithIsWebsocket(true))
		if err != nil {
			t.Fatal(err)
		}
		for {
			err = conn.Request("", 1, []byte("hello"))
			if err != nil {
				log.Infow("request", "error", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
		conn.Close()
		return io.EOF
	})
	err = errGroup.Wait()
	if err != nil {
		log.Infow("errGroup", "error", err)

	}
	return
}

type GateHandler3 struct {
}

func (self *GateHandler3) OnClientStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	req, err = s.Recv(context.TODO())
	if err != nil {
		log.Infow("recv", "error", err)
	} else {
		log.Infow("recv", "data", string(req.Payload()))
	}
}

// 客户端连接上服务器后，每隔1秒发送1次消息，然后由服务器主动关闭链接（通过函数返回的方式）
func TestClientClose3(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode(true), WithIsWebsocket(true))
	if err != nil {
		t.Fatal(err)
	}
	errGroup, errCtx := errgroup.WithContext(context.TODO())
	errGroup.Go(func() error {
		go func() {
			gate.Serve(&GateHandler3{})

		}()
		select {
		case <-errCtx.Done():
			gate.Shutdown()
		}
		return nil
	})
	errGroup.Go(func() error {
		time.Sleep(1 * time.Millisecond)
		var conn gira.GateClient
		conn, err = client.Dial("127.0.0.1:1234", client.WithDebugMode(), client.WithIsWebsocket(true))
		if err != nil {
			t.Fatal(err)
		}
		for {
			err = conn.Request("", 1, []byte("hello"))
			if err != nil {
				log.Infow("request", "error", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
		conn.Close()
		return io.EOF
	})
	err = errGroup.Wait()
	if err != nil {
		log.Infow("errGroup", "error", err)

	}
	return
}

type GateHandler4 struct {
}

func (self *GateHandler4) OnClientStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	req, err = s.Recv(context.TODO())
	if err != nil {
		log.Infow("recv", "error", err)

	} else {
		log.Infow("recv", "data", string(req.Payload()))

	}
}

// 客户端连接上服务器后，不进行握手
func TestClientClose4(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode(true))
	if err != nil {
		t.Fatal(err)
	}
	errGroup, errCtx := errgroup.WithContext(context.TODO())
	errGroup.Go(func() error {
		go func() {
			gate.Serve(&GateHandler4{})
		}()
		select {
		case <-errCtx.Done():
			time.Sleep(1 * time.Second)
			log.Infow("gate stat", "ActiveSessionCount", gate.Stat.ActiveSessionCount)
			log.Infow("gate stat", "CumulativeSessionCount", gate.Stat.CumulativeSessionCount)
			gate.Shutdown()
		}
		return nil
	})
	errGroup.Go(func() error {
		time.Sleep(1 * time.Millisecond)
		var conn net.Conn
		conn, err = net.Dial("tcp", "127.0.0.1:1234")
		if err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		log.Info(n, err)
		return io.EOF
	})
	err = errGroup.Wait()
	if err != nil {
		log.Infow("errGroup", "error", err)

	}
	return
}
