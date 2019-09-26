package adapter

import "context"

type Type int

const (
	_ Type = iota
	Ep6v2
)

type MeasureData struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
}

type CtrlData struct {
	Values map[string]interface{}
	Error  chan error
}

type Client interface {
	Open(ctx context.Context, option map[string]interface{}) error
	Close() error

	//创建数据接收通道
	Create(option map[string]interface{}) (<-chan *MeasureData, error)

	//关闭数据接收通道
	Drop(<-chan *MeasureData)

	//插入控制通道
	Plug(<-chan *CtrlData) error

	//报告状态
	Stats() (map[string]interface{}, error)

	//设置参数
	Set(params map[string]interface{}) error

	//获取状态参数当前值
	Get(keys ...string) (map[string]interface{}, error)
}
