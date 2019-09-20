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
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/perm"
	"github.com/maritimusj/centrum/web/response"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		var params []helper.OptionFN
		var orgID int64
		if perm.IsDefaultAdminUser(ctx) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = perm.AdminUser(ctx).OrganizationID()
		}
		if orgID > 0 {
			params = append(params, helper.Organization(orgID))
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		params = append(params, helper.Page(page, pageSize))

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

func Create(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			OrgID int64  `json:"org"`
			Name  string `json:"name"`
			Title string `json:"title"`
		}

		if err := ctx.ReadJSON(&form); err != nil || form.Name == "" {
			return lang.ErrInvalidRequestData
		}

		if exists, err := s.IsRoleExists(form.Name); err != nil {
			return err
		} else if exists {
			return lang.ErrRoleExists
		}
		var org interface{}
		if perm.IsDefaultAdminUser(ctx) {
			if form.OrgID > 0 {
				org = form.OrgID
			} else {
				org = cfg.DefaultOrganization()
			}
		} else {
			org = perm.AdminUser(ctx).OrganizationID()
		}

		role, err := s.CreateRole(org, form.Name, form.Title)
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
			Invoke        *bool `json:"invoke"`
			View          *bool `json:"view"`
			Ctrl          *bool `json:"ctrl"`
		}

		var form struct {
			Title    string `json:"title"`
			Policies []P    `json:"policies"`
		}

		if err = ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if len(form.Policies) > 0 {
			for _, p := range form.Policies {
				res, err := s.GetResource(resource.Class(p.ResourceClass), p.ResourceID)
				if err != nil {
					return err
				}
				if p.Invoke != nil {
					effect := util.If(*p.Invoke, resource.Allow, resource.Deny).(resource.Effect)
					_, err = role.SetPolicy(res, resource.Invoke, effect, make(map[model.Resource]struct{}))
					if err != nil {
						return err
					}
				}
				if p.View != nil {
					effect := util.If(*p.View, resource.Allow, resource.Deny).(resource.Effect)
					_, err = role.SetPolicy(res, resource.View, effect, make(map[model.Resource]struct{}))
					if err != nil {
						return err
					}
				}
				if p.Ctrl != nil {
					effect := util.If(*p.Ctrl, resource.Allow, resource.Deny).(resource.Effect)
					_, err = role.SetPolicy(res, resource.Ctrl, effect, make(map[model.Resource]struct{}))
					if err != nil {
						return err
					}
				}
			}
		}

		if form.Title != "" {
			role.SetTitle(form.Title)
		}

		err = role.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(roleID int64, ctx iris.Context, s store.Store) hero.Result {
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
