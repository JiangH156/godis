package dict

type Consumer func(key string, value any) bool

type Dict interface {
	Get(key string) (val any, exists bool)
	Put(key string, val any) (result int)
	PutIfAbsent(key string, val any) (result int)
	PutIfExists(key string, val any) (result int)
	Remove(key string) (result int)
	Len() int
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}
