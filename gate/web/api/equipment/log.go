package equipment

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

func LogList(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		return log.GetLogList(ctx, equipment.OrganizationID(), equipment.UID())
	})
}

func LogDelete(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = app.LogDBStore.Delete(equipment.OrganizationID(), equipment.UID())
		if err != nil {
			return err
		}

		logStr := lang.Str(lang.DeviceLogDeletedByUser, admin.Name(), equipment.Title())

		log2.WithField("src", logStore.SystemLog).Info(logStr)
		equipment.Logger().Infoln(logStr)

		return lang.Ok
	})
}
