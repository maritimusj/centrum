package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	log2 "github.com/maritimusj/centrum/gate/web/api/log"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
)

func LogList(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, device, resource2.View) {
			return lang2.ErrNoPermission
		}

		return log2.GetLogList(ctx, device.OrganizationID(), device.UID())
	})
}

func LogDelete(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, device, resource2.Ctrl) {
			return lang2.ErrNoPermission
		}

		return log2.DeleteLog(ctx, device.OrganizationID(), device.UID())
	})
}
