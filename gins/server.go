package gins

/// 参考 https://juejin.cn/post/6844903833273892871

import (
	"context"
	"net/http"
	"time"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
)

type HttpServer struct {
	config  gira.HttpConfig
	Handler http.Handler
	server  *http.Server
	ctx     context.Context
}

func NewConfigHttpServer(config gira.HttpConfig, router http.Handler) (*HttpServer, error) {
	self := &HttpServer{
		config:  config,
		Handler: router,
	}
	server := &http.Server{
		Addr:           self.config.Addr,
		Handler:        self.Handler,
		ReadTimeout:    time.Duration(self.config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(self.config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	self.server = server
	return self, nil
}

func (self *HttpServer) OnStart(ctx context.Context) error {
	self.ctx = ctx
	return nil
}

func (self *HttpServer) OnStop() error {
	log.Debugw("http server on stop")
	return nil
}
func (self *HttpServer) Serve() error {
	log.Debugw("http server started", "addr", self.config.Addr)
	go func() {
		<-self.ctx.Done()
		self.server.Close()
	}()
	var err error
	if self.config.Ssl && len(self.config.CertFile) > 0 && len(self.config.KeyFile) > 0 {
		if err = self.server.ListenAndServeTLS(self.config.CertFile, self.config.KeyFile); err == http.ErrServerClosed {
			err = nil
		}
	} else {
		if err = self.server.ListenAndServe(); err == http.ErrServerClosed {
			err = nil
		}
	}
	return err
}