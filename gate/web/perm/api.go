package perm

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/util"
)

func CheckApiPerm(ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	checkFN := func() interface{} {
		if app.IsDefaultAdminUser(admin) {
			return nil
		}

		router := ctx.GetCurrentRoute()
		res, err := s.GetApiResource(router.Name())
		if err != nil {
			if err != lang.ErrApiResourceNotFound.Error() {
				return err
			}
			return util.If(app.Config.DefaultEffect() == resource.Allow, nil, lang.ErrNoPermission)
		}

		if app.Allow(admin, res, resource.Invoke) {
			return nil
		}
		return lang.ErrNoPermission
	}

	if err := checkFN(); err != nil {
		return response.Wrap(err)
	}

	ctx.Next()
	return nil
}
