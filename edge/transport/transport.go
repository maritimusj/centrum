package transport

import (
	"github.com/maritimusj/centrum/edge/adapter"

	"context"
	"io"
)

type Type int

const (
	_ Type = iota
	InfluxDB
)

type Client interface {
	Open(context.Context, map[string]interface{}) error
	Close() error

	Add(string, adapter.Client) (io.Closer, error)
	Stats() (map[string]interface{}, error)
}
