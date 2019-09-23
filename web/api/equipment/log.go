package equipment

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/api/web"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
)

func LogList(equipmentID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		return web.GetLogList(ctx, equipment.LogUID())
	})
}

func LogDelete(equipmentID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, equipment, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		return web.DeleteLog(ctx, equipment.LogUID())
	})
}

