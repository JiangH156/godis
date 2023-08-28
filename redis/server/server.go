package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/jiangh156/godis/database"
	"github.com/jiangh156/godis/interface/db"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/lib/sync/atomic"
	"github.com/jiangh156/godis/lib/sync/wait"
	"github.com/jiangh156/godis/redis/connection"
	"github.com/jiangh156/godis/redis/parser"
	"github.com/jiangh156/godis/redis/protocol"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RedisHandler struct {
	ActiveConn     sync.Map
	DB             db.DataBase
	closing        atomic.AtomicBool
	waitingClients wait.Wait
}

func (handler *RedisHandler) Handle(ctx context.Context, conn net.Conn) {
	// Handler已关闭，拒绝处理新的连接
	if handler.closing.Get() {
		return
	}
	defer func() {
		_ = handler.Close()
	}()
	client := connection.NewConn(conn)
	handler.ActiveConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		//TODO
		fmt.Printf("%v", payload)
		if payload.Err != nil {
			// conn close
			if payload.Err == io.EOF || errors.Is(payload.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				handler.Close()
				logger.Info("connection closed: " + client.RemoteAddr())
				return
			}
			errReply := protocol.MakeErrReply(payload.Err.Error())
			// write back errMsg
			_, err := client.Write(errReply.ToBytes())
			if err != nil {
				handler.Close()
				logger.Error("connection closed: ", err.Error())
				return
			}
		} else {
			// payload.Data is null
			if payload.Data == nil {
				logger.Info("empty reply")
				continue
			}
			args := [][]byte{}
			// require bulk reply
			switch payload.Data.(type) {
			case *protocol.StatusReply: // PING
				status := payload.Data.(*protocol.StatusReply).Status
				if strings.ToUpper(status) != "PING" {
					logger.Info("statusReply must is 'PING'")
					continue
				}
				args[0] = []byte(status)
			case *protocol.MultiBulkReply:
				args = payload.Data.(*protocol.MultiBulkReply).Args
			default:
				logger.Info("require Bulk or multiBulk")
				continue
			}
			result := handler.DB.Exec(client, args)
			if result != nil {
				_, _ = client.Write(result.ToBytes())
			} else {
				_, _ = client.Write(unknownErrReplyBytes)
			}
		}
	}

}

func (handler *RedisHandler) Close() error {
	logger.Info("Handler shutting down...")
	if !handler.closing.Get() {
		handler.closing.Set(true)
	}
	handler.ActiveConn.Range(func(key, value any) bool {
		go func() {
			handler.waitingClients.Add(1)
			defer handler.waitingClients.Done()
			client := key.(*connection.Connection)
			client.Close() // 这里会阻塞
			handler.ActiveConn.Delete(key)
		}()
		return true
	})
	handler.waitingClients.Wait()
	handler.DB.Close()
	return nil
}

func MakeRedisHandler() *RedisHandler {
	return &RedisHandler{
		DB: database.NewSingleServer(),
	}
}
