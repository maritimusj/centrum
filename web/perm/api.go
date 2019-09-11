package perm

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
	"log"
)

func Check(ctx iris.Context, store store.Store, cfg config.Config) hero.Result {
	result := MustAdminOk(ctx, func(admin model.User) interface{} {
		//总是允许系统默认用户
		if admin.Name() == cfg.DefaultUserName() {
			return nil
		}

		router := ctx.GetCurrentRoute()
		res, err := store.GetApiResource(router.Name())
		if err != nil {
			return err
		}

		allowed, err := admin.IsAllowed(res, resource.Invoke)
		if allowed {
			return nil
		}

		if cfg.DefaultEffect() == resource.Allow && err == lang.Error(lang.ErrPolicyNotFound) {
			log.Printf("user: %s, router: %s allowed by default policy\r\n", admin.Name(), router.Name())
			return nil
		}
		return lang.Error(lang.ErrNoPermission)
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
