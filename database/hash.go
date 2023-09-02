package database

import (
	"github.com/jiangh156/godis/datastruct/dict"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
)

func (db *DB) getAsHash(key string) (dict.Dict, redis.ErrReply) {
	entity, exists := db.Get(key)
	if !exists {
		return nil, nil
	}
	hash, ok := entity.Data.(dict.Dict)
	if !ok {
		return nil, protocol.MakeWrongTypeErrReply()
	}
	return hash, nil
}
func (db *DB) getOrInitHash(key string) (dict.Dict, redis.ErrReply) {
	hash, errReply := db.getAsHash(key)
	if errReply != nil {
		return nil, errReply
	}
	if hash == nil {
		hash = dict.MakeSyncDict()
		db.Put(key, &DataEntity{Data: hash})
	}
	return hash, nil
}

// HSET key field value
func execHSet(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'HSet' command")
	}
	key := string(args[0])
	field := string(args[1])
	value := args[2]
	hash, err := db.getOrInitHash(key)
	if err != nil {
		return err
	}
	result := hash.Put(field, value)
	return protocol.MakeIntReply(int64(result))
}

// HGET key field
func execHGet(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'HGet' command")
	}
	key := string(args[0])
	field := string(args[1])
	hash, err := db.getAsHash(key)
	if err != nil {
		return err
	}
	if hash == nil {
		return protocol.MakeNullBulkReply()
	}
	val, exists := hash.Get(field)
	if !exists {
		return protocol.MakeNullBulkReply()
	}
	return protocol.MakeBulkReply(val.([]byte))
}

// HDEL key field [field ...]
func execHDel(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'HDel' command")
	}
	key := string(args[0])
	fields := args[1:]
	hash, err := db.getAsHash(key)
	if err != nil {
		return err
	}
	var removed int
	for _, field := range fields {
		result := hash.Remove(string(field))
		if result == 1 {
			removed++
		}
	}
	return protocol.MakeIntReply(int64(removed))
}

// HKEYS key
func execHKeys(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'HKeys' command")
	}
	key := string(args[0])
	hash, err := db.getAsHash(key)
	if err != nil {
		return err
	}
	keyStrs := hash.Keys()
	keys := make([][]byte, len(keyStrs))
	for i, str := range keyStrs {
		keys[i] = []byte(str)
	}
	return protocol.MakeMultiBulkReply(keys)
}

func init() {
	RegisterCommand("HSet", execHSet, 4)
	RegisterCommand("HGet", execHGet, 3)
	RegisterCommand("HDel", execHDel, -3)
	RegisterCommand("HKeys", execHKeys, 2)
}
