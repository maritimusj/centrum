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
}

func New() *data {
	return &data{
		lockerMap: make(map[interface{}]*sync.Mutex),
		done:      make(chan struct{}),
	}
}

func (data *data) Close() {
	close(data.done)
	data.wg.Wait()
}

func (data *data) Do(obj interface{}, fn func() interface{}) <-chan interface{} {
	data.mu.Lock()
	defer func() {
		data.mu.Unlock()
	}()

	v, ok := data.lockerMap[obj]
	if !ok {
		v = &sync.Mutex{}
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
