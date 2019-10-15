package ep6v2

import (
	"context"
	"errors"
	"github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/modbus"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
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

	sync.RWMutex
}


func New() *Device {
	return &Device{}
}

func (device *Device) SetConnector(connector Connector) {
	device.connector = connector
}

func (device *Device) onDisconnected(err error) {
	device.Lock()
	defer device.Unlock()

	device.model = nil
	device.addr = nil
	device.client = nil
	device.handler = nil
	device.status = Disconnected
}

func (device *Device) IsConnected() bool {
	device.RLock()
	defer device.RUnlock()

	return device != nil && device.client != nil && device.status == Connected
}

func (device *Device) getModbusClient() (modbusClient, error) {
	if device.client != nil && device.status == Connected {
		return device.client, nil
	}
	return nil, errors.New("device not connected")
}

func (device *Device) Connect(ctx context.Context, address string) error {
	device.Lock()
	{
		device.status = Connecting
		if device.connector == nil {
			device.connector = NewTCPConnector()
		}
	}
	device.Unlock()

	conn, err := device.connector.Try(ctx, address)
	if err != nil {
		device.Lock()
		{
			device.status = Disconnected
		}
		device.Unlock()
		return err
	}

	device.Lock()
	{
		handler := modbus.NewTCPClientHandlerFrom(conn)
		client := modbus.NewClient(handler)

		device.chAI = make(map[int]*AI)
		device.chAO = make(map[int]*AO)
		device.chDI = make(map[int]*DI)
		device.chDO = make(map[int]*DO)

		device.handler = handler
		device.client = client
		device.status = Connected
	}
	device.Unlock()

	return nil
}

func (device *Device) Close() {
	device.Lock()
	defer device.Unlock()

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

func (device *Device) GetStatusTitle() string {
	return lang.Str(lang.StrIndex(device.status))
}

func (device *Device) GetModel() (*Model, error) {
	device.Lock()

	if device.model == nil {
		if client, err := device.getModbusClient(); err != nil {
			device.Unlock()
			return nil, err
		} else {
			device.model = &Model{}
			device.Unlock()
			if err := device.model.fetchData(client); err != nil {
				return device.model, err
			}
		}
	} else {
		device.Unlock()
	}

	return device.model, nil
}

func (device *Device) GetAddr() (*Addr, error) {
	device.Lock()

	if device.addr == nil {
		if client, err := device.getModbusClient(); err != nil {
			device.Unlock()
			return device.addr, err
		} else {
			device.addr = &Addr{}
			device.Unlock()
			if err := device.addr.fetchData(client); err != nil {
				return device.addr, err
			}
		}
	} else {
		device.Unlock()
	}

	return device.addr, nil
}

func (device *Device) GetCHNum() (*CHNum, error) {
	chNum := &CHNum{}
	if client, err := device.getModbusClient(); err != nil {
		return chNum, err
	} else {
		if err := chNum.fetchData(client); err != nil {
			return chNum, err
		}
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

	if client, err := device.getModbusClient(); err != nil {
		return nil, err
	} else {
		err = data.fetchData(client)
		if err != nil {
			return nil, err
		}
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

	client, err := device.getModbusClient()
	if err != nil {
		return nil, err
	}

	config := &AIConfig{}
	if err := config.fetchData(client, index); err != nil {
		return nil, err
	}

	alarm := &AIAlarmConfig{}
	if err := alarm.fetchData(client, index); err != nil {
		return nil, err
	}

	ai := &AI{
		Index:       index,
		config:      config,
		alarmConfig: alarm,
		conn:        client,
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

	client, err := device.getModbusClient()
	if err != nil {
		return nil, err
	}

	config := &AOConfig{}
	if err := config.fetchData(client, index); err != nil {
		return nil, err
	}
	ao := &AO{
		Index:  index,
		config: config,
		conn:   client,
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

	client, err := device.getModbusClient()
	if err != nil {
		return nil, err
	}

	config := &DIConfig{}
	if err := config.fetchData(client, index); err != nil {
		return nil, err
	}

	di := &DI{
		Index:  index,
		config: config,
		conn:   client,
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

	client, err := device.getModbusClient()
	if err != nil {
		return nil, err
	}

	config := &DOConfig{}
	if err := config.fetchData(client, index); err != nil {
		return nil, err
	}

	do := &DO{
		config: config,
		Index:  index,
		conn:   client,
	}

	device.chDO[index] = do
	return do, nil
}
