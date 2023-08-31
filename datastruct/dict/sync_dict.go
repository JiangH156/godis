package dict

import "sync"

// 使用sync.Map 实现的简单的Dict
type SyncDict struct {
	m sync.Map
}

func (dict *SyncDict) Get(key string) (val any, exists bool) {
	return dict.m.Load(key)
}

func (dict *SyncDict) Put(key string, val any) (result int) {
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val any) (result int) {
	_, exists := dict.m.Load(key)
	if exists { // 存在
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIfExists(key string, val any) (result int) {
	oldVal, exists := dict.m.Load(key)
	if exists { // 存在
		ok := dict.m.CompareAndSwap(key, oldVal, val)
		if ok { // 修改成功
			return 1
		}
		return 0
	}
	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	dict.m.Delete(key)
	return 1
}

func (dict *SyncDict) Len() int {
	i := 0
	dict.m.Range(func(key, value any) bool {
		i++
		return true
	})
	return i
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(k, value any) bool {
		key := k.(string)
		result := consumer(key, value)
		return result
	})
}

func (dict *SyncDict) Keys() []string {
	keys := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(k, value any) bool {
		key := k.(string)
		keys[i] = key
		i++
		return true
	})
	return keys
}

func (dict *SyncDict) RandomKeys(limit int) []string {
	keys := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(k, value any) bool {
			key := k.(string)
			keys[i] = key
			return false
		})
	}
	return keys
}

func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	keys := make([]string, limit)
	i := 0
	dict.m.Range(func(k, value any) bool {
		key := k.(string)
		keys[i] = key
		i++
		if i == limit {
			return false
		}
		return true
	})
	return keys
}

func (dict *SyncDict) Clear() {
	*dict = SyncDict{}
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}
