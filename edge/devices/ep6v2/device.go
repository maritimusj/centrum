package ep6v2

import (
	"context"
	"errors"
	"github.com/maritimusj/modbus"
	"io"
	"net"
	"strconv"
	"strings"
)

type Connector interface {
	Try(ctx context.Context, addr string) (net.Conn, error)
}

const (
	Disconnected = iota
	Connecting
	Connected
)

type Device struct {
	model *Model
	addr  *Addr

	status int

	connector Connector

	handler io.Closer
	client  modbus.Client

	chAI map[int]*AI
	chAO map[int]*AO
	chDI map[int]*DI
	chDO map[int]*DO
}

func New() *Device {
	return &Device{}
}

func (device *Device) SetConnector(connector Connector) {
	device.connector = connector
}

func (device *Device) onDisconnected(err error) {
	device.client = nil
	device.handler = nil
	device.status = Disconnected
}

func (device *Device) Connect(ctx context.Context, address string) error {
	device.status = Connecting
	if device.connector == nil {
		device.connector = NewTCPConnector()
	}
	conn, err := device.connector.Try(ctx, address)
	if err != nil {
		device.status = Disconnected
		return err
	}

	handler := modbus.NewTCPClientHandlerFrom(conn)
	client := modbus.NewClient(handler)

	device.chAI = make(map[int]*AI)
	device.chAO = make(map[int]*AO)
	device.chDI = make(map[int]*DI)
	device.chDO = make(map[int]*DO)

	device.handler = handler
	device.client = client
	device.status = Connected

	return nil
}

func (device *Device) Close() {
	device.status = Disconnected

	device.model = nil
	device.addr = nil

	device.chAI = nil
	device.chAO = nil
	device.chDI = nil
	device.chDO = nil

	device.client = nil
	if device.handler != nil {
		_ = device.handler.Close()
		device.handler = nil
	}
}

func (device *Device) GetStatus() int {
	return device.status
}

func (device *Device) GetModel() (*Model, error) {
	if device.model == nil {
		device.model = &Model{}
		if err := device.model.fetchData(device.client); err != nil {
			return device.model, err
		}
	}

	return device.model, nil
}

func (device *Device) GetAddr() (*Addr, error) {
	if device.addr != nil {
		device.addr = &Addr{}
		if err := device.addr.fetchData(device.client); err != nil {
			return device.addr, err
		}
	}

	return device.addr, nil
}

func (device *Device) GetCHNum() (*CHNum, error) {
	chNum := &CHNum{}
	if err := chNum.fetchData(device.client); err != nil {
		return nil, err
	}

	return chNum, nil
}

func (device *Device) GetRealTimeData() (*RealTimeData, error) {
	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}

	data := &RealTimeData{
		chNum: chNum,
	}

	err = data.fetchData(device.client)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (device *Device) SetCHValue(tag string, value interface{}) error {
	if strings.HasPrefix(tag, "DO") {
		do, err := device.GetDOFromTag(tag)
		if err != nil {
			return err
		}

		_, err = do.SetValue(IsOn(value))
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("invalid ch index")
}

func (device *Device) GetCHValue(tag string) (value map[string]interface{}, err error) {
	seg := strings.SplitN(tag, "-", 2)
	if len(seg) != 2 {
		err = errors.New("invalid ch")
		return
	}
	var index int64
	index, err = strconv.ParseInt(seg[1], 10, 0)
	if err != nil {
		return
	}

	//index is base on 0
	index -= 1

	switch strings.ToUpper(seg[0]) {
	case "AI":
		var ai *AI
		ai, err = device.GetAI(int(index))
		if err != nil {
			return
		}
		var v float32
		v, err = ai.GetValue()
		if err != nil {
			return
		}
		var av AlarmValue
		av, err = ai.GetAlarmState()
		if err != nil {
			return
		}
		return map[string]interface{}{
			"title": ai.GetConfig().Title,
			"tag":   ai.GetConfig().TagName,
			"unit":  ai.GetConfig().Uint,
			"alarm": AlarmDesc(av),
			"value": v,
		}, nil
	case "AO":
		var ao *AO
		ao, err = device.GetAO(int(index))
		if err != nil {
			return
		}
		var v float32
		v, err = ao.GetValue()
		if err != nil {
			return
		}

		return map[string]interface{}{
			"title": ao.GetConfig().Title,
			"tag":   ao.GetConfig().TagName,
			"unit":  ao.GetConfig().Uint,
			"value": v,
		}, nil
	case "DI":
		var di *DI
		di, err = device.GetDI(int(index))
		if err != nil {
			return
		}
		var v bool
		v, err = di.GetValue()
		if err != nil {
			return
		}
		return map[string]interface{}{
			"title": di.GetConfig().Title,
			"tag":   di.GetConfig().TagName,
			"value": v,
		}, nil
	case "DO":
		var do *DO
		do, err = device.GetDO(int(index))
		if err != nil {
			return
		}
		var v bool
		v, err = do.GetValue()
		if err != nil {
			return
		}
		return map[string]interface{}{
			"title": do.GetConfig().Title,
			"tag":   do.GetConfig().TagName,
			"value": v,
		}, nil
	default:
		err = errors.New("invalid ch")
		return
	}
}

func (device *Device) GetAI(index int) (*AI, error) {
	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}
	if index >= chNum.AI {
		return nil, errors.New("invalid AI index")
	}

	if ai, ok := device.chAI[index]; ok {
		return ai, nil
	}

	config := &AIConfig{}
	if err := config.fetchData(device.client, index); err != nil {
		return nil, err
	}

	alarm := &AIAlarmConfig{}
	if err := alarm.fetchData(device.client, index); err != nil {
		return nil, err
	}

	ai := &AI{
		Index:       index,
		config:      config,
		alarmConfig: alarm,
		conn:        device.client,
	}

	device.chAI[index] = ai

	return ai, nil
}

func (device *Device) GetAO(index int) (*AO, error) {
	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}
	if index >= chNum.AO {
		return nil, errors.New("invalid AO index")
	}

	if ao, ok := device.chAO[index]; ok {
		return ao, nil
	}

	config := &AOConfig{}
	if err := config.fetchData(device.client, index); err != nil {
		return nil, err
	}
	ao := &AO{
		Index:  index,
		config: config,
		conn:   device.client,
	}
	device.chAO[index] = ao
	return ao, nil
}

func (device *Device) GetDI(index int) (*DI, error) {
	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}
	if index >= chNum.DI {
		return nil, errors.New("invalid DI index")
	}

	if di, ok := device.chDI[index]; ok {
		return di, nil
	}

	config := &DIConfig{}
	if err := config.fetchData(device.client, index); err != nil {
		return nil, err
	}

	di := &DI{
		Index:  index,
		config: config,
		conn:   device.client,
	}

	device.chDI[index] = di
	return di, nil
}

func (device *Device) GetDOFromTag(tag string) (*DO, error) {
	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}

	for index := 0; index < chNum.DO; index++ {
		do, err := device.GetDO(index)
		if err != nil {
			continue
		}
		if do.GetConfig().TagName == tag {
			return do, nil
		}
	}

	return nil, errors.New("invalid DO index")
}

func (device *Device) GetDO(index int) (*DO, error) {
	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}

	if index >= chNum.DO {
		return nil, errors.New("invalid DO index")
	}

	if do, ok := device.chDO[index]; ok {
		return do, nil
	}

	config := &DOConfig{}
	if err := config.fetchData(device.client, index); err != nil {
		return nil, err
	}

	do := &DO{
		config: config,
		Index:  index,
		conn:   device.client,
	}

	device.chDO[index] = do
	return do, nil
}
