package device

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	edge2 "github.com/maritimusj/centrum/gate/web/edge"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/global"
	"net"
)

func Reset(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		device, err := app2.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}
		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, device, resource2.View) {
			return lang2.ErrNoPermission
		}

		edge2.ResetConfig(device)
		return lang2.Ok
	})
}

func Status(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		device, err := app2.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, device, resource2.View) {
			return lang2.ErrNoPermission
		}

		if ctx.URLParamExists("simple") {
			index, title := global.GetDeviceStatus(device)
			return iris.Map{
				"index": index,
				"title": title,
			}
		}

		baseInfo, err := edge2.GetStatus(device)
		if err != nil {
			return err
		}
		return baseInfo
	})
}

func Data(deviceID int64, ctx iris.Context) hero.Result {
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

		data, err := edge2.GetData(device)
		if err != nil {
			if netErr, ok := err.(net.Error); ok {
				return lang2.Error(lang2.ErrNetworkFail, netErr.Error())
			}
			return err
		}

		testPerm := func(measure model2.Measure, action resource2.Action) bool {
			if app2.IsDefaultAdminUser(admin) {
				return true
			}
			return app2.Allow(admin, measure, action)
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
					if testPerm(measure, resource2.View) {
						e["perm"] = map[string]bool{
							"view": true,
							"ctrl": testPerm(measure, resource2.Ctrl),
						}
						result = append(result, entry)
					}
				}
			}
		}

		return data
	})
}

func Ctrl(deviceID int64, chTagName string, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Val bool `form:"value" json:"value"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		device, err := app2.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		measure, err := app2.Store().GetMeasureFromTagName(device.GetID(), chTagName)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, measure, resource2.Ctrl) {
			return lang2.ErrNoPermission
		}

		err = edge2.SetCHValue(device, chTagName, form.Val)
		if err != nil {
			return err
		}

		val, err := edge2.GetCHValue(device, chTagName)
		if err != nil {
			return err
		}

		return val
	})
}

func GetCHValue(deviceID int64, chTagName string, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		device, err := app2.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		measure, err := app2.Store().GetMeasureFromTagName(device.GetID(), chTagName)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, measure, resource2.View) {
			return lang2.ErrNoPermission
		}

		val, err := edge2.GetCHValue(device, chTagName)
		if err != nil {
			return err
		}
		return val
	})
}
