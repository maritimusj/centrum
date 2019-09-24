package event

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

type CallbackFN func(ctx context.Context, code int, values map[string]interface{})

type Data struct {
	Code   int
	Values map[string]interface{}
}

func (data *Data) Set(key string, v interface{}) {
	data.Values[key] = v
}

func (data *Data) Get(key string) interface{} {
	if v, ok := data.Values[key]; ok {
		return v
	}
	return nil
}

func (data *Data) Pop(key string) interface{} {
	if v, ok := data.Values[key]; ok {
		delete(data.Values, key)
		return v
	}
	return nil
}

func (data *Data) Clone() (*Data, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var clone Data
	err = json.Unmarshal(jsonData, &clone)
	if err != nil {
		return nil, err
	}
	return &clone, nil
}

type Event interface {
	Fire(data *Data)
	Sub(ctx context.Context, codes ...int) <-chan *Data
	Register(ctx context.Context, fn CallbackFN, codes ...int)
	Wait()
}

type eventX struct {
	callbackFNMap map[int][]*fnPair
	channelsMap   map[int][]*chPair
	mutex         sync.RWMutex
	wg            sync.WaitGroup
}

type fnPair struct {
	cb  CallbackFN
	ctx context.Context
}

type chPair struct {
	ch  chan *Data
	ctx context.Context
}

func New() Event {
	return &eventX{
		callbackFNMap: map[int][]*fnPair{},
		channelsMap:   map[int][]*chPair{},
	}
}

func (e *eventX) Wait() {
	e.wg.Wait()
}

func (e *eventX) Sub(ctx context.Context, codes ...int) <-chan *Data {
	if len(codes) == 0 {
		panic(errors.New("event register: zero code"))
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	ch := make(chan *Data, 10)

	for _, code := range codes {
		e.channelsMap[code] = append(e.channelsMap[code], &chPair{
			ch:  ch,
			ctx: ctx,
		})
	}

	return ch
}

func (e *eventX) Register(ctx context.Context, fn CallbackFN, codes ...int) {
	if len(codes) == 0 {
		panic(errors.New("event register: zero code"))
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, code := range codes {
		e.callbackFNMap[code] = append(e.callbackFNMap[code], &fnPair{
			cb:  fn,
			ctx: ctx,
		})
	}
}

func (e *eventX) Fire(data *Data) {
	if data != nil {
		e.mutex.RLock()
		e.mutex.RUnlock()

		if v, ok := e.channelsMap[data.Code]; ok && len(v) > 0 {
			for _, pair := range v {
				e.wg.Add(1)
				go func() {
					defer e.wg.Done()
					data, err := data.Clone()
					if err != nil {
						log.Errorf("event fire, data clone failed: %s", err)
					} else {
						select {
						case <-pair.ctx.Done():
							return
						default:
							pair.ch <- data
						}
					}
				}()
			}
		}

		if v, ok := e.callbackFNMap[data.Code]; ok && len(v) > 0 {
			for _, pair := range v {
				if pair.cb == nil {
					continue
				}

				clone, err := data.Clone()
				if err != nil {
					log.Errorf("event fire, data clone failed: %s", err)
					continue
				}

				select {
				case <-pair.ctx.Done():
					return
				default:
					pair.cb(pair.ctx, clone.Code, clone.Values)
				}
			}
		}
	}
}
