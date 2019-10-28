package organization

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	"github.com/sirupsen/logrus"
)

func List(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		admin := s.MustGetUserFromContext(ctx)

		if !app2.IsDefaultAdminUser(admin) {
			return lang2.ErrNoPermission
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())

		var params = []helper2.OptionFN{
			helper2.Page(page, pageSize),
		}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper2.Keyword(keyword))
		}

		organizations, total, err := s.GetOrganizationList(params...)
		if err != nil {
			return err
		}

		var result = make([]model2.Map, 0, len(organizations))
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
	return response2.Wrap(func() interface{} {
		return app2.TransactionDo(func(s store2.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if !app2.IsDefaultAdminUser(admin) {
				return lang2.ErrNoPermission
			}

			var form struct {
				Name  string `json:"name" valid:"required"`
				Title string `json:"title"`
			}

			if err := ctx.ReadJSON(&form); err != nil {
				return lang2.ErrInvalidRequestData
			}

			if _, err := govalidator.ValidateStruct(&form); err != nil {
				return lang2.ErrInvalidRequestData
			}

			if exists, err := s.IsOrganizationExists(form.Name); err != nil {
				return err
			} else if exists {
				return lang2.ErrOrganizationExists
			}

			org, err := s.CreateOrganization(form.Name, form.Title)
			if err != nil {
				go admin.Logger().WithFields(logrus.Fields{
					"name":  form.Name,
					"title": form.Title,
				}).Info(lang2.Str(lang2.CreateOrgFail, form.Name, form.Title, err))
				return err
			} else {
				go admin.Logger().WithFields(logrus.Fields(org.Brief())).Info(lang2.Str(lang2.CreateOrgOk, org.Title(), org.Name()))
			}

			return org.Simple()
		})
	})
}

func Detail(orgID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()

		admin := s.MustGetUserFromContext(ctx)
		if !app2.IsDefaultAdminUser(admin) {
			return lang2.ErrNoPermission
		}

		org, err := s.GetOrganization(orgID)
		if err != nil {
			return err
		}

		return org.Detail()
	})
}

func Update(orgID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		return app2.TransactionDo(func(s store2.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if !app2.IsDefaultAdminUser(admin) {
				return lang2.ErrNoPermission
			}

			org, err := s.GetOrganization(orgID)
			if err != nil {
				return err
			}

			var form struct {
				Title *string `json:"title"`
			}

			if err = ctx.ReadJSON(&form); err != nil {
				return lang2.ErrInvalidRequestData
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
			return lang2.Ok
		})
	})
}

func Delete(orgID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		return app2.TransactionDo(func(s store2.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if !app2.IsDefaultAdminUser(admin) {
				return lang2.ErrNoPermission
			}

			org, err := s.GetOrganization(orgID)
			if err != nil {
				return err
			}

			if org.Name() == app2.Config.DefaultOrganization() {
				return lang2.ErrFailedRemoveDefaultOrganization
			}

			err = org.Destroy()
			if err != nil {
				return err
			}

			return lang2.Ok
		})
	})
}
