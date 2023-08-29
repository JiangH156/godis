package server

import (
	"context"
	"errors"
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
		_ = conn.Close()
	}()
	client := connection.NewConn(conn)
	handler.ActiveConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
label:
	for payload := range ch {
		if payload.Err != nil {
			// conn close
			if payload.Err == io.EOF || errors.Is(payload.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				handler.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr())
				return
			}
			errReply := protocol.MakeErrReply(payload.Err.Error())
			// write back errMsg
			_, err := client.Write(errReply.ToBytes())
			if err != nil {
				handler.closeClient(client)
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
					continue label
				}
				args[0] = []byte(status)
			case *protocol.MultiBulkReply:
				args = payload.Data.(*protocol.MultiBulkReply).Args
			default:
				logger.Info("require Bulk or multiBulk")
				continue label
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

func (handler *RedisHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	handler.DB.AfterClientClose(client)
	handler.ActiveConn.Delete(client)
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
			handler.closeClient(client) //这一步最多阻塞5秒
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
