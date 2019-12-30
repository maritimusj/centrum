package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/logStore"
	"github.com/maritimusj/centrum/gate/web/api/log"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"

	log2 "github.com/sirupsen/logrus"
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

		err = app.LogDBStore.Delete(device.OrganizationID(), device.UID())
		if err != nil {
			return err
		}

		logStr := lang.Str(lang.DeviceLogDeletedByUser, admin.Name(), device.Title())

		log2.WithField("src", logStore.SystemLog).Info(logStr)
		device.Logger().Infoln(logStr)

		return lang.Ok
	})
}
