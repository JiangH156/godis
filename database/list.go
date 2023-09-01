package database

import (
	List "github.com/jiangh156/godis/datastruct/list"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
	"strconv"
)

func (db *DB) getAsList(key string) (list List.List, errReply redis.ErrReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	list, ok = entity.Data.(List.List)
	if !ok {
		return nil, protocol.MakeWrongTypeErrReply()
	}
	return list, nil
}
func (db *DB) getOrInitList(key string) (list List.List, isNew bool, errReply redis.ErrReply) {
	list, err := db.getAsList(key)
	if err != nil {
		return nil, false, err
	}
	if list == nil {
		list = List.Make()
		db.Put(key, &DataEntity{
			Data: list,
		})
		isNew = true
	}
	return list, isNew, nil
}

// LPUSH key value [value ...]
func execLPush(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'LPush' command")
	}
	key := string(args[0])
	values := args[1:]
	list, _, err := db.getOrInitList(key)
	if err != nil {
		return err
	}
	for _, val := range values {
		list.Insert(0, val)
	}
	return protocol.MakeIntReply(int64(list.Len()))
}

// RPUSH key value [value ...]
func execRPush(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'RPush' command")
	}
	key := string(args[0])
	values := args[1:]
	list, _, err := db.getOrInitList(key)
	if err != nil {
		return err
	}
	for _, val := range values {
		list.Add(val)
	}
	return protocol.MakeIntReply(int64(list.Len()))
}

// LPOP key
func execLPop(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'LPop' command")
	}
	key := string(args[0])
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeNullBulkReply()
	}
	val, _ := list.Remove(0).([]byte)
	if list.Len() == 0 {
		db.Remove(key)
	}
	return protocol.MakeBulkReply(val)
}

// RPOP key
func execRPop(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'RPop' command")
	}
	key := string(args[0])
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeNullBulkReply()
	}
	val, _ := list.RemoveLast().([]byte)
	if list.Len() == 0 {
		db.Remove(key)
	}
	return protocol.MakeBulkReply(val)
}

// LLEN key
func execLLen(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'LLen' command")
	}
	key := string(args[0])
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeIntReply(0)
	}
	size := list.Len()
	return protocol.MakeIntReply(int64(size))
}

// LINDEX key index
func execLIndex(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'LIndex' command")
	}
	key := string(args[0])
	index64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	index := int(index64)
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeNullBulkReply()
	}
	size := list.Len()
	if index < -1*size {
		return protocol.MakeNullBulkReply()
	} else if index < 0 {
		index = size + index
	} else if index >= size {
		return protocol.MakeNullBulkReply()
	}
	return protocol.MakeBulkReply(list.Get(index).([]byte))
}

// LRANGE key start stop
func execLRange(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'LRange' command")
	}
	key := string(args[0])
	start64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	start := int(start64)
	stop64, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	stop := int(stop64)
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeEmptyMultiBulkReply()
	}
	// compute index
	size := list.Len() // assert: size > 0
	if start < -1*size {
		start = 0
	} else if start < 0 {
		start = size + start
	} else if start >= size {
		return protocol.MakeEmptyMultiBulkReply()
	}
	if stop < -1*size {
		stop = 0
	} else if stop < 0 {
		stop = size + stop + 1
	} else if stop < size {
		stop = stop + 1
	} else {
		stop = size
	}
	if stop < start {
		stop = start
	}
	slice := list.Range(start, stop)
	result := make([][]byte, len(slice))
	for i, raw := range slice {
		result[i] = raw.([]byte)
	}
	return protocol.MakeMultiBulkReply(result)
}

// LSET key index value
func execLSet(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'LSet' command")
	}
	key := string(args[0])
	index64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	index := int(index64)
	value := args[2]
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeErrReply("ERR no such key")
	}
	size := list.Len()
	if index < -1*size {
		return protocol.MakeErrReply("ERR index out of range")
	} else if index < 0 {
		index = size + index
	} else if index >= size {
		return protocol.MakeErrReply("ERR index out of range")
	}
	list.Set(index, value)
	return protocol.MakeOkReply()
}

// LREM key count value 从列表中删除指定数量的匹配元素。
func execLRme(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'LRem' command")
	}
	key := string(args[0])
	count64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	count := int(count64)
	value := args[2]
	list, errReply := db.getAsList(key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeIntReply(0)
	}
	var removed int
	if count < 0 {
		count = -count
	}
	removed = list.RemoveByVal(value, count)
	return protocol.MakeIntReply(int64(removed))
}

// LTRIM key start stop 修剪列表，保留指定范围内的元素。
func execLTrim(db *DB, args [][]byte) redis.Reply {
	return nil
}

func init() {
	RegisterCommand("LPush", execLPush, -3)
	RegisterCommand("RPush", execRPush, -3)
	RegisterCommand("LPop", execLPop, 2)
	RegisterCommand("RPop", execRPop, 2)
	RegisterCommand("LLen", execLLen, 2)
	RegisterCommand("LIndex", execLIndex, 3)
	RegisterCommand("LRange", execLRange, 4)
	RegisterCommand("LSet", execLSet, 4)
	RegisterCommand("LRme", execLRme, 4)
}
