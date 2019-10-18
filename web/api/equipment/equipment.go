package equipment

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/event"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
	"github.com/maritimusj/centrum/web/store"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()

		var params []helper.OptionFN
		var orgID int64

		admin := s.MustGetUserFromContext(ctx)
		if app.IsDefaultAdminUser(admin) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = admin.OrganizationID()
		}
		if orgID > 0 {
			params = append(params, helper.Organization(orgID))
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
		params = append(params, helper.Page(page, pageSize))

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		if ctx.URLParamExists("group") {
			groupID := ctx.URLParamInt64Default("group", 0)
			params = append(params, helper.Group(groupID))
		}

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
			params = append(params, helper.User(admin.GetID()))
		}

		equipments, total, err := s.GetEquipmentList(params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(equipments))
		for _, equipment := range equipments {
			brief := equipment.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app.Allow(admin, equipment, resource.Ctrl),
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

func Create(ctx iris.Context, validate *validator.Validate) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			OrgID  int64   `json:"org"`
			Title  string  `json:"title" validate:"required"`
			Desc   string  `json:"desc"`
			Groups []int64 `json:"groups"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		result := app.TransactionDo(func(s store.Store) interface{} {
			var org interface{}

			admin := s.MustGetUserFromContext(ctx)
			if app.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app.Config.DefaultOrganization()
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

			err = app.SetAllow(admin, equipment, resource.View, resource.Ctrl)
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
			app.Event.Publish(event.EquipmentCreated, data.Get("userID"), data.Get("equipmentID"))
			return data.Pop("result")
		}

		return result
	})
}

func Detail(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		equipment, err := s.GetEquipment(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		return equipment.Detail()
	})
}

func Update(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Title  *string  `json:"title"`
			Desc   *string  `json:"desc"`
			Groups *[]int64 `json:"groups"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		result := app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}

			if !app.Allow(admin, equipment, resource.Ctrl) {
				return lang.ErrNoPermission
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
			app.Event.Publish(event.EquipmentUpdated, data.Get("userID"), data.Get("equipmentID"))
			return lang.Ok
		}

		return result
	})
}

func Delete(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		result := app.TransactionDo(func(s store.Store) interface{} {
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app.Allow(admin, equipment, resource.Ctrl) {
				return lang.ErrNoPermission
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
			app.Event.Publish(event.EquipmentDeleted, data.Get("userID"), data.Get("title"))
			return lang.Ok
		}
		return result
	})
}
