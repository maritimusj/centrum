package synchronized

import (
	"sync"
	"sync/atomic"
)

var (
	defaultSynchronized = New()
)

type mutexWrapper struct {
	mu    sync.Mutex
	total int32
}

func (w *mutexWrapper) Lock() {
	atomic.AddInt32(&w.total, 1)
	w.mu.Lock()
}

func (w *mutexWrapper) Unlock() {
	w.mu.Unlock()
	atomic.AddInt32(&w.total, -1)
}

func (w *mutexWrapper) IsIdle() bool {
	return atomic.LoadInt32(&w.total) < 0
}

func Close() {
	defaultSynchronized.Close()
}

func Do(obj interface{}, fn func() interface{}) <-chan interface{} {
	return defaultSynchronized.Do(obj, fn)
}

type data struct {
	lockerMap map[interface{}]*mutexWrapper
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func New() *data {
	return &data{
		lockerMap: make(map[interface{}]*mutexWrapper),
	}
}

func (data *data) Close() {
	data.wg.Wait()
}

func (data *data) Do(obj interface{}, fn func() interface{}) <-chan interface{} {
	data.mu.Lock()
	defer data.mu.Unlock()

	v, ok := data.lockerMap[obj]
	if !ok {
		v = &mutexWrapper{}
		data.lockerMap[obj] = v
	}

	resultChan := make(chan interface{}, 1)
	data.wg.Add(1)

	go func() {
		v.Lock()
		defer func() {
			close(resultChan)
			v.Unlock()
			data.mu.Lock()
			if v.IsIdle() {
				delete(data.lockerMap, obj)
			}
			data.mu.Unlock()
			data.wg.Done()
		}()

		resultChan <- fn()
	}()

	return resultChan
}
