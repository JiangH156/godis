package cluster

import (
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
)

func flushdb(cluster *ClusterServer, c redis.Connection, cmdArgs [][]byte) redis.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply redis.ErrReply
	for _, r := range replies {
		if protocol.IsErrorReply(r) {
			errReply = r.(redis.ErrReply)
			break
		}
	}
	if errReply == nil {
		return protocol.MakeOkReply()
	}
	return protocol.MakeErrReply("error: " + errReply.Error())
}
