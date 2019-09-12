package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
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
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())

		var params = []store.OptionFN{store.Page(page, pageSize), store.User(perm.AdminUser(ctx).GetID())}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, store.Keyword(keyword))
		}

		groupID := ctx.URLParamInt64Default("group", -1)
		if groupID != -1 {
			params = append(params, store.Group(groupID))
		}

		devices, total, err := s.GetDeviceList(params...)
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(devices))
		for _, device := range devices {
			result = append(result, device.Brief())
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context, s store.Store, validate *validator.Validate) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Title   string `json:"title" validate:"required"`
			ConnStr string `json:"params.connStr" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		device, err := s.CreateDevice(form.Title, map[string]interface{}{
			"params": map[string]interface{}{
				"connStr": form.ConnStr,
			},
		})
		if err != nil {
			return err
		}

		return device.Simple()
	})
}

func Detail(deviceID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, device, resource.View) {
			return lang.ErrNoPermission
		}

		return device.Detail()
	})
}

func Update(deviceID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, device, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		var form struct {
			Title   *string `json:"title"`
			ConnStr *string `json:"params.connStr"`
		}
		if err = ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.Title != nil {
			err = device.SetTitle(*form.Title)
			if err != nil {
				return err
			}
		}
		if form.ConnStr != nil {
			err = device.SetOption("params.connStr", form.ConnStr)
			if err != nil {
				return err
			}
		}
		err = device.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(deviceID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, device, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = device.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func MeasureList(deviceID int64, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")
		kind := ctx.URLParamIntDefault("kind", int(model.AllKind))

		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, device, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		measures, total, err := s.GetMeasureList(store.Device(device.GetID()), store.Keyword(keyword), store.Kind(model.MeasureKind(kind)), store.Page(page, pageSize))
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(measures))
		for _, measure := range measures {
			result = append(result, measure.Brief())
		}
		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func MeasureDetail(measureID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, measure, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		return measure.Detail()
	})
}
