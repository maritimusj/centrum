package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/edge"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
)

func Reset(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := app.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}
		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		edge.ResetConfig(device)
		return lang.Ok
	})
}

func Status(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := app.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		if ctx.URLParamExists("simple") {
			index, title := global.GetDeviceStatus(device)
			return iris.Map{
				"index": index,
				"title": title,
			}
		}

		baseInfo, err := edge.GetStatus(device)
		if err != nil {
			return err
		}
		return baseInfo
	})
}

func Data(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := app.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		data, err := edge.GetData(device)
		if err != nil {
			return err
		}
		return data
	})
}

func Ctrl(deviceID int64, chTagName string, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Val bool `form:"value" json:"value"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		device, err := app.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), chTagName)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, measure, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = edge.SetCHValue(device, chTagName, form.Val)
		if err != nil {
			return err
		}

		val, err := edge.GetCHValue(device, chTagName)
		if err != nil {
			return err
		}

		return val
	})
}

func GetCHValue(deviceID int64, chTagName string, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := app.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), chTagName)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		val, err := edge.GetCHValue(device, chTagName)
		if err != nil {
			return err
		}
		return val
	})
}
