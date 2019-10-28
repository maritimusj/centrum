package equipment

import (
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	store2 "github.com/maritimusj/centrum/gate/web/store"
)

func StateList(equipmentID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		equipment, err := s.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, equipment, resource2.View) {
			return lang2.ErrNoPermission
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())
		kind := ctx.URLParamIntDefault("kind", int(resource2.AllKind))

		var params = []helper2.OptionFN{
			helper2.Page(page, pageSize),
			helper2.Kind(resource2.MeasureKind(kind)),
			helper2.Equipment(equipment.GetID()),
		}

		if !app2.IsDefaultAdminUser(admin) {
			params = append(params, helper2.DefaultEffect(app2.Config.DefaultEffect()))
			params = append(params, helper2.User(admin.GetID()))
		}

		states, total, err := s.GetStateList(params...)
		if err != nil {
			return err
		}

		var result = make([]model2.Map, 0, len(states))
		for _, state := range states {
			brief := state.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app2.Allow(admin, state, resource2.Ctrl),
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
	return response2.Wrap(func() interface{} {
		var form struct {
			Title           string `json:"title" valid:"required"`
			Desc            string `json:"desc"`
			MeasureID       int64  `json:"measure_id" valid:"required"`
			TransformScript string `json:"script"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}
		if _, err := govalidator.ValidateStruct(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		return app2.TransactionDo(func(s store2.Store) interface{} {
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app2.Allow(admin, equipment, resource2.Ctrl) {
				return lang2.ErrNoPermission
			}

			state, err := equipment.CreateState(form.Title, form.Desc, form.MeasureID, form.TransformScript)
			if err != nil {
				return err
			}
			err = app2.SetAllow(admin, state, resource2.View, resource2.Ctrl)
			if err != nil {
				return err
			}

			return state.Simple()
		})
	})
}

func StateDetail(stateID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, state, resource2.View) {
			return lang2.ErrNoPermission
		}

		return state.Detail()
	})
}

func UpdateState(stateID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, state, resource2.Ctrl) {
			return lang2.ErrNoPermission
		}

		var form struct {
			Title           *string `json:"title"`
			MeasureID       *int64  `json:"measure_id"`
			TransformScript *string `json:"script"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
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

		return lang2.Ok
	})
}

func DeleteState(stateID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, state, resource2.Ctrl) {
			return lang2.ErrNoPermission
		}

		err = state.Destroy()
		if err != nil {
			return err
		}
		return lang2.Ok
	})
}
