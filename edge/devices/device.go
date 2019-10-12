package devices

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"github.com/maritimusj/centrum/edge/devices/ep6v2"
	"github.com/maritimusj/centrum/edge/lang"
	httpLogger "github.com/maritimusj/centrum/edge/logStore/http"
	"github.com/maritimusj/centrum/json_rpc"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
)

type Adapter struct {
	client *ep6v2.Device
	conf   *json_rpc.Conf
	ch     chan *MeasureData
	logger *log.Logger
	done   chan struct{}
	wg     sync.WaitGroup
}

func (adapter *Adapter) Close() {
	if adapter != nil && adapter.done != nil {
		close(adapter.done)
		adapter.wg.Wait()
	}
}

func (adapter *Adapter) OnDeviceStatusChange(index lang.StrIndex) {
	if adapter.conf.CallbackURL != "" {
		data, _ := json.Marshal(map[string]string{
			"uid":    adapter.conf.UID,
			"status": lang.Str(index),
		})
		req, err := http.NewRequest("post", adapter.conf.CallbackURL, bytes.NewReader(data))
		if err != nil {
			log.Errorf("[OnDeviceStatusChange] %s", err)
			return
		}
		_, err = httpLogger.DefaultHttpClient().Do(req)
		if err != nil {
			log.Errorf("[OnDeviceStatusChange] %s", err)
		}
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

func (runner *Runner) GetBaseInfo(uid string) (map[string]interface{}, error) {
	if v, ok := runner.adapters.Load(uid); ok {
		info := make(map[string]interface{})

		adapter := v.(*Adapter)
		model, err := adapter.client.GetModel()
		if err != nil {
			return nil, err
		}

		info["model"] = model.ID
		info["version"] = model.Version

		addr, err := adapter.client.GetAddr()
		if err != nil {
			return info, err
		}

		info["addr"] = addr.Ip.String() + "/" + addr.Mask.String()
		info["mac"] = addr.Mask.String()

		info["status"] = adapter.client.GetStatus()
		return info, nil
	}

	return nil, errors.New("device not exists")
}

func (runner *Runner) Active(conf *json_rpc.Conf) error {
	if v, ok := runner.adapters.Load(conf.UID); ok {
		adapter := v.(*Adapter)
		if adapter.conf.Address != conf.Address {
			adapter.Close()
			runner.adapters.Delete(conf.UID)
		} else {
			adapter.conf.Interval = conf.Interval
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

		loggerHook := httpLogger.New()
		err = loggerHook.Open(runner.ctx, conf.CallbackURL)
		if err != nil {
			return err
		}

		logger.SetLevel(level)
		logger.AddHook(loggerHook)
	}

	adapter := &Adapter{
		client: ep6v2.New(),
		conf:   conf,
		logger: logger,
		ch:     make(chan *MeasureData, 100),
		done:   make(chan struct{}),
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
	return nil, errors.New("device not exists")
}

func (runner *Runner) SetValue(val *json_rpc.Value) error {
	if v, ok := runner.adapters.Load(val.UID); ok {
		adapter := v.(*Adapter)
		return adapter.client.SetCHValue(val.Tag, val.V)
	}
	return errors.New("device not exists")
}

func (runner *Runner) GetRealtimeData(uid string) ([]interface{}, error) {
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)
		r, err := adapter.client.GetRealTimeData()
		if err != nil {
			return nil, err
		}
		values := make([]interface{}, 0)
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
				})
			}
		}

		return values, nil
	}
	return nil, errors.New("device not exists")
}

func (runner *Runner) Remove(uid string) {
	fmt.Println("stop device: ", uid)
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)
		adapter.Close()
		runner.adapters.Delete(uid)
	}
}

func (runner *Runner) Serve(adapter *Adapter) error {
	adapter.OnDeviceStatusChange(lang.AdapterInitializing)

	fmt.Printf("%# v", pretty.Formatter(adapter.conf))
	log.Info("start influx http client")

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

	log.Info("create influx db: ", adapter.conf.DB)

	_, err = c.Query(influx.Query{
		Database: adapter.conf.DB,
		Command:  fmt.Sprintf("CREATE DATABASE \"%s\"", adapter.conf.DB),
	})
	if err != nil {
		return err
	}

	getDataFN := func() error {
		bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
			Precision: "s",
			Database:  adapter.conf.DB,
		})
		for {
			select {
			case data := <-adapter.ch:
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
			log.Warnln("influx routine exit!")
			adapter.wg.Done()
		}()

		for {
			select {
			case <-runner.ctx.Done():
				return
			case <-adapter.done:
				return
			default:
				err := getDataFN()
				if err != nil {
					log.Error(err)
					return
				}
			}
		}
	}()

	go func() {
		defer func() {
			close(adapter.ch)
			adapter.wg.Done()
			log.Warnln("fetch data routine exit!")
		}()

	makeConnection:
		client := adapter.client
		for {
			log.Info("try connect to :", adapter.conf.Address)
			adapter.OnDeviceStatusChange(lang.Connecting)

			err := client.Connect(runner.ctx, adapter.conf.Address)
			if err != nil {
				if err == runner.ctx.Err() {
					return
				}
				select {
				case <-adapter.done:
					return
				case <-time.After(6 * time.Second):
					continue
				}
			} else {
				break
			}
		}

		adapter.OnDeviceStatusChange(lang.Connecting)

		for {
			select {
			case <-runner.ctx.Done():
				return
			case <-adapter.done:
				return
			case <-time.After(adapter.conf.Interval):
				log.Info("start fetch data from: ", adapter.conf.Address)
				err := runner.fetchData(adapter)
				if err != nil {
					log.Error(err)
					adapter.OnDeviceStatusChange(lang.Disconnected)
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
				adapter.ch <- data
				log.Printf("%s => %#v", ai.GetConfig().TagName, v)
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
				adapter.ch <- data
				log.Printf("%s => %#v", di.GetConfig().TagName, v)
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
				adapter.ch <- data
				log.Printf("%s => %#v", ao.GetConfig().TagName, v)
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
				adapter.ch <- data
				log.Printf("%s => %#v", do.GetConfig().TagName, v)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
