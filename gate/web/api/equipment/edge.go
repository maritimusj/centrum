package equipment

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/edge"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/global"

	edgeLang "github.com/maritimusj/centrum/edge/lang"
	_ "github.com/maritimusj/centrum/edge/lang/zhCN"
)

func rangeEquipmentStates(user model.User, equipment model.Equipment, fn func(device model.Device, measure model.Measure, state model.State) error) error {
	var params []helper.OptionFN
	if user != nil && !app.IsDefaultAdminUser(user) {
		params = append(params, helper.User(user.GetID()))
	}

	states, _, err := equipment.GetStateList(params...)
	if err != nil {
		return err
	}

	var (
		device  model.Device
		measure model.Measure
	)

	for _, state := range states {
		measure = state.Measure()
		if measure != nil {
			device = measure.Device()
		} else {
			device = nil
		}
		if err := fn(device, measure, state); err != nil {
			return err
		}
	}
	return nil
}

func getEquipmentSimpleStatus(user model.User, equipment model.Equipment) interface{} {
	res := map[string]interface{}{
		"index": edgeLang.Connected,
		"title": edgeLang.Str(edgeLang.Connected),
	}
	_ = rangeEquipmentStates(user, equipment, func(device model.Device, measure model.Measure, state model.State) error {
		if device == nil {
			res["index"] = edgeLang.MalFunctioned
			return lang.Error(lang.ErrDeviceNotFound)
		} else if measure == nil {
			res["index"] = edgeLang.MalFunctioned
			return lang.Error(lang.ErrMeasureNotFound)
		}
		index, title := global.GetDeviceStatus(device)
		if index != int(edgeLang.Connected) {
			res["index"] = index
			res["title"] = title
			return errors.New(edgeLang.Str(edgeLang.StrIndex(index)))
		}
		return nil
	})
	return res
}

func Status(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := app.Store().GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		if ctx.URLParamExists("simple") {
			return getEquipmentSimpleStatus(admin, equipment)
		}

		devices := make([]map[string]interface{}, 0)
		err = rangeEquipmentStates(admin, equipment, func(device model.Device, measure model.Measure, state model.State) error {
			dataMap := map[string]interface{}{
				"id":    state.GetID(),
				"title": state.Title(),
			}

			if device != nil {
				baseInfo, err := edge.GetStatus(device)
				if err != nil {
					index, title := global.GetDeviceStatus(device)
					if index != 0 {
						dataMap["edge"] = map[string]interface{}{
							"status": map[string]interface{}{
								"index": index,
								"title": title,
							},
						}
					} else {
						dataMap["error"] = err.Error()
					}
				} else {
					dataMap["edge"] = baseInfo
				}
			}

			if device == nil {
				dataMap["error"] = lang.Error(lang.ErrDeviceNotFound)
			} else if measure == nil {
				dataMap["error"] = lang.Error(lang.ErrMeasureNotFound)
			}

			devices = append(devices, dataMap)
			return nil
		})

		if err != nil {
			return err
		}

		return devices
	})
}

func Data(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := app.Store().GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		testCtrlPerm := func(state model.State) bool {
			if app.IsDefaultAdminUser(admin) {
				return true
			}
			return app.Allow(admin, state, resource.Ctrl)
		}

		result := make([]interface{}, 0)
		err = rangeEquipmentStates(admin, equipment, func(device model.Device, measure model.Measure, state model.State) error {
			dataMap := map[string]interface{}{
				"id":    state.GetID(),
				"title": state.Title(),
			}

			if device == nil {
				dataMap["error"] = lang.Error(lang.ErrDeviceNotFound)
			} else if measure == nil {
				dataMap["error"] = lang.Error(lang.ErrMeasureNotFound)
			}

			if device != nil && measure != nil {
				data, err := edge.GetCHValue(device, measure.TagName())
				if err != nil {
					dataMap["error"] = err.Error()
				} else {
					for k, v := range data {
						dataMap[k] = v
					}
				}

				dataMap["perm"] = map[string]bool{
					"view": true,
					"ctrl": testCtrlPerm(state),
				}
			}

			result = append(result, dataMap)
			return nil
		})

		if err != nil {
			return err
		}

		return result
	})
}

func Ctrl(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		//equipmentID, err := ctx.Params().GetInt64("id")
		//if err != nil {
		//	return err
		//}
		//equipment, err := app.Store().GetEquipment(equipmentID)
		//if err != nil {
		//	return err
		//}

		stateID, err := ctx.Params().GetInt64("stateID")
		if err != nil {
			return err
		}

		var form struct {
			Val bool `form:"value" json:"value"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		state, err := app.Store().GetState(stateID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, state, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		measure := state.Measure()
		if measure == nil {
			return lang.Error(lang.ErrMeasureNotFound)
		}

		device := measure.Device()
		if device == nil {
			return lang.Error(lang.ErrDeviceNotFound)
		}

		err = edge.SetCHValue(device, measure.TagName(), form.Val)
		if err != nil {
			return err
		}

		val, err := edge.GetCHValue(device, measure.TagName())
		if err != nil {
			return err
		}
		return val
	})
}

func GetCHValue(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		//equipmentID, err := ctx.Params().GetInt64("id")
		//if err != nil {
		//	return lang.ErrInvalidRequestData
		//}
		//
		//equipment, err := app.Store().GetEquipment(equipmentID)
		//if err != nil {
		//	return err
		//}

		stateID, err := ctx.Params().GetInt64("stateID")
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		state, err := app.Store().GetState(stateID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, state, resource.View) {
			return lang.ErrNoPermission
		}

		measure := state.Measure()
		if measure == nil {
			return lang.Error(lang.ErrMeasureNotFound)
		}

		device := measure.Device()
		if device == nil {
			return lang.Error(lang.ErrDeviceNotFound)
		}

		val, err := edge.GetCHValue(device, measure.TagName())
		if err != nil {
			return err
		}
		return val
	})
}
