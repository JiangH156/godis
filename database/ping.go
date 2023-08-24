package database

import (
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
)

func execPing(db *DB, args [][]byte) redis.Reply {
	return protocol.MakePingReply()
}

func init() {
	RegisterCommand("PING", execPing, 1)
}
