package devices

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/maritimusj/centrum/edge/logStore"

	"github.com/maritimusj/centrum/edge/devices/InverseServer"
	"github.com/maritimusj/centrum/edge/devices/ep6v2"
	"github.com/maritimusj/centrum/edge/devices/measure"
	"github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/centrum/json_rpc"
	log "github.com/sirupsen/logrus"

	influx "github.com/influxdata/influxdb1-client/v2"
	httpLogStore "github.com/maritimusj/centrum/edge/logStore/http"
)

type Runner struct {
	ctx      context.Context
	adapters sync.Map

	RestartMainFN func()
}

func New() *Runner {
	runner := &Runner{
		ctx: context.Background(),
	}
	return runner
}

func (runner *Runner) Close() {
	runner.adapters.Range(func(key, v interface{}) bool {
		adapter := v.(*Adapter)
		adapter.Close()
		return true
	})
}

func (runner *Runner) StartInverseServer(conf *json_rpc.InverseConf) error {
	return InverseServer.Start(runner.ctx, conf.Address, conf.Port)
}

func (runner *Runner) GetBaseInfo(uid string) (map[string]interface{}, error) {
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)

		baseInfo := make(map[string]interface{})

		model, err := adapter.device.GetModel()
		if err != nil {
			return nil, err
		}

		baseInfo["model"] = model.ID
		baseInfo["version"] = model.Version
		baseInfo["title"] = model.Title

		addr, err := adapter.device.GetAddr()
		if err != nil {
			return nil, err
		}

		baseInfo["addr"] = addr.Ip.String() + "/" + addr.Mask.String()
		baseInfo["mac"] = addr.Mac.String()

		baseInfo["status"] = map[string]interface{}{
			"index": adapter.device.GetStatus(),
			"title": adapter.device.GetStatusTitle(),
		}

		return baseInfo, nil
	}

	return nil, lang.Error(lang.ErrDeviceNotExists)
}

func (runner *Runner) needRestartAdapter(conf *json_rpc.Conf, newConf *json_rpc.Conf) bool {
	return conf.Address != newConf.Address ||
		conf.InfluxDBUrl != newConf.InfluxDBUrl ||
		conf.InfluxDBUserName != newConf.InfluxDBUserName ||
		conf.InfluxDBPassword != newConf.InfluxDBPassword ||
		conf.DB != newConf.DB ||
		conf.CallbackURL != newConf.CallbackURL
}

func (runner *Runner) Restart() {
	if runner.RestartMainFN != nil {
		runner.RestartMainFN()
	}
}

func (runner *Runner) Active(conf *json_rpc.Conf) error {
	log.Traceln("active:", conf.UID, conf.Address)

	select {
	case <-runner.ctx.Done():
		return runner.ctx.Err()
	default:
	}

	if v, ok := runner.adapters.Load(conf.UID); ok {
		adapter := v.(*Adapter)
		if !adapter.IsAlive() || runner.needRestartAdapter(adapter.conf, conf) {
			log.Warningln("restart adapter, isAlive:", adapter.IsAlive(), "conf changed:", runner.needRestartAdapter(adapter.conf, conf))
			runner.adapters.Delete(conf.UID)
			adapter.Close()
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

			adapter.OnDeviceStatusChanged(adapter.device.GetStatus())
			return nil
		}
	}

	var (
		logger     = log.New()
		loggerHook logStore.Store
	)

	if conf.CallbackURL != "" {
		if conf.LogLevel == "" {
			conf.LogLevel = "error"
		}

		level, err := log.ParseLevel(conf.LogLevel)
		if err != nil {
			return err
		}

		loggerHook = httpLogStore.New()
		loggerHook.SetUID(conf.UID)

		logger.SetLevel(level)
		logger.AddHook(loggerHook)

		err = loggerHook.Open(runner.ctx, conf.CallbackURL, level)
		if err != nil {
			return err
		}
	}

	adapter := &Adapter{
		device:         ep6v2.New(),
		conf:           conf,
		logger:         logger,
		loggerStore:    loggerHook,
		lastActiveTime: time.Now(),
		measureDataCH:  make(chan *measure.Data, 60),
		done:           make(chan struct{}),
	}

	err := runner.Serve(adapter)
	if err != nil {
		return err
	}

	if exists, ok := runner.adapters.LoadOrStore(conf.UID, adapter); ok {
		exists.(*Adapter).Close()
	}

	return nil
}

