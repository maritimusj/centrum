package role

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
	"github.com/maritimusj/centrum/web/store"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		admin := s.MustGetUserFromContext(ctx)

		var params []helper.OptionFN
		var orgID int64
		if app.IsDefaultAdminUser(admin) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = admin.OrganizationID()
		}
		if orgID > 0 {
			params = append(params, helper.Organization(orgID))
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
		params = append(params, helper.Page(page, pageSize))

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		var (
			userID = ctx.URLParamInt64Default("user", 0)
			user   model.User
			err    error
		)

		matchRoles := hashset.New()
		if userID > 0 {
			user, err = s.GetUser(userID)
			if err != nil {
				return err
			}

			if app.IsDefaultAdminUser(user) {
				return lang.ErrFailedEditDefaultUser
			}

			roles, err := user.GetRoles()
			if err != nil {
				return err
			}
			for _, p := range roles {
				matchRoles.Add(p.GetID())
			}
		}
		roles, total, err := s.GetRoleList(params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(roles))
		for _, role := range roles {
			//普通用户无法查看__sys__角色
			if role.Name() != lang.RoleSystemAdminName {
				brief := role.Brief()
				if userID > 0 {
					brief["matched"] = matchRoles.Contains(role.GetID())
				}
				result = append(result, brief)
			}
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			OrgID int64  `json:"org"`
			Name  string `json:"name"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
		}

		if err := ctx.ReadJSON(&form); err != nil || form.Name == "" {
			return lang.ErrInvalidRequestData
		}

		return app.TransactionDo(func(s store.Store) interface{} {
			if exists, err := s.IsRoleExists(form.Name); err != nil {
				return err
			} else if exists {
				return lang.ErrRoleExists
			}

			var org interface{}

			admin := s.MustGetUserFromContext(ctx)
			if app.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app.Config.DefaultOrganization()
				}
			} else {
				org = admin.OrganizationID()
			}

			role, err := s.CreateRole(org, form.Name, form.Title, form.Desc)
			if err != nil {
				return err
			}
			return role.Brief()
		})
	})
}

func Detail(roleID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		role, err := app.Store().GetRole(roleID)
		if err != nil {
			return err
		}
		admin := app.Store().MustGetUserFromContext(ctx)
		if app.IsDefaultAdminUser(admin) || role.Name() != lang.RoleSystemAdminName {
			return role.Detail()
		}
		return lang.ErrNoPermission
	})
}

func Update(roleID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		role, err := s.GetRole(roleID)
		if err != nil {
			return err
		}

		if role.Name() == lang.RoleSystemAdminName {
			return lang.ErrFailedEditDefaultUser
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

func Delete(roleID int64) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
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
	})
}
