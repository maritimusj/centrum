package global

import (
	"fmt"

	"github.com/maritimusj/centrum/gate/web/model"
)

func UpdateDeviceStatus(device model.Device, index int, title string) {
	path := fmt.Sprintf("device.%d.stats", device.GetID())
	_ = Stats.Set(path, map[string]interface{}{
		"index": index,
		"title": title,
	})
}

func GetDeviceStatus(device model.Device) (int, string) {
	key := fmt.Sprintf("device.%d.stats", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.(map[string]interface{}); ok {
			return int(vv["index"].(float64)), vv["title"].(string)
		}
	}
	return 0, ""
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
