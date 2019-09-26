package adapter

import (
	"context"
	"github.com/maritimusj/centrum/edge"
)

type Type int

const (
	_ Type = iota
	Ep6v2
)

type Client interface {
	Open(ctx context.Context, option map[string]interface{}) error
	Close() error

	//创建数据接收通道
	Create(option map[string]interface{}) (<-chan *edge.MeasureData, error)

	//关闭数据接收通道
	Drop(<-chan *edge.MeasureData)

	//插入控制通道
	Plug(<-chan *edge.CtrlData) error

	//报告状态
	Stats() (map[string]interface{}, error)

	//设置参数
	Set(params map[string]interface{}) error

	//获取状态参数当前值
	Get(keys ...string) (map[string]interface{}, error)
}
