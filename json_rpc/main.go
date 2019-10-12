package json_rpc

import (
	"net/http"
	"time"
)

type Sink interface {
	GetBaseInfo(uid string) (map[string]interface{}, error)
	Remove(uid string)
	Active(conf *Conf) error
	SetValue(val *Value) error
	GetValue(ch *CH) (interface{}, error)
	GetRealtimeData(uid string) ([]interface{}, error)
}

type Edge struct {
	sink Sink
}

type Conf struct {
	UID              string
	Inverse          bool
	Address          string
	Interval         time.Duration
	DB               string
	InfluxDBAddress  string
	InfluxDBUserName string
	InfluxDBPassword string
	CallbackURL      string
	LogLevel         string
}

type CH struct {
	UID string
	Tag string
}

type Value struct {
	CH
	V interface{}
}

type Result struct {
	Code int
	Msg  string
	Data interface{}
}

func New(sink Sink) *Edge {
	return &Edge{
		sink: sink,
	}
}

//GetBaseInfo 获取设备基本信息
func (e *Edge) GetBaseInfo(r *http.Request, uid *string, result *Result) error {
	data, err := e.sink.GetBaseInfo(*uid)
	if err != nil {
		return err
	}
	result.Data = data
	return nil
}

//Active 用于激活一个设备
func (e *Edge) Active(r *http.Request, conf *Conf, result *Result) error {
	return e.sink.Active(conf)
}

//Remove 移除一个设备，不再读取相关数据
func (e *Edge) Remove(r *http.Request, uid *string, result *Result) error {
	e.sink.Remove(*uid)
	return nil
}

//SetValue 设置设备一个点位值，目前只支持DO
func (e *Edge) SetValue(r *http.Request, val *Value, result *Result) error {
	return e.sink.SetValue(val)
}

//GetValue 获取设备的一个点位值
func (e *Edge) GetValue(r *http.Request, ch *CH, result *Result) error {
	v, err := e.sink.GetValue(ch)
	if err != nil {
		return err
	}
	result.Data = v
	return nil
}

//GetRealtimeData 获取设备全部的即时数据
func (e *Edge) GetRealtimeData(r *http.Request, uid *string, result *Result) error {
	data, err := e.sink.GetRealtimeData(*uid)
	if err != nil {
		return err
	}
	result.Data = data
	return nil
}
