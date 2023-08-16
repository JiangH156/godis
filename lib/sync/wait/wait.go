package wait

import (
	"sync"
	"time"
)

type Wait struct {
	wg sync.WaitGroup
}

func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

func (w *Wait) Done() {
	w.wg.Done()
}

func (w *Wait) Wait() {
	w.wg.Wait()
}
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		w.wg.Wait()
		c <- struct{}{}
	}()
	select {
	case <-c: // 正常等待结束
		return true
	case <-time.After(timeout):
		return false
	}
}
