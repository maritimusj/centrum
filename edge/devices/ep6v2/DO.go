package ep6v2

import (
	"fmt"
	"time"

	"github.com/maritimusj/centrum/edge/devices/modbus"
	"github.com/maritimusj/centrum/edge/devices/util"
)

const (
	DOCHStartAddress = 20480
	ON               = 0xFF00
	OFF              = 0x0000
)

type DO struct {
	Index  int
	config *DOConfig

	value        bool
	lastReadTime time.Time

	conn modbus.Client
}

type DOConfig struct {
	Enabled      bool //是否启用
	AlarmEnabled bool

	TagName string //频道名称
	Title   string //中文名称

	AutoControl bool //自动控制是否开启
	LogEnabled  bool //记录启用时间
	Reverse     bool //反向输出
	IsManual    bool //手动控制是否开启

	EnableSwitch bool
	OnTime       int
	OffTime      int
}

func (do *DO) expired() bool {
	return time.Now().Sub(do.lastReadTime) > 1*time.Second
}

func (do *DO) GetValue() (bool, error) {
	if do.expired() {
		data, _, err := do.conn.ReadCoils(uint16(do.Index), 1)
		if err != nil {
			return false, err
		}
		do.value = data[0] > 0
		do.lastReadTime = time.Now()
	}

	return do.value, nil
}

func (do *DO) SetValue(v bool) (bool, error) {
	var data uint16
	if v {
		data = ON
	} else {
		data = OFF
	}
	res, _, err := do.conn.WriteSingleCoil(uint16(do.Index), data)
	if err != nil {
		return false, err
	}
	return res[0] > 0, nil
}

func (do *DO) GetConfig() *DOConfig {
	if do.config == nil {
		config := &DOConfig{}
		if err := config.fetchData(do.conn, do.Index); err != nil {
			return config
		}
		do.config = config
	}
	return do.config
}

func (c *DOConfig) fetchData(conn modbus.Client, index int) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("unexpect error: %#v", err)
			return
		}
	}()

	start := DOCHStartAddress + uint16(index)*CHBlockSize
	data, _, err := conn.ReadHoldingRegisters(start, 15)
	if err != nil {
		return err
	}

	c.TagName = fmt.Sprintf("DO-%d", index+1)
	c.Title = util.DecodeUtf16String(data[0:])

	data, _, err = conn.ReadHoldingRegisters(start+32, 5)
	if err != nil {
		return err
	}

	c.AlarmEnabled = true
	c.Enabled = data[1] > 0
	c.AutoControl = data[3] > 0
	c.LogEnabled = data[5] > 0
	c.Reverse = data[7] > 0
	c.IsManual = data[9] > 0

	return nil
}
