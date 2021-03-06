package ep6v2

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/maritimusj/centrum/edge/devices/modbus"
	"github.com/maritimusj/centrum/edge/devices/util"
)

const (
	AOCHStartAddress    = 28672
	AOValueStartAddress = 0
)

type AO struct {
	Index  int
	config *AOConfig
	conn   modbus.Client
}

type AOConfig struct {
	Enabled bool //是否启用

	CTLMode        int
	CTLSource      int
	CTLMin         float32
	CTLMax         float32
	CTLTarget      float32
	CTLKp          float32
	CTLKi          float32
	CTLKd          float32
	CTLManualValue float32

	TagName string //频道名称
	Title   string //中文名称
	Point   int    //小位数
	Uint    string //单位名称
}

func (ao *AO) GetValue() (float32, error) {
	return 0, errors.New("implement me")
}

func (ao *AO) GetConfig() *AOConfig {
	if ao.config == nil {
		config := &AOConfig{}
		if err := config.fetchData(ao.conn, ao.Index); err != nil {
			return config
		}
		ao.config = config
	}
	return ao.config
}

func (c *AOConfig) fetchData(conn modbus.Client, index int) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			retErr = fmt.Errorf("unexpect error: %#v", err)
			return
		}
	}()

	start := AOCHStartAddress + uint16(index)*CHBlockSize
	data, _, err := conn.ReadHoldingRegisters(start, 16)
	if err != nil {
		return err
	}

	c.Title = util.DecodeUtf16String(data[0:32])
	//英文名称
	c.TagName = fmt.Sprintf("AO-%d", index+1)
	c.Uint = "mA"
	c.Point = 3

	data, _, err = conn.ReadHoldingRegisters(start+32, 1)
	if err != nil {
		return err
	}

	c.Enabled = data[1] > 0

	data, _, err = conn.ReadHoldingRegisters(start+42, 9)
	if err != nil {
		return err
	}

	c.CTLMode = int(binary.BigEndian.Uint16(data[0:]))

	c.CTLSource = int(binary.BigEndian.Uint16(data[2:]))
	c.CTLMin = util.ToFloat32(float32(binary.BigEndian.Uint16(data[4:]))/100, 2)
	c.CTLMax = util.ToFloat32(float32(binary.BigEndian.Uint16(data[6:]))/100, 2)
	c.CTLTarget = util.ToFloat32(float32(binary.BigEndian.Uint16(data[8:]))/1000, 2)
	c.CTLKp = util.ToFloat32(float32(binary.BigEndian.Uint16(data[10:]))/1000, 2)
	c.CTLKi = util.ToFloat32(float32(binary.BigEndian.Uint16(data[12:]))/1000, 2)
	c.CTLKd = util.ToFloat32(float32(binary.BigEndian.Uint16(data[14:]))/1000, 2)
	c.CTLManualValue = util.ToFloat32(float32(binary.BigEndian.Uint16(data[16:]))/1000, 2)

	return nil
}
