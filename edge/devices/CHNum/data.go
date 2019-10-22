package CHNum

import (
	"encoding/binary"
	"github.com/maritimusj/centrum/edge/devices/modbus"
	"sync"
)

var (
	defaultCHNumPool = &sync.Pool{
		New: func() interface{} {
			return &Data{}
		},
	}
)

type Data struct {
	AI   int
	DI   int
	DO   int
	AO   int
	VO   int
	pool *sync.Pool
}

func New() *Data {
	data := defaultCHNumPool.New().(*Data)
	data.pool = defaultCHNumPool
	return data
}

func (ch *Data) Release() {
	ch.AI = 0
	ch.AO = 0
	ch.DI = 0
	ch.DO = 0
	ch.VO = 0
	ch.pool.Put(ch)
}

func (ch *Data) Sum() int {
	return ch.AI + ch.DI + ch.DO + ch.AO + ch.VO
}

func (ch *Data) Clone() *Data {
	data := New()
	data.AI = ch.AI
	data.AO = ch.AO
	data.DI = ch.DI
	data.DO = ch.DO
	data.VO = ch.VO
	return data
}

func (ch *Data) FetchData(conn modbus.Client) error {
	data, err := conn.ReadHoldingRegisters(16, 5)
	if err != nil {
		return err
	}

	ch.AI = int(binary.BigEndian.Uint16(data[0:]))
	ch.DI = int(binary.BigEndian.Uint16(data[2:]))
	ch.DO = int(binary.BigEndian.Uint16(data[4:]))
	ch.AO = int(binary.BigEndian.Uint16(data[6:]))
	ch.VO = int(binary.BigEndian.Uint16(data[8:]))

	return nil
}
