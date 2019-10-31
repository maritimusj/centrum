package global

import (
	"fmt"
	"github.com/maritimusj/centrum/gate/web/model"
)

func UpdateDeviceStatus(device model.Device, index int, title string) {
	path := fmt.Sprintf("device.stats.%d", device.GetID())
	_ = Stats.Set(path, map[string]interface{}{
		"index": index,
		"title": title,
	})
}

func GetDeviceStatus(device model.Device) (int, string) {
	key := fmt.Sprintf("device.stats.%d", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.(map[string]interface{}); ok {
			return int(vv["index"].(float64)), vv["title"].(string)
		}
	}
	return 0, ""
}
