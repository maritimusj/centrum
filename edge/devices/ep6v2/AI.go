package ep6v2

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/maritimusj/centrum/edge/devices/modbus"
	"github.com/maritimusj/centrum/edge/devices/util"
)

const (
	CHBlockSize         = 256
	AIValueStartAddress = 96
	AIAlarmStartAddress = 48
)

type AlarmValue int

const (
	AlarmInvalid AlarmValue = 0xFF
	AlarmError   AlarmValue = 0xFE
	AlarmNormal  AlarmValue = 0x0
	AlarmHF      AlarmValue = 0x10
	AlarmHH      AlarmValue = 0x0c
	AlarmHI      AlarmValue = 0x04
	AlarmLO      AlarmValue = 0x01
	AlarmLL      AlarmValue = 0x03
	AlarmLF      AlarmValue = 0x20
)

var (
	alarmMap = map[AlarmValue]string{
		AlarmInvalid: "READ..",
		AlarmError:   "Err..",
		AlarmNormal:  "",
		AlarmHF:      "HF",
		AlarmHH:      "HH",
		AlarmHI:      "HI",
		AlarmLO:      "LO",
		AlarmLL:      "LL",
		AlarmLF:      "LF",
	}
)

func AlarmDesc(alarm AlarmValue) string {
	if alarm == 0 {
		return ""
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
	Index       int
	config      *AIConfig
	alarmConfig *AIAlarmConfig

	value             float32
	lastValueReadTime time.Time

	alarmState             AlarmValue
	lastAlarmStateReadTime time.Time

	conn modbus.Client
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

func (alarm *AIAlarmConfig) fetchData(conn modbus.Client, index int) error {
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

	alarm.PrimalMaxValue = util.ToSingle(data[0:])
	alarm.PrimalMinValue = util.ToSingle(data[4:])
	alarm.MaxValue = util.ToSingle(data[8:])
	alarm.MinValue = util.ToSingle(data[12:])
	alarm.LowCut = int(math.Round(float64(util.ToSingle(data[24:]))))
	alarm.HiHi.Value = util.ToSingle(data[32:])
	alarm.HI.Value = util.ToSingle(data[36:])
	alarm.LO.Value = util.ToSingle(data[40:])
	alarm.LoLo.Value = util.ToSingle(data[44:])
	alarm.HF.Value = util.ToSingle(data[48:])
	alarm.LF.Value = util.ToSingle(data[52:])
	alarm.DeadBand = util.ToSingle(data[56:])

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

func (ai *AI) valueExpired() bool {
	return time.Now().Sub(ai.lastValueReadTime) > 1*time.Second
}

func (ai *AI) stateExpired() bool {
	return time.Now().Sub(ai.lastAlarmStateReadTime) > 1*time.Second
}

func (ai *AI) fetchValue() error {
	var address, quantity uint16 = uint16(AIValueStartAddress + ai.Index*2), 2
	data, err := ai.conn.ReadInputRegisters(address, quantity)
	if err != nil {
		return err
	}

	ai.value = util.ToFloat32(util.ToSingle(data), ai.config.Point)
	ai.lastValueReadTime = time.Now()
	return nil
}

func (ai *AI) GetValue() (float32, error) {
	if ai.valueExpired() {
		if err := ai.fetchValue(); err != nil {
			return 0, nil
		}
	}
	return ai.value, nil
}

func (ai *AI) GetConfig() *AIConfig {
	if ai.config == nil {
		config := &AIConfig{}
		if err := config.fetchData(ai.conn, ai.Index); err != nil {
			return config
		}
		ai.config = config
	}
	return ai.config
}

func (ai *AI) GetAlarmState() (AlarmValue, error) {
	if ai.stateExpired() {
		address := uint16(AIAlarmStartAddress + ai.Index)
		data, err := ai.conn.ReadInputRegisters(address, 1)
		if err != nil {
			return 0, err
		}
		ai.alarmState = AlarmValue(data[1])
		ai.lastAlarmStateReadTime = time.Now()
	}

	return ai.alarmState, nil
}

func (ai *AI) getAlarmConfig() (*AIAlarmConfig, error) {
	if ai.alarmConfig == nil {
		alarmConfig := &AIAlarmConfig{}
		err := alarmConfig.fetchData(ai.conn, ai.Index)
		if err != nil {
			return nil, err
		}
		ai.alarmConfig = alarmConfig
	}
	return ai.alarmConfig, nil
}

func (ai *AI) CheckAlarm(val float32) (AlarmValue, float32) {
	cfg, err := ai.getAlarmConfig()
	if err != nil {
		return AlarmError, 0
	}

	if val > cfg.HF.Value && cfg.HF.Style == Alarm {
		return AlarmHF, cfg.HF.Value
	}

	if val >= cfg.HiHi.Value-cfg.DeadBand && cfg.HiHi.Style == Alarm {
		return AlarmHH, cfg.HiHi.Value
	}

	if val >= cfg.HI.Value-cfg.DeadBand && cfg.HI.Style == Alarm {
		return AlarmHI, cfg.HI.Value
	}

	if val < cfg.LO.Value+cfg.DeadBand && cfg.LO.Style == Alarm {
		return AlarmLO, cfg.LO.Value
	}

	if val < cfg.LoLo.Value+cfg.DeadBand && cfg.LoLo.Style == Alarm {
		return AlarmLL, cfg.LoLo.Value
	}

	if val < cfg.LF.Value {
		return AlarmLF, cfg.LF.Value
	}
	return AlarmNormal, 0
}

func (c *AIConfig) fetchData(conn modbus.Client, index int) error {
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
	c.Title = util.DecodeUtf16String(data[0:32])
	c.Uint = util.DecodeUtf16String(data[32:64])
	c.Enabled = binary.BigEndian.Uint16(data[64:]) > 0
	c.Point = int(binary.BigEndian.Uint16(data[66:]))

	return nil
}
