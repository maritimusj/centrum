package organization

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/perm"
	"github.com/maritimusj/centrum/web/response"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		if !perm.IsDefaultAdminUser(ctx) {
			return lang.ErrNoPermission
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())

		var params = []helper.OptionFN{
			helper.Page(page, pageSize),
		}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		organizations, total, err := s.GetOrganizationList(params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(organizations))
		for _, org := range organizations {
			brief := org.Brief()
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context, s store.Store, validate *validator.Validate) hero.Result {
	return response.Wrap(func() interface{} {
		if !perm.IsDefaultAdminUser(ctx) {
			return lang.ErrNoPermission
		}

		var form struct {
			Name  string `json:"name" validate:"required"`
			Title string `json:"title"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if _, err := s.GetOrganization(form.Name); err != lang.Error(lang.ErrOrganizationNotFound) {
			return lang.ErrOrganizationExists
		}

		org, err := s.CreateOrganization(form.Name, form.Title)
		if err != nil {
			go perm.AdminUser(ctx).Logger().WithFields(logrus.Fields{
				"name":  form.Name,
				"title": form.Title,
			}).Info(lang.Str(lang.CreateOrgFail, form.Name, form.Title, err))
			return err
		} else {
			go perm.AdminUser(ctx).Logger().WithFields(logrus.Fields(org.Brief())).Info(lang.Str(lang.CreateOrgOk, org.Title(), org.Name()))
		}

		return org.Simple()
	})
}

func Detail(orgID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		if !perm.IsDefaultAdminUser(ctx) {
			return lang.ErrNoPermission
		}

		org, err := s.GetOrganization(orgID)
		if err != nil {
			return err
		}

		return org.Detail()
	})
}

func Update(orgID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		if !perm.IsDefaultAdminUser(ctx) {
			return lang.ErrNoPermission
		}

		org, err := s.GetOrganization(orgID)
		if err != nil {
			return err
		}

		var form struct {
			Title *string `json:"title"`
		}

		if err = ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		logFields := make(map[string]interface{})

		if form.Title != nil {
			org.SetTitle(*form.Title)
			logFields["title"] = form.Title
		}

		err = org.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(orgID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		if !perm.IsDefaultAdminUser(ctx) {
			return lang.ErrNoPermission
		}

		org, err := s.GetOrganization(orgID)
		if err != nil {
			return err
		}

		err = org.Destroy()
		if err != nil {
			return err
		}

		return lang.Ok
	})
}
