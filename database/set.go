package database

import (
	Set "github.com/jiangh156/godis/datastruct/set"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
	"strconv"
)

func (db *DB) getAsSet(key string) (*Set.Set, redis.ErrReply) {
	entity, exists := db.Get(key)
	if !exists {
		return nil, nil
	}
	set, ok := entity.Data.(*Set.Set)
	if !ok {
		return nil, protocol.MakeWrongTypeErrReply()
	}
	return set, nil
}
func (db *DB) getOrInitSet(key string) (*Set.Set, redis.ErrReply) {
	set, err := db.getAsSet(key)
	if err != nil {
		return nil, err
	}
	if set == nil {
		set = Set.Make()
		db.Put(key, &DataEntity{
			Data: set,
		})
	}
	return set, nil
}

// SADD key member [member ...]
func execSAdd(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SAdd' command")
	}
	key := string(args[0])
	members := args[1:]
	set, err := db.getOrInitSet(key)
	if err != nil {
		return err
	}
	for _, member := range members {
		set.Add(string(member))
	}
	return protocol.MakeIntReply(int64(len(members)))
}

// SMEMBERS key
func execSMembers(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SMembers' command")
	}
	key := string(args[0])
	set, err := db.getOrInitSet(key)
	if err != nil {
		return err
	}
	if set == nil {
		return protocol.MakeEmptyMultiBulkReply()
	}
	slice := set.ToSlice()
	result := make([][]byte, len(slice))
	for i, s := range slice {
		result[i] = []byte(s)
	}
	return protocol.MakeMultiBulkReply(result)
}

// SREM key member [member ...]
func execSRem(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SRem' command")
	}
	key := string(args[0])
	members := args[1:]
	set, err := db.getOrInitSet(key)
	if err != nil {
		return err
	}
	if set == nil {
		return protocol.MakeIntReply(0)
	}
	var removed int
	for _, member := range members {
		res := set.Remove(string(member))
		if res == 1 {
			removed++
		}
	}
	return protocol.MakeIntReply(int64(removed))
}

// SISMEMBER key member
func execSIsMember(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SIsMember' command")
	}
	key := string(args[0])
	member := string(args[1])
	set, err := db.getOrInitSet(key)
	if err != nil {
		return err
	}
	if set == nil {
		return protocol.MakeIntReply(0)
	}
	ok := set.Has(member)
	result := 0
	if ok {
		result = 1
	}
	return protocol.MakeIntReply(int64(result))
}

// SCARD key
func execSCard(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SCard' command")
	}
	key := string(args[0])
	set, err := db.getOrInitSet(key)
	if err != nil {
		return err
	}
	if set == nil {
		return protocol.MakeIntReply(0)
	}
	return protocol.MakeIntReply(int64(set.Len()))
}

// SRANDMEMBER key [count]
func execSRandMember(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 && len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SRandMember' command")
	}
	key := string(args[0])
	set, err := db.getOrInitSet(key)
	if err != nil {
		return err
	}
	if set == nil {
		return protocol.MakeIntReply(0)
	}
	if len(args) == 1 {
		member := set.RandomMembers(1)
		return protocol.MakeBulkReply([]byte(member[0]))
	} else {
		count64, err := strconv.ParseInt(string(args[1]), 10, 64)
		if err != nil {
			return protocol.MakeErrReply("ERR value is not an integer or out of range")
		}
		count := int(count64)
		if count > 0 {
			members := set.RandomMembers(count)
			result := make([][]byte, len(members))
			for i, v := range members {
				result[i] = []byte(v)
			}
			return protocol.MakeMultiBulkReply(result)
		} else if count < 0 {
			members := set.RandomDistinctMembers(-count)
			result := make([][]byte, len(members))
			for i, v := range members {
				result[i] = []byte(v)
			}
			return protocol.MakeMultiBulkReply(result)
		} else {
			return protocol.MakeEmptyMultiBulkReply()
		}
	}

}

// SDIFF key [key ...]
func execSDiff(db *DB, args [][]byte) redis.Reply {
	if len(args) < 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'sdiff' command")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = string(arg)
	}

	var result *Set.Set
	for i, key := range keys {
		set, errReply := db.getAsSet(key)
		if errReply != nil {
			return errReply
		}
		if set == nil {
			if i == 0 {
				// early termination
				return protocol.MakeEmptyMultiBulkReply()
			} else {
				continue
			}
		}
		if result == nil {
			// init
			result = Set.MakeFromVals(set.ToSlice()...)
		} else {
			result = result.Diff(set)
			if result.Len() == 0 {
				// early termination
				return protocol.MakeEmptyMultiBulkReply()
			}
		}
	}

	if result == nil {
		// all keys are nil
		return protocol.MakeEmptyMultiBulkReply()
	}
	arr := make([][]byte, result.Len())
	i := 0
	result.ForEach(func(member string) bool {
		arr[i] = []byte(member)
		i++
		return true
	})
	return protocol.MakeMultiBulkReply(arr)
}

// SUNION key [key ...]
func execSUnion(db *DB, args [][]byte) redis.Reply {
	if len(args) < 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'SUnion' command")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = string(arg)
	}

	var result *Set.Set
	for i, key := range keys {
		set, errReply := db.getAsSet(key)
		if errReply != nil {
			return errReply
		}
		if set == nil {
			if i == 0 {
				// early termination
				return protocol.MakeEmptyMultiBulkReply()
			} else {
				continue
			}
		}
		if result == nil {
			// init
			result = Set.MakeFromVals(set.ToSlice()...)
		} else {
			result = result.Union(set)
			if result.Len() == 0 {
				// early termination
				return protocol.MakeEmptyMultiBulkReply()
			}
		}
	}
	arr := make([][]byte, result.Len())
	i := 0
	result.ForEach(func(member string) bool {
		arr[i] = []byte(member)
		i++
		return true
	})
	return protocol.MakeMultiBulkReply(arr)
}

func init() {
	RegisterCommand("SAdd", execSAdd, -3)
	RegisterCommand("SMembers", execSMembers, 2)
	RegisterCommand("SIsMember", execSIsMember, 3)
	RegisterCommand("SRandMember", execSIsMember, -2)
	RegisterCommand("SCard", execSCard, 2)
	RegisterCommand("SRandMember", execSRandMember, -2)
	RegisterCommand("SDiff", execSDiff, -2)
	RegisterCommand("SUnion", execSUnion, -2)
}