func (runner *Runner) GetValue(ch *json_rpc.CH) (retVal interface{}, err error) {
	if v, ok := runner.adapters.Load(ch.UID); ok {
		adapter := v.(*Adapter)
		return adapter.device.GetCHValue(ch.Tag)
	}
	return nil, lang.Error(lang.ErrDeviceNotExists)
}

func (runner *Runner) SetValue(val *json_rpc.Value) error {
	if v, ok := runner.adapters.Load(val.UID); ok {
		adapter := v.(*Adapter)
		return adapter.device.SetCHValue(val.Tag, val.V)
	}
	return lang.Error(lang.ErrDeviceNotExists)
}

func (runner *Runner) GetRealtimeData(uid string) ([]map[string]interface{}, error) {
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)
		r, err := adapter.device.GetRealTimeData()
		if err != nil {
			return nil, err
		}

		defer r.Release()

		values := make([]map[string]interface{}, 0)
		for i := 0; i < r.AINum(); i++ {
			ai, err := adapter.device.GetAI(i)
			if err != nil {
				return values, err
			}

			entry := map[string]interface{}{
				"tag":   ai.GetConfig().TagName,
				"title": ai.GetConfig().Title,
				"unit":  ai.GetConfig().Uint,
			}

			if v, ok := r.GetAIValue(i, ai.GetConfig().Point); ok {
				av, x := ai.CheckAlarm(v)

				entry["alarm"] = ep6v2.AlarmDesc(av)
				entry["value"] = v

				if av != ep6v2.AlarmNormal {
					entry["threshold"] = x
				}
			}

			values = append(values, entry)
			adapter.OnMeasureDiscovered(ai.GetConfig().TagName, ai.GetConfig().Title)
		}

		for i := 0; i < r.AONum(); i++ {
			ao, err := adapter.device.GetAO(i)
			if err != nil {
				return values, err
			}

			entry := map[string]interface{}{
				"tag":   ao.GetConfig().TagName,
				"title": ao.GetConfig().Title,
				"unit":  ao.GetConfig().Uint,
			}

			if v, ok := r.GetAOValue(i); ok {
				entry["value"] = v
			}

			values = append(values, entry)
			adapter.OnMeasureDiscovered(ao.GetConfig().TagName, ao.GetConfig().Title)
		}

		for i := 0; i < r.DINum(); i++ {
			di, err := adapter.device.GetDI(i)
			if err != nil {
				return values, err
			}

			entry := map[string]interface{}{
				"tag":   di.GetConfig().TagName,
				"title": di.GetConfig().Title,
			}

			if v, ok := r.GetDIValue(i); ok {
				entry["value"] = v
			}

			values = append(values, entry)
			adapter.OnMeasureDiscovered(di.GetConfig().TagName, di.GetConfig().Title)
		}

		for i := 0; i < r.DONum(); i++ {
			do, err := adapter.device.GetDO(i)
			if err != nil {
				return values, err
			}

			entry := map[string]interface{}{
				"tag":   do.GetConfig().TagName,
				"title": do.GetConfig().Title,
				"ctrl":  do.GetConfig().IsManual,
			}

			if v, ok := r.GetDOValue(i); ok {
				entry["value"] = v
			}

			values = append(values, entry)
			adapter.OnMeasureDiscovered(do.GetConfig().TagName, do.GetConfig().Title)
		}

		return values, nil
	}

	return nil, lang.Error(lang.ErrDeviceNotExists)
}

func (runner *Runner) Reset(uid string) {
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)
		adapter.device.Reset()
	}
}

func (runner *Runner) Remove(uid string) {
	if v, ok := runner.adapters.Load(uid); ok {
		adapter := v.(*Adapter)

		runner.adapters.Delete(uid)
		adapter.Close()
	}
}

func (runner *Runner) InitInfluxDB(conf *json_rpc.Conf) (influx.Client, error) {
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     conf.InfluxDBUrl,
		Username: conf.InfluxDBUserName,
		Password: conf.InfluxDBPassword,
	})
	if err != nil {
		return nil, err
	}

	if _, _, err = c.Ping(3 * time.Second); err != nil {
		return nil, err
	}

	_, err = c.Query(influx.Query{
		Database: conf.DB,
		Command:  fmt.Sprintf("CREATE DATABASE \"%s\"", conf.DB),
	})

	if err != nil {
		return nil, err
	}
	return c, nil
}

