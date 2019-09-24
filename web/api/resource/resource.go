package resource

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
)

func GroupList() hero.Result {
	return response.Wrap(func() interface{} {
		return app.Store().GetResourceGroupList()
	})
}

func getUserPerm(user model.User, res model.Resource) interface{} {
	perm := map[string]bool{}
	switch res.ResourceClass() {
	case resource.Api:
		perm["invoke"] = app.Allow(user, res, resource.Invoke)
	default:
		perm["view"] = app.Allow(user, res, resource.View)
		perm["ctrl"] = app.Allow(user, res, resource.Ctrl)
	}
	return perm
}

func getRolePerm(role model.Role, res model.Resource) (interface{}, error) {
	policies, err := role.GetPolicy(res)
	if err != nil {
		return nil, err
	}

	perm := map[string]interface{}{}
	switch res.ResourceClass() {
	case resource.Api:
		perm["invoke"] = util.If(role.Name() == lang.RoleSystemAdminName, true, func() interface{} {
			if v, ok := policies[resource.Invoke]; ok {
				return v.Effect() == resource.Allow
			} else {
				return app.Config.DefaultEffect() == resource.Allow
			}
		})
	default:
		perm["view"] = util.If(role.Name() == lang.RoleSystemAdminName, true, func() interface{} {
			if v, ok := policies[resource.View]; ok {
				return v.Effect() == resource.Allow
			} else {
				return app.Config.DefaultEffect()
			}
		})
		perm["ctrl"] = util.If(role.Name() == lang.RoleSystemAdminName, true, func() interface{} {
			if v, ok := policies[resource.Ctrl]; ok {
				return v.Effect() == resource.Allow
			} else {
				return app.Config.DefaultEffect() == resource.Allow
			}
		})
	}
	return perm, nil
}

func List(classID int, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
		keyword := ctx.URLParam("keyword")
		roleID := ctx.URLParamInt64Default("role", 0)
		userID := ctx.URLParamInt64Default("user", 0)

		var err error

		var role model.Role
		var user model.User

		s := app.Store()

		if roleID > 0 {
			role, err = s.GetRole(roleID)
			if err != nil {
				return err
			}
		} else if userID > 0 {
			user, err = s.GetUser(userID)
			if err != nil {
				return err
			}
		}

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

		resources, total, err := s.GetResourceList(class, params...)
		if err != nil {
			return err
		}

		var list = make([]model.Map, 0, len(resources))
		for _, res := range resources {
			entry := model.Map{
				"id":          res.ResourceID(),
				"title":       res.ResourceTitle(),
				"desc":        res.ResourceDesc(),
				"class":       class,
				"class_title": lang.ResourceClassTitle(class),
			}

			if role != nil {
				perm, err := getRolePerm(role, res)
				if err != nil {
					return err
				}
				entry["perm"] = perm
			} else if user != nil {
				entry["perm"] = getUserPerm(user, res)
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
