package ep6v2

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	CHBlockSize         = 256
	AIValueStartAddress = 96
	AIAlarmStartAddress = 48
)

const (
	AlarmHF = 0x10
	AlarmHH = 0x0c
	AlarmHI = 0x04
	AlarmLO = 0x01
	AlarmLL = 0x03
	AlarmLF = 0x20
)

var (
	alarmMap = map[int]string{
		AlarmHF: "HF",
		AlarmHH: "HH",
		AlarmHI: "HI",
		AlarmLO: "LO",
		AlarmLL: "LL",
		AlarmLF: "LF",
	}
)

func FormatAlarm(alarm int) string {
	if alarm == 0 {
		return "Ok"
	} else if v, ok := alarmMap[alarm]; ok {
		return v
	}
	return "<unknown>"
}

const (
	None = iota
	Control
	Alarm
)

type AI struct {
	Index  int
	config *AIConfig

	conn modbusClient
}

type AlarmItem struct {
	Style int
	Value float32
}

type AIAlarmConfig struct {
	Enabled        bool
	PrimalMaxValue float32
	PrimalMinValue float32
	MaxValue       float32
	MinValue       float32
	LowCut         int
	DeadBand       float32
	HiHi           AlarmItem
	HI             AlarmItem
	LO             AlarmItem
	LoLo           AlarmItem
	HF             AlarmItem
	LF             AlarmItem

	Delay int
}

func (alarm *AIAlarmConfig) fetchData(conn modbusClient, index int) error {
	var address, quantity uint16 = uint16(index+1)*CHBlockSize + 47, 11
	data, err := conn.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}

	alarm.HiHi.Style = int(binary.BigEndian.Uint16(data[0:]))
	alarm.HI.Style = int(binary.BigEndian.Uint16(data[2:]))
	alarm.LO.Style = int(binary.BigEndian.Uint16(data[4:]))
	alarm.LoLo.Style = int(binary.BigEndian.Uint16(data[6:]))
	alarm.HF.Style = int(binary.BigEndian.Uint16(data[8:]))
	alarm.LF.Style = int(binary.BigEndian.Uint16(data[10:]))

	alarm.Delay = int(binary.BigEndian.Uint16(data[20:]))

	address, quantity = uint16(index+1)*CHBlockSize+80, 30
	data, err = conn.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}

	alarm.PrimalMaxValue = ToSingle(data[0:])
	alarm.PrimalMinValue = ToSingle(data[4:])
	alarm.MaxValue = ToSingle(data[8:])
	alarm.MinValue = ToSingle(data[12:])
	alarm.LowCut = int(math.Round(float64(ToSingle(data[24:]))))
	alarm.HiHi.Value = ToSingle(data[32:])
	alarm.HI.Value = ToSingle(data[36:])
	alarm.LO.Value = ToSingle(data[40:])
	alarm.LoLo.Value = ToSingle(data[44:])
	alarm.HF.Value = ToSingle(data[48:])
	alarm.LF.Value = ToSingle(data[52:])
	alarm.DeadBand = ToSingle(data[56:])

	return nil
}

type AIConfig struct {
	Enabled      bool //是否启用
	AlarmEnabled bool

	TagName string //频道名称
	Title   string //中文名称
	Point   int    //小位数
	Uint    string //单位名称
	Gain    float32
	Offset  float32

	Alarm *AIAlarmConfig //警报设置
}

func (ai *AI) GetValue() (float32, int, error) {
	var address, quantity uint16 = uint16(AIValueStartAddress + ai.Index*2), 2
	data, err := ai.conn.ReadInputRegisters(address, quantity)
	if err != nil {
		return 0, 0, err
	}

	v := ToFloat32(ToSingle(data), ai.config.Point)

	//alarm
	address = uint16(AIAlarmStartAddress + ai.Index)
	data, err = ai.conn.ReadInputRegisters(address, 1)
	if err != nil {
		return v, 0, err
	}

	return v, int(data[1]), nil
}

func (c *AIConfig) fetchData(conn modbusClient, index int) error {
	var address, quantity uint16 = uint16(index+1) * CHBlockSize, 34
	data, err := conn.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return err
	}

	c.AlarmEnabled = true
	c.Gain = 1
	c.Offset = 0
	//英文名称
	c.TagName = fmt.Sprintf("AI-%d", index+1)
	//中文名称
	c.Title = DecodeUtf16String(data[0:32])
	c.Uint = DecodeUtf16String(data[32:64])
	c.Enabled = binary.BigEndian.Uint16(data[64:]) > 0
	c.Point = int(binary.BigEndian.Uint16(data[66:]))

	return nil
}
