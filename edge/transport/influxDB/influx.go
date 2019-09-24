package influxDB

import (
	"context"
	"github.com/maritimusj/centrum/edge/adapter"
	"github.com/maritimusj/centrum/edge/transport"
	"io"
)

type client struct{}

func New() transport.Client {
	return &client{}
}

func (c *client) Open(context.Context, transport.Option) error {
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
