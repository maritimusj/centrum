package perm

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/util"
)

func CheckApiPerm(ctx iris.Context) hero.Result {
	s := app2.Store()
	admin := s.MustGetUserFromContext(ctx)

	checkFN := func() interface{} {
		if app2.IsDefaultAdminUser(admin) {
			return nil
		}

		router := ctx.GetCurrentRoute()
		res, err := s.GetApiResource(router.Name())
		if err != nil {
			if err != lang2.Error(lang2.ErrApiResourceNotFound) {
				return err
			}
			return util.If(app2.Config.DefaultEffect() == resource2.Allow, nil, lang2.ErrNoPermission)
		}

		if app2.Allow(admin, res, resource2.Invoke) {
			return nil
		}
		return lang2.ErrNoPermission
	}

	if err := checkFN(); err != nil {
		return response2.Wrap(err)
	}

	ctx.Next()
	return nil
}
