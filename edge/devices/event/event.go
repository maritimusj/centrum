package event

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/EventBus"
	"github.com/maritimusj/centrum/edge/devices/measure"
	"github.com/maritimusj/centrum/edge/lang"
	httpLoggerStore "github.com/maritimusj/centrum/edge/logStore/http"
	"github.com/maritimusj/centrum/json_rpc"
	log "github.com/sirupsen/logrus"
)

var (
	Event = EventBus.New()
)

var (
	__httpRequestCH = make(chan *HttpRequest, 1000)
)

type HttpRequest struct {
	url  string
	data []byte
}

const (
	DeviceStatusChanged = "device:status::changed"
	DevicePerfChanged   = "device:perf::changed"
	MeasureDiscovered   = "measure::discovered"
	MeasureAlarm        = "measure::alarm"
)

func isHttpTooBusy() bool {
	return len(__httpRequestCH) > 600
}

func Init(ctx context.Context) {
	eventsMap := map[string]interface{}{
		DeviceStatusChanged: OnDeviceStatusChanged,
		DevicePerfChanged:   OnDevicePerfChanged,
		MeasureDiscovered:   OnMeasureDiscovered,
		MeasureAlarm:        OnMeasureAlarm,
	}

	for e, fn := range eventsMap {
		_ = Event.SubscribeAsync(e, fn, false)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case request := <-__httpRequestCH:
				if request != nil {
					_, err := doHttpPost(request.url, request.data)
					if err != nil {
						log.Traceln("doHttpRequest:", err)
					}
				}
			}
		}
	}()
}

func Publish(topic string, args ...interface{}) {
	Event.Publish(topic, args...)
}

func HttpPost(url string, data interface{}) {
	x, err := json.Marshal(data)
	if err != nil {
		log.Traceln("[httpPost]", err)
		return
	}

	__httpRequestCH <- &HttpRequest{
		url:  url,
		data: x,
	}
}

func doHttpPost(url string, data []byte) ([]byte, error) {
	defer func() {
		recover()
	}()

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	resp, err := httpLoggerStore.DefaultHttpClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func OnDeviceStatusChanged(conf *json_rpc.Conf, status lang.StrIndex) {
	if conf.CallbackURL != "" && !isHttpTooBusy() {
		HttpPost(conf.CallbackURL, map[string]interface{}{
			"status": map[string]interface{}{
				"uid":   conf.UID,
				"index": status,
				"title": lang.Str(status),
			},
		})
	}
}

func OnDevicePerfChanged(conf *json_rpc.Conf, perf map[string]interface{}) {
	if conf.CallbackURL != "" && !isHttpTooBusy() {
		perf["uid"] = conf.UID
		HttpPost(conf.CallbackURL, map[string]interface{}{
			"perf": perf,
		})
	}
}

func OnMeasureDiscovered(conf *json_rpc.Conf, tagName, title string) {
	if conf.CallbackURL != "" && !isHttpTooBusy() {
		HttpPost(conf.CallbackURL, map[string]interface{}{
			"measure": map[string]interface{}{
				"uid":   conf.UID,
				"tag":   tagName,
				"title": title,
			},
		})
	}
}

func OnMeasureAlarm(conf *json_rpc.Conf, measureData *measure.Data) {
	defer measureData.Release()

	if conf.CallbackURL != "" {
		HttpPost(conf.CallbackURL, map[string]interface{}{
			"alarm": measureData,
		})
	}
}
