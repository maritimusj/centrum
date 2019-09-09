package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")

		devices, total, err := s.GetDeviceList(store.Keyword(keyword), store.Page(page, pageSize))
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
			ConnStr string `json:"conn_str" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		device, err := s.CreateDevice(form.Title, map[string]interface{}{
			"connStr": form.ConnStr,
		})
		if err != nil {
			return err
		}

		return device.Simple()
	})
}

func Detail(deviceID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
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
		var form struct {
			Title   *string `json:"title"`
			ConnStr *string `json:"conn_str"`
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
			err = device.SetOption(map[string]interface{}{
				"connStr": *form.ConnStr,
			})
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

func Delete(deviceID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
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

		measures, total, err := s.GetMeasureList(device.GetID(), store.Keyword(keyword), store.Kind(model.MeasureKind(kind)), store.Page(page, pageSize))
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

func MeasureDetail(measureID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}
		return measure.Detail()
	})
}
