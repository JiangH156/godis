package database

import (
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/interface/db"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/lib/sync/atomic"
	"github.com/jiangh156/godis/redis/protocol"
	"strconv"
	"strings"
)

type SingleServer struct {
	DBSet      []*DB
	aofLoading atomic.AtomicBool
}

var RedisServerInstance *SingleServer

var _ db.DataBase = (*SingleServer)(nil)

// redis节点Datebase
func NewSingleServer() *SingleServer {
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	server := &SingleServer{DBSet: make([]*DB, 16)}
	for i := range server.DBSet {
		db := MakeDB()
		db.index = i
		server.DBSet[i] = db
	}
	RedisServerInstance = server
	//TODO AOF
	if config.Properties.AppendOnly {
		server.aofLoading.Set(true)
		LoadAof()
	}
	return server
}

func (s *SingleServer) Exec(conn redis.Connection, args [][]byte) redis.Reply {
	// 处理select命令
	cmdName := args[0]
	if strings.ToLower(string(cmdName)) == "select" {
		return execSelect(conn, args)
	}
	index := conn.GetDBIndex()
	return s.DBSet[index].Exec(conn, args)
}

func execSelect(conn redis.Connection, args [][]byte) redis.Reply {
	dbNum, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR illegal number: " + string(args[1]))
	}
	if dbNum < 0 || dbNum > 15 {
		return protocol.MakeErrReply("ERR invalid DB index: " + strconv.Itoa(int(dbNum)))
	}
	conn.SelectDB(int(dbNum))
	if config.Properties.AppendOnly {
		aofReply := RedisServerInstance.DBSet[0].makeAofCmd("select", args[1:])
		RedisServerInstance.DBSet[0].addAof(aofReply)
	}
	return protocol.MakeOkReply()
}

func (s *SingleServer) Close() {
	for _, db := range s.DBSet {
		db.Close()
	}
}

func (s *SingleServer) AfterClientClose(conn redis.Connection) {
	//TODO implement me
}
