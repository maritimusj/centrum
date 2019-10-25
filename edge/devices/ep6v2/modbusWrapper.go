package ep6v2

import (
	"github.com/maritimusj/modbus"
	"sync"
)

type modbusWrapper struct {
	client modbus.Client
	sync.Mutex
}

func (w *modbusWrapper) ReadCoils(address, quantity uint16) (results []byte, err error) {
	w.Lock()
	defer w.Unlock()
	return w.client.ReadCoils(address, quantity)
}

func (w *modbusWrapper) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	w.Lock()
	defer w.Unlock()
	return w.client.ReadDiscreteInputs(address, quantity)
}

func (w *modbusWrapper) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	w.Lock()
	defer w.Unlock()
	return w.client.WriteSingleCoil(address, value)
}

func (w *modbusWrapper) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	w.Lock()
	defer w.Unlock()
	return w.client.ReadInputRegisters(address, quantity)
}

func (w *modbusWrapper) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	w.Lock()
	defer w.Unlock()
	return w.client.ReadHoldingRegisters(address, quantity)
}
