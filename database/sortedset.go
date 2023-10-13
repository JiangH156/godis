package database

import (
	"github.com/jiangh156/godis/datastruct/sortedset"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
	"strconv"
)

func (db *DB) getAsSortedSet(key string) (*sortedset.SortedSet, redis.ErrReply) {
	entity, exists := db.Get(key)
	if !exists {
		return nil, nil
	}
	sortedSet, ok := entity.Data.(*sortedset.SortedSet)
	if !ok {
		return nil, &protocol.WrongTypeErrReply{}
	}
	return sortedSet, nil
}

func (db *DB) getOrInitSortedSet(key string) (sortedSet *sortedset.SortedSet, inited bool, errReply redis.ErrReply) {
	sortedSet, errReply = db.getAsSortedSet(key)
	if errReply != nil {
		return nil, false, errReply
	}
	inited = false
	if sortedSet == nil {
		sortedSet = sortedset.Make()
		db.Put(key, &DataEntity{
			Data: sortedSet,
		})
		inited = true
	}
	return sortedSet, inited, nil
}

// ZADD key score member [score member ...]
func execZAdd(db *DB, args [][]byte) redis.Reply {
	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZAdd' command")
	}
	key := string(args[0])
	zSet, _, err := db.getOrInitSortedSet(key)
	if err != nil {
		return err
	}
	kvLen := len(args)

	cnt := 0
	// 下标从1开始，args[0] 为key
	for i := 1; i < kvLen; i += 2 {
		member := string(args[i])
		score, err := strconv.ParseFloat(string(args[i+1]), 64)
		if err != nil {
			return protocol.MakeErrReply("ERR value is not a valid float64")
		}
		add := zSet.Add(member, score)
		if add {
			cnt++
		}
	}
	return protocol.MakeIntReply(int64(cnt))
}

// ZSCORE key member
func execZScore(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZScore' command")
	}
	key := string(args[0])
	zSet, _, err := db.getOrInitSortedSet(key)
	if err != nil {
		return err
	}
	member := string(args[1])
	element, ok := zSet.Get(member)
	if !ok {
		return protocol.MakeNullBulkReply()
	}
	return protocol.MakeBulkReply([]byte(strconv.FormatFloat(element.Score, 'f', -1, 64)))
}

// ZRANK key member
func execZRank(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRank' command")
	}
	key := string(args[0])
	zSet, _, err := db.getOrInitSortedSet(key)
	if err != nil {
		return err
	}
	if zSet == nil {
		return protocol.MakeNullBulkReply()
	}
	member := string(args[1])
	rank := zSet.GetRank(member, false)
	return protocol.MakeIntReply(rank)
}

// ZREVRANK key member
func execZRevRank(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRevRank' command")
	}
	key := string(args[0])
	zSet, _, err := db.getOrInitSortedSet(key)
	if err != nil {
		return err
	}
	if zSet == nil {
		return protocol.MakeNullBulkReply()
	}
	member := string(args[1])
	rank := zSet.GetRank(member, true)
	return protocol.MakeIntReply(rank)
}

// ZCARD key
func execZCard(db *DB, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZCard' command")
	}
	key := string(args[0])
	zSet, _, err := db.getOrInitSortedSet(key)
	if err != nil {
		return err
	}
	if zSet == nil {
		return protocol.MakeNullBulkReply()
	}
	return protocol.MakeIntReply(zSet.Len())
}

// ZRANGE key start stop
func execZRange(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRange' command")
	}
	key := string(args[0])
	start, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	stop, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeEmptyMultiBulkReply()
	}
	elements := zSet.Range(start, stop, false)
	members := make([][]byte, len(elements))
	for i, e := range elements {
		members[i] = []byte(e.Member)
	}
	return protocol.MakeMultiBulkReply(members)
}

// ZREVRANGE key start stop
func execZRevRange(db *DB, args [][]byte) redis.Reply {
	if len(args) != 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRevRange' command")
	}
	key := string(args[0])
	start, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	stop, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeEmptyMultiBulkReply()
	}
	elements := zSet.Range(start, stop, true)
	members := make([][]byte, len(elements))
	for i, e := range elements {
		members[i] = []byte(e.Member)
	}
	return protocol.MakeMultiBulkReply(members)
}

