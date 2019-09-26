package influxDB

import (
	"context"
	"github.com/maritimusj/centrum/edge/adapter"
	"io"

	db "github.com/influxdata/influxdb1-client/v2"
)

type client struct {
	db db.Client
}

func New() *client {
	return &client{}
}

func (c *client) Open(context.Context, context.Context, map[string]interface{}) error {
	panic("implement me")
}

func (c *client) Close() error {
	panic("implement me")
}

func (c *client) Add(string, adapter.Client) (io.Closer, error) {
	panic("implement me")
}

func (c *client) Stats() (map[string]interface{}, error) {
	panic("implement me")
}
