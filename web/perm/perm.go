package perm

import (
	"github.com/kataras/iris"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/util"
)

const (
	AdminUserKey        = "__admin__"
	DefaultAdminUserKey = "__defaultAdminUser__"
	DefaultEffect       = "__defaultEffect__"
)

func IsDefaultAdminUser(ctx iris.Context) bool {
	v := ctx.Values().Get(DefaultAdminUserKey)
	return v != nil
}

func AdminUser(ctx iris.Context) model.User {
	v := ctx.Values().Get(AdminUserKey)
	if v != nil {
		return v.(model.User)
	}
	return nil
}

func IsAdminOk(ctx iris.Context, fn func(admin model.User) interface{}) (interface{}, bool) {
	user := AdminUser(ctx)
	if user != nil && user.IsEnabled() {
		if fn != nil {
			return fn(user), true
		}
		return nil, true
	}
	return nil, false
}

func MustAdminOk(ctx iris.Context, fn func(admin model.User) interface{}) interface{} {
	result, ok := IsAdminOk(ctx, fn)
	return util.If(ok, result, lang.Error(lang.ErrInvalidUser))
}

func IsAllowed(ctx iris.Context, resource resource.Resource, action resource.Action) (bool, error) {
	//if IsDefaultAdminUser(ctx) {
	//	return true, nil
	//}

	user := AdminUser(ctx)
	if user != nil && user.IsEnabled() {
		return user.IsAllow(resource, action)
	}

	return false, lang.Error(lang.ErrInvalidUser)
}

func Deny(ctx iris.Context, res resource.Resource, action resource.Action) bool {
	return !Allow(ctx, res, action)
}

func Allow(ctx iris.Context, res resource.Resource, action resource.Action) bool {
	allowed, err := IsAllowed(ctx, res, action)
	if allowed {
		return true
	}

	if err == lang.Error(lang.ErrPolicyNotFound) {
		println("没找到策略，默认通过！resource title: ", res.ResourceTitle(), "id: ", res.ResourceID(), " class: ", res.ResourceClass())
		return ctx.Values().Get(DefaultEffect).(resource.Effect) == resource.Allow
	}
	return false
}
