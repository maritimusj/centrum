package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/api/log"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
)

func LogList(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		return log.GetLogList(ctx, device.OrganizationID(), device.UID())
	})
}

func LogDelete(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		return log.DeleteLog(ctx, device.OrganizationID(), device.UID())
	})
}
