package ep6v2

import (
	"net"
	"sync"
	"time"

	"github.com/maritimusj/centrum/util"

	"github.com/maritimusj/modbus"
)

type modbusWrapper struct {
	client modbus.Client
	sync.Mutex
}

func (w *modbusWrapper) retry(fn func() ([]byte, error)) (result []byte, used time.Duration, err error) {
	for i := 0; i < 3; i++ {
		begin := time.Now()
		result, err = fn()
		if err != nil {
			if e, ok := err.(net.Error); ok && (e.Temporary() || e.Timeout()) {
				time.Sleep(time.Duration(util.Exponent(10, uint64(i+1))) * time.Millisecond)
				continue
			}

			if e, ok := err.(modbus.Error); ok && e.Temporary() {
				time.Sleep(time.Duration(util.Exponent(10, uint64(i+1))) * time.Millisecond)
				continue
			}
		}
		return result, time.Now().Sub(begin), nil
	}
	used = -1
	return
}

func (w *modbusWrapper) ReadCoils(address, quantity uint16) (results []byte, duration time.Duration, err error) {
	w.Lock()
	defer w.Unlock()
	return w.retry(func() (bytes []byte, err error) {
		return w.client.ReadCoils(address, quantity)
	})
}

func (w *modbusWrapper) ReadDiscreteInputs(address, quantity uint16) (results []byte, duration time.Duration, err error) {
	w.Lock()
	defer w.Unlock()
	return w.retry(func() (bytes []byte, err error) {
		return w.client.ReadDiscreteInputs(address, quantity)
	})
}

func (w *modbusWrapper) WriteSingleCoil(address, value uint16) (results []byte, duration time.Duration, err error) {
	w.Lock()
	defer w.Unlock()
	return w.retry(func() (bytes []byte, err error) {
		return w.client.WriteSingleCoil(address, value)
	})
}

func (w *modbusWrapper) ReadInputRegisters(address, quantity uint16) (results []byte, duration time.Duration, err error) {
	w.Lock()
	defer w.Unlock()
	return w.retry(func() (bytes []byte, err error) {
		return w.client.ReadInputRegisters(address, quantity)
	})
}

func (w *modbusWrapper) ReadHoldingRegisters(address, quantity uint16) (results []byte, duration time.Duration, err error) {
	w.Lock()
	defer w.Unlock()
	return w.retry(func() (bytes []byte, err error) {
		return w.client.ReadHoldingRegisters(address, quantity)
	})
}
