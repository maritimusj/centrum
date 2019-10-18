package synchronized

import (
	"errors"
	"sync"
)

var (
	defaultSynchronized = New()
)

func Close() {
	defaultSynchronized.Close()
}

func Do(obj interface{}, fn func() interface{}) <-chan interface{} {
	return defaultSynchronized.Do(obj, fn)
}

type data struct {
	lockerMap map[interface{}]*sync.Mutex
	mu        sync.Mutex
	done      chan struct{}
	wg        sync.WaitGroup
	pool      *sync.Pool
}

func New() *data {
	return &data{
		lockerMap: make(map[interface{}]*sync.Mutex),
		done:      make(chan struct{}),
		pool: &sync.Pool{
			New: func() interface{} {
				return &sync.Mutex{}
			},
		},
	}
}

func (data *data) Close() {
	close(data.done)
	data.wg.Wait()
}

func (data *data) Do(obj interface{}, fn func() interface{}) <-chan interface{} {
	data.mu.Lock()
	defer data.mu.Unlock()

	v, ok := data.lockerMap[obj]
	if !ok {
		v = data.pool.Get().(*sync.Mutex)
		data.lockerMap[obj] = v
	}

	resultChan := make(chan interface{}, 1)
	go func() {
		data.wg.Add(1)
		v.Lock()

		defer func() {
			close(resultChan)
			data.wg.Done()
			v.Unlock()

			data.mu.Lock()
			{
				delete(data.lockerMap, obj)
				data.pool.Put(v)
			}
			data.mu.Unlock()
		}()

		select {
		case <-data.done:
			resultChan <- errors.New("synchronized Do() exit")
			return
		default:
			resultChan <- fn()
			return
		}
	}()

	return resultChan
}
