package user

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/api/log"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/response"
)

func LogList(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := app.Store().GetUser(userID)
		if err != nil {
			return err
		}

		return log.GetLogList(ctx, user.OrganizationID(), user.UID())
	})
}

func LogDelete(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := app.Store().GetUser(userID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.IsDefaultAdminUser(admin) {
			return lang.ErrNoPermission
		}

		return log.DeleteLog(ctx, user.OrganizationID(), user.UID())
	})
}
