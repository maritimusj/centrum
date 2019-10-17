package devices

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/kr/pretty"
	"github.com/maritimusj/centrum/edge/devices/InverseServer"
	"github.com/maritimusj/centrum/edge/devices/ep6v2"
	"github.com/maritimusj/centrum/edge/devices/event"
	"github.com/maritimusj/centrum/edge/lang"
	httpLoggerStore "github.com/maritimusj/centrum/edge/logStore/http"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/json_rpc"
	log "github.com/sirupsen/logrus"
)

type Adapter struct {
	client        *ep6v2.Device
	conf          *json_rpc.Conf
	measureDataCH chan *MeasureData
	logger        *log.Logger
	done          chan struct{}

	wg sync.WaitGroup
}

func (adapter *Adapter) IsDone() bool {
	select {
	case <-adapter.done:
		return true
	default:
		return false
	}
}

func (adapter *Adapter) Close() {
	if adapter != nil {
		if adapter.client != nil {
			adapter.client.Close()
			adapter.client = nil
		}

		close(adapter.done)
		adapter.wg.Wait()
	}
}

func (adapter *Adapter) OnDeviceStatusChanged(index lang.StrIndex) {
	event.Publish(event.DeviceStatusChanged, adapter.conf, index)
}

func (adapter *Adapter) OnMeasureDiscovered(tagName, title string) {
	key := "tag:" + adapter.conf.UID + ":" + tagName
	if v, ok := global.Params.Get(key); !ok || v.(string) != title {
		global.Params.Set(key, title)
		event.Publish(event.MeasureDiscovered, adapter.conf, tagName, title)
	}
}

type Runner struct {
	ctx      context.Context
	adapters sync.Map
}

func New() *Runner {
	runner := &Runner{
		ctx: context.Background(),
	}
	return runner
}

func (runner *Runner) StartInverseServer(conf *json_rpc.InverseConf) error {
	return InverseServer.Start(runner.ctx, conf.Address, conf.Port)
}

func (runner *Runner) GetBaseInfo(uid string) (map[string]interface{}, error) {
	if v, ok := runner.adapters.Load(uid); ok {
		baseInfo := make(map[string]interface{})

		adapter := v.(*Adapter)
		model, err := adapter.client.GetModel()
		if err != nil {
			return nil, err
		}

		baseInfo["model"] = model.ID
		baseInfo["version"] = model.Version
		baseInfo["title"] = model.Title

		addr, err := adapter.client.GetAddr()
		if err != nil {
			return baseInfo, err
		}

		baseInfo["addr"] = addr.Ip.String() + "/" + addr.Mask.String()
		baseInfo["mac"] = addr.Mac.String()

		baseInfo["status"] = map[string]interface{}{
			"index": adapter.client.GetStatus(),
			"title": adapter.client.GetStatusTitle(),
		}
		return baseInfo, nil
	}

	return nil, lang.Error(lang.ErrDeviceNotExists)
}

func needRestartAdapter(conf *json_rpc.Conf, newConf *json_rpc.Conf) bool {
	return conf.Address != newConf.Address ||
		conf.InfluxDBAddress != newConf.InfluxDBAddress ||
		conf.InfluxDBUserName != newConf.InfluxDBUserName ||
		conf.InfluxDBPassword != newConf.InfluxDBPassword ||
		conf.DB != newConf.DB ||
		conf.CallbackURL != newConf.CallbackURL
}

func (runner *Runner) Active(conf *json_rpc.Conf) error {
	if v, ok := runner.adapters.Load(conf.UID); ok {
		adapter := v.(*Adapter)
		if needRestartAdapter(adapter.conf, conf) {
			adapter.Close()
			runner.adapters.Delete(conf.UID)
		} else {
			adapter.conf.Interval = conf.Interval
			if adapter.conf.LogLevel != conf.LogLevel {
				adapter.conf.LogLevel = conf.LogLevel

				level, err := log.ParseLevel(conf.LogLevel)
				if err != nil {
					return err
				}
				adapter.logger.SetLevel(level)
			}

			adapter.OnDeviceStatusChanged(adapter.client.GetStatus())
			return nil
		}
	}

	logger := log.New()
	if conf.CallbackURL != "" {
		if conf.LogLevel == "" {
			conf.LogLevel = "error"
		}

		level, err := log.ParseLevel(conf.LogLevel)
		if err != nil {
			return err
		}

		loggerHook := httpLoggerStore.New()
		loggerHook.SetUID(conf.UID)

		logger.SetLevel(level)
		logger.AddHook(loggerHook)

		err = loggerHook.Open(runner.ctx, conf.CallbackURL)
		if err != nil {
			return err
		}
	}

	adapter := &Adapter{
		client:        ep6v2.New(),
		conf:          conf,
		logger:        logger,
		measureDataCH: make(chan *MeasureData, 100),
		done:          make(chan struct{}),
	}

	if _, ok := runner.adapters.LoadOrStore(conf.UID, adapter); !ok {
		return runner.Serve(adapter)
	}

	return nil
}

