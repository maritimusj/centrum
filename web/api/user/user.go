package user

import (
	"github.com/maritimusj/centrum/app"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/api/web"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/status"
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/response"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context) hero.Result {
	s := app.Store()
	admin := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {

		println("get")
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
		params = append(params, helper.Keyword(keyword))

		users, total, err := s.GetUserList(params...)
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(users))
		for _, user := range users {
			result = append(result, user.Brief())
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context, validate *validator.Validate) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			OrgID    int64  `json:"org"`
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
			RoleID   *int64 `json:"role"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		return app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)

			if exists, err := s.IsUserExists(form.Username); err != nil {
				return err
			} else if exists {
				return lang.ErrUserExists
			}

			var roles []interface{}
			var err error

			if app.Config.IsRoleEnabled() {
				if form.RoleID != nil {
					role, err := s.GetRole(*form.RoleID)
					if err != nil {
						return err
					}
					roles = append(roles, role)
				}

				if len(roles) == 0 {
					return lang.ErrRoleNotFound
				}
			} else {
				roles = append(roles, lang.RoleGuestName)
				//创建用户同名的role
				role, err := s.CreateRole(app.Config.DefaultOrganization(), form.Username, form.Username)
				if err != nil {
					return err
				}
				roles = append(roles, role)
			}

			var org interface{}
			if app.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app.Config.DefaultOrganization()
				}
			} else {
				org = admin.OrganizationID()
			}

			user, err := s.CreateUser(org, form.Username, []byte(form.Password), roles...)
			if err != nil {
				return err
			}

			return user.Simple()
		})
	})
}

func Detail(userID int64) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}
			return user.Detail()
		})
	})
}

func Update(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)

			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}

			if app.IsDefaultAdminUser(user) && !app.IsDefaultAdminUser(admin) {
				return lang.ErrNoPermission
			}

			var form struct {
				Enable   *bool   `json:"enable"`
				Password *string `json:"password"`
				Title    *string `json:"title"`
				Mobile   *string `json:"mobile"`
				Email    *string `json:"email"`
				Roles    []int64 `json:"roles"`
			}

			err = ctx.ReadJSON(&form)
			if err != nil {
				return lang.ErrInvalidRequestData
			}

			if app.Config.IsRoleEnabled() && form.Roles != nil {
				roles := make([]interface{}, 0, len(form.Roles))
				for _, role := range form.Roles {
					roles = append(roles, role)
				}
				err = user.SetRoles(roles...)
				if err != nil {
					return err
				}
			}

			if form.Password != nil && *form.Password != "" {
				user.ResetPassword(*form.Password)
			}

			var data = model.Map{}
			if form.Enable != nil {
				if false == *form.Enable {
					if app.IsDefaultAdminUser(user) {
						return lang.ErrFailedDisableDefaultUser
					}
					if user.Name() == admin.Name() {
						return lang.ErrFailedDisableUserSelf
					}
				}
				data["enable"] = util.If(*form.Enable, status.Enable, status.Disable)
			}
			if form.Title != nil {
				data["title"] = *form.Title
			}
			if form.Mobile != nil {
				data["mobile"] = *form.Mobile
			}
			if form.Email != nil {
				data["email"] = *form.Email
			}

			if len(data) > 0 {
				user.Update(data)
			}

			err = user.Save()
			if err != nil {
				return err
			}
			return lang.Ok
		})
	})
}

func Delete(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)

			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}

			if app.IsDefaultAdminUser(user) {
				return lang.ErrFailedRemoveDefaultUser
			}

			if user.Name() == admin.Name() {
				return lang.ErrFailedRemoveUserSelf
			}

			if !app.Config.IsRoleEnabled() {
				role, err := s.GetRole(user.Name())
				if err != nil {
					return err
				}
				err = role.Destroy()
				if err != nil {
					return err
				}
			}

			err = user.Destroy()
			if err != nil {
				return err
			}

			return lang.Ok
		})
	})
}

func UpdatePerm(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		return app.TransactionDo(func(s store.Store) interface{} {
			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}

			if app.IsDefaultAdminUser(user) {
				return lang.ErrFailedEditDefaultUserPerm
			}

			roles, err := user.GetRoles()
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
				Policies []P `json:"policies"`
			}
			if err = ctx.ReadJSON(&form); err != nil {
				return lang.ErrInvalidRequestData
			}

			update := func(role model.Role) interface{} {
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

				return lang.Ok
			}

			for _, role := range roles {
				if role.Title() == user.Name() {
					return update(role)
				}
			}

			return lang.ErrRoleNotFound
		})
	})
}

func LogList(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := app.Store().GetUser(userID)
		if err != nil {
			return err
		}

		return web.GetLogList(ctx, user.LogUID())
	})
}

func LogDelete(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := app.Store().GetUser(userID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.IsDefaultAdminUser(admin) {
			return lang.ErrNoPermission
		}

		return web.DeleteLog(ctx, user.LogUID())
	})
}
