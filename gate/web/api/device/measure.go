package device

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	store2 "github.com/maritimusj/centrum/gate/web/store"
)

func MeasureList(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()

		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, device, resource2.View) {
			return lang2.ErrNoPermission
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())
		kind := ctx.URLParamIntDefault("kind", int(resource2.AllKind))

		var params = []helper2.OptionFN{
			helper2.Page(page, pageSize),
			helper2.Kind(resource2.MeasureKind(kind)),
			helper2.Device(device.GetID()),
		}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper2.Keyword(keyword))
		}

		if !app2.IsDefaultAdminUser(admin) {
			params = append(params, helper2.User(admin.GetID()))
			params = append(params, helper2.DefaultEffect(app2.Config.DefaultEffect()))
		}

		measures, total, err := s.GetMeasureList(params...)
		if err != nil {
			return err
		}

		var result = make([]model2.Map, 0, len(measures))
		for _, measure := range measures {
			brief := measure.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app2.Allow(admin, measure, resource2.Ctrl),
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
	return response2.Wrap(func() interface{} {
		var form struct {
			Title string `json:"title" valid:"required"`
			Tag   string `json:"tag" valid:"required"`
			Kind  int8   `json:"kind" valid:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		if _, err := govalidator.ValidateStruct(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		return app2.TransactionDo(func(s store2.Store) interface{} {
			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app2.Allow(admin, device, resource2.Ctrl) {
				return lang2.ErrNoPermission
			}

			measure, err := s.CreateMeasure(device.GetID(), form.Title, form.Tag, resource2.MeasureKind(form.Kind))
			if err != nil {
				return err
			}

			err = app2.SetAllow(admin, measure, resource2.View, resource2.Ctrl)
			if err != nil {
				return err
			}

			return measure.Simple()
		})
	})
}

func MeasureDetail(measureID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, measure, resource2.View) {
			return lang2.ErrNoPermission
		}

		return measure.Detail()
	})
}

func DeleteMeasure(measureID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, measure, resource2.View) {
			return lang2.ErrNoPermission
		}

		err = measure.Destroy()
		if err != nil {
			return err
		}
		return lang2.Ok
	})
}