func (runner *Runner) getMeasureData(client influx.Client, db string, ch <-chan *measure.Data) error {
	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Precision: "ns",
		Database:  db,
	})
	for {
		select {
		case data := <-ch:
			if data == nil {
				return errors.New("ch closed")
			}

			point, err := influx.NewPoint(data.Name, data.Tags, data.Fields, data.Time)

			data.Release()

			if err != nil {
				log.Errorln(err)
				continue
			} else {
				bp.AddPoint(point)
			}

		case <-time.After(1 * time.Second):
			if len(bp.Points()) > 0 {
				err := client.Write(bp)
				if err != nil {
					return err
				} else {
					return err
				}
			}
		}
	}
}
func (runner *Runner) Serve(adapter *Adapter) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = lang.InternalError(errors.New("unknown error"))
			}
		}

		adapter.OnDeviceStatusChanged(lang.MalFunctioned)
	}()

	adapter.OnDeviceStatusChanged(lang.AdapterInitializing)

	c, err := runner.InitInfluxDB(adapter.conf)
	if err != nil {
		adapter.OnDeviceStatusChanged(lang.InfluxDBError)
		return err
	}

	defer func() {
		_ = c.Close()
	}()

	adapter.wg.Add(2)
	go func() {
		defer func() {
			adapter.wg.Done()
		}()

		for {
			select {
			case <-runner.ctx.Done():
				return
			case <-adapter.done:
				return
			default:
				err := runner.getMeasureData(c, adapter.conf.DB, adapter.measureDataCH)
				if err != nil {
					adapter.logger.Error(err)
					//return
				}
			}
		}
	}()

	go func() {
		defer func() {
			close(adapter.measureDataCH)
			adapter.wg.Done()
		}()

		const delay = 10 * time.Second

	tryConnectToDevice:
		device := adapter.device
		for {
			adapter.heartBeat()

			adapter.OnDeviceStatusChanged(lang.Connecting)

			err := device.Connect(runner.ctx, adapter.conf.Address)
			if err != nil {
				if err == runner.ctx.Err() {
					return
				}

				adapter.OnDeviceStatusChanged(lang.Disconnected)

				select {
				case <-adapter.done:
					return
				case <-time.After(delay):
					continue
				}
			} else {
				if device.IsConnected() {
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
				adapter.heartBeat()

				err := runner.gatherData(adapter)
				if err != nil {
					adapter.logger.Errorln(err)
					device.Close()

					go adapter.OnDeviceStatusChanged(lang.Disconnected)
					go adapter.OnDevicePerfChanged(map[string]interface{}{
						"delay": -1,
					})
					goto tryConnectToDevice
				}
			}
		}
	}()

	return nil
}

func (runner *Runner) gatherData(adapter *Adapter) error {
	client := adapter.device

	data, err := client.GetRealTimeData()
	if err != nil {
		return err
	}

	defer data.Release()

	adapter.OnDevicePerfChanged(map[string]interface{}{
		"delay": data.TimeUsed(),
	})

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
				av, x := ai.CheckAlarm(v)

				data := measure.New(ai.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", ai.GetConfig().TagName)
				data.AddTag("title", ai.GetConfig().Title)
				data.AddTag("alarm", ep6v2.AlarmDesc(av))
				data.AddField("val", v)

				if av != ep6v2.AlarmNormal {
					data.AddTag("unit", ai.GetConfig().Uint)
					data.AddField("threshold", x)
				}

				adapter.measureDataCH <- data
				if av != 0 {
					adapter.OnMeasureAlarm(data.Clone())
				}
			}

			adapter.OnMeasureDiscovered(ai.GetConfig().TagName, ai.GetConfig().Title)
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
				data := measure.New(di.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", di.GetConfig().TagName)
				data.AddTag("title", di.GetConfig().Title)
				data.AddField("val", v)
				adapter.measureDataCH <- data
			}
			adapter.OnMeasureDiscovered(di.GetConfig().TagName, di.GetConfig().Title)
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
				data := measure.New(ao.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", ao.GetConfig().TagName)
				data.AddTag("title", ao.GetConfig().Title)
				data.AddField("val", v)
				adapter.measureDataCH <- data
			}
			adapter.OnMeasureDiscovered(ao.GetConfig().TagName, ao.GetConfig().Title)
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
				data := measure.New(do.GetConfig().TagName)
				data.AddTag("uid", adapter.conf.UID)
				data.AddTag("address", adapter.conf.Address)
				data.AddTag("tag", do.GetConfig().TagName)
				data.AddTag("title", do.GetConfig().Title)
				data.AddField("val", v)
				adapter.measureDataCH <- data
			}
			adapter.OnMeasureDiscovered(do.GetConfig().TagName, do.GetConfig().Title)
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
