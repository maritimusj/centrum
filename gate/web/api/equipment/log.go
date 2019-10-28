package equipment

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	log2 "github.com/maritimusj/centrum/gate/web/api/log"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
)

func LogList(equipmentID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, equipment, resource2.View) {
			return lang2.ErrNoPermission
		}

		return log2.GetLogList(ctx, equipment.OrganizationID(), equipment.UID())
	})
}

func LogDelete(equipmentID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, equipment, resource2.Ctrl) {
			return lang2.ErrNoPermission
		}

		return log2.DeleteLog(ctx, equipment.OrganizationID(), equipment.UID())
	})
}
