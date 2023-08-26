package database

import (
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/datastruct/dict"
	"github.com/jiangh156/godis/interface/redis"
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
		db := &DB{
			index: i,
			Data:  dict.NewSyncDict(),
		}
		server.DBSet[i] = db
	}
	return server
}

func (s *Server) Exec(conn redis.Connection, args [][]byte) redis.Reply {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Close() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) AfterClientClose(conn redis.Connection) {
	//TODO implement me
	panic("implement me")
}
