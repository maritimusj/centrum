package group

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
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

		parentGroupID := ctx.URLParamInt64Default("parent", 0)
		if parentGroupID > 0 {
			_, err := s.GetGroup(parentGroupID)
			if err != nil {
				return err
			}
			params = append(params, helper.Parent(parentGroupID))
		}

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.User(admin.GetID()))
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
		}

		groups, total, err := s.GetGroupList(params...)
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(groups))
		for _, group := range groups {
			brief := group.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app.Allow(admin, group, resource.Ctrl),
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
			OrgID         int64  `json:"org"`
			Title         string `json:"title" validate:"required"`
			Desc          string `json:"desc"`
			ParentGroupID int64  `json:"parent" validate:"min=0"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		return app.TransactionDo(func(s store.Store) interface{} {
			if form.ParentGroupID > 0 {
				_, err := s.GetGroup(form.ParentGroupID)
				if err != nil {
					return err
				}
			}

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

			group, err := s.CreateGroup(org, form.Title, form.Desc, form.ParentGroupID)
			if err != nil {
				return err
			}

			err = app.SetAllow(admin, group, resource.View, resource.Ctrl)
			if err != nil {
				return err
			}

			return group.Simple()
		})
	})
}

func Detail(groupID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, group, resource.View) {
			return lang.ErrNoPermission
		}

		return group.Detail()
	})
}

func Update(groupID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, group, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		var form struct {
			ParentGroupID *int64  `json:"parent"`
			Title         *string `json:"title"`
			Desc          *string `json:"desc"`
			Devices       []int64 `json:"devices"`
			Equipments    []int64 `json:"equipments"`
		}

		err = ctx.ReadJSON(&form)
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.ParentGroupID != nil {
			group.SetParent(*form.ParentGroupID)
		}

		if form.Title != nil && *form.Title != "" {
			group.SetTitle(*form.Title)
		}

		if form.Desc != nil {
			group.SetDesc(*form.Desc)
		}

		if len(form.Devices) > 0 {
			devices := make([]interface{}, 0, len(form.Devices))
			for _, deviceID := range form.Devices {
				device, err := s.GetDevice(deviceID)
				if err != nil {
					return err
				}

				if !app.Allow(admin, device, resource.Ctrl) {
					return lang.ErrNoPermission
				}

				if device.OrganizationID() != group.OrganizationID() {
					return lang.ErrDeviceOrganizationDifferent
				}

				devices = append(devices, device)
			}

			err = group.AddDevice(devices...)
			if err != nil {
				return err
			}
		}
		if len(form.Equipments) > 0 {
			equipments := make([]interface{}, 0, len(form.Equipments))
			for _, equipmentID := range form.Equipments {
				equipment, err := s.GetEquipment(equipmentID)
				if err != nil {
					return err
				}

				if !app.Allow(admin, equipment, resource.Ctrl) {
					return lang.ErrNoPermission
				}

				if equipment.OrganizationID() != group.OrganizationID() {
					return lang.ErrEquipmentOrganizationDifferent
				}

				equipments = append(equipments, equipment)
			}
			err = group.AddEquipment(equipments...)
			if err != nil {
				return err
			}
		}
		err = group.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(groupID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, group, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = group.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
