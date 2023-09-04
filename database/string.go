package database

import (
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
	"strconv"
	"time"
)

func (db *DB) getAsString(key string) ([]byte, redis.ErrReply) {
	entity, exists := db.Get(key)
	if !exists {
		return nil, nil
	}
	bytes, ok := entity.Data.([]byte)
	if !ok {
		return nil, protocol.MakeWrongTypeErrReply()
	}
	return bytes, nil
}

func execGet(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'get' command")
	}
	key := string(args[0])
	val, err := db.getAsString(key)
	if err != nil {
		return err
	}
	if len(val) == 0 {
		return protocol.MakeNullBulkReply()
	}
	return protocol.MakeBulkReply(val)
}

func execSet(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'set' command")
	}
	key := string(args[0])
	val := args[1]
	entity := &DataEntity{
		Data: val,
	}
	result := db.Put(key, entity)
	if result == 0 {
		return protocol.MakeErrReply("ERR fail to store data")
	}
	return protocol.MakeOkReply()
}

func execSetNX(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'setNX' command")
	}
	key := string(args[0])
	val := args[1]
	result := db.PutIfAbsent(key, &DataEntity{
		Data: val,
	})
	return protocol.MakeIntReply(int64(result))
}

func execGetSet(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'getset' command")
	}
	key := string(args[0])
	val := args[1]

	entity, exists := db.Get(key)
	oldVal := entity.Data.([]byte)
	if !exists {
		return protocol.MakeErrReply("ERR no such key")
	}
	db.Put(key, &DataEntity{
		Data: val,
	})
	return protocol.MakeBulkReply(oldVal)
}
func execStrlen(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'strlen' command")
	}
	key := string(args[0])
	bytes, err := db.getAsString(key)
	if err != nil {
		return err
	}
	return protocol.MakeIntReply(int64(len(bytes)))
}

// SETEX <key> <seconds> <value>
func execSetEX(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SetEX' command")
	}
	key := string(args[0])
	val := args[2]
	ttlArg, err := strconv.ParseInt(string(args[1]), 10, 64)
	ttlArg *= 10
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	if ttlArg <= 0 {
		return protocol.MakeErrReply("ERR invalid expire time")
	}
	db.Put(key, &DataEntity{Data: val})
	expireTime := time.Now().Add(time.Duration(ttlArg) * time.Millisecond)
	db.Expire(key, expireTime)
	return protocol.MakeOkReply()
}

func init() {
	RegisterCommand("Get", execGet, 2)
	RegisterCommand("Set", execSet, 3)
	RegisterCommand("SetNX", execSetNX, 3)
	RegisterCommand("GetSet", execGetSet, 3)
	RegisterCommand("Strlen", execStrlen, 2)
	RegisterCommand("SetEX", execSetEX, 4)
}
