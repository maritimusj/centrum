package ep6v2

import (
	"context"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/maritimusj/centrum/edge/devices/CHNum"
	"github.com/maritimusj/centrum/edge/devices/InverseServer"
	"github.com/maritimusj/centrum/edge/devices/modbus"
	"github.com/maritimusj/centrum/edge/devices/realtime"
	"github.com/maritimusj/centrum/edge/devices/util"
	"github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/synchronized"
	rawModbus "github.com/maritimusj/modbus"
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
}

func New() *Device {
	return &Device{
		status: lang.Disconnected,
	}
}

func (device *Device) SetConnector(connector Connector) {
	if device != nil {
		device.connector = connector
	}
}

func (device *Device) onDisconnected(err error) {
	device.Reset(func() {
		device.status = lang.Disconnected
	})
}

func (device *Device) IsConnected() bool {
	if device != nil {
		return device.client != nil && device.status == lang.Connected
	}
	return false
}

func (device *Device) getModbusClient() (modbus.Client, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}

	if device.client != nil && device.status == lang.Connected {
		return device.client, nil
	}
	return nil, lang.Error(lang.ErrDeviceNotConnected)
}

func (device *Device) Connect(ctx context.Context, address string) error {
	if device == nil {
		return lang.Error(lang.ErrDeviceNotExists)
	}
	<-synchronized.Do(device, func() interface{} {
		device.status = lang.Connecting
		if device.connector == nil {
			if govalidator.IsMAC(address) {
				device.connector = InverseServer.DefaultConnector()
			} else {
				device.connector = NewTCPConnector()
			}
		}
		return nil
	})

	conn, err := device.connector.Try(ctx, address)
	if err != nil {
		device.status = lang.Disconnected
		return err
	}

	device.Reset(func() {
		device.status = lang.Connected
		handler := rawModbus.NewTCPClientHandlerFrom(conn)
		client := rawModbus.NewClient(handler)
		device.handler = handler
		device.client = &modbusWrapper{client: client}
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
	if device != nil {
		<-synchronized.Do(device, func() interface{} {
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
			return nil
		})
	}
}

func (device *Device) GetStatus() lang.StrIndex {
	if device != nil {
		return device.status
	}
	return lang.Disconnected
}

func (device *Device) GetStatusTitle() string {
	if device != nil {
		return lang.Str(device.status)
	}
	return lang.Str(lang.Disconnected)
}

func (device *Device) GetModel() (*Model, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	result := <-synchronized.Do(device, func() interface{} {
		if device.model == nil {
			if client, err := device.getModbusClient(); err != nil {
				return err
			} else {
				model := &Model{}
				if err := model.fetchData(client); err != nil {
					return err
				}
				device.model = model
			}
		}

		return device.model
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*Model), nil
}

func (device *Device) GetAddr() (*Addr, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	result := <-synchronized.Do(device, func() interface{} {
		if device.addr == nil {
			if client, err := device.getModbusClient(); err != nil {
				return err
			} else {
				addr := &Addr{}
				if err := addr.fetchData(client); err != nil {
					return err
				}
				device.addr = addr
			}
		}
		return device.addr
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*Addr), nil
}

func (device *Device) GetCHNum(flush bool) (*CHNum.Data, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	result := <-synchronized.Do(device, func() interface{} {
		if device.chNum == nil || flush {
			chNum := CHNum.New()
			if client, err := device.getModbusClient(); err != nil {
				chNum.Release()
				return err
			} else {
				if err := chNum.FetchData(client); err != nil {
					chNum.Release()
					return err
				}
			}
			device.chNum = chNum
		}

		return device.chNum
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*CHNum.Data), nil
}

func (device *Device) GetRealTimeData() (*realtime.Data, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}

	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum(true)
	if err != nil {
		return nil, err
	}

	result := <-synchronized.Do(device, func() interface{} {
		if device.readTimeData == nil {
			device.readTimeData = realtime.New(chNum)
		}

		if client, err := device.getModbusClient(); err != nil {
			return err
		} else {
			err = device.readTimeData.FetchData(client)
			if err != nil {
				return err
			}
		}

		return device.readTimeData.Clone()
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*realtime.Data), nil
}

func (device *Device) SetCHValue(tag string, value interface{}) error {
	if device == nil {
		return lang.Error(lang.ErrDeviceNotExists)
	}
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
	if device == nil {
		err = lang.Error(lang.ErrDeviceNotExists)
		return
	}
	if !device.IsConnected() {
		err = lang.Error(lang.ErrDeviceNotConnected)
		return
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

		if device.readTimeData != nil {
			device.readTimeData.SetDOValue(do.Index, v)
		}

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
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum(false)
	if err != nil {
		return nil, err
	}

	if index >= chNum.AI {
		return nil, errors.New("invalid AI index")
	}

	result := <-synchronized.Do(device, func() interface{} {
		if ai, ok := device.chAI[index]; ok {
			return ai
		}

		client, err := device.getModbusClient()
		if err != nil {
			return err
		}

		config := &AIConfig{}
		if err := config.fetchData(client, index); err != nil {
			return err
		}

		alarm := &AIAlarmConfig{}
		if err := alarm.fetchData(client, index); err != nil {
			return err
		}

		ai := &AI{
			Index:       index,
			config:      config,
			alarmConfig: alarm,
			conn:        client,
		}

		device.chAI[index] = ai
		return ai
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*AI), nil
}

func (device *Device) GetAO(index int) (*AO, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum(false)
	if err != nil {
		return nil, err
	}

	if index >= chNum.AO {
		return nil, errors.New("invalid AO index")
	}

	result := <-synchronized.Do(device, func() interface{} {
		if ao, ok := device.chAO[index]; ok {
			return ao
		}

		client, err := device.getModbusClient()
		if err != nil {
			return err
		}

		config := &AOConfig{}
		if err := config.fetchData(client, index); err != nil {
			return err
		}
		ao := &AO{
			Index:  index,
			config: config,
			conn:   client,
		}
		device.chAO[index] = ao
		return ao
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*AO), nil
}

func (device *Device) GetDI(index int) (*DI, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum(false)
	if err != nil {
		return nil, err
	}
	if index >= chNum.DI {
		return nil, errors.New("invalid DI index")
	}

	result := <-synchronized.Do(device, func() interface{} {
		if di, ok := device.chDI[index]; ok {
			return di
		}

		client, err := device.getModbusClient()
		if err != nil {
			return err
		}

		config := &DIConfig{}
		if err := config.fetchData(client, index); err != nil {
			return err
		}

		di := &DI{
			Index:  index,
			config: config,
			conn:   client,
		}

		device.chDI[index] = di
		return di
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*DI), nil
}

func (device *Device) GetDOFromTag(tag string) (*DO, error) {
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum(false)
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
	if device == nil {
		return nil, lang.Error(lang.ErrDeviceNotExists)
	}
	if !device.IsConnected() {
		return nil, lang.Error(lang.ErrDeviceNotConnected)
	}

	chNum, err := device.GetCHNum(false)
	if err != nil {
		return nil, err
	}

	if index >= chNum.DO {
		return nil, errors.New("invalid DO index")
	}

	result := <-synchronized.Do(device, func() interface{} {
		if do, ok := device.chDO[index]; ok {
			return do
		}

		client, err := device.getModbusClient()
		if err != nil {
			return err
		}

		config := &DOConfig{}
		if err := config.fetchData(client, index); err != nil {
			return err
		}

		do := &DO{
			config: config,
			Index:  index,
			conn:   client,
		}

		device.chDO[index] = do
		return do
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*DO), nil
}
