package redis

type Connection interface {
	Write([]byte)
}
