package database

import (
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/lib/wildcard"
	"github.com/jiangh156/godis/redis/protocol"
	"strconv"
	"time"
)

// DEL k1 k2 k3
func execDel(db *DB, args [][]byte) redis.Reply {
	if len(args) < 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'del' command")
	}
	var keys []string
	for _, arg := range args {
		keys = append(keys, string(arg))
	}
	result := db.Removes(keys...)
	db.CleanExpire()
	return protocol.MakeIntReply(int64(result))
}

// EXISTS k1 k2 k3
func execExists(db *DB, args [][]byte) redis.Reply {
	if len(args) < 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'exists' command")
	}
	result := 0
	for _, arg := range args {
		_, exists := db.Get(string(arg))
		if !exists {
			continue
		}
		//expire
		if !db.IsExpire(string(arg)) {
			result++
		} else {
			// delete key
			db.Remove(string(arg))
			db.Persist(string(arg))
		}
	}
	return protocol.MakeIntReply(int64(result))
}

func execKeys(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'keys' command")
	}
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.Data.ForEach(func(key string, value any) bool {
		if pattern.IsMatch(key) && !db.IsExpire(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return protocol.MakeMultiBulkReply(result)
}
func execFlushDB(db *DB, args [][]byte) redis.Reply {
	if len(args) != 0 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'flushDB' command")
	}
	db.Flush()
	return protocol.MakeOkReply()
}
func execType(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'type' command")
	}
	key := string(args[0])
	entity, exists := db.Get(key)
	if !exists {
		return protocol.MakeStatusReply("none")
	}
	if db.IsExpire(key) {
		db.Persist(key)
		return protocol.MakeStatusReply("none")
	}
	// TODO 目前只有string，后续完善
	switch entity.Data.(type) {
	case []byte:
		return protocol.MakeStatusReply("string")
	}
	return protocol.MakeUnknownErrReply()
}
func execRename(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	oldKye := string(args[0])
	newKey := string(args[1])
	entity, exists := db.Get(oldKye)
	if !exists {
		return protocol.MakeErrReply("ERR no such key")
	}
	// key expired
	if db.IsExpire(oldKye) {
		db.Persist(oldKye)
		return protocol.MakeStatusReply("ERR no such key")
	}
	db.Remove(oldKye)
	db.Put(newKey, entity)
	return protocol.MakeStatusReply("OK")
}
func execRenameNX(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	oldKye := string(args[0])
	newKey := string(args[1])
	entity, exists := db.Get(oldKye)
	if !exists {
		return protocol.MakeErrReply("ERR no such key")
	}
	// key expired
	if db.IsExpire(oldKye) {
		db.Persist(oldKye)
		return protocol.MakeStatusReply("ERR no such key")
	}
	_, exists = db.Get(newKey)
	if exists {
		return protocol.MakeErrReply("ERR target key name already exists")
	}
	db.Put(newKey, entity)
	return protocol.MakeOkReply()
}

// EXPIRE <key> <seconds>
func execExpire(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'Expire' command")
	}
	key := string(args[0])
	ttlArg, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	if ttlArg <= 0 {
		return protocol.MakeErrReply("ERR invalid expire time")
	}
	_, exists := db.Get(key)
	if !exists {
		return protocol.MakeIntReply(0)
	}
	ttlArg *= 10
	expireTime := time.Now().Add(time.Duration(ttlArg) * time.Millisecond)
	db.Expire(key, expireTime)
	return protocol.MakeIntReply(1)
}
func init() {
	RegisterCommand("Del", execDel, -2)
	RegisterCommand("Exists", execExists, -2)
	RegisterCommand("Keys", execKeys, 2)
	RegisterCommand("FlushDB", execFlushDB, 1)
	RegisterCommand("Type", execType, 2)
	RegisterCommand("Rename", execRename, 3)
	RegisterCommand("RenameNX", execRenameNX, 3)
	RegisterCommand("Expire", execExpire, 3)

	RegisterSingleCommand("FLUSHDB")
}
