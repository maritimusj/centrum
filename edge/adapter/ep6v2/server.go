package ep6v2

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/maritimusj/modbus"
	log "github.com/sirupsen/logrus"

	"github.com/maritimusj/chuanyan/gate/config"
)

type InverseServer struct {
	address string
	port    int

	lsr     net.Listener
	connMap sync.Map

	wg sync.WaitGroup
}

type MAC [6]uint8

func (mac *MAC) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

var (
	defaultInverseServer = NewInverseServer("", Config.DefaultInverseServerPort)
)

func GetInverseServer() *InverseServer {
	return defaultInverseServer
}

func StartDefaultServer(ctx context.Context, addr string, port int) error {
	defaultInverseServer.address = addr
	defaultInverseServer.port = port
	return defaultInverseServer.Start(ctx)
}

func NewInverseServer(addr string, port int) *InverseServer {
	return &InverseServer{
		address: addr,
		port:    port,
	}
}

func (server *InverseServer) Close() {
	_ = server.lsr.Close()

	server.connMap.Range(func(mac, v interface{}) bool {
		_ = v.(net.Conn).Close()
		server.connMap.Delete(mac)
		return true
	})
}

func (server *InverseServer) Wait() {
	server.wg.Wait()
	//log.Trace("ep6v2 inverse server server shutdown.")
}

func (server *InverseServer) Try(mac string) (net.Conn, error) {
	var conn net.Conn
	server.connMap.Range(func(key, value interface{}) bool {
		if key.(string) == mac {
			server.connMap.Delete(key)
			conn = value.(net.Conn)
			return false
		}
		return true
	})

	if conn != nil {
		return conn, nil
	}

	return nil, fmt.Errorf("mac not found: %s", mac)
}

func (server *InverseServer) Start(ctx context.Context) error {
	var err error
	server.lsr, err = net.Listen("tcp", fmt.Sprintf("%s:%d", server.address, server.port))
	if err != nil {
		return err
	}

	server.wg.Add(2)
	go func() {
		defer server.wg.Done()
		log.Trace("ep6v2 inverse server start listening on :", server.lsr.Addr().String())
		for {
			select {
			case <-ctx.Done():
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
		defer server.wg.Done()
		for {
			select {
			case <-ctx.Done():
				if server.lsr != nil {
					_ = server.lsr.Close()
				}
				return
			}
		}
	}()

	return nil
}

func (server *InverseServer) handler(ctx context.Context, conn net.Conn) {
	defer server.wg.Done()

	log.Debug("ep6v2 inverse server: handler new conn:", conn.RemoteAddr().String())

	var buf [64]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		log.Println(err)
		return
	}

	log.Debug("ep6v2 inverse server: read:", string(buf[0:n]))

	handler := modbus.NewTCPClientHandlerFrom(conn)
	handler.IdleTimeout = 0

	client := modbus.NewClient(handler)

	data, err := client.ReadHoldingRegisters(44, 6)
	if err != nil {
		log.Error(err)
		return
	}

	var mac MAC
	for i := range mac {
		mac[i] = byte(binary.BigEndian.Uint16(data[i*2:]))
	}

	log.Debug("ep6v2 inverse server: 6mac: ", mac.String())

	old, ok := server.connMap.LoadOrStore(mac.String(), conn)
	if ok {
		_ = old.(net.Conn).Close()
	}

	//for {
	//	select {
	//	case <-ctx.Done():
	//		return
	//	case <-time.After(10 * time.Second):
	//		log.Println("ep6v2 inverse server: heart beat:", conn.RemoteAddr().String())
	//		_, err := client.ReadHoldingRegisters(0, 1)
	//		if err != nil {
	//			log.Println(err)
	//			return
	//		}
	//	}
	//}
}
