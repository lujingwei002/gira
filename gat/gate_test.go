package gat

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/lujingwei/gira"
	"github.com/lujingwei/gira/gat/client"
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

type GateHandler1 struct {
}

func (self *GateHandler1) OnGateStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	for {
		req, err = s.Recv(context.TODO())
		if err != nil {
			gira.Infow("recv", "error", err)
			break
		} else {
			gira.Infow("recv", "data", string(req.Payload()))
		}
	}
}

// 客户端连接上服务器后，每隔1秒发送1次消息，然后主动关闭
func TestClientClose1(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode(), WithIsWebsocket(true))
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
		gira.Infow("errGroup", "error", err)
	}
	return
}

type GateHandler2 struct {
}

func (self *GateHandler2) OnGateStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	req, err = s.Recv(context.TODO())
	if err != nil {
		gira.Infow("recv", "error", err)

	} else {
		gira.Infow("recv", "data", string(req.Payload()))

	}
	s.Close()
}

// 客户端连接上服务器后，每隔1秒发送1次消息，然后由服务器主动关闭链接
func TestClientClose2(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode(), WithIsWebsocket(true))
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
				gira.Infow("request", "error", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
		conn.Close()
		return io.EOF
	})
	err = errGroup.Wait()
	if err != nil {
		gira.Infow("errGroup", "error", err)

	}
	return
}

type GateHandler3 struct {
}

func (self *GateHandler3) OnGateStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	req, err = s.Recv(context.TODO())
	if err != nil {
		gira.Infow("recv", "error", err)
	} else {
		gira.Infow("recv", "data", string(req.Payload()))
	}
}

// 客户端连接上服务器后，每隔1秒发送1次消息，然后由服务器主动关闭链接（通过函数返回的方式）
func TestClientClose3(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode(), WithIsWebsocket(true))
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
				gira.Infow("request", "error", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
		conn.Close()
		return io.EOF
	})
	err = errGroup.Wait()
	if err != nil {
		gira.Infow("errGroup", "error", err)

	}
	return
}

type GateHandler4 struct {
}

func (self *GateHandler4) OnGateStream(s gira.GateConn) {
	var req gira.GateRequest
	var err error
	req, err = s.Recv(context.TODO())
	if err != nil {
		gira.Infow("recv", "error", err)

	} else {
		gira.Infow("recv", "data", string(req.Payload()))

	}
}

// 客户端连接上服务器后，不进行握手
func TestClientClose4(t *testing.T) {
	var err error
	var gate *Gate
	gate, err = Listen(context.TODO(), ":1234", WithDebugMode())
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
		log.Println(n, err)
		return io.EOF
	})
	err = errGroup.Wait()
	if err != nil {
		gira.Infow("errGroup", "error", err)

	}
	return
}
