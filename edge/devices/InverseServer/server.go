package InverseServer

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/modbus"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

type Server struct {
	addr string
	port int

	lsr     net.Listener
	connMap sync.Map

	done chan struct{}
	wg   sync.WaitGroup
}

type MAC [6]uint8

func (mac *MAC) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func New() *Server {
	return &Server{
		done: make(chan struct{}),
	}
}

func (server *Server) Start(ctx context.Context, addr string, port int) error {
	var err error
	server.lsr, err = net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return lang.InternalError(err)
	}

	log.Tracef("[inverse] start at %s:%d", addr, port)

	server.wg.Add(2)
	go func() {
		defer func() {
			server.wg.Done()
			log.Trace("[inverse] listen routine exit")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-server.done:
				return
			default:
				conn, err := server.lsr.Accept()
				if err != nil {
					continue
				}
				server.wg.Add(1)
				go server.handler(ctx, conn)
			}
		}
	}()

	go func() {
		defer func() {
			server.connMap.Range(func(mac, v interface{}) bool {
				_ = v.(net.Conn).Close()
				server.connMap.Delete(mac)
				return true
			})

			if server.lsr != nil {
				_ = server.lsr.Close()
				server.lsr = nil
			}

			server.wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-server.done:
				return
			}
		}
	}()

	return nil
}

func (server *Server) handler(ctx context.Context, conn net.Conn) {
	defer server.wg.Done()

	log.Debug("[inverse] handler new conn:", conn.RemoteAddr().String())

	var buf [64]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		log.Println(err)
		return
	}

	log.Debug("[inverse] read:", string(buf[0:n]))

	handler := modbus.NewTCPClientHandlerFrom(conn)
	handler.IdleTimeout = 0

	client := modbus.NewClient(handler)

	data, err := client.ReadHoldingRegisters(44, 6)
	if err != nil {
		log.Errorln(err)
		return
	}

	var mac MAC
	for i := range mac {
		mac[i] = byte(binary.BigEndian.Uint16(data[i*2:]))
	}

	log.Debug("[inverse] mac: ", mac.String())

	old, ok := server.connMap.LoadOrStore(mac.String(), conn)
	if ok {
		_ = old.(net.Conn).Close()
		server.connMap.Store(mac.String(), conn)
	}
}

func (server *Server) Close() {
	close(server.done)
	server.wg.Wait()
}

func (server *Server) Wait() {
	server.wg.Wait()
}

func (server *Server) Try(_ context.Context, mac string) (net.Conn, error) {
	var conn net.Conn
	server.connMap.Range(func(key, value interface{}) bool {
		if key.(string) == mac {
			server.connMap.Delete(key)
			conn = value.(net.Conn)
			log.Trace("[inverse] new connection: ", mac, conn.RemoteAddr().String())
			return false
		}
		return true
	})

	if conn != nil {
		return conn, nil
	}

	return nil, fmt.Errorf("[inverse]mac not found: %s", mac)
}
