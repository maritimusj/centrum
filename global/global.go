package global

import "sync"

var (
	Params = New()
	Stats = New()
)

type stats struct {
	data sync.Map
}

func New() *stats {
	return &stats{}
}

func (stats *stats) Set(name string, val interface{}) *stats {
	stats.data.Store(name, val)
	return stats
}

func (stats *stats) Get(name string)(interface{}, bool) {
	return stats.data.Load(name)
}

func(stats *stats)  MustGet(name string)interface{} {
	if v, ok := stats.data.Load(name); ok {
		return v
	}
	panic("stats not exists")
}

func (stats *stats) Remove(name string) *stats {
	stats.data.Delete(name)
	return stats
}
