package organization

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/gate/web/store"
	"github.com/sirupsen/logrus"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var (
			s     = app.Store()
			admin = s.MustGetUserFromContext(ctx)
		)

		if !app.IsDefaultAdminUser(admin) {
			return lang.ErrNoPermission
		}

		var (
			page     = ctx.URLParamInt64Default("page", 1)
			pageSize = ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())

			params = []helper.OptionFN{
				helper.Page(page, pageSize),
			}
		)

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

func Create(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if !app.IsDefaultAdminUser(admin) {
				return lang.ErrNoPermission
			}

			var form struct {
				Name  string `json:"name" valid:"required"`
				Title string `json:"title"`
			}

			if err := ctx.ReadJSON(&form); err != nil {
				return lang.ErrInvalidRequestData
			}

			if _, err := govalidator.ValidateStruct(&form); err != nil {
				return lang.ErrInvalidRequestData
			}

			if exists, err := s.IsOrganizationExists(form.Name); err != nil {
				return err
			} else if exists {
				return lang.ErrOrganizationExists
			}

			org, err := s.CreateOrganization(form.Name, form.Title)
			if err != nil {
				go admin.Logger().WithFields(logrus.Fields{
					"name":  form.Name,
					"title": form.Title,
				}).Info(lang.Str(lang.CreateOrgFail, form.Name, form.Title, err))
				return err
			} else {
				go admin.Logger().WithFields(logrus.Fields(org.Brief())).Info(lang.Str(lang.CreateOrgOk, org.Title(), org.Name()))
			}

			return org.Simple()
		})
	})
}

func Detail(orgID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()

		admin := s.MustGetUserFromContext(ctx)
		if !app.IsDefaultAdminUser(admin) {
			return lang.ErrNoPermission
		}

		org, err := s.GetOrganization(orgID)
		if err != nil {
			return err
		}

		return org.Detail()
	})
}

func Update(orgID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if !app.IsDefaultAdminUser(admin) {
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
	})
}

func Delete(orgID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if !app.IsDefaultAdminUser(admin) {
				return lang.ErrNoPermission
			}

			org, err := s.GetOrganization(orgID)
			if err != nil {
				return err
			}

			if org.Name() == app.Config.DefaultOrganization() {
				return lang.ErrFailedRemoveDefaultOrganization
			}

			err = org.Destroy()
			if err != nil {
				return err
			}

			return lang.Ok
		})
	})
}
