package cluster

import "github.com/jiangh156/godis/interface/redis"

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc // exists k1
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["ping"] = ping
	routerMap["rename"] = rename
	routerMap["renamenx"] = rename
	routerMap["flushdb"] = flushdb
	routerMap["del"] = del
	routerMap["Select"] = execSelect
	return routerMap
}

// GET key // set k1 v1
// 只跟key有关
func defaultFunc(cluster *ClusterServer, c redis.Connection, cmdArgs [][]byte) redis.Reply {
	key := string(cmdArgs[1])
	peer := cluster.peerPicker.Get(key)
	return cluster.relay(peer, c, cmdArgs)
}
