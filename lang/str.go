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

	ApiResourcesMap = [][3]string{
		{resource.UserList, "用户列表", "查看用户列表"},
		{resource.UserCreate, "创建用户", "创建新用户"},
		{resource.UserUpdate, "编辑用户", "修改用户资料"},
		{resource.UserDetail, "查看资料", "查看用户资料"},
		{resource.UserDelete, "删除用户", "删除用户"},

		{resource.RoleList, "角色列表", "查看角色列表"},
		{resource.RoleCreate, "创建角色", "创建新的角色"},
		{resource.RoleUpdate, "编辑角色", "修改角色权限"},
		{resource.RoleDetail, "查看角色", "查看角色详情"},
		{resource.RoleDelete, "删除角色", "删除角色"},
	}
)

func ResourceClassTitle(class resource.Class) string {
	if v, ok := resourceGroupsMap[class]; ok {
		return v
	}
	panic(errors.New("unknown resource class"))
}
