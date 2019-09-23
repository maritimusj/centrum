package ep6v2

import (
	"context"
	"github.com/maritimusj/centrum/edge/adapter"
)

type ep6v2Adapter struct {
}

func New() adapter.Adapter {
	return &ep6v2Adapter{}
}

func (e *ep6v2Adapter) Open(ctx context.Context, option adapter.Option) error {
	panic("implement me")
}

func (e *ep6v2Adapter) Close() error {
	panic("implement me")
}

func (e *ep6v2Adapter) Create(option adapter.Option) (<-chan *adapter.MeasureData, error) {
	panic("implement me")
}

func (e *ep6v2Adapter) Drop(<-chan *adapter.MeasureData) {
	panic("implement me")
}

func (e *ep6v2Adapter) Plug(<-chan *adapter.CtrlData) error {
	panic("implement me")
}

func (e *ep6v2Adapter) Stats() (map[string]interface{}, error) {
	panic("implement me")
}

func (e *ep6v2Adapter) Set(params adapter.Option) error {
	panic("implement me")
}

func (e *ep6v2Adapter) Get(keys ...string) (map[string]interface{}, error) {
	panic("implement me")
}
