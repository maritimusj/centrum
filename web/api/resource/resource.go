package resource

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

func GroupList(store store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		return store.GetResourceGroupList()
	})
}

func List(classID int, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")
		roleID := ctx.URLParamInt64Default("role", 0)

		var params = []helper.OptionFN{
			helper.Page(page, pageSize),
		}

		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		class := resource.Class(classID)
		if !resource.IsValidClass(class) {
			return lang.Error(lang.ErrInvalidResourceClassID)
		}

		var role model.Role
		var err error
		if roleID > 0 {
			role, err = s.GetRole(roleID)
			if err != nil {
				return err
			}
		}

		if class == resource.Api {
			sub := ctx.URLParam("sub")
			if sub != "" {
				params = append(params, helper.Name(sub))
			}
		} else {
			sub := ctx.URLParamInt64Default("sub", -1)
			if sub != -1 {
				switch class {
				case resource.Group:
					params = append(params, helper.Parent(sub))
				case resource.Device:
					params = append(params, helper.Group(sub))
				case resource.Measure:
					params = append(params, helper.Device(sub))
				case resource.Equipment:
					params = append(params, helper.Group(sub))
				case resource.State:
					params = append(params, helper.Equipment(sub))
				}
			}
		}

		resources, total, err := s.GetResourceList(resource.Class(classID), params...)
		if err != nil {
			return err
		}

		var list = make([]model.Map, 0, len(resources))

		for _, res := range resources {
			entry := model.Map{
				"id":          res.ResourceID(),
				"title":       res.ResourceTitle(),
				"desc":        res.ResourceDesc(),
				"class":       classID,
				"class_title": lang.ResourceClassTitle(class),
			}

			if role != nil {
				perm := model.Map{}
				policies, err := role.GetPolicy(res)
				if err != nil {
					return err
				}
				if v, ok := policies[resource.View]; ok {
					perm["view"] = v.Effect()
				} else {
					perm["view*"] = cfg.DefaultEffect()
				}
				if v, ok := policies[resource.Ctrl]; ok {
					perm["ctrl"] = v.Effect()
				} else {
					perm["ctrl*"] = cfg.DefaultEffect()
				}
				entry["perm"] = perm
			}

			list = append(list, entry)
		}

		result := iris.Map{
			"total": total,
			"list":  list,
		}
		if role != nil {
			result["role"] = role.Brief()
		}
		return result
	})
}