func (runner *Runner) GetValue(ch *json_rpc.CH) (interface{}, error) {
	if v, ok := runner.adapters.Load(ch.UID); ok {
		adapter := v.(*Adapter)
		v, err := adapter.client.GetCHValue(ch.Tag)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil,lang.Error(lang.ErrDeviceNotExists)
}

func (runner *Runner) SetValue(val *json_rpc.Value) error {
	if v, ok := runner.adapters.Load(val.UID); ok {
		adapter := v.(*Adapter)
		return adapter.client.SetCHValue(val.Tag, val.V)
	}
	return lang.Error(lang.ErrDeviceNotExists)
}

func (runner *Runner) GetRealtimeData(uid string) ([]map[string]interface{}, error) {
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)
		r, err := adapter.client.GetRealTimeData()
		if err != nil {
			return nil, err
		}
		adapter.logger.Info("GetRealtimeData: ", uid)
		values := make([]map[string]interface{}, 0)
		for i := 0; i < r.AINum(); i++ {
			ai, err := adapter.client.GetAI(i)
			if err != nil {
				return values, err
			}
			if v, ok := r.GetAIValue(i, ai.GetConfig().Point); ok {
				av := ai.CheckAlarm(v)
				values = append(values, map[string]interface{}{
					"tag":   ai.GetConfig().TagName,
					"title": ai.GetConfig().Title,
					"unit":  ai.GetConfig().Uint,
					"alarm": ep6v2.AlarmDesc(av),
					"value": v,
				})

				adapter.OnMeasureDiscovered(ai.GetConfig().TagName, ai.GetConfig().Title)
			}
		}

		for i := 0; i < r.AONum(); i++ {
			ao, err := adapter.client.GetAO(i)
			if err != nil {
				return values, err
			}
			if v, ok := r.GetAOValue(i); ok {
				values = append(values, map[string]interface{}{
					"tag":   ao.GetConfig().TagName,
					"title": ao.GetConfig().Title,
					"unit":  ao.GetConfig().Uint,
					"value": v,
				})
				adapter.OnMeasureDiscovered(ao.GetConfig().TagName, ao.GetConfig().Title)
			}
		}

		for i := 0; i < r.DINum(); i++ {
			di, err := adapter.client.GetDI(i)
			if err != nil {
				return values, err
			}
			if v, ok := r.GetDIValue(i); ok {
				values = append(values, map[string]interface{}{
					"tag":   di.GetConfig().TagName,
					"title": di.GetConfig().Title,
					"value": v,
				})
				adapter.OnMeasureDiscovered(di.GetConfig().TagName, di.GetConfig().Title)
			}
		}

		for i := 0; i < r.DONum(); i++ {
			do, err := adapter.client.GetDO(i)
			if err != nil {
				return values, err
			}
			if v, ok := r.GetDOValue(i); ok {
				values = append(values, map[string]interface{}{
					"tag":   do.GetConfig().TagName,
					"title": do.GetConfig().Title,
					"value": v,
					"ctrl":  do.GetConfig().IsManual,
				})
				adapter.OnMeasureDiscovered(do.GetConfig().TagName, do.GetConfig().Title)
			}
		}

		return values, nil
	}

	return nil, lang.Error(lang.ErrDeviceNotExists)
}

func (runner *Runner) Remove(uid string) {
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)

		adapter.Close()
		runner.adapters.Delete(uid)
	}
}

