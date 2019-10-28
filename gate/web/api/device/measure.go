package device

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
	"github.com/maritimusj/centrum/web/store"
)

func MeasureList(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()

		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
		kind := ctx.URLParamIntDefault("kind", int(resource.AllKind))

		var params = []helper.OptionFN{
			helper.Page(page, pageSize),
			helper.Kind(resource.MeasureKind(kind)),
			helper.Device(device.GetID()),
		}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.User(admin.GetID()))
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
		}

		measures, total, err := s.GetMeasureList(params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(measures))
		for _, measure := range measures {
			brief := measure.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app.Allow(admin, measure, resource.Ctrl),
			}
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func CreateMeasure(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Title string `json:"title" valid:"required"`
			Tag   string `json:"tag" valid:"required"`
			Kind  int8   `json:"kind" valid:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if _, err := govalidator.ValidateStruct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		return app.TransactionDo(func(s store.Store) interface{} {
			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app.Allow(admin, device, resource.Ctrl) {
				return lang.ErrNoPermission
			}

			measure, err := s.CreateMeasure(device.GetID(), form.Title, form.Tag, resource.MeasureKind(form.Kind))
			if err != nil {
				return err
			}

			err = app.SetAllow(admin, measure, resource.View, resource.Ctrl)
			if err != nil {
				return err
			}

			return measure.Simple()
		})
	})
}

func MeasureDetail(measureID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		return measure.Detail()
	})
}

func DeleteMeasure(measureID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		err = measure.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
