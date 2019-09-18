package perm

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/response"
)

func Check(ctx iris.Context, store store.Store, cfg config.Config) hero.Result {
	result := MustAdminOk(ctx, func(admin model.User) interface{} {
		//总是允许系统默认用户
		if IsDefaultAdminUser(ctx) {
			return nil
		}

		router := ctx.GetCurrentRoute()
		res, err := store.GetApiResource(router.Name())
		if err != nil {
			if err != lang.Error(lang.ErrApiResourceNotFound) {
				return err
			}
			return util.If(cfg.DefaultEffect() == resource.Allow, nil, lang.Error(lang.ErrNoPermission))
		}

		allowed, err := admin.IsAllow(res, resource.Invoke)
		if err != nil {
			if err != lang.Error(lang.ErrPolicyNotFound) {
				return err
			}
			return util.If(cfg.DefaultEffect() == resource.Allow, nil, lang.Error(lang.ErrNoPermission))
		}

		return util.If(allowed, nil, lang.Error(lang.ErrNoPermission))
	})

	if result != nil {
		if err := result.(error); err != nil {
			ctx.StopExecution()
			return response.Wrap(result)
		}
	}

	ctx.Next()
	return nil
}
