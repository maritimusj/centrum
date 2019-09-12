package equipment

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

		var params = []store.OptionFN{store.Page(page, pageSize)}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, store.Keyword(keyword))
		}

		groupID := ctx.URLParamInt64Default("group", -1)
		if groupID != -1 {
			params = append(params, store.Group(groupID))
		}

		equipments, total, err := s.GetEquipmentList(params...)
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(equipments))
		for _, equipment := range equipments {
			result = append(result, equipment.Brief())
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
			err = equipment.SetTitle(*form.Title)
			if err != nil {
				return err
			}
		}
		if form.Desc != nil {
			err = equipment.SetDesc(*form.Desc)
			if err != nil {
				return err
			}
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
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")
		kind := ctx.URLParamIntDefault("kind", int(model.AllKind))

		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if perm.Deny(ctx, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		states, total, err := s.GetStateList(store.Equipment(equipment.GetID()), store.Keyword(keyword), store.Kind(model.MeasureKind(kind)), store.Page(page, pageSize))
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(states))
		for _, state := range states {
			result = append(result, state.Brief())
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
			MeasureID       int64  `json:"measure_id" validate:"required"`
			TransformScript string `json:"script"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}
		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		state, err := equipment.CreateState(form.Title, form.MeasureID, form.TransformScript)
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
			err = state.SetTitle(*form.Title)
			if err != nil {
				return err
			}
		}
		if form.MeasureID != nil {
			err = state.SetMeasure(*form.MeasureID)
			if err != nil {
				return err
			}
		}
		if form.TransformScript != nil {
			err = state.SetScript(*form.TransformScript)
			if err != nil {
				return err
			}
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
