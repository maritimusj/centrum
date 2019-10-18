package ep6v2

import (
	"bytes"
	"encoding/binary"
	"time"
)

const (
	realtimeDataStartAddress  = 4106
	realtimeStateStartAddress = 8202
)

type RealTimeData struct {
	chNum *CHNum
	data  bytes.Buffer
	ready bytes.Buffer

	lastReadTime time.Time
}

func (r *RealTimeData) Clone() *RealTimeData {
	rt := &RealTimeData{
		chNum:        r.chNum,
		lastReadTime: r.lastReadTime,
	}
	rt.data.Write(r.data.Bytes())
	rt.ready.Write(r.ready.Bytes())
	return rt
}

func (r *RealTimeData) expired() bool {
	return time.Now().Sub(r.lastReadTime) > 1*time.Second
}

func (r *RealTimeData) AINum() int {
	return r.chNum.AI
}

func (r *RealTimeData) DINum() int {
	return r.chNum.DI
}

func (r *RealTimeData) AONum() int {
	return r.chNum.AO
}

func (r *RealTimeData) DONum() int {
	return r.chNum.DO
}

func (r *RealTimeData) VONum() int {
	return r.chNum.VO
}

func (r *RealTimeData) GetAIValue(index int, point int) (float32, bool) {
	v, ok := r.getFloat32(index)
	return ToFloat32(v, point), ok
}

func (r *RealTimeData) GetDIValue(index int) (bool, bool) {
	v, ready := r.getBool(r.chNum.AI + index)
	if ready {
		return v, true
	}
	return false, false
}

func (r *RealTimeData) GetDOValue(index int) (bool, bool) {
	v, ready := r.getBool(r.chNum.AI + r.chNum.DI + index)
	if ready {
		return v, true
	}
	return false, false
}

func (r *RealTimeData) SetDOValue(index int, value bool) {
	var v uint32
	if value {
		v = 1
	} else {
		v = 0
	}
	r.setInt(r.chNum.AI + r.chNum.DI + index, v)
}

func (r *RealTimeData) GetAOValue(index int) (float32, bool) {
	return r.getFloat32(r.chNum.AI + r.chNum.DI + r.chNum.DO + index)
}

func (r *RealTimeData) GetVOValue(index int) (bool, bool) {
	v, ready := r.getBool(r.chNum.AI + r.chNum.DI + r.chNum.DO + r.chNum.AO + index)
	if ready {
		return v, true
	}
	return false, false
}

func (r *RealTimeData) getFloat32(index int) (float32, bool) {
	pos := index * 4
	if r.data.Len() > pos+4 && index < r.ready.Len() && r.ready.Bytes()[index] == 0 {
		return ToSingle(r.data.Bytes()[pos:]), true
	}
	return 0, false
}

func (r *RealTimeData) getInt(index int) (int, bool) {
	pos := index * 4
	if r.data.Len() > pos+4 && index < r.ready.Len() && r.ready.Bytes()[index] == 0 {
		return int(binary.BigEndian.Uint32(r.data.Bytes()[pos:])), true
	}
	return 0, false
}

func (r *RealTimeData) setInt(index int, value uint32) {
	pos := index * 4
	if r.data.Len() > pos+4 && index < r.ready.Len() && r.ready.Bytes()[index] == 0 {
		binary.BigEndian.PutUint32(r.data.Bytes()[pos:], value)
	}
}

func (r *RealTimeData) getBool(index int) (bool, bool) {
	v, ready := r.getInt(index)
	if ready {
		return v > 0, ready
	}
	return false, false
}

func (r *RealTimeData) fetchData(conn modbusClient) error {
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

	r.data.Truncate(0)
	r.ready.Truncate(0)

	var address uint16 = realtimeDataStartAddress
	var quantity uint16
	var amount = total

	for amount > 0 {
		if amount > 124 {
			quantity = 124
		} else {
			quantity = uint16(amount)
		}
		data, err := conn.ReadInputRegisters(address, quantity)
		if err != nil {
			return err
		}

		r.data.Write(data)

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
		state, err := conn.ReadInputRegisters(address, quantity)
		if err != nil {
			return err
		}

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
