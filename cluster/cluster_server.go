package cluster

import (
	"context"
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/database"
	"github.com/jiangh156/godis/interface/db"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/lib/consistenthash"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/redis/protocol"
	pool "github.com/jolestar/go-commons-pool/v2"
	"strings"
)

type ClusterServer struct {
	self string

	nodes          []string
	peerPicker     *consistenthash.Map
	peerConnection map[string]*pool.ObjectPool
	db             db.DataBase
}

func MakeClusterServer() *ClusterServer {
	cluster := &ClusterServer{
		self:           config.Properties.Self,
		db:             database.NewSingleServer(),
		peerPicker:     consistenthash.New(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	// all nodes
	nodes = append(nodes, config.Properties.Self)
	cluster.nodes = nodes
	// 初始化一致性hash节点
	cluster.peerPicker.Add(nodes...)
	ctx := context.Background()
	// 使用连接池初始化peer结点
	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{Peer: peer})
	}
	return cluster
}

type CmdFunc func(cluster *ClusterServer, c redis.Connection, cmdArgs [][]byte) redis.Reply

var router = makeRouter()

func (cluster *ClusterServer) Exec(conn redis.Connection, args [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = protocol.MakeUnknownErrReply()
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return protocol.MakeErrReply("not supported cmd")
	}
	result = cmdFunc(cluster, conn, args)
	return
}

func (cluster *ClusterServer) Close() {
	cluster.db.Close()
}

func (cluster *ClusterServer) AfterClientClose(c redis.Connection) {
	cluster.db.AfterClientClose(c)
}
