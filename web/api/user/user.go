package user

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/maritimusj/centrum/event"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/store"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/response"
	"github.com/maritimusj/centrum/web/status"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var params []helper.OptionFN
		var orgID int64

		s := app.Store()
		admin := s.MustGetUserFromContext(ctx)
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
			OrgID    int64   `json:"org"`
			Username string  `json:"username" validate:"required"`
			Password string  `json:"password" validate:"required"`
			RoleIDs  []int64 `json:"roles"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		result := app.TransactionDo(func(s store.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if exists, err := s.IsUserExists(form.Username); err != nil {
				return err
			} else if exists {
				return lang.ErrUserExists
			}

			var roles []interface{}
			var err error

			if len(form.RoleIDs) > 0 {
				for _, roleID := range form.RoleIDs {
					role, err := s.GetRole(roleID)
					if err != nil {
						return err
					}
					roles = append(roles, role)
				}
			}

			//创建用户同名的role，并设置guest权限
			role, err := s.CreateRole(app.Config.DefaultOrganization(), form.Username, form.Username, lang.Str(lang.UserDefaultRoleDesc))
			if err != nil {
				return err
			}
			for _, res := range resource.Guest {
				if res == resource.Unknown {
					continue
				}
				res, err := s.GetApiResource(res)
				if err != nil {
					return err
				}
				_, err = role.SetPolicy(res, resource.Invoke, resource.Allow, nil)
				if err != nil {
					return err
				}
			}

			roles = append(roles, role)

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

			data := event.Data{
				"userID":  user.GetID(),
				"adminID": admin.GetID(),
				"result":  user.Simple(),
			}
			return data
		})

		if data, ok := result.(*event.Data); ok {
			app.Event.Publish(event.UserCreated, data.Get("adminID"), data.Get("userID"))
			return data.Pop("result")
		}

		return result
	})
}

func Detail(userID int64) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := app.Store().GetUser(userID)
		if err != nil {
			return err
		}
		return user.Detail()
	})
}

func Update(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		result := app.TransactionDo(func(s store.Store) interface{} {
			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if app.IsDefaultAdminUser(user) && !app.IsDefaultAdminUser(admin) {
				return lang.ErrNoPermission
			}

			var form struct {
				Enable   *bool    `json:"enable"`
				Password *string  `json:"password"`
				Title    *string  `json:"title"`
				Mobile   *string  `json:"mobile"`
				Email    *string  `json:"email"`
				Roles    *[]int64 `json:"roles"`
			}

			err = ctx.ReadJSON(&form)
			if err != nil {
				return lang.ErrInvalidRequestData
			}

			if form.Roles != nil {
				roles := make([]interface{}, 0, len(*form.Roles))
				for _, role := range *form.Roles {
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

			eventData := event.Data{
				"userID":  user.GetID(),
				"adminID": admin.GetID(),
			}

			return eventData
		})

		if data, ok := result.(event.Data); ok {
			app.Event.Publish(event.UserUpdated, data.Get("adminID"), data.Get("userID"))
			return lang.Ok
		}

		return result
	})
}

func Delete(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		result := app.TransactionDo(func(s store.Store) interface{} {
			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if app.IsDefaultAdminUser(user) {
				return lang.ErrFailedRemoveDefaultUser
			}

			if user.Name() == admin.Name() {
				return lang.ErrFailedRemoveUserSelf
			}

			//删除用户同名角色
			role, err := s.GetRole(user.Name())
			if err != nil {
				return err
			}

			data := event.Data{
				"name":    user.Name(),
				"adminID": admin.GetID(),
			}

			err = role.Destroy()
			if err != nil {
				return err
			}

			err = user.Destroy()
			if err != nil {
				return err
			}

			return data
		})

		if data, ok := result.(event.Data); ok {
			app.Event.Publish(event.UserDeleted, data.Get("adminID"), data.Get("name"))
			return lang.Ok
		}
		return result
	})
}

func UpdatePerm(userID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		result := app.TransactionDo(func(s store.Store) interface{} {
			user, err := s.GetUser(userID)
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

			type P struct {
				ResourceClass int   `json:"class"`
				ResourceID    int64 `json:"id"`
				View          *bool `json:"view"`   //是否可以观察资源
				Ctrl          *bool `json:"ctrl"`   //是否可以控制资源
				Enable        *bool `json:"enable"` //角色是否启用
			}
			var form struct {
				Policies []P `json:"policies"`
			}
			if err = ctx.ReadJSON(&form); err != nil {
				return lang.ErrInvalidRequestData
			}

			newRoles := hashset.New()
			for _, role := range roles {
				newRoles.Add(role.GetID())
			}

			admin := s.MustGetUserFromContext(ctx)

			//先处理角色设定
			for _, p := range form.Policies {
				if p.Enable != nil {
					role, err := s.GetRole(p.ResourceID)
					if err != nil {
						return err
					}
					if app.IsDefaultAdminUser(admin) || role.Name() != lang.RoleSystemAdminName {
						if *p.Enable {
							newRoles.Add(role.GetID())
						} else {
							newRoles.Remove(role.GetID())
						}
					}
				}
			}
			err = user.SetRoles(newRoles.Values()...)
			if err != nil {
				return err
			}

			update := func(role model.Role) interface{} {
				for _, p := range form.Policies {
					//角色设置，则跳过
					if p.Enable != nil {
						continue
					}

					res, err := s.GetResource(resource.Class(p.ResourceClass), p.ResourceID)
					if err != nil {
						return err
					}

					//Api权限不允许单独分配（只能通过角色分配）
					if res.ResourceClass() == resource.Api {
						return lang.ErrNoPermission
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

				return nil
			}

			data := event.Data{
				"userID":  user.GetID(),
				"adminID": admin.GetID(),
			}

			for _, role := range roles {
				if role.Name() == user.Name() {
					err := update(role)
					if err != nil {
						return err
					}
					return data
				}
			}

			return lang.ErrRoleNotFound
		})

		if data, ok := result.(event.Data); ok {
			app.Event.Publish(event.UserUpdated, data.Get("adminID"), data.Get("userID"))
			return lang.Ok
		}
		return result
	})
}
