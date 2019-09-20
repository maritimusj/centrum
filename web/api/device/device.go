package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/app"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/api/web"
	"github.com/maritimusj/centrum/web/response"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		var params []helper.OptionFN
		var orgID int64
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
		pageSize := ctx.URLParamInt64Default("pagesize", app.Cfg.DefaultPageSize())
		params = append(params, helper.Page(page, pageSize))

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		groupID := ctx.URLParamInt64Default("group", -1)
		if groupID != -1 {
			params = append(params, helper.Group(groupID))
		}

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.User(admin.GetID()))
			params = append(params, helper.DefaultEffect(app.Cfg.DefaultEffect()))
		}

		devices, total, err := s.GetDeviceList(params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(devices))
		for _, device := range devices {
			brief := device.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app.Allow(admin, device, resource.Ctrl),
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
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		var form struct {
			OrgID    int64  `json:"org"`
			Title    string `json:"title" validate:"required"`
			ConnStr  string `json:"params.connStr" validate:"required"`
			Interval int64  `json:"params.interval" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		var org interface{}
		if app.IsDefaultAdminUser(admin) {
			if form.OrgID > 0 {
				org = form.OrgID
			} else {
				org = app.Cfg.DefaultOrganization()
			}
		} else {
			org = admin.OrganizationID()
		}

		device, err := app.Store().CreateDevice(org, form.Title, map[string]interface{}{
			"params": map[string]interface{}{
				"connStr":  form.ConnStr,
				"interval": form.Interval,
			},
		})
		if err != nil {
			go admin.Logger().WithFields(logrus.Fields{
				"title":    form.Title,
				"connStr":  form.ConnStr,
				"interval": form.Interval,
			}).Info(lang.Str(lang.CreateDeviceFail, err))
			return err
		} else {
			go admin.Logger().WithFields(logrus.Fields(device.Brief())).Info(lang.Str(lang.CreateDeviceOk, device.Title()))
		}

		return device.Simple()
	})
}

func Detail(deviceID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		device, err := app.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		return device.Detail()
	})
}

func Update(deviceID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, device, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		var form struct {
			Title    *string  `json:"title"`
			ConnStr  *string  `json:"params.connStr"`
			Interval *int64   `json:"params.interval"`
			Groups   *[]int64 `json:"groups"`
		}

		if err = ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		logFields := make(map[string]interface{})

		if form.Title != nil {
			device.SetTitle(*form.Title)
			logFields["title"] = form.Title
		}

		if form.ConnStr != nil {
			err = device.SetOption("params.connStr", form.ConnStr)
			if err != nil {
				return err
			}
			logFields["connStr"] = form.ConnStr
		}

		if form.Interval != nil {
			err = device.SetOption("params.interval", form.Interval)
			if err != nil {
				return err
			}
			logFields["Interval"] = form.Interval
		}

		if form.Groups != nil && len(*form.Groups) > 0 {
			var groups []interface{}
			for _, g := range *form.Groups {
				groups = append(groups, g)
			}
			err = device.SetGroups(groups...)
			if err != nil {
				return err
			}
		}

		err = device.Save()
		if err != nil {
			go admin.Logger().WithFields(logFields).Info(lang.Str(lang.UpdateDeviceFail, device.Title(), err))
			return err
		} else {
			go admin.Logger().WithFields(logFields).Info(lang.Str(lang.UpdateDeviceOk, device.Title()))
		}

		return lang.Ok
	})
}

func Delete(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)

			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}

			if !app.Allow(admin, device, resource.Ctrl) {
				return lang.ErrNoPermission
			}

			err = device.Destroy()
			if err != nil {
				go admin.Logger().Info(lang.Str(lang.DeleteDeviceFail, device.Title(), err))
				return err
			} else {
				go admin.Logger().Info(lang.Str(lang.DeleteDeviceOk, device.Title()))
			}

			return lang.Ok
		})
	})
}

func MeasureList(deviceID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Cfg.DefaultPageSize())
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
			params = append(params, helper.DefaultEffect(app.Cfg.DefaultEffect()))
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

func CreateMeasure(deviceID int64, ctx iris.Context, validate *validator.Validate) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, device, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		var form struct {
			Title string `json:"title" validate:"required"`
			Tag   string `json:"tag" validate:"required"`
			Kind  int8   `json:"kind" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		measure, err := s.CreateMeasure(device.GetID(), form.Title, form.Tag, resource.MeasureKind(form.Kind))
		if err != nil {
			return err
		}

		return measure.Simple()
	})
}

func MeasureDetail(measureID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		measure, err := app.Store().GetMeasure(measureID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		return measure.Detail()
	})
}

func DeleteMeasure(measureID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		measure, err := app.Store().GetMeasure(measureID)
		if err != nil {
			return err
		}

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

func LogList(deviceID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		return web.GetLogList(ctx, device.LogUID())
	})
}

func LogDelete(deviceID int64, ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, device, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		return web.DeleteLog(ctx, device.LogUID())
	})
}
