package database

import (
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/datastruct/dict"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/redis/protocol"
	"os"
	"time"
)

var (
	aofQueueSize = 1 << 4
)

type ExecFunc func(db *DB, args [][]byte) redis.Reply
type CmdLine [][]byte
type DataEntity struct {
	Data any
}
type DB struct {
	index  int
	Data   dict.Dict
	TTLMap dict.Dict

	aofChan     chan *protocol.MultiBulkReply
	aofFile     *os.File
	aofFilename string
}

func MakeDB() *DB {
	db := &DB{
		Data:   dict.MakeSyncDict(),
		TTLMap: dict.MakeSyncDict(),
	}
	if config.Properties.AppendOnly {
		aofFile, err := os.OpenFile(config.Properties.AppendFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			logger.Warn(err.Error())
		} else {
			db.aofFile = aofFile
			db.aofChan = make(chan *protocol.MultiBulkReply, aofQueueSize)
		}
		go db.handleAof()
	}
	return db
}
func (db *DB) Expire(key string, expireTime time.Time) {
	db.TTLMap.Put(key, expireTime)
}
func (db *DB) IsExpire(key string) bool {
	rawExpireTime, ok := db.TTLMap.Get(key)
	if !ok {
		return false
	}
	expireTime := rawExpireTime.(time.Time)
	expired := time.Now().After(expireTime)
	if expired {
		db.TTLMap.Remove(key)
	}
	return expired
}
func (db *DB) Persist(key string) {
	db.TTLMap.Remove(key)
}
func (db *DB) CleanExpire() {
	keys := db.TTLMap.Keys()
	for _, key := range keys {
		rawExpireTime, ok := db.TTLMap.Get(key)
		// key is deleted when range
		if !ok {
			continue
		}
		_, exists := db.Get(key)
		if !exists {
			db.TTLMap.Remove(key)
			continue
		}
		expireTime := rawExpireTime.(time.Time)
		expired := time.Now().After(expireTime)
		if expired {
			db.TTLMap.Remove(key)
			db.Data.Remove(key)
		}
	}
}

func (db *DB) Close() {
	db.Data.Clear()
}

func (db *DB) AfterClientClose(c redis.Connection) {
	//TODO pub/sub 处理
}

func (db *DB) GetEntity(key string) (entity *DataEntity, ok bool) {
	entity, ok = db.Get(key)
	if !ok {
		return nil, false
	}
	return entity, true
}
func (db *DB) Exec(conn redis.Connection, cmdLine CmdLine) redis.Reply {
	cmdName := string(cmdLine[0])
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return protocol.MakeWrongTypeErrReply()
	}
	ok = validateCommand(cmd, cmdLine)
	if !ok {
		return protocol.MakeArgNumErrReply(cmdName)
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
func (db *DB) Put(key string, val *DataEntity) (result int) {
	result = db.Data.Put(key, val)
	return result
}
func (db *DB) PutIfExists(key string, val *DataEntity) (result int) {
	result = db.Data.PutIfExists(key, val)
	return result
}
func (db *DB) PutIfAbsent(key string, val *DataEntity) (result int) {
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
			db.Data.Remove(key)
			result++
		}
	}
	return result
}
func (db *DB) Flush() {
	db.Data.Clear()
	db.TTLMap.Clear()
}
