package devices

import (
	_ "github.com/influxdata/influxdb1-client"
	"github.com/maritimusj/centrum/edge/devices/ep6v2"
	"github.com/maritimusj/centrum/edge/devices/event"
	"github.com/maritimusj/centrum/edge/devices/measure"
	"github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/json_rpc"
	"github.com/maritimusj/centrum/synchronized"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Adapter struct {
	device *ep6v2.Device
	conf   *json_rpc.Conf

	measureDataCH chan *measure.Data

	logger *log.Logger

	done chan struct{}
	wg   sync.WaitGroup
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
		if adapter.device != nil {
			adapter.device.Close()
			adapter.device = nil
		}

		close(adapter.done)
		adapter.wg.Wait()
		return nil
	})
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

func (adapter *Adapter) OnMeasureAlarm(data *measure.Data) {
	event.Publish(event.MeasureAlarm, adapter.conf, data)
}
