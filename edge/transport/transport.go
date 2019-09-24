package transport

import (
	"github.com/maritimusj/centrum/edge/adapter"

	"context"
	"io"
)

type Option map[string]interface{}

type Client interface {
	Open(context.Context, Option) error
	Close() error

	Add(string, adapter.Client) (io.Closer, error)
	Stats() (map[string]interface{}, error)
}
