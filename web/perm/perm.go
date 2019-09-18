package perm

import (
	"github.com/kataras/iris"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/util"
	log "github.com/sirupsen/logrus"
)

const (
	AdminUserKey        = "__admin__"
	DefaultAdminUserKey = "__defaultAdminUser__"
	DefaultEffect       = "__defaultEffect__"
)

func IsDefaultAdminUser(ctx iris.Context) bool {
	v := ctx.Values().Get(DefaultAdminUserKey)
	return v != nil && v.(bool)
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
	if IsDefaultAdminUser(ctx) {
		return true, nil
	}

	user := AdminUser(ctx)
	if user == nil {
		return false, lang.Error(lang.ErrInvalidUser)
	}

	if resource.OrganizationID() != 0 && resource.OrganizationID() != user.OrganizationID() {
		return false, lang.Error(lang.ErrNoPermission)
	}

	if !user.IsEnabled() {
		return false, lang.Error(lang.ErrUserDisabled)
	}

	return user.IsAllow(resource, action)
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
		log.Trace("没找到策略，使用默认设置：resource title: ", res.ResourceTitle(), "id: ", res.ResourceID(), " class: ", res.ResourceClass())
		return ctx.Values().Get(DefaultEffect).(resource.Effect) == resource.Allow
	}

	return false
}
