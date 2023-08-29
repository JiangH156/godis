package database

import (
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
	"strings"
)

type Server struct {
	DBSet []*DB
}

// redis节点Datebase
func NewSingleServer() *Server {
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	server := &Server{DBSet: make([]*DB, 16)}
	for i := range server.DBSet {
		db := MakeDB()
		db.index = i
		server.DBSet[i] = db
	}
	return server
}

func (s *Server) Exec(conn redis.Connection, args [][]byte) redis.Reply {
	// 处理select命令
	cmdName := args[0]
	if strings.ToLower(string(cmdName)) == "ping" {
		return protocol.MakePingReply()
	}
	index := conn.GetDBIndex()
	return s.DBSet[index].Exec(conn, args)
}

func (s *Server) Close() {
	for _, db := range s.DBSet {
		db.Close()
	}
}

func (s *Server) AfterClientClose(conn redis.Connection) {
	//TODO implement me
}
