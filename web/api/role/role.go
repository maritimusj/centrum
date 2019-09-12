package role

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())

		var params = []helper.OptionFN{helper.Page(page, pageSize)}

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		userID := ctx.URLParamInt64Default("user", -1)
		if userID != -1 {
			params = append(params, helper.User(userID))
		}

		roles, total, err := s.GetRoleList(params...)
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(roles))
		for _, role := range roles {
			result = append(result, role.Brief())
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Title string `json:"title"`
		}

		if err := ctx.ReadJSON(&form); err != nil || form.Title == "" {
			return lang.ErrInvalidRequestData
		}

		role, err := s.CreateRole(form.Title)
		if err != nil {
			return err
		}
		return role.Brief()
	})
}

func Detail(roleID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		role, err := s.GetRole(roleID)
		if err != nil {
			return err
		}
		return role.Detail()
	})
}

func Update(roleID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		role, err := s.GetRole(roleID)
		if err != nil {
			return err
		}

		type P struct {
			ResourceClass int   `json:"class"`
			ResourceID    int64 `json:"id"`
			Action        int   `json:"action"`
			Effect        int   `json:"effect"`
		}
		var form struct {
			Title   *string `json:"title"`
			Polices []P     `json:"policies"`
		}
		if err = ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.Title != nil && *form.Title != "" {
			role.SetTitle(*form.Title)
		}

		if len(form.Polices) > 0 {
			for _, p := range form.Polices {
				res, err := s.GetResource(resource.Class(p.ResourceClass), p.ResourceID)
				if err != nil {
					return err
				}
				_, err = role.SetPolicy(res, resource.Action(p.Action), resource.Effect(p.Effect))
				if err != nil {
					return err
				}
			}
		}

		err = role.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(roleID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		role, err := s.GetRole(roleID)
		if err != nil {
			return err
		}

		err = role.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
