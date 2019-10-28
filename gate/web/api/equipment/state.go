package equipment

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/gate/web/store"
)

func StateList(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, equipment, resource.View) {
			return lang.ErrNoPermission
		}

		var (
			page     = ctx.URLParamInt64Default("page", 1)
			pageSize = ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
			kind     = ctx.URLParamIntDefault("kind", int(resource.AllKind))

			params = []helper.OptionFN{
				helper.Page(page, pageSize),
				helper.Kind(resource.MeasureKind(kind)),
				helper.Equipment(equipment.GetID()),
			}
		)

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
			params = append(params, helper.User(admin.GetID()))
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
				"ctrl": app.Allow(admin, state, resource.Ctrl),
			}
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func CreateState(equipmentID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Title           string `json:"title" valid:"required"`
			Desc            string `json:"desc"`
			MeasureID       int64  `json:"measure_id" valid:"required"`
			TransformScript string `json:"script"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}
		if _, err := govalidator.ValidateStruct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		return app.TransactionDo(func(s store.Store) interface{} {
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app.Allow(admin, equipment, resource.Ctrl) {
				return lang.ErrNoPermission
			}

			state, err := equipment.CreateState(form.Title, form.Desc, form.MeasureID, form.TransformScript)
			if err != nil {
				return err
			}
			err = app.SetAllow(admin, state, resource.View, resource.Ctrl)
			if err != nil {
				return err
			}

			return state.Simple()
		})
	})
}

func StateDetail(stateID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, state, resource.View) {
			return lang.ErrNoPermission
		}

		return state.Detail()
	})
}

func UpdateState(stateID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, state, resource.Ctrl) {
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

func DeleteState(stateID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, state, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = state.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
