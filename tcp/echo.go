//package tcp
//
//import (
//	"bufio"
//	"context"
//	"fmt"
//	"github.com/jiangh156/godis/lib/logger"
//	"github.com/jiangh156/godis/lib/sync/atomic"
//	"github.com/jiangh156/godis/lib/sync/wait"
//	"io"
//	"net"
//	"sync"
//	"time"
//)
//
//// 客户端
//type Client struct {
//	// tcp 连接
//	Conn net.Conn
//	// 当服务端开始发送数据时进入waiting, 阻止其它goroutine关闭连接
//	Wait wait.Wait
//}
//
//// 维护客户端
//type EchoHandler struct {
//	// 保存所有工作状态client的集合(把map当set用)
//	// 需使用并发安全的容器
//	ActiveConn sync.Map
//	// 关闭状态标识位
//	closing atomic.AtomicBool
//}
//
//func (e *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
//	// handler已经关闭，拒绝处理业务
//	if e.closing.Get() {
//		conn.Close()
//		return
//	}
//	//启动一个客户端处理连接
//	client := &Client{Conn: conn}
//	//维护存活的客户端连接
//	e.ActiveConn.Store(client, struct{}{})
//	reader := bufio.NewReader(conn)
//	for {
//		msg, err := reader.ReadString('\n')
//		if err != nil {
//			// 这里要考虑是否为连接关闭的错误
//			if err == io.EOF {
//				logger.Info("connection close")
//				e.ActiveConn.Delete(client)
//
//			} else {
//				logger.Warn(err)
//			}
//			// 关闭客户端
//			client.Close()
//			return
//		}
//		// 对于每个消息，进行回写
//		//go func() {
//		//	client.Wait.Add(1)
//		//	b := []byte(msg)
//		//	fmt.Println(msg)
//		//	conn.Write(b)
//		//	client.Wait.Done()
//		//}()
//		client.Wait.Add(1)
//		//logger.Info("sleeping")
//		//time.Sleep(10 * time.Second)
//		b := []byte(msg)
//		_, _ = conn.Write(b)
//		client.Wait.Done()
//	}
//}
//
//func (c *Client) Close() error {
//	c.Wait.WaitWithTimeout(5 * time.Second) // 允许5秒钟完成剩余任务，超时强制关闭
//	c.Conn.Close()
//	return nil
//}
//
//// 客户端管理模块关闭，由上层服务端控制客户端关闭
//func (e *EchoHandler) Close() error {
//	fmt.Println("handler shutting down...")
//	if e.closing.Get() { //客户端未关闭
//		e.closing.Set(true)
//	}
//	e.ActiveConn.Range(func(k, v any) bool {
//		client := k.(*Client)
//		client.Close()
//		return true
//	})
//	return nil
//}

package tcp

/**
 * A echo server to test whether the server is functioning normally
 */

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/lib/sync/atomic"
	"github.com/jiangh156/godis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

// EchoHandler echos received line to client, using for test
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.AtomicBool
}

// MakeEchoHandler creates EchoHandler
func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// EchoClient is client for EchoHandler, using for test
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close close connection
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// Handle echos received line to client
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		// may occurs: client EOF, client timeout, server early close
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		//logger.Info("sleeping")
		//time.Sleep(10 * time.Second)
		fmt.Println(msg)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

// Close stops echo handler
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
