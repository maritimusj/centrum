package resource

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/util"
	"strconv"
	"strings"
)

func GroupList() hero.Result {
	return response2.Wrap(func() interface{} {
		return app2.Store().GetResourceGroupList()
	})
}

func getUserPerm(user model2.User, res model2.Resource) interface{} {
	perm := map[string]bool{}
	switch res.ResourceClass() {
	case resource2.Api:
		perm["invoke"] = app2.Allow(user, res, resource2.Invoke)
	default:
		perm["view"] = app2.Allow(user, res, resource2.View)
		perm["ctrl"] = app2.Allow(user, res, resource2.Ctrl)
	}
	return perm
}

func getRolePerm(role model2.Role, res model2.Resource) (interface{}, error) {
	policies, err := role.GetPolicy(res)
	if err != nil {
		return nil, err
	}

	perm := map[string]interface{}{}
	switch res.ResourceClass() {
	case resource2.Api:
		perm["invoke"] = util.If(role.Name() == lang2.RoleSystemAdminName, true, func() interface{} {
			if v, ok := policies[resource2.Invoke]; ok {
				return v.Effect() == resource2.Allow
			} else {
				return app2.Config.DefaultEffect() == resource2.Allow
			}
		})
	default:
		perm["view"] = util.If(role.Name() == lang2.RoleSystemAdminName, true, func() interface{} {
			if v, ok := policies[resource2.View]; ok {
				return v.Effect() == resource2.Allow
			} else {
				return app2.Config.DefaultEffect()
			}
		})
		perm["ctrl"] = util.If(role.Name() == lang2.RoleSystemAdminName, true, func() interface{} {
			if v, ok := policies[resource2.Ctrl]; ok {
				return v.Effect() == resource2.Allow
			} else {
				return app2.Config.DefaultEffect() == resource2.Allow
			}
		})
	}
	return perm, nil
}

func List(classID int, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())
		keyword := ctx.URLParam("keyword")
		roleID := ctx.URLParamInt64Default("role", 0)
		userID := ctx.URLParamInt64Default("user", 0)

		var (
			err error

			role model2.Role
			user model2.User
		)

		s := app2.Store()

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

		var params = []helper2.OptionFN{
			helper2.Page(page, pageSize),
		}

		if keyword != "" {
			params = append(params, helper2.Keyword(keyword))
		}

		class := resource2.Class(classID)
		if !resource2.IsValidClass(class) {
			return lang2.Error(lang2.ErrInvalidResourceClassID)
		}

		if class == resource2.Api {
			sub := ctx.URLParam("sub")
			if sub != "" {
				params = append(params, helper2.Name(sub))
			}
		} else {
			sub := ctx.URLParamInt64Default("sub", -1)
			if sub != -1 {
				switch class {
				case resource2.Group:
					params = append(params, helper2.Parent(sub))
				case resource2.Device:
					params = append(params, helper2.Group(sub))
				case resource2.Measure:
					params = append(params, helper2.Device(sub))
				case resource2.Equipment:
					params = append(params, helper2.Group(sub))
				case resource2.State:
					params = append(params, helper2.Equipment(sub))
				}
			}
		}

		resources, total, err := s.GetResourceList(class, params...)
		if err != nil {
			return err
		}

		var list = make([]model2.Map, 0, len(resources))
		for _, res := range resources {
			entry := model2.Map{
				"id":          res.ResourceID(),
				"title":       res.ResourceTitle(),
				"desc":        res.ResourceDesc(),
				"class":       res.ResourceClass(),
				"class_title": lang2.ResourceClassTitle(class),
				"seg":         fmt.Sprintf("%d_%d", res.ResourceClass(), res.ResourceID()),
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

func GetList(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())

		roleID := ctx.URLParamInt64Default("role", 0)
		userID := ctx.URLParamInt64Default("user", 0)
		keyword := ctx.URLParam("keyword")

		var params = []helper2.OptionFN{
			helper2.Page(page, pageSize),
		}

		if keyword != "" {
			params = append(params, helper2.Keyword(keyword))
		}

		s := app2.Store()

		var (
			err       error
			resources []model2.Resource
			total     int64
		)

		seg := ctx.URLParam("seg")
		if len(seg) == 0 {
			resources, total, err = s.GetResourceList(resource2.Group, params...)
		} else {
			pair := strings.SplitN(seg, "_", 2)
			if len(pair) < 2 {
				return lang2.ErrInvalidRequestData
			}

			classID, err := strconv.ParseInt(pair[0], 10, 0)
			if err != nil {
				return lang2.ErrInvalidRequestData
			}

			resourceID, err := strconv.ParseInt(pair[1], 10, 0)
			if err != nil {
				return lang2.ErrInvalidRequestData
			}

			res, err := s.GetResource(resource2.Class(classID), resourceID)
			if err != nil {
				return err
			}

			resources, total, err = res.GetChildrenResources(params...)
			if err != nil {
				return err
			}
		}

		var (
			role model2.Role
			user model2.User
		)

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
			if app2.IsDefaultAdminUser(user) {
				return lang2.ErrFailedEditDefaultUser
			}
		}

		var list = make([]model2.Map, 0, len(resources))
		for _, res := range resources {
			classID := res.ResourceClass()
			entry := model2.Map{
				"id":          res.ResourceID(),
				"title":       res.ResourceTitle(),
				"desc":        res.ResourceDesc(),
				"class":       classID,
				"class_title": lang2.ResourceClassTitle(classID),
				"seg":         fmt.Sprintf("%d_%d", classID, res.ResourceID()),
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
