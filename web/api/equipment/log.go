package equipment

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/api/log"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
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

		return log.GetLogList(ctx, equipment.OrganizationID(), equipment.LogUID())
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

		return log.DeleteLog(ctx, equipment.OrganizationID(), equipment.LogUID())
	})
}
