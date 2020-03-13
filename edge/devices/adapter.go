package devices

import (
	"sync"
	"time"

	"github.com/maritimusj/centrum/edge/logStore"

	_ "github.com/influxdata/influxdb1-client"
	"github.com/maritimusj/centrum/edge/devices/ep6v2"
	"github.com/maritimusj/centrum/edge/devices/event"
	"github.com/maritimusj/centrum/edge/devices/measure"
	"github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/json_rpc"
	"github.com/maritimusj/centrum/synchronized"
	log "github.com/sirupsen/logrus"
)

type Adapter struct {
	device *ep6v2.Device
	conf   *json_rpc.Conf

	measureDataCH chan *measure.Data

	logger      *log.Logger
	loggerStore logStore.Store

	lastActiveTime time.Time

	done chan struct{}
	wg   sync.WaitGroup
}

func (adapter *Adapter) heartBeat() {
	adapter.lastActiveTime = time.Now()
}

func (adapter *Adapter) IsAlive() bool {
	if adapter == nil || adapter.device == nil {
		return false
	}

	select {
	case <-adapter.done:
		return false
	default:
	}

	//log.Debugln("IsAlive:", time.Now().Sub(adapter.lastActiveTime).String())

	return time.Now().Sub(adapter.lastActiveTime) < adapter.conf.Interval*2
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
	<-synchronized.Do(adapter, func() interface{} {
		defer func() {
			recover()
		}()

		if adapter.device != nil {
			adapter.device.Close()
			adapter.device = nil
		}

		close(adapter.done)
		adapter.wg.Wait()

		if adapter.loggerStore != nil {
			adapter.loggerStore.Close()
		}

		return nil
	})
}

func (adapter *Adapter) OnDeviceStatusChanged(index lang.StrIndex) {
	event.Publish(event.DeviceStatusChanged, adapter.conf, index)
}

func (adapter *Adapter) OnDevicePerfChanged(perf map[string]interface{}) {
	event.Publish(event.DevicePerfChanged, adapter.conf, perf)
}

func (adapter *Adapter) OnMeasureDiscovered(tagName, title string) {
	path := "tag." + adapter.conf.UID + "." + tagName
	if v, ok := global.Params.Get(path); !ok || v.(string) != title {
		_ = global.Params.Set(path, title)
		event.Publish(event.MeasureDiscovered, adapter.conf, tagName, title)
	}
}

func (adapter *Adapter) OnMeasureAlarm(data *measure.Data) {
	event.Publish(event.MeasureAlarm, adapter.conf, data)
}
