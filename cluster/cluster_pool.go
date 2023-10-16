package cluster

import (
	"context"
	"errors"
	"github.com/jiangh156/godis/redis/client"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type connectionFactory struct {
	Peer string
}

func (cluster *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	// 存放每一个客户端地址
	c, err := client.MakeClient(cluster.Peer)
	if err != nil {
		return nil, err
	}
	c.Start()
	return pool.NewPooledObject(c), nil
}

func (cluster *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	c, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	c.Close()
	return nil
}

func (cluster *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (cluster *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (cluster *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
