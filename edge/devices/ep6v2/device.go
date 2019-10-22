package ep6v2

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/maritimusj/centrum/edge/devices/CHNum"
	"github.com/maritimusj/centrum/edge/devices/InverseServer"
	"github.com/maritimusj/centrum/edge/devices/realtime"
	"github.com/maritimusj/centrum/edge/devices/util"
	"github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/centrum/global"
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

type Device struct {
	model *Model
	addr  *Addr

	status lang.StrIndex

	connector Connector

	handler io.Closer
	client  modbus.Client

	chAI map[int]*AI
	chAO map[int]*AO
	chDI map[int]*DI
	chDO map[int]*DO

	chNum        *CHNum.Data
	readTimeData *realtime.Data

	sync.RWMutex
}

func New() *Device {
	return &Device{}
}

func (device *Device) SetConnector(connector Connector) {
	device.connector = connector
}

func (device *Device) onDisconnected(err error) {
	device.Reset(func() {
		device.status = lang.Disconnected
	})
}

func (device *Device) IsConnected() bool {
	if device != nil {
		device.RLock()
		defer device.RUnlock()

		return device.client != nil && device.status == lang.Connected
	}
	return false
}

func (device *Device) getModbusClient() (modbus.Client, error) {
	if device.client != nil && device.status == lang.Connected {
		return device.client, nil
	}
	return nil, errors.New("device not connected")
}

func (device *Device) Connect(ctx context.Context, address string) error {
	device.Lock()
	device.status = lang.Connecting

	if device.connector == nil {
		if govalidator.IsMAC(address) {
			println("InverseServer.DefaultConnector")
			device.connector = InverseServer.DefaultConnector()
		} else {
			println("NewTCPConnector")
			device.connector = NewTCPConnector()
		}
	}

	device.Unlock()

	conn, err := device.connector.Try(ctx, address)
	if err != nil {
		device.status = lang.Disconnected
		return err
	}

	device.Reset(func() {
		device.status = lang.Connected
		handler := modbus.NewTCPClientHandlerFrom(conn)
		client := modbus.NewClient(handler)
		device.handler = handler
		device.client = client
	})

	return nil
}

func (device *Device) Close() {
	device.Reset(func() {
		device.status = lang.Disconnected

		device.client = nil
		if device.handler != nil {
			_ = device.handler.Close()
			device.handler = nil
		}
	})
}

func (device *Device) Reset(otherFN ...func()) {
	device.Lock()
	defer device.Unlock()

	device.model = nil
	device.addr = nil

	if device.chNum != nil {
		device.chNum.Release()
		device.chNum = nil
	}
	if device.readTimeData != nil {
		device.readTimeData.Release()
		device.readTimeData = nil
	}

	device.chAI = make(map[int]*AI)
	device.chAO = make(map[int]*AO)
	device.chDI = make(map[int]*DI)
	device.chDO = make(map[int]*DO)

	global.Params.Reset()

	for _, fn := range otherFN {
		if fn != nil {
			fn()
		}
	}
}

func (device *Device) GetStatus() lang.StrIndex {
	device.RLock()
	defer device.RUnlock()

	return device.status
}

func (device *Device) GetStatusTitle() string {
	device.RLock()
	defer device.RUnlock()

	return lang.Str(device.status)
}

func (device *Device) GetModel() (*Model, error) {
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	device.Lock()
	defer device.Unlock()

	if device.model == nil {
		if client, err := device.getModbusClient(); err != nil {
			return nil, err
		} else {
			model := &Model{}
			if err := model.fetchData(client); err != nil {
				return nil, err
			}
			device.model = model
		}
	}

	return device.model, nil
}

func (device *Device) GetAddr() (*Addr, error) {
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	device.Lock()
	defer device.Unlock()

	if device.addr == nil {
		if client, err := device.getModbusClient(); err != nil {
			return nil, err
		} else {
			addr := &Addr{}
			if err := addr.fetchData(client); err != nil {
				return nil, err
			}
			device.addr = addr
		}
	}
	return device.addr, nil
}

func (device *Device) GetCHNum() (*CHNum.Data, error) {
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	device.Lock()
	defer device.Unlock()

	if device.chNum == nil {
		chNum := CHNum.New()
		if client, err := device.getModbusClient(); err != nil {
			chNum.Release()
			return nil, err
		} else {
			if err := chNum.FetchData(client); err != nil {
				chNum.Release()
				return nil, err
			}
		}
		device.chNum = chNum
	}

	return device.chNum, nil
}

func (device *Device) GetRealTimeData() (*realtime.Data, error) {
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}

	device.Lock()
	defer device.Unlock()

	if device.readTimeData == nil {
		device.readTimeData = realtime.New(chNum)
	}

	if client, err := device.getModbusClient(); err != nil {
		return nil, err
	} else {
		err = device.readTimeData.FetchData(client)
		if err != nil {
			return nil, err
		}
	}

	return device.readTimeData.Clone(), nil
}

func (device *Device) SetCHValue(tag string, value interface{}) error {
	if !device.IsConnected() {
		return lang.Error(lang.ErrDeviceNotConnected)
	}

	if strings.HasPrefix(tag, "DO") {
		do, err := device.GetDOFromTag(tag)
		if err != nil {
			return err
		}

		_, err = do.SetValue(util.IsOn(value))
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("invalid ch index")
}

func (device *Device) GetCHValue(tag string) (value map[string]interface{}, err error) {
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

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
		device.Lock()
		if device.readTimeData != nil {
			device.readTimeData.SetDOValue(do.Index, v)
		}
		device.Unlock()
		return map[string]interface{}{
			"title": do.GetConfig().Title,
			"tag":   do.GetConfig().TagName,
			"value": v,
			"ctrl":  do.GetConfig().IsManual,
		}, nil
	default:
		err = errors.New("invalid ch")
		return
	}
}

func (device *Device) GetAI(index int) (*AI, error) {
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}

	if index >= chNum.AI {
		return nil, errors.New("invalid AI index")
	}

	device.Lock()
	defer device.Unlock()

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
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}

	if index >= chNum.AO {
		return nil, errors.New("invalid AO index")
	}

	device.Lock()
	defer device.Unlock()

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
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}
	if index >= chNum.DI {
		return nil, errors.New("invalid DI index")
	}

	device.Lock()
	defer device.Unlock()

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
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

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
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum()
	if err != nil {
		return nil, err
	}

	if index >= chNum.DO {
		return nil, errors.New("invalid DO index")
	}

	device.Lock()
	defer device.Unlock()

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
