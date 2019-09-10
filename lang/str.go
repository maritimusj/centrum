package lang

import (
	"errors"
	"github.com/maritimusj/centrum/resource"
)

var (
	resourceGroupsMap = map[resource.Class]string{
		resource.Default:   "默认分组",
		resource.Api:       "后台权限",
		resource.Group:     "设备分组",
		resource.Device:    "设备",
		resource.Equipment: "自定义设备",
		resource.Measure:   "设备点位",
		resource.State:     "自定义点位",
	}
)

func ResourceClassTitle(group resource.Class) string {
	if v, ok := resourceGroupsMap[group]; ok {
		return v
	}
	panic(errors.New("unknown resource group"))
}