/*
 * param limit: limit < 0 means no limit
 */
// ZCOUNT key min max
func execZCount(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZCount' command")
	}
	key := string(args[0])
	min, err := strconv.ParseFloat(string(args[1]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	max, err := strconv.ParseFloat(string(args[2]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeIntReply(0)
	}
	elements := zSet.RangeByScore(&sortedset.ScoreBorder{Value: min}, &sortedset.ScoreBorder{Value: max}, 0, -1, false)
	members := make([][]byte, len(elements))
	return protocol.MakeIntReply(int64(len(members)))
}

// ZRANGEBYSCORE key min max
func execZRangeByScore(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRangeByScore' command")
	}
	key := string(args[0])
	min, err := strconv.ParseFloat(string(args[1]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	max, err := strconv.ParseFloat(string(args[2]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeEmptyMultiBulkReply()
	}
	elements := zSet.RangeByScore(&sortedset.ScoreBorder{Value: min}, &sortedset.ScoreBorder{Value: max}, 0, -1, false)
	members := make([][]byte, len(elements))
	for i, e := range elements {
		members[i] = []byte(e.Member)
	}
	return protocol.MakeMultiBulkReply(members)
}

func execZRevRangeByScore(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRevRangeByScore' command")
	}
	key := string(args[0])
	min, err := strconv.ParseFloat(string(args[1]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	max, err := strconv.ParseFloat(string(args[2]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeEmptyMultiBulkReply()
	}
	elements := zSet.RangeByScore(&sortedset.ScoreBorder{Value: min}, &sortedset.ScoreBorder{Value: max}, 0, -1, true)
	members := make([][]byte, len(elements))
	for i, e := range elements {
		members[i] = []byte(e.Member)
	}
	return protocol.MakeMultiBulkReply(members)
}

// ZREMRANGEBYSCORE key min max
func execZRemRangeByScore(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRemRangeByScore' command")
	}
	key := string(args[0])
	min, err := strconv.ParseFloat(string(args[1]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	max, err := strconv.ParseFloat(string(args[2]), 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeIntReply(0)
	}
	removed := zSet.RemoveByScore(&sortedset.ScoreBorder{Value: min}, &sortedset.ScoreBorder{Value: max})
	return protocol.MakeIntReply(removed)
}

// ZREMRANGEBYRANK key start stop
func execZRemRangeByRank(db *DB, args [][]byte) redis.Reply {
	if len(args) != 3 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRemRangeByRank' command")
	}
	key := string(args[0])
	start, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	stop, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeIntReply(0)
	}
	removed := zSet.RemoveByRank(start, stop)
	return protocol.MakeIntReply(removed)
}

// ZREM key member [member ...]
func execZRem(db *DB, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ZRem' command")
	}
	key := string(args[0])
	members := args[1:]
	zSet, _, errReply := db.getOrInitSortedSet(key)
	if errReply != nil {
		return errReply
	}
	if zSet == nil {
		return protocol.MakeIntReply(0)
	}
	removed := int64(0)
	for _, member := range members {
		ok := zSet.Remove(string(member))
		if ok {
			removed++
		}
	}
	return protocol.MakeIntReply(removed)
}

func init() {
	RegisterCommand("ZAdd", execZAdd, -4)
	RegisterCommand("ZScore", execZScore, -3)
	RegisterCommand("ZRank", execZRank, 3)
	RegisterCommand("ZRevRank", execZRevRank, 3)
	RegisterCommand("ZCard", execZCard, 2)
	RegisterCommand("ZRange", execZRange, 4)
	RegisterCommand("ZRevRange", execZRevRange, 4)
	RegisterCommand("ZCount", execZCount, 4)
	RegisterCommand("ZRangeByScore", execZRangeByScore, 4)
	RegisterCommand("ZRevRangeByScore", execZRevRangeByScore, 4)
	RegisterCommand("ZRemRangeByScore", execZRemRangeByScore, 4)
	RegisterCommand("ZRemRangeByRank", execZRemRangeByRank, 4)
	RegisterCommand("ZRem", execZRem, -3)
}
