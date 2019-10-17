package ep6v2

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

type Model struct {
	ID      string
	Version string
	Title   string
}

func (model *Model) fetchData(conn modbusClient) error {
	data, err := conn.ReadHoldingRegisters(0, 4)
	if err != nil {
		return err
	}

	model.ID = string([]byte{data[1], data[0], data[3], data[2], data[5], data[4]})
	model.Version = fmt.Sprintf("v%.2f", (float32(data[6])*100+float32(data[7]))/100)

	//读取设备名称
	data, err = conn.ReadHoldingRegisters(0x0040, 32)
	if err != nil {
		return err
	}

	var buf [32]uint16
	var index int
	for index = range buf {
		buf[index] = binary.BigEndian.Uint16(data[index*2:])
		if buf[index] == 0 {
			break
		}
	}

	model.Title = string(utf16.Decode(buf[0:index]))
	return nil
}
