package global

import (
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	Params = New()
	Stats  = New()
)

type stats struct {
	data []byte
	sync.RWMutex
}

func New() *stats {
	return &stats{}
}

func (stats *stats) Set(path string, val interface{}) error {
	stats.Lock()
	defer stats.Unlock()

	data, err := sjson.SetBytes(stats.data, path, val)
	if err != nil {
		return err
	}
	stats.data = data
	return nil
}

func (stats *stats) Exists(path string) bool {
	stats.RLock()
	defer stats.RUnlock()

	return gjson.GetBytes(stats.data, path).Exists()
}

func (stats *stats) Get(path string) (interface{}, bool) {
	stats.RLock()
	defer stats.RUnlock()

	v := gjson.GetBytes(stats.data, path)
	return v.Value(), v.Exists()
}

func (stats *stats) MustGet(path string) interface{} {
	stats.RLock()
	defer stats.RUnlock()

	v := gjson.GetBytes(stats.data, path)
	if !v.Exists() {
		panic("stats not exists")
	}

	return v.Value()
}

func (stats *stats) Remove(path string) {
	stats.Lock()
	defer stats.Unlock()

	v, err := sjson.DeleteBytes(stats.data, path)
	if err != nil {
		return
	}
	stats.data = v
}

func (stats *stats) Reset() {
	stats.data = []byte{}
}
