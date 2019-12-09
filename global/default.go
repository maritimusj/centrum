package global

import (
	"fmt"
	"time"

	edgeLang "github.com/maritimusj/centrum/edge/lang"

	"github.com/maritimusj/centrum/gate/web/model"
)

func UpdateDeviceStatus(device model.Device, index int, title string) {
	path := fmt.Sprintf("device.%d.stats", device.GetID())
	data := map[string]interface{}{
		"index": index,
		"title": title,
	}
	orgIndex, _, from := GetDeviceStatus(device)
	if orgIndex != int(edgeLang.Connected) && index == int(edgeLang.Connected) {
		data["from"] = time.Now().Format(time.RFC3339)
	} else {
		data["from"] = from
	}
	_ = Stats.Set(path, data)
}

func GetDeviceStatus(device model.Device) (int, string, time.Time) {
	key := fmt.Sprintf("device.%d.stats", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.(map[string]interface{}); ok {
			from, _ := time.Parse(time.RFC3339, vv["from"].(string))
			return int(vv["index"].(float64)), vv["title"].(string), from
		}
	}
	return int(edgeLang.EdgeUnknownState), edgeLang.Str(edgeLang.EdgeUnknownState), time.Now()
}

func UpdateDevicePerf(device model.Device, data map[string]interface{}) {
	path := fmt.Sprintf("device.%d.perf", device.GetID())
	_ = Stats.Set(path, data)
}

func GetDevicePerf(device model.Device) map[string]interface{} {
	key := fmt.Sprintf("device.%d.perf", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.(map[string]interface{}); ok {
			return vv
		}
	}
	return map[string]interface{}{}
}
