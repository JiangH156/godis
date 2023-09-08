package database

import (
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/redis/protocol"
)

//func LoadAof() {
//	currentDB := 0
//
//}

func (db *DB) addAof(args *protocol.MultiBulkReply) {
	if config.Properties.AppendOnly {
		db.aofChan <- args
	}
}

func (db *DB) handleAof() {
	if !config.Properties.AppendOnly {
		return
	}
	select {
	case reply := <-db.aofChan:
		_, err := db.aofFile.Write(reply.ToBytes())
		if err != nil {
			logger.Warn(err.Error())
		}
	}
}

func (db *DB) makeAofCmd(cmd string, args [][]byte) *protocol.MultiBulkReply {
	result := make([][]byte, len(args)+1)
	copy(result[1:], args)
	result[0] = []byte(cmd)
	return protocol.MakeMultiBulkReply(result)
}
