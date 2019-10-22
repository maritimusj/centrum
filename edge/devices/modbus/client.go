package modbus

type Client interface {
	ReadCoils(address, quantity uint16) (results []byte, err error)
	ReadDiscreteInputs(address, quantity uint16) (results []byte, err error)
	WriteSingleCoil(address, value uint16) (results []byte, err error)
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
}
