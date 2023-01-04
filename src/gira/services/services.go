package services

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type Interface interface {
	Awake() error
	Start() error
	Destory() error
}

type HttpServer interface {
	NewHttpServer(serverName string) (error, *http.Server)
}

func Main(service Interface) error {
	if err := service.Awake(); err != nil {
		return err
	}
	if _, ok := service.(HttpServer); ok {
		log.Println("is a http")
	}
	rand.Seed(time.Now().UnixNano())
	ctx, _ := context.WithCancel(context.Background())
	group, errCtx := errgroup.WithContext(ctx)

	ctrlFunc := func() error {
		quit := make(chan os.Signal)
		defer close(quit)
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		for {
			select {
			case s := <-quit:
				switch s {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					log.Println("ctrl shutdown.")
					return nil
				case syscall.SIGUSR1:
					log.Println("sigusr1.")
				case syscall.SIGUSR2:
					log.Println("sigusr2.")
				default:
					log.Println("single x")
				}
			case <-errCtx.Done():
				log.Println("recv ctx:", errCtx.Err().Error())
				return nil
			}
		}
		log.Printf("ctrlServer shutdown.")
		return nil
	}
	group.Go(ctrlFunc)
	if err := service.Start(); err != nil {
		return err
	}
	if err := group.Wait(); err != nil {
		log.Println("all goroutine done. err:", err.Error())
	} else {
		log.Println("all goroutine done.")
	}
	return nil
}
