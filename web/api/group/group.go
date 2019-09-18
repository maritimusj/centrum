package group

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/perm"
	"github.com/maritimusj/centrum/web/response"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		var params []helper.OptionFN
		var orgID int64
		if perm.IsDefaultAdminUser(ctx) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = perm.AdminUser(ctx).OrganizationID()
		}
		if orgID > 0 {
			params = append(params, helper.Organization(orgID))
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
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

		if !perm.IsDefaultAdminUser(ctx) {
			params = append(params, helper.User(perm.AdminUser(ctx).GetID()))
			params = append(params, helper.DefaultEffect(cfg.DefaultEffect()))
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
				"ctrl": perm.Allow(ctx, group, resource.Ctrl),
			}
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context, s store.Store, cfg config.Config, validate *validator.Validate) hero.Result {
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

		if form.ParentGroupID > 0 {
			_, err := s.GetGroup(form.ParentGroupID)
			if err != nil {
				return err
			}
		}

		var org interface{}
		if perm.IsDefaultAdminUser(ctx) {
			if form.OrgID > 0 {
				org = form.OrgID
			} else {
				org = cfg.DefaultOrganization()
			}
		} else {
			org = perm.AdminUser(ctx).OrganizationID()
		}

		group, err := s.CreateGroup(org, form.Title, form.Desc, form.ParentGroupID)
		if err != nil {
			return err
		}

		return group.Simple()
	})
}

func Detail(groupID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, group, resource.View) {
			return lang.ErrNoPermission
		}

		return group.Detail()
	})
}

func Update(groupID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, group, resource.Ctrl) {
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

				if perm.Deny(ctx, device, resource.Ctrl) {
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

				if perm.Deny(ctx, equipment, resource.Ctrl) {
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

func Delete(groupID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, group, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = group.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
