package resource

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/app"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/web/response"
)

func GroupList() hero.Result {
	return response.Wrap(func() interface{} {
		return app.Store().GetResourceGroupList()
	})
}

func List(classID int, ctx iris.Context) hero.Result {
	s := app.Store()

	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
		keyword := ctx.URLParam("keyword")
		roleID := ctx.URLParamInt64Default("role", 0)
		userID := ctx.URLParamInt64Default("user", 0)

		var err error

		var role model.Role
		var user model.User

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
				if resource.Class(classID) == resource.Api {
					if v, ok := policies[resource.Invoke]; ok {
						perm["invoke"] = v.Effect() == resource.Allow
					} else {
						perm["invoke"] = app.Config.DefaultEffect() == resource.Allow
					}
				} else {
					if v, ok := policies[resource.View]; ok {
						perm["view"] = v.Effect() == resource.Allow
					} else {
						perm["view"] = app.Config.DefaultEffect()
					}
					if v, ok := policies[resource.Ctrl]; ok {
						perm["ctrl"] = v.Effect() == resource.Allow
					} else {
						perm["ctrl"] = app.Config.DefaultEffect() == resource.Allow
					}
				}
				entry["perm"] = perm
			} else if user != nil {
				perm := model.Map{}
				if resource.Class(classID) == resource.Api {
					allowed, err := user.IsAllow(res, resource.Invoke)
					if err != nil {
						switch err {
						case lang.Error(lang.ErrPolicyNotFound):
							perm["invoke"] = app.Config.DefaultEffect() == resource.Allow
						default:
							perm["invoke"] = false
						}
					} else {
						perm["invoke"] = allowed
					}
				} else {
					allowView, err := user.IsAllow(res, resource.View)
					if err != nil {
						switch err {
						case lang.Error(lang.ErrPolicyNotFound):
							perm["view"] = app.Config.DefaultEffect() == resource.Allow
						default:
							perm["view"] = false
						}
					} else {
						perm["view"] = allowView
					}
					allowCtrl, err := user.IsAllow(res, resource.Ctrl)
					if err != nil {
						switch err {
						case lang.Error(lang.ErrPolicyNotFound):
							perm["ctrl"] = app.Config.DefaultEffect() == resource.Allow
						default:
							perm["ctrl"] = false
						}
					} else {
						perm["ctrl"] = allowCtrl
					}
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
