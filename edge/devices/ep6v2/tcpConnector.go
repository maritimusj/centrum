package ep6v2

import (
	"context"
	"net"
	"time"
)

type TCPConnector struct {
	timeout time.Duration
}

func NewTCPConnector() Connector {
	return &TCPConnector{
		timeout: 6 * time.Second,
	}
}

func (c *TCPConnector) Try(ctx context.Context, addr string) (net.Conn, error) {
	dialer := net.Dialer{Timeout: c.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
