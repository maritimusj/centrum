package ep6v2

import "encoding/binary"

type CHNum struct {
	AI int
	DI int
	DO int
	AO int
	VO int
}

func (ch *CHNum) Sum() int {
	return ch.AI + ch.DI + ch.DO + ch.AO + ch.VO
}

func (ch *CHNum) Clone() *CHNum {
	return &CHNum{
		AI: ch.AI,
		DI: ch.DI,
		DO: ch.DO,
		AO: ch.AO,
		VO: ch.VO,
	}
}

func (ch *CHNum) fetchData(conn modbusClient) error {
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
