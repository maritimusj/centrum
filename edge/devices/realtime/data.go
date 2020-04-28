package realtime

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/maritimusj/centrum/edge/devices/CHNum"
	"github.com/maritimusj/centrum/edge/devices/modbus"
	"github.com/maritimusj/centrum/edge/devices/util"
)

const (
	realtimeDataStartAddress  = 4106
	realtimeStateStartAddress = 8202
)

var (
	defaultRealTimeDataPool = &sync.Pool{
		New: func() interface{} {
			return &Data{}
		},
	}
)

type Data struct {
	chNum *CHNum.Data
	data  bytes.Buffer
	ready bytes.Buffer

	pool *sync.Pool

	timeUsed     time.Duration
	lastReadTime time.Time
}

func New(chNum *CHNum.Data) *Data {
	data := defaultRealTimeDataPool.Get().(*Data)
	data.chNum = chNum.Clone()
	data.pool = defaultRealTimeDataPool
	return data
}

func (r *Data) Release() {
	if r.chNum != nil {
		r.chNum.Release()
	}

	r.data.Reset()
	r.ready.Reset()

	r.pool.Put(r)
}

func (r *Data) Clone() *Data {
	rt := New(r.chNum)
	rt.lastReadTime = r.lastReadTime
	rt.timeUsed = r.timeUsed
	rt.data.Write(r.data.Bytes())
	rt.ready.Write(r.ready.Bytes())
	return rt
}

func (r *Data) expired() bool {
	return time.Now().Sub(r.lastReadTime) > 1*time.Second
}

func (r *Data) CHNum() *CHNum.Data {
	return r.chNum
}

func (r *Data) TimeUsed() time.Duration {
	return r.timeUsed
}

func (r *Data) AINum() int {
	return r.chNum.AI
}

func (r *Data) DINum() int {
	return r.chNum.DI
}

func (r *Data) AONum() int {
	return r.chNum.AO
}

func (r *Data) DONum() int {
	return r.chNum.DO
}

func (r *Data) VONum() int {
	return r.chNum.VO
}

func (r *Data) GetAIValue(index int, point int) (float32, bool) {
	if v, ready := r.getFloat32(index); ready {
		return util.ToFloat32(v, point), true
	}

	return 0, false
}

func (r *Data) GetDIValue(index int) (bool, bool) {
	v, ready := r.getBool(r.chNum.AI + index)
	if ready {
		return v, true
	}
	return false, false
}

func (r *Data) GetDOValue(index int) (bool, bool) {
	v, ready := r.getBool(r.chNum.AI + r.chNum.DI + index)
	if ready {
		return v, true
	}
	return false, false
}

func (r *Data) SetDOValue(index int, value bool) {
	var v uint32
	if value {
		v = 1
	} else {
		v = 0
	}
	r.setInt(r.chNum.AI+r.chNum.DI+index, v)
}

func (r *Data) GetAOValue(index int) (float32, bool) {
	return r.getFloat32(r.chNum.AI + r.chNum.DI + r.chNum.DO + index)
}

func (r *Data) GetVOValue(index int) (bool, bool) {
	v, ready := r.getBool(r.chNum.AI + r.chNum.DI + r.chNum.DO + r.chNum.AO + index)
	if ready {
		return v, true
	}
	return false, false
}

func (r *Data) getFloat32(index int) (float32, bool) {
	pos := index * 4
	if r.data.Len() >= pos+4 && r.ready.Len() > index && r.ready.Bytes()[index] == 0 {
		return util.ToSingle(r.data.Bytes()[pos:]), true
	}

	return 0, false
}

func (r *Data) getInt(index int) (int, bool) {
	pos := index * 4
	if r.data.Len() >= pos+4 && index < r.ready.Len() && r.ready.Bytes()[index] == 0 {
		return int(binary.BigEndian.Uint32(r.data.Bytes()[pos:])), true
	}
	return 0, false
}

func (r *Data) setInt(index int, value uint32) {
	pos := index * 4
	if r.data.Len() >= pos+4 && index < r.ready.Len() && r.ready.Bytes()[index] == 0 {
		binary.BigEndian.PutUint32(r.data.Bytes()[pos:], value)
	}
}

func (r *Data) getBool(index int) (bool, bool) {
	v, ready := r.getInt(index)
	if ready {
		return v > 0, ready
	}
	return false, false
}

func (r *Data) FetchData(conn modbus.Client) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("unexpect error: %#v", err)
			return
		}
	}()

	if !r.expired() {
		return nil
	}

	total := r.chNum.Sum() * 2
	if r.data.Len() < total {
		r.data.Grow(total)
	}

	if r.ready.Len() < r.chNum.Sum() {
		r.ready.Grow(r.chNum.Sum())
	}

	r.timeUsed = 0
	r.data.Reset()
	r.ready.Reset()

	var (
		address  uint16 = realtimeDataStartAddress
		quantity uint16
		amount   = total
	)

	for amount > 0 {
		if amount > 124 {
			quantity = 124
		} else {
			quantity = uint16(amount)
		}

		data, used, err := conn.ReadInputRegisters(address, quantity)
		if err != nil {
			return err
		}

		r.data.Write(data)
		r.timeUsed = r.timeUsed + used

		amount -= int(quantity)
		address += quantity
	}

	//读取数据有效状态
	address = realtimeStateStartAddress
	amount = r.chNum.Sum()
	for amount > 0 {
		if amount > 124 {
			quantity = 124
		} else {
			quantity = uint16(amount)
		}

		state, used, err := conn.ReadInputRegisters(address, quantity)
		if err != nil {
			return err
		}

		r.timeUsed = r.timeUsed + used

		var i uint16
		for i = 0; i < quantity; i++ {
			reading := byte(binary.BigEndian.Uint16(state[i*2:]))
			r.ready.Write([]byte{reading})
		}

		amount -= int(quantity)
		address += quantity
	}

	r.lastReadTime = time.Now()
	return nil
}