func (runner *Runner) Serve(adapter *Adapter) error {
	adapter.OnDeviceStatusChanged(lang.AdapterInitializing)

	fmt.Printf("%# v", pretty.Formatter(adapter.conf))
	adapter.logger.Info("start influx http client")

	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     adapter.conf.InfluxDBAddress,
		Username: adapter.conf.InfluxDBUserName,
		Password: adapter.conf.InfluxDBPassword,
	})
	if err != nil {
		return err
	}

	if _, _, err = c.Ping(3 * time.Second); err != nil {
		return err
	}

	adapter.logger.Info("create influx db: ", adapter.conf.DB)

	_, err = c.Query(influx.Query{
		Database: adapter.conf.DB,
		Command:  fmt.Sprintf("CREATE DATABASE \"%s\"", adapter.conf.DB),
	})
	if err != nil {
		adapter.logger.Error(err)
		return err
	}

	getMeasureDataFN := func() error {
		bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
			Precision: "ns",
			Database:  adapter.conf.DB,
		})
		for {
			select {
			case data := <-adapter.measureDataCH:
				if data == nil {
					return errors.New("got nil data")
				}
				point, err := influx.NewPoint(data.Name, data.Tags, data.Fields, data.Time)
				if err != nil {
					log.Error(err)
					continue
				} else {
					bp.AddPoint(point)
				}
			case <-time.After(1 * time.Second):
				if len(bp.Points()) > 0 {
					err := c.Write(bp)
					if err != nil {
						return err
					} else {
						return nil
					}
				}
			}
		}
	}

	adapter.wg.Add(2)
	go func() {
		defer func() {
			adapter.logger.Warnln("influx routine exit!")
			adapter.wg.Done()
		}()

		for {
			select {
			case <-runner.ctx.Done():
				return
			case <-adapter.done:
				return
			default:
				err := getMeasureDataFN()
				if err != nil {
					adapter.logger.Error(err)
					return
				}
			}
		}
	}()

	go func() {
		defer func() {
			close(adapter.measureDataCH)
			adapter.wg.Done()
			adapter.logger.Warnln("fetch data routine exit!")
		}()

	makeConnection:
		client := adapter.client
		for {
			adapter.logger.Info("try connect to :", adapter.conf.Address)
			adapter.OnDeviceStatusChanged(lang.Connecting)

			err := client.Connect(runner.ctx, adapter.conf.Address)
			if err != nil {
				if err == runner.ctx.Err() {
					return
				}

				adapter.OnDeviceStatusChanged(lang.Disconnected)

				select {
				case <-adapter.done:
					return
				case <-time.After(adapter.conf.Interval):
					continue
				}
			} else {
				if client.IsConnected() {
					break
				}
			}
		}

		adapter.OnDeviceStatusChanged(lang.Connected)

		for {
			select {
			case <-runner.ctx.Done():
				return
			case <-adapter.done:
				return
			case <-time.After(adapter.conf.Interval):
				adapter.logger.Info("start fetch data from: ", adapter.conf.Address)
				err := runner.fetchData(adapter)
				if err != nil {
					adapter.logger.Error(err)
					go adapter.OnDeviceStatusChanged(lang.Disconnected)
					adapter.client.Close()
					goto makeConnection
				}
			}
		}
	}()

	return nil
}

func (runner *Runner) fetchData(adapter *Adapter) error {
	client := adapter.client
	data, err := client.GetRealTimeData()
	if err != nil {
		return err
	}

	getData := func(fn func() error) error {
		select {
		case <-runner.ctx.Done():
			return runner.ctx.Err()
		case <-adapter.done:
			return errors.New("adapter closed")
		default:
			if err := fn(); err != nil {
				return err
			}
		}
		return nil
	}

	for i := 0; i < data.AINum(); i++ {
		err := getData(func() error {
			ai, err := client.GetAI(i)
			if err != nil {
				return err
			}
			v, ok := data.GetAIValue(i, ai.GetConfig().Point)
			if ok {
				data := NewMeasureData(ai.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", ai.GetConfig().TagName)
				data.AddTag("title", ai.GetConfig().Title)
				data.AddTag("alarm", ep6v2.AlarmDesc(ai.CheckAlarm(v)))
				data.AddField("val", v)
				adapter.measureDataCH <- data

				adapter.OnMeasureDiscovered(ai.GetConfig().TagName, ai.GetConfig().Title)
				log.Tracef("%s => %#v", ai.GetConfig().TagName, v)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	for i := 0; i < data.DINum(); i++ {
		err := getData(func() error {
			di, err := client.GetDI(i)
			if err != nil {
				return err
			}
			v, ok := data.GetDIValue(i)
			if ok {
				data := NewMeasureData(di.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", di.GetConfig().TagName)
				data.AddTag("title", di.GetConfig().Title)
				data.AddField("val", v)
				adapter.measureDataCH <- data

				adapter.OnMeasureDiscovered(di.GetConfig().TagName, di.GetConfig().Title)
				log.Tracef("%s => %#v", di.GetConfig().TagName, v)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	for i := 0; i < data.AONum(); i++ {
		err := getData(func() error {
			ao, err := client.GetAO(i)
			if err != nil {
				return err
			}
			v, ok := data.GetAOValue(i)
			if ok {
				data := NewMeasureData(ao.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", ao.GetConfig().TagName)
				data.AddTag("title", ao.GetConfig().Title)
				data.AddField("val", v)
				adapter.measureDataCH <- data

				adapter.OnMeasureDiscovered(ao.GetConfig().TagName, ao.GetConfig().Title)
				log.Tracef("%s => %#v", ao.GetConfig().TagName, v)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	for i := 0; i < data.DONum(); i++ {
		err := getData(func() error {
			do, err := client.GetDO(i)
			if err != nil {
				return err
			}
			v, ok := data.GetDOValue(i)
			if ok {
				data := NewMeasureData(do.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", do.GetConfig().TagName)
				data.AddTag("title", do.GetConfig().Title)
				data.AddField("val", v)
				adapter.measureDataCH <- data

				adapter.OnMeasureDiscovered(do.GetConfig().TagName, do.GetConfig().Title)
				log.Tracef("%s => %#v", do.GetConfig().TagName, v)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
