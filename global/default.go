package global

import (
	"fmt"
	"github.com/maritimusj/centrum/web/model"
)

func UpdateDeviceStatus(device model.Device, index int, title string) {
	key := fmt.Sprintf("device:%d", device.GetID())
	Stats.Set(key, [2]interface{}{index, title})
}

func GetDeviceStatus(device model.Device) (int, string) {
	key := fmt.Sprintf("device:%d", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.([2]interface{}); ok {
			return vv[0].(int), vv[1].(string)
		}
	}
	return 0, ""
}

func UpdateEquipmentStatus(equipment model.Equipment, index int, title string) {
	key := fmt.Sprintf("equipment:%d", equipment.GetID())
	Stats.Set(key, [2]interface{}{index, title})
}

func GetEquipmentStatus(equipment model.Equipment) (int, string) {
	key := fmt.Sprintf("equipment:%d", equipment.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.([2]interface{}); ok {
			return vv[0].(int), vv[1].(string)
		}
	}
	return 0, ""
}
