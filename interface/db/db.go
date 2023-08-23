package db

import "github.com/jiangh156/godis/interface/redis"

type DataBase interface {
	Exec(conn redis.Connection, args [][]byte) redis.Reply
	Close()
	AfterClientClose(c redis.Connection)
}
