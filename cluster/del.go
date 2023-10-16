package cluster

import (
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
)

// del k1 k2 k3 k4 ...
func del(cluster *ClusterServer, c redis.Connection, cmdArgs [][]byte) redis.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply redis.ErrReply
	var deleted int64 = 0
	for _, r := range replies {
		if protocol.IsErrorReply(r) {
			errReply = r.(redis.ErrReply)
			break
		}
		intReply, ok := r.(*protocol.IntReply)
		if !ok {
			errReply = protocol.MakeErrReply("error")
		}
		deleted += intReply.Code
	}
	if errReply == nil {
		return protocol.MakeIntReply(deleted)
	}
	return protocol.MakeErrReply("error: " + errReply.Error())
}
