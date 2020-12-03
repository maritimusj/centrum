package measure

import (
	"sync"
	"time"
)

var (
	defaultMeasureDataPool = &sync.Pool{
		New: func() interface{} {
			return &Data{
				Tags:   map[string]string{},
				Fields: map[string]interface{}{},
			}
		},
	}
)

type Data struct {
	Name   string                 `json:"name"`
	Tags   map[string]string      `json:"tags"`
	Fields map[string]interface{} `json:"fields"`
	Time   time.Time              `json:"time"`
	pool   *sync.Pool
}

func New(name string) *Data {
	data := defaultMeasureDataPool.Get().(*Data)
	data.Time = time.Now()
	data.Name = name
	data.pool = defaultMeasureDataPool
	return data
}

func (measure *Data) Clone() *Data {
	data := New(measure.Name)
	data.Time = measure.Time

	for k, v := range measure.Tags {
		data.Tags[k] = v
	}
	for k, v := range measure.Fields {
		data.Fields[k] = v
	}
	return data
}

func (measure *Data) Release() {
	measure.Name = ""
	measure.Tags = map[string]string{}
	measure.Fields = map[string]interface{}{}
	measure.pool.Put(measure)
}

func (measure *Data) AddTag(name, val string) *Data {
	measure.Tags[name] = val
	return measure
}

func (measure *Data) GetTag(name string) (interface{}, bool) {
	v, ok := measure.Tags[name]
	return v, ok
}

func (measure *Data) AddField(name string, val interface{}) *Data {
	measure.Fields[name] = val
	return measure
}

func (measure *Data) GetField(name string) (interface{}, bool) {
	v, ok := measure.Fields[name]
	return v, ok
}
