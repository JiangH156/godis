package database

import (
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/redis/parser"
	"github.com/jiangh156/godis/redis/protocol"
	"io"
	"strconv"
	"strings"
)

// 加载aof文件
func LoadAof() {
	currentDB := 0
	var server = RedisServerInstance

	aofFile := server.DBSet[0].aofFile

	ch := parser.ParseStream(aofFile)
	for payload := range ch {
		if payload.Err != nil {
			// conn close
			logger.Warn(payload.Err.Error())
			if payload.Err == io.EOF {
				RedisServerInstance.aofLoading.Set(false)
				return
			}
		} else {
			// payload.Data is null
			if payload.Data == nil {
				logger.Info("empty reply")
				continue
			}
			args := [][]byte{}
			// require bulk reply
			switch payload.Data.(type) {
			case *protocol.StatusReply: // PING
				status := payload.Data.(*protocol.StatusReply).Status
				args = append(args, []byte(status))
			case *protocol.MultiBulkReply:
				args = payload.Data.(*protocol.MultiBulkReply).Args
			}
			// handle select
			if strings.ToLower(string(args[0])) == "select" {
				dbNum, err := strconv.ParseInt(string(args[1]), 10, 64)
				if err != nil {
					logger.Warn(err.Error())
					continue
				}
				currentDB = int(dbNum)
			} else {
				// handle common
				server.DBSet[currentDB].Exec(nil, args)
			}
		}
	}
}

func (db *DB) addAof(args *protocol.MultiBulkReply) {
	if config.Properties.AppendOnly && db.aofChan != nil && !RedisServerInstance.aofLoading.Get() {
		db.aofChan <- args
	}
}

func (db *DB) handleAof() {
	if !config.Properties.AppendOnly {
		return
	}
	for {
		select {
		case reply := <-db.aofChan:
			_, err := db.aofFile.Write(reply.ToBytes())
			if err != nil {
				logger.Warn(err.Error())
			}
			db.aofFile.Write([]byte("\r\n"))
			db.aofFile.Sync()
		}
	}
}

func (db *DB) makeAofCmd(cmd string, args [][]byte) *protocol.MultiBulkReply {
	result := make([][]byte, len(args)+1)
	copy(result[1:], args)
	result[0] = []byte(cmd)
	return protocol.MakeMultiBulkReply(result)
}
