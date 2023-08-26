package redis

type Connection interface {
	Write([]byte) (int, error)
	GetDBIndex() int
	SelectDB(int)
}
