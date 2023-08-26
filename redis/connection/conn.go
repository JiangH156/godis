package connection

import (
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/lib/sync/wait"
	"net"
	"sync"
)

type Connection struct {
	Conn       net.Conn
	SelectedDB int

	sendWait wait.Wait
}

var connPool = sync.Pool{
	New: func() any {
		return &Connection{}
	},
}

func NewConnection() *Connection {
	conn, ok := connPool.Get().(*Connection)
	if !ok {
		logger.Error("connection pool make wrong type")
		return &Connection{}
	}
	return conn
}

func (c *Connection) Write(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	c.sendWait.Add(1)
	defer c.sendWait.Done()
	return c.Conn.Write(bytes)
}

func (c *Connection) GetDBIndex() int {
	return c.SelectedDB
}

func (c *Connection) SelectDB(dbNum int) {
	c.SelectedDB = dbNum
}
