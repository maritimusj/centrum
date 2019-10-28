package global

import (
	"fmt"
	"github.com/maritimusj/centrum/gate/web/model"
)

func UpdateDeviceStatus(device model.Device, index int, title string) {
	key := fmt.Sprintf("device:%d", device.GetID())
	println("UpdateDeviceStatus: ", device.GetID(), index, title)
	Stats.Set(key, [2]interface{}{index, title})
}

func GetDeviceStatus(device model.Device) (int, string) {
	key := fmt.Sprintf("device:%d", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.([2]interface{}); ok {
			println("GetDeviceStatus: ", vv[0].(int), vv[1].(string))
			return vv[0].(int), vv[1].(string)
		}
	}
	return 0, ""
}
