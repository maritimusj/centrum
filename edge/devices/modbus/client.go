package modbus

import "time"

type Client interface {
	ReadCoils(address, quantity uint16) (results []byte, duration time.Duration, err error)
	ReadDiscreteInputs(address, quantity uint16) (results []byte, duration time.Duration, err error)
	WriteSingleCoil(address, value uint16) (results []byte, duration time.Duration, err error)
	ReadInputRegisters(address, quantity uint16) (results []byte, duration time.Duration, err error)
	ReadHoldingRegisters(address, quantity uint16) (results []byte, duration time.Duration, err error)
}
