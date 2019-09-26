package ep6v2

import (
	"context"
	"github.com/maritimusj/centrum/edge"
)

type ep6v2Adapter struct {
}

func New() *ep6v2Adapter {
	return &ep6v2Adapter{}
}

func (e *ep6v2Adapter) Open(ctx context.Context, option map[string]interface{}) error {
	panic("implement me")
}

func (e *ep6v2Adapter) Close() error {
	panic("implement me")
}

func (e *ep6v2Adapter) Create(option map[string]interface{}) (<-chan *edge.MeasureData, error) {
	panic("implement me")
}

func (e *ep6v2Adapter) Drop(<-chan *edge.MeasureData) {
	panic("implement me")
}

func (e *ep6v2Adapter) Plug(<-chan *edge.CtrlData) error {
	panic("implement me")
}

func (e *ep6v2Adapter) Stats() (map[string]interface{}, error) {
	panic("implement me")
}

func (e *ep6v2Adapter) Set(params map[string]interface{}) error {
	panic("implement me")
}

func (e *ep6v2Adapter) Get(keys ...string) (map[string]interface{}, error) {
	panic("implement me")
}
