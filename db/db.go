package db

import (
	"github.com/jiangh156/godis/datastruct/dict"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
)

type ExecFunc func(db *DB, args [][]byte) redis.Reply
type CmdLine [][]byte

type DataEntity struct {
	Data any
}
type DB struct {
	index int
	Data  dict.Dict
}

func (db *DB) Close() {
	//TODO aof处理
}

func (db *DB) AfterClientClose(c redis.Connection) {
	//TODO pub/sub 处理
}

func MakeDB() *DB {
	db := &DB{
		Data: dict.MakeSyncDict(),
	}
	return db
}

func (db *DB) Exec(conn redis.Connection, cmdLine CmdLine) redis.Reply {
	cmdName := string(cmdLine[0])
	cmd := cmdTable[cmdName]
	ok := validateCommand(cmd, cmdLine)
	if !ok {
		return protocol.MakeArgNumErrReply()
	}
	return cmd.exector(db, cmdLine[1:])
}

func validateCommand(c *command, cmdLine CmdLine) bool {
	if c.arity < 0 { // -3 至少3个字段, cmdLine的数量大于等于-arity
		return len(cmdLine) >= -c.arity
	}
	return len(cmdLine) == c.arity // cmdLine的数量必须等于arity
}

func (db *DB) Get(key string) (entity *DataEntity, exists bool) {
	raw, exists := db.Data.Get(key)
	if !exists {
		return nil, false
	}
	entity, _ = raw.(*DataEntity)
	return entity, true
}
func (db *DB) Put(key string, val any) (result int) {
	result = db.Data.Put(key, val)
	return result
}
func (db *DB) PutIfExists(key string, val any) (result int) {
	result = db.Data.PutIfExists(key, val)
	return result
}
func (db *DB) PutIfAbsent(key string, val any) (result int) {
	result = db.Data.PutIfAbsent(key, val)
	return result
}
func (db *DB) Remove(key string) (result int) {
	_, exists := db.Data.Get(key)
	if !exists {
		return 0
	}
	result = db.Data.Remove(key)
	return result
}
func (db *DB) Removes(keys ...string) (result int) {
	for _, key := range keys {
		_, exists := db.Data.Get(key)
		if exists {
			result++
		}
	}
	return result
}
func (db *DB) Flush() {
	db.Data.Clear()
}
