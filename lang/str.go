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
		resource.Measure:   "点位",
		resource.Equipment: "自定义设备",
		resource.State:     "自定义点位",
	}

	ApiResourcesMap = [...][3]string{
		{resource.UserList, "用户列表", "查看用户列表"},
		{resource.UserCreate, "创建用户", "创建新用户"},
		{resource.UserUpdate, "编辑用户", "修改用户资料"},
		{resource.UserDetail, "查看资料", "查看用户资料"},
		{resource.UserDelete, "删除用户", "删除用户"},

		/*
			{resource.RoleList, "角色列表", "查看角色列表"},
			{resource.RoleCreate, "创建角色", "创建新的角色"},
			{resource.RoleUpdate, "编辑角色", "修改角色权限"},
			{resource.RoleDetail, "查看角色", "查看角色详情"},
			{resource.RoleDelete, "删除角色", "删除角色"},
		*/

		{resource.GroupList, "分组列表", "查看全部设备分组"},
		{resource.GroupCreate, "创建分组", "创建新的设备分组"},
		{resource.GroupDetail, "查看分组", "查看分组详情"},
		{resource.GroupUpdate, "编辑分组", "修改分组信息"},
		{resource.GroupDelete, "删除分级", "删除设备分组"},

		{resource.DeviceList, "设备列表", "查看设备列表"},
		{resource.DeviceCreate, "添加设备", "添加新设备"},
		{resource.DeviceDetail, "查看设备", "查看设备详情"},
		{resource.DeviceUpdate, "编辑设备", "修改设备详情"},
		{resource.DeviceDelete, "删除设备", "移除一个设备"},

		{resource.MeasureList, "点位列表", "查看点位列表"},
		{resource.MeasureCreate, "添加点位", "添加一个新的点位"},
		{resource.MeasureDetail, "查看点位", "查看点位详情"},
		{resource.MeasureUpdate, "编辑点位", "修改点位"},
		{resource.MeasureDelete, "删除点位", "删除点位"},

		{resource.EquipmentList, "自定义设备列表", "查看自定义设备列表"},
		{resource.EquipmentCreate, "创建自定义设备", "创建一个新的自定义设备"},
		{resource.EquipmentDetail, "查看自定义设备", "查看自定义设备详情"},
		{resource.EquipmentUpdate, "编辑自定义设备", "修改自定义设备"},
		{resource.EquipmentDelete, "删除自定义设备", "删除一个自定设备"},

		{resource.StateList, "自定义点位列表", "查看自定义点位列表"},
		{resource.StateCreate, "添加自定义点位", "添加一个新的自定义点位"},
		{resource.StateDetail, "查看自定义点位", "查看片定义点位详情"},
		{resource.StateUpdate, "编辑自定义点位", "修改自定义点位"},
		{resource.StateDelete, "删除自定义点位", "删除一个自定义点位"},
	}
)

func ResourceClassTitle(class resource.Class) string {
	if v, ok := resourceGroupsMap[class]; ok {
		return v
	}
	panic(errors.New("unknown resource class"))
}
