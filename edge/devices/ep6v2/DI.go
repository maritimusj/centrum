package ep6v2

import (
	"encoding/binary"
	"fmt"
)

const (
	DICHStartAddress = 12288
)

type DI struct {
	Index  int
	config *DIConfig

	conn modbusClient
}

type DIConfig struct {
	Enabled      bool //是否启用
	Inverse      bool
	AlarmEnabled bool

	TagName string //频道名称
	Title   string //中文名称

	AlarmDelay int //警报延迟(秒)
}

func (di *DI) GetValue() (bool, error) {
	data, err := di.conn.ReadDiscreteInputs(uint16(di.Index), 1)
	if err != nil {
		return false, err
	}
	return data[0] > 0, nil
}

func (di *DI) GetConfig() *DIConfig {
	return di.config
}

func (c *DIConfig) fetchData(conn modbusClient, index int) error {
	var address, quantity uint16 = DICHStartAddress + uint16(index)*CHBlockSize, 16
	data, err := conn.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}

	c.Title = DecodeUtf16String(data[0:32])
	c.TagName = fmt.Sprintf("DI-%d", index+1)

	address, quantity = DICHStartAddress+uint16(index)*CHBlockSize+32, 5
	data, err = conn.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}

	c.Enabled = data[1] > 0
	c.AlarmDelay = int(binary.BigEndian.Uint16(data[2:]))
	c.Inverse = data[7] > 0
	c.AlarmEnabled = data[9] > 0

	return nil
}
