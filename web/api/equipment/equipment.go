package equipment

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/logStore"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/api/web"
	"github.com/maritimusj/centrum/web/perm"
	"github.com/maritimusj/centrum/web/response"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())

		var params = []helper.OptionFN{
			helper.Page(page, pageSize),
		}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		groupID := ctx.URLParamInt64Default("group", -1)
		if groupID != -1 {
			params = append(params, helper.Group(groupID))
		}

		if !perm.IsDefaultAdminUser(ctx) {
			params = append(params, helper.DefaultEffect(cfg.DefaultEffect()))
			params = append(params, helper.User(perm.AdminUser(ctx).GetID()))
		}

		equipments, total, err := s.GetEquipmentList(params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(equipments))
		for _, equipment := range equipments {
			brief := equipment.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": perm.Allow(ctx, equipment, resource.Ctrl),
			}
			result = append(result, brief)
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
			Title string `json:"title" validate:"required"`
			Desc  string `json:"desc"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		equipment, err := s.CreateEquipment(form.Title, form.Desc)
		if err != nil {
			return err
		}

		return equipment.Simple()
	})
}

func Detail(deviceID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(deviceID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		return equipment.Detail()
	})
}

func Update(equipmentID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		var form struct {
			Title *string `json:"title"`
			Desc  *string `json:"desc"`
		}
		if err = ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}
		if form.Title != nil {
			equipment.SetTitle(*form.Title)
		}
		if form.Desc != nil {
			equipment.SetDesc(*form.Desc)
		}

		err = equipment.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(equipmentID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = equipment.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func StateList(equipmentID int64, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		kind := ctx.URLParamIntDefault("kind", int(resource.AllKind))

		var params = []helper.OptionFN{
			helper.Page(page, pageSize),
			helper.Kind(resource.MeasureKind(kind)),
			helper.Equipment(equipment.GetID()),
		}

		if !perm.IsDefaultAdminUser(ctx) {
			params = append(params, helper.DefaultEffect(cfg.DefaultEffect()))
			params = append(params, helper.User(perm.AdminUser(ctx).GetID()))
		}

		states, total, err := s.GetStateList(params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(states))
		for _, state := range states {
			brief := state.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": perm.Allow(ctx, state, resource.Ctrl),
			}
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func CreateState(equipmentID int64, ctx iris.Context, s store.Store, validate *validator.Validate) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		var form struct {
			Title           string `json:"title" validate:"required"`
			Desc            string `json:"desc"`
			MeasureID       int64  `json:"measure_id" validate:"required"`
			TransformScript string `json:"script"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}
		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		state, err := equipment.CreateState(form.Title, form.Desc, form.MeasureID, form.TransformScript)
		if err != nil {
			return err
		}
		return state.Simple()
	})
}

func StateDetail(stateID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, state, resource.View) {
			return lang.ErrNoPermission
		}

		return state.Detail()
	})
}

func UpdateState(stateID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, state, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		var form struct {
			Title           *string `json:"title"`
			MeasureID       *int64  `json:"measure_id"`
			TransformScript *string `json:"script"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.Title != nil {
			state.SetTitle(*form.Title)
		}

		if form.MeasureID != nil {
			state.SetMeasure(*form.MeasureID)
		}

		if form.TransformScript != nil {
			state.SetScript(*form.TransformScript)
		}

		err = state.Save()
		if err != nil {
			return err
		}

		return lang.Ok
	})
}

func DeleteState(stateID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, state, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = state.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func LogList(equipmentID int64, ctx iris.Context, s store.Store, store logStore.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		return web.GetLogList(equipment.LogUID(), ctx, store, cfg)
	})
}

func LogDelete(equipmentID int64, ctx iris.Context, s store.Store, store logStore.Store) hero.Result {
	return response.Wrap(func() interface{} {
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.Ctrl) {
			return lang.ErrNoPermission
		}
		return web.DeleteLog(equipment.LogUID(), ctx, store)
	})
}
