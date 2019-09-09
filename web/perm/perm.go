package perm

import (
	"github.com/kataras/iris"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/util"
)

func AdminUser(ctx iris.Context) model.User {
	return ctx.Values().Get("__admin__").(model.User)
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
