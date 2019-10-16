package equipment

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/edge"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"

	edgeLang "github.com/maritimusj/centrum/edge/lang"
)

func RangeEquipmentStates(equipment model.Equipment, fn func(device model.Device, measure model.Measure, state model.State) error) error {
	states, _, err := equipment.GetStateList()
	if err != nil {
		return err
	}
	var device model.Device
	var measure model.Measure

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
			res := map[string]interface{}{
				"index": edgeLang.Connected,
				"title": edgeLang.Str(edgeLang.Connected),
			}
			_ = RangeEquipmentStates(equipment, func(device model.Device, measure model.Measure, state model.State) error {
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

		devices := make([]map[string]interface{}, 0)
		err = RangeEquipmentStates(equipment, func(device model.Device, measure model.Measure, state model.State) error {
			dataMap := map[string]interface{}{
				"id":    state.GetID(),
				"title": state.Title(),
			}

			if device != nil {
				baseInfo, err := edge.GetStatus(device)
				if err != nil {
					dataMap["error"] = err.Error()
				} else {
					dataMap["device"] = baseInfo
					devices = append(devices, dataMap)
				}
			}

			if device == nil {
				dataMap["error"] = lang.Error(lang.ErrDeviceNotFound)
			} else if measure == nil {
				dataMap["error"] = lang.Error(lang.ErrMeasureNotFound)
			}

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

		states, _, err := equipment.GetStateList()
		if err != nil {
			return err
		}

		devices := make([]interface{}, 0, len(states))
		for _, state := range states {
			dataMap := map[string]interface{}{
				"id":    state.GetID(),
				"title": state.Title(),
			}
			measure := state.Measure()
			if measure == nil {
				dataMap["error"] = lang.Error(lang.ErrMeasureNotFound)
				continue
			}

			device := measure.Device()
			if device == nil {
				dataMap["error"] = lang.Error(lang.ErrDeviceNotFound)
				continue
			}

			data, err := edge.GetCHValue(device, measure.TagName())
			if err != nil {
				dataMap["error"] = err.Error()
				continue
			}
			dataMap["data"] = data
			devices = append(devices, dataMap)
		}

		return devices
	})
}

func Ctrl(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		equipmentID, err := ctx.URLParamInt64("id")
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		stateID, err := ctx.URLParamInt64("stateID")
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		equipment, err := app.Store().GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.Ctrl) {
			return lang.ErrNoPermission
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
		equipmentID, err := ctx.URLParamInt64("id")
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		stateID, err := ctx.URLParamInt64("stateID")
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		equipment, err := app.Store().GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		state, err := app.Store().GetState(stateID)
		if err != nil {
			return err
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
