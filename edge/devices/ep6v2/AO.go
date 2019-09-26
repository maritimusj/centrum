package ep6v2

import (
	"encoding/binary"
	"fmt"
)

const (
	AOCHStartAddress    = 28672
	AOValueStartAddress = 0
)

type AO struct {
	Index  int
	config *AOConfig
	conn   modbusClient
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
	panic("implement me")
}

func (c *AOConfig) fetchData(conn modbusClient, index int) error {
	start := AOCHStartAddress + uint16(index)*CHBlockSize
	data, err := conn.ReadHoldingRegisters(start, 16)
	if err != nil {
		return err
	}

	c.Title = DecodeUtf16String(data[0:32])
	//英文名称
	c.TagName = fmt.Sprintf("AO-%d", index+1)
	c.Uint = "mA"
	c.Point = 3

	data, err = conn.ReadHoldingRegisters(start+32, 1)
	if err != nil {
		return err
	}

	c.Enabled = data[1] > 0

	data, err = conn.ReadHoldingRegisters(start+42, 9)
	if err != nil {
		return err
	}

	c.CTLMode = int(binary.BigEndian.Uint16(data[0:]))

	c.CTLSource = int(binary.BigEndian.Uint16(data[2:]))
	c.CTLMin = ToFloat32(float32(binary.BigEndian.Uint16(data[4:]))/100, 2)
	c.CTLMax = ToFloat32(float32(binary.BigEndian.Uint16(data[6:]))/100, 2)
	c.CTLTarget = ToFloat32(float32(binary.BigEndian.Uint16(data[8:]))/1000, 2)
	c.CTLKp = ToFloat32(float32(binary.BigEndian.Uint16(data[10:]))/1000, 2)
	c.CTLKi = ToFloat32(float32(binary.BigEndian.Uint16(data[12:]))/1000, 2)
	c.CTLKd = ToFloat32(float32(binary.BigEndian.Uint16(data[14:]))/1000, 2)
	c.CTLManualValue = ToFloat32(float32(binary.BigEndian.Uint16(data[16:]))/1000, 2)

	return nil
}
