package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

// 一致性哈希节点映射
type Map struct {
	hashFunc HashFunc
	keys     []int // sorted
	hashMap  map[int]string
}

func New(fn HashFunc) *Map {
	m := &Map{
		hashFunc: fn,
		hashMap:  make(map[int]string), // 虚拟节点 hash 值到物理节点地址的映射
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := m.hashFunc([]byte(key))
		m.keys = append(m.keys, int(hash))
		m.hashMap[int(hash)] = key
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if key == "" {
		return ""
	}
	hash := m.hashFunc([]byte(key))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= int(hash)
	})
	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}
