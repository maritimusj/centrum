package ep6v2

import (
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

func (c *TCPConnector) Try(addr string) (net.Conn, error) {
	dialer := net.Dialer{Timeout: c.timeout}
	println("addr:", addr)
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
