package lang

import (
	"errors"
	"github.com/maritimusj/centrum/model"
)

var (
	resourceGroupsMap = map[model.ResourceClass]string{
		model.DefaultResClass:   "默认分组",
		model.ApiResClass:       "后台权限",
		model.GroupResClass:     "设备分组",
		model.DeviceResClass:    "设备",
		model.EquipmentResClass: "自定义设备",
		model.MeasureResClass:   "设备点位",
		model.StateResClass:     "自定义点位",
	}
)

func ResourceClassTitle(group model.ResourceClass) string {
	if v, ok := resourceGroupsMap[group]; ok {
		return v
	}
	panic(errors.New("unknown resource group"))
}
