package InverseServer

import (
	"context"
)

var (
	defaultServer = New()
)

func DefaultConnector() *Server {
	return defaultServer
}

func Start(ctx context.Context, addr string, port int) error {
	if defaultServer.lsr != nil {
		defaultServer.Close()
		defaultServer = New()
	}

	return defaultServer.Start(ctx, addr, port)
}

func Close() {
	defaultServer.Close()
}
