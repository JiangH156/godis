package cluster

import "github.com/jiangh156/godis/interface/redis"

func execSelect(cluster *ClusterServer, c redis.Connection, cmdArgs [][]byte) redis.Reply {
	return cluster.db.Exec(c, cmdArgs)
}
