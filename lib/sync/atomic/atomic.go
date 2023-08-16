package atomic

import "sync/atomic"

type AtomicBool uint32

func (bool *AtomicBool) Get() bool {
	return atomic.LoadUint32((*uint32)(bool)) != 0
}

func (bool *AtomicBool) Set(val bool) {
	if val {
		atomic.StoreUint32((*uint32)(bool), 1)
	} else {
		atomic.StoreUint32((*uint32)(bool), 0)
	}
}
