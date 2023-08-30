package database

import (
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/lib/wildcard"
	"github.com/jiangh156/godis/redis/protocol"
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
		if exists {
			result++
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
		if pattern.IsMatch(key) {
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
	return protocol.MakeStatusReply("OK")
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
	_, exists = db.Get(newKey)
	if exists {
		return protocol.MakeErrReply("ERR target key name already exists")
	}
	db.Put(newKey, entity)
	return protocol.MakeStatusReply("OK")
}
func init() {
	RegisterCommand("Del", execDel, -2)
	RegisterCommand("Exists", execExists, -2)
	RegisterCommand("Keys", execKeys, 2)
	RegisterCommand("FlushDB", execFlushDB, 1)
	RegisterCommand("Type", execType, 2)
	RegisterCommand("Rename", execRename, 3)
	RegisterCommand("RenameNX", execRenameNX, 3)

	RegisterSingleCommand("FLUSHDB")
}
