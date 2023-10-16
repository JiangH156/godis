package cluster

import (
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
)

// rename k1 k2
func rename(cluster *ClusterServer, c redis.Connection, cmdArgs [][]byte) redis.Reply {
	if len(cmdArgs) != 3 {
		return protocol.MakeErrReply("ERR wrong number args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	srcPeer := cluster.peerPicker.Get(src)
	destPeer := cluster.peerPicker.Get(dest)
	// 目前只允许在同一结点rename
	if srcPeer != destPeer {
		return protocol.MakeErrReply("ERR rename must within on peer")
	}
	return cluster.relay(srcPeer, c, cmdArgs)
}
