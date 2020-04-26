package ep6v2

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/maritimusj/centrum/edge/devices/modbus"
	"github.com/maritimusj/centrum/edge/devices/util"
)

const (
	DICHStartAddress = 12288
)

type DI struct {
	Index  int
	config *DIConfig

	value        bool
	lastReadTime time.Time
	conn         modbus.Client
}

type DIConfig struct {
	Enabled      bool //是否启用
	Inverse      bool
	AlarmEnabled bool

	TagName string //频道名称
	Title   string //中文名称

	AlarmConfig int
	AlarmDelay  int //警报延迟(秒)
}

func (di *DI) expired() bool {
	return time.Now().Sub(di.lastReadTime) > 1*time.Second
}

func (di *DI) GetValue() (bool, error) {
	if di.expired() {
		data, _, err := di.conn.ReadDiscreteInputs(uint16(di.Index), 1)
		if err != nil {
			return false, err
		}
		di.value = data[0] > 0
		di.lastReadTime = time.Now()
	}
	return di.value, nil
}

func (di *DI) GetConfig() *DIConfig {
	if di.config == nil {
		config := &DIConfig{}
		if err := config.fetchData(di.conn, di.Index); err != nil {
			return config
		}
		di.config = config
	}
	return di.config
}

func (c *DIConfig) fetchData(conn modbus.Client, index int) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("unexpect error: %#v", err)
			return
		}
	}()

	var address, quantity uint16 = DICHStartAddress + uint16(index)*CHBlockSize, 16
	data, _, err := conn.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}

	c.Title = util.DecodeUtf16String(data[0:32])
	c.TagName = fmt.Sprintf("DI-%d", index+1)

	address, quantity = DICHStartAddress+uint16(index)*CHBlockSize+32, 7
	data, _, err = conn.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}

	c.Enabled = data[1] > 0
	c.AlarmDelay = int(binary.BigEndian.Uint16(data[2:]))
	c.Inverse = data[7] > 0
	c.AlarmEnabled = data[9] > 0
	c.AlarmConfig = int(binary.BigEndian.Uint16(data[12:]))
	return nil
}
