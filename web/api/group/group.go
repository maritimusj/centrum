package group

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")
		parentGroupID := ctx.URLParamInt64Default("parent", 0)

		if parentGroupID > 0 {
			_, err := s.GetGroup(parentGroupID)
			if err != nil {
				return err
			}
		}

		groups, total, err := s.GetGroupList(parentGroupID, store.Keyword(keyword), store.Page(page, pageSize))
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(groups))
		for _, group := range groups {
			result = append(result, group.Brief())
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
			Title         string `json:"title" validate:"required"`
			ParentGroupID int64  `json:"parent" validate:"min=0"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.ParentGroupID > 0 {
			_, err := s.GetGroup(form.ParentGroupID)
			if err != nil {
				return err
			}
		}

		group, err := s.CreateGroup(form.Title, form.ParentGroupID)
		if err != nil {
			return err
		}

		return group.Simple()
	})
}

func Detail(groupID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}
		return group.Detail()
	})
}

func Update(groupID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}

		var form struct {
			Title      *string `json:"title"`
			Devices    []int64 `json:"roles"`
			Equipments []int64 `json:"equipments"`
		}

		err = ctx.ReadJSON(&form)
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.Title != nil && *form.Title != "" {
			err = group.SetTitle(*form.Title)
			if err != nil {
				return err
			}
		}

		if len(form.Devices) > 0 {
			devices := make([]interface{}, 0, len(form.Devices))
			for _, device := range form.Devices {
				devices = append(devices, device)
			}
			err = group.AddDevice(devices...)
			if err != nil {
				return err
			}
		}
		if len(form.Equipments) > 0 {
			equipments := make([]interface{}, 0, len(form.Equipments))
			for _, equipment := range form.Equipments {
				equipments = append(equipments, equipment)
			}
			err = group.AddEquipment(equipments...)
			if err != nil {
				return err
			}
		}

		err = group.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(groupID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return err
		}
		err = group.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
