package server

import (
	"context"
	"github.com/jiangh156/godis/database"
	"github.com/jiangh156/godis/interface/db"
	"github.com/jiangh156/godis/lib/sync/atomic"
	"net"
	"sync"
)

type Handler struct {
	ActiveConn sync.Map
	DB         db.DataBase
	closing    atomic.AtomicBool
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) Close() error {
	//TODO implement me
	panic("implement me")
}

func MakeHandler() *Handler {
	return &Handler{
		DB: database.NewSingleServer(),
	}
}
