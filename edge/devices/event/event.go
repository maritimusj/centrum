package event

import (
	"bytes"
	"encoding/json"
	"github.com/asaskevich/EventBus"
	"github.com/maritimusj/centrum/edge/lang"
	httpLoggerStore "github.com/maritimusj/centrum/edge/logStore/http"
	"github.com/maritimusj/centrum/json_rpc"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var (
	Event = EventBus.New()
)

const (
	DeviceStatusChanged = "device:status::changed"
	MeasureDiscovered   = "measure::discovered"
)

func init() {
	eventsMap := map[string]interface{}{
		DeviceStatusChanged: OnDeviceStatusChanged,
		MeasureDiscovered:   OnMeasureDiscovered,
	}

	for e, fn := range eventsMap {
		_ = Event.SubscribeAsync(e, fn, false)
	}
}

func Publish(topic string, args ...interface{}) {
	Event.Publish(topic, args...)
}

func HttpPost(url string, data interface{}) ([]byte, error) {
	b, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
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
	if conf.CallbackURL != "" {
		data, err := HttpPost(conf.CallbackURL, map[string]interface{}{
			"status": map[string]interface{}{
				"uid":   conf.UID,
				"index": status,
				"title": lang.Str(status),
			},
		})
		if err != nil {
			log.Errorf("[OnDeviceStatusChanged] %s", err)
			return
		}

		println("[OnDeviceStatusChanged]", conf.CallbackURL, string(data))
	}
}

func OnMeasureDiscovered(conf *json_rpc.Conf, tagName, title string) {
	if conf.CallbackURL != "" {
		data, err := HttpPost(conf.CallbackURL, map[string]interface{}{
			"measure": map[string]interface{}{
				"uid":   conf.UID,
				"tag":   tagName,
				"title": title,
			},
		})
		if err != nil {
			log.Errorf("[OnMeasureDiscovered] %s", err)
			return
		}

		println("[OnMeasureDiscovered]", conf.CallbackURL, string(data))
	}
}
