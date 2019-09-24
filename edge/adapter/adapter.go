package adapter

import "context"

type MeasureData struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
}

type CtrlData struct {
	Values map[string]interface{}
	Error  chan error
}

type Option map[string]interface{}

type Client interface {
	Open(ctx context.Context, option Option) error
	Close() error

	//创建数据接收通道
	Create(option Option) (<-chan *MeasureData, error)

	//关闭数据接收通道
	Drop(<-chan *MeasureData)

	//插入控制通道
	Plug(<-chan *CtrlData) error

	//报告状态
	Stats() (map[string]interface{}, error)

	//设置参数
	Set(params Option) error

	//获取状态参数当前值
	Get(keys ...string) (map[string]interface{}, error)
}
