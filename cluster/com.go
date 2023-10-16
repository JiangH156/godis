package cluster

import (
	"context"
	"errors"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/lib/utils"
	"github.com/jiangh156/godis/redis/client"
	"github.com/jiangh156/godis/redis/protocol"
	"strconv"
)

func (cluster *ClusterServer) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return nil, errors.New("connection not found")
	}
	ctx := context.Background()
	object, err := pool.BorrowObject(ctx)
	if err != nil {
		return nil, err
	}
	c, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("wrong type")
	}
	return c, nil
}
func (cluster *ClusterServer) returnPeerClient(peer string, returnClient *client.Client) error {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return errors.New("connection not found")
	}
	return pool.ReturnObject(context.Background(), returnClient)
}

func (cluster *ClusterServer) relay(peer string, c redis.Connection, args [][]byte) redis.Reply {
	if peer == cluster.self {
		return cluster.db.Exec(c, args)
	}
	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return protocol.MakeErrReply("ERR" + err.Error())
	}
	defer cluster.returnPeerClient(peer, peerClient)
	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(c.GetDBIndex())))
	return peerClient.Send(args)
}

func (cluster *ClusterServer) broadcast(c redis.Connection, args [][]byte) map[string]redis.Reply {
	results := make(map[string]redis.Reply)
	for _, node := range cluster.nodes {
		result := cluster.relay(node, c, args)
		results[node] = result
	}
	return results
}
