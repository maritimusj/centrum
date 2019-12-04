package device

import (
	"net"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/edge"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/global"
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
		s := app.Store()
		device, err := s.GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		data, err := edge.GetRealTimeData(device)
		if err != nil {
			if netErr, ok := err.(net.Error); ok {
				return lang.Error(lang.ErrNetworkFail, netErr.Error())
			}
			return err
		}

		testPerm := func(measure model.Measure, action resource.Action) bool {
			if app.IsDefaultAdminUser(admin) {
				return true
			}
			return app.Allow(admin, measure, action)
		}

		//过滤掉没有权限的measure
		var result = make([]interface{}, 0, len(data))
		for _, entry := range data {
			if e, ok := entry.(map[string]interface{}); ok {
				if chTagName, ok := e["tag"].(string); ok {
					measure, err := s.GetMeasureFromTagName(deviceID, chTagName)
					if err != nil {
						continue
					}
					if testPerm(measure, resource.View) {
						e["perm"] = map[string]bool{
							"view": true,
							"ctrl": testPerm(measure, resource.Ctrl),
						}
						result = append(result, entry)
					}
				}
			}
		}

		return result
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
