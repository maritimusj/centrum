package user

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	log2 "github.com/maritimusj/centrum/gate/web/api/log"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	response2 "github.com/maritimusj/centrum/gate/web/response"
)

func LogList(userID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		user, err := app2.Store().GetUser(userID)
		if err != nil {
			return err
		}

		return log2.GetLogList(ctx, user.OrganizationID(), user.UID())
	})
}

func LogDelete(userID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		user, err := app2.Store().GetUser(userID)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.IsDefaultAdminUser(admin) {
			return lang2.ErrNoPermission
		}

		return log2.DeleteLog(ctx, user.OrganizationID(), user.UID())
	})
}
