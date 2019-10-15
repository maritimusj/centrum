package global

import "sync"

var (
	Params = New()
	Stats  = New()
)

type stats struct {
	data sync.Map
}

func New() *stats {
	return &stats{}
}

func (stats *stats) Set(name interface{}, val interface{}) *stats {
	stats.data.Store(name, val)
	return stats
}

func (stats *stats) Get(name interface{}) (interface{}, bool) {
	return stats.data.Load(name)
}

func (stats *stats) MustGet(name interface{}) interface{} {
	if v, ok := stats.data.Load(name); ok {
		return v
	}
	panic("stats not exists")
}

func (stats *stats) Remove(name interface{}) *stats {
	stats.data.Delete(name)
	return stats
}
