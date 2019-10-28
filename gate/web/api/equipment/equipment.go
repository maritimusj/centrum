package equipment

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/event"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	store2 "github.com/maritimusj/centrum/gate/web/store"
)

func List(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()

		var params []helper2.OptionFN
		var orgID int64

		admin := s.MustGetUserFromContext(ctx)
		if app2.IsDefaultAdminUser(admin) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = admin.OrganizationID()
		}
		if orgID > 0 {
			params = append(params, helper2.Organization(orgID))
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())
		params = append(params, helper2.Page(page, pageSize))

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper2.Keyword(keyword))
		}

		if ctx.URLParamExists("group") {
			groupID := ctx.URLParamInt64Default("group", 0)
			params = append(params, helper2.Group(groupID))
		}

		if !app2.IsDefaultAdminUser(admin) {
			params = append(params, helper2.DefaultEffect(app2.Config.DefaultEffect()))
			params = append(params, helper2.User(admin.GetID()))
		}

		equipments, total, err := s.GetEquipmentList(params...)
		if err != nil {
			return err
		}

		var (
			result = make([]model2.Map, 0, len(equipments))
		)

		for _, equipment := range equipments {
			brief := equipment.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app2.Allow(admin, equipment, resource2.Ctrl),
			}
			brief["edge"] = iris.Map{
				"status": getEquipmentSimpleStatus(admin, equipment),
			}
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			OrgID  int64   `json:"org"`
			Title  string  `json:"title" valid:"required"`
			Desc   string  `json:"desc"`
			Groups []int64 `json:"groups"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		if _, err := govalidator.ValidateStruct(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		result := app2.TransactionDo(func(s store2.Store) interface{} {
			var org interface{}

			admin := s.MustGetUserFromContext(ctx)
			if app2.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app2.Config.DefaultOrganization()
				}
			} else {
				org = admin.OrganizationID()
			}

			equipment, err := s.CreateEquipment(org, form.Title, form.Desc)
			if err != nil {
				return err
			}

			if len(form.Groups) > 0 {
				var groups []interface{}
				for _, g := range form.Groups {
					groups = append(groups, g)
				}
				err = equipment.SetGroups(groups...)
				if err != nil {
					return err
				}
			}

			err = app2.SetAllow(admin, equipment, resource2.View, resource2.Ctrl)
			if err != nil {
				return err
			}

			data := event.Data{
				"equipmentID": equipment.GetID(),
				"userID":      admin.GetID(),
				"result":      equipment.Simple(),
			}
			return data
		})

		if data, ok := result.(event.Data); ok {
			app2.Event.Publish(event.EquipmentCreated, data.Get("userID"), data.Get("equipmentID"))
			return data.Pop("result")
		}

		return result
	})
}

func Detail(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		equipment, err := s.GetEquipment(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, equipment, resource2.View) {
			return lang2.ErrNoPermission
		}

		return equipment.Detail()
	})
}

func Update(equipmentID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Title  *string  `json:"title"`
			Desc   *string  `json:"desc"`
			Groups *[]int64 `json:"groups"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		result := app2.TransactionDo(func(s store2.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}

			if !app2.Allow(admin, equipment, resource2.Ctrl) {
				return lang2.ErrNoPermission
			}

			if form.Title != nil {
				equipment.SetTitle(*form.Title)
			}

			if form.Desc != nil {
				equipment.SetDesc(*form.Desc)
			}

			if form.Groups != nil && len(*form.Groups) > 0 {
				var groups []interface{}
				for _, g := range *form.Groups {
					groups = append(groups, g)
				}
				err = equipment.SetGroups(groups...)
				if err != nil {
					return err
				}
			}

			err = equipment.Save()
			if err != nil {
				return err
			}

			data := event.Data{
				"equipmentID": equipment.GetID(),
				"userID":      admin.GetID(),
			}

			return data
		})

		if data, ok := result.(*event.Data); ok {
			app2.Event.Publish(event.EquipmentUpdated, data.Get("userID"), data.Get("equipmentID"))
			return lang2.Ok
		}

		return result
	})
}

func Delete(equipmentID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		result := app2.TransactionDo(func(s store2.Store) interface{} {
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app2.Allow(admin, equipment, resource2.Ctrl) {
				return lang2.ErrNoPermission
			}

			data := event.Data{
				"title":  equipment.Title(),
				"userID": admin.GetID(),
			}

			err = equipment.Destroy()
			if err != nil {
				return err
			}
			return data
		})
		if data, ok := result.(event.Data); ok {
			app2.Event.Publish(event.EquipmentDeleted, data.Get("userID"), data.Get("title"))
			return lang2.Ok
		}
		return result
	})
}
