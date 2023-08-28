package connection

import (
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	Conn        net.Conn
	SelectedDB  int
	sendingWait wait.Wait

	password string
}

var connPool = sync.Pool{
	New: func() any {
		return &Connection{}
	},
}

func NewConn(conn net.Conn) *Connection {
	c, ok := connPool.Get().(*Connection)
	if !ok {
		logger.Error("connection pool make wrong type")
		return &Connection{
			Conn: conn,
		}
	}
	c.Conn = conn
	return c
}

func (c *Connection) RemoteAddr() string {
	return c.Conn.RemoteAddr().String()
}
func (c *Connection) Name() string {
	if c.Conn != nil {
		return c.Conn.RemoteAddr().String()
	}
	return ""
}

func (c *Connection) SetPassword(password string) {
	c.password = password
}
func (c *Connection) GetPassword() string {
	return c.password
}

func (c *Connection) Close() error {
	_ = c.sendingWait.WaitWithTimeout(5 * time.Second)
	c.Conn.Close()
	return nil
}

func (c *Connection) Write(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	c.sendingWait.Add(1)
	defer c.sendingWait.Done()
	return c.Conn.Write(bytes)
}

func (c *Connection) GetDBIndex() int {
	return c.SelectedDB
}

func (c *Connection) SelectDB(dbNum int) {
	c.SelectedDB = dbNum
}
