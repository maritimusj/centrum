package equipment

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	edge2 "github.com/maritimusj/centrum/gate/web/edge"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/global"

	edgeLang "github.com/maritimusj/centrum/edge/lang"
	_ "github.com/maritimusj/centrum/edge/lang/zhCN"
)

func rangeEquipmentStates(user model2.User, equipment model2.Equipment, fn func(device model2.Device, measure model2.Measure, state model2.State) error) error {
	var params []helper2.OptionFN
	if user != nil && !app2.IsDefaultAdminUser(user) {
		params = append(params, helper2.User(user.GetID()))
	}

	states, _, err := equipment.GetStateList(params...)
	if err != nil {
		return err
	}

	var (
		device  model2.Device
		measure model2.Measure
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

func getEquipmentSimpleStatus(user model2.User, equipment model2.Equipment) interface{} {
	res := map[string]interface{}{
		"index": edgeLang.Connected,
		"title": edgeLang.Str(edgeLang.Connected),
	}
	_ = rangeEquipmentStates(user, equipment, func(device model2.Device, measure model2.Measure, state model2.State) error {
		if device == nil {
			res["index"] = edgeLang.MalFunctioned
			return lang2.Error(lang2.ErrDeviceNotFound)
		} else if measure == nil {
			res["index"] = edgeLang.MalFunctioned
			return lang2.Error(lang2.ErrMeasureNotFound)
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
	return response2.Wrap(func() interface{} {
		equipment, err := app2.Store().GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, equipment, resource2.View) {
			return lang2.ErrNoPermission
		}

		if ctx.URLParamExists("simple") {
			return getEquipmentSimpleStatus(admin, equipment)
		}

		devices := make([]map[string]interface{}, 0)
		err = rangeEquipmentStates(admin, equipment, func(device model2.Device, measure model2.Measure, state model2.State) error {
			dataMap := map[string]interface{}{
				"id":    state.GetID(),
				"title": state.Title(),
			}

			if device != nil {
				baseInfo, err := edge2.GetStatus(device)
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
				dataMap["error"] = lang2.Error(lang2.ErrDeviceNotFound)
			} else if measure == nil {
				dataMap["error"] = lang2.Error(lang2.ErrMeasureNotFound)
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
	return response2.Wrap(func() interface{} {
		equipment, err := app2.Store().GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, equipment, resource2.View) {
			return lang2.ErrNoPermission
		}

		testCtrlPerm := func(state model2.State) bool {
			if app2.IsDefaultAdminUser(admin) {
				return true
			}
			return app2.Allow(admin, state, resource2.Ctrl)
		}

		result := make([]interface{}, 0)
		err = rangeEquipmentStates(admin, equipment, func(device model2.Device, measure model2.Measure, state model2.State) error {
			dataMap := map[string]interface{}{
				"id":    state.GetID(),
				"title": state.Title(),
			}

			if device == nil {
				dataMap["error"] = lang2.Error(lang2.ErrDeviceNotFound)
			} else if measure == nil {
				dataMap["error"] = lang2.Error(lang2.ErrMeasureNotFound)
			}

			if device != nil && measure != nil {
				data, err := edge2.GetCHValue(device, measure.TagName())
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
	return response2.Wrap(func() interface{} {
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
			return lang2.ErrInvalidRequestData
		}

		state, err := app2.Store().GetState(stateID)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, state, resource2.Ctrl) {
			return lang2.ErrNoPermission
		}

		measure := state.Measure()
		if measure == nil {
			return lang2.Error(lang2.ErrMeasureNotFound)
		}

		device := measure.Device()
		if device == nil {
			return lang2.Error(lang2.ErrDeviceNotFound)
		}

		err = edge2.SetCHValue(device, measure.TagName(), form.Val)
		if err != nil {
			return err
		}

		val, err := edge2.GetCHValue(device, measure.TagName())
		if err != nil {
			return err
		}
		return val
	})
}

func GetCHValue(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
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
			return lang2.ErrInvalidRequestData
		}

		state, err := app2.Store().GetState(stateID)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, state, resource2.View) {
			return lang2.ErrNoPermission
		}

		measure := state.Measure()
		if measure == nil {
			return lang2.Error(lang2.ErrMeasureNotFound)
		}

		device := measure.Device()
		if device == nil {
			return lang2.Error(lang2.ErrDeviceNotFound)
		}

		val, err := edge2.GetCHValue(device, measure.TagName())
		if err != nil {
			return err
		}
		return val
	})
}
