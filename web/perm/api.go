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
	log "github.com/sirupsen/logrus"
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
			if err == lang.Error(lang.ErrApiResourceNotFound) && cfg.DefaultEffect() == resource.Allow {
				log.Tracef("没找到API资源，user: %s, router: %s 默认通过！", admin.Name(), router.Name())
				return nil
			}
			return err
		}

		allowed, err := admin.IsAllow(res, resource.Invoke)
		if allowed {
			return nil
		}

		if err == lang.Error(lang.ErrPolicyNotFound) && cfg.DefaultEffect() == resource.Allow {
			log.Tracef("没找到策略，user: %s, router: %s 默认通过！", admin.Name(), router.Name())
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
