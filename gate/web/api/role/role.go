package role

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	"github.com/maritimusj/centrum/util"
)

func List(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		admin := s.MustGetUserFromContext(ctx)

		var params []helper2.OptionFN
		var orgID int64
		if app2.IsDefaultAdminUser(admin) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = admin.OrganizationID()
		}
		if orgID > 0 {
			params = append(params, helper2.Organization(orgID))
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())
		params = append(params, helper2.Page(page, pageSize))

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper2.Keyword(keyword))
		}

		var (
			userID = ctx.URLParamInt64Default("user", 0)
			user   model2.User
			err    error
		)

		matchRoles := hashset.New()
		if userID > 0 {
			user, err = s.GetUser(userID)
			if err != nil {
				return err
			}

			if app2.IsDefaultAdminUser(user) {
				return lang2.ErrFailedEditDefaultUser
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

		var result = make([]model2.Map, 0, len(roles))
		for _, role := range roles {
			//普通用户无法查看__sys__角色
			if role.Name() != lang2.RoleSystemAdminName {
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
	return response2.Wrap(func() interface{} {
		var form struct {
			OrgID int64  `json:"org"`
			Name  string `json:"name"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
		}

		if err := ctx.ReadJSON(&form); err != nil || form.Name == "" {
			return lang2.ErrInvalidRequestData
		}

		return app2.TransactionDo(func(s store2.Store) interface{} {
			if exists, err := s.IsRoleExists(form.Name); err != nil {
				return err
			} else if exists {
				return lang2.ErrRoleExists
			}

			var org interface{}

			admin := s.MustGetUserFromContext(ctx)
			if app2.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app2.Config.DefaultOrganization()
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
	return response2.Wrap(func() interface{} {
		role, err := app2.Store().GetRole(roleID)
		if err != nil {
			return err
		}
		admin := app2.Store().MustGetUserFromContext(ctx)
		if app2.IsDefaultAdminUser(admin) || role.Name() != lang2.RoleSystemAdminName {
			return role.Detail()
		}
		return lang2.ErrNoPermission
	})
}

func Update(roleID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		role, err := s.GetRole(roleID)
		if err != nil {
			return err
		}

		if role.Name() == lang2.RoleSystemAdminName {
			return lang2.ErrFailedEditDefaultUser
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
			return lang2.ErrInvalidRequestData
		}

		if len(form.Policies) > 0 {
			for _, p := range form.Policies {
				res, err := s.GetResource(resource2.Class(p.ResourceClass), p.ResourceID)
				if err != nil {
					return err
				}
				if p.Invoke != nil {
					effect := util.If(*p.Invoke, resource2.Allow, resource2.Deny).(resource2.Effect)
					_, err = role.SetPolicy(res, resource2.Invoke, effect, make(map[model2.Resource]struct{}))
					if err != nil {
						return err
					}
				}
				if p.View != nil {
					effect := util.If(*p.View, resource2.Allow, resource2.Deny).(resource2.Effect)
					_, err = role.SetPolicy(res, resource2.View, effect, make(map[model2.Resource]struct{}))
					if err != nil {
						return err
					}
				}
				if p.Ctrl != nil {
					effect := util.If(*p.Ctrl, resource2.Allow, resource2.Deny).(resource2.Effect)
					_, err = role.SetPolicy(res, resource2.Ctrl, effect, make(map[model2.Resource]struct{}))
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
		return lang2.Ok
	})
}

func Delete(roleID int64) hero.Result {
	return response2.Wrap(func() interface{} {
		return app2.TransactionDo(func(s store2.Store) interface{} {
			role, err := s.GetRole(roleID)
			if err != nil {
				return err
			}

			err = role.Destroy()
			if err != nil {
				return err
			}
			return lang2.Ok
		})
	})
}
