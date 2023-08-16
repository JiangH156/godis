package tcp

import (
	"context"
	"fmt"
	"github.com/jiangh156/godis/interface/tcp"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
	Port    int
	MaxConn int
}

// 信号控制
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
	if err != nil {
		log.Fatalln("listen err", err)
		return err
	}
	fmt.Printf("Bind:%s, start listen...\n", cfg.Address)
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	defer func() {
		_ = handler.Close()
		_ = listener.Close()
	}()
	go func() {
		<-closeChan
		_ = handler.Close()
		_ = listener.Close()
	}()
	ctx := context.Background()
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept err:", err)
			break
		}
		wg.Add(1)
		fmt.Println("accept link")
		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
}
