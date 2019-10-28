package user

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/event"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	status2 "github.com/maritimusj/centrum/gate/web/status"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/util"
)

func List(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var params []helper2.OptionFN
		var orgID int64

		s := app2.Store()
		admin := s.MustGetUserFromContext(ctx)
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
		params = append(params, helper2.Keyword(keyword))

		users, total, err := s.GetUserList(params...)
		if err != nil {
			return err
		}
		var result = make([]model2.Map, 0, len(users))
		for _, user := range users {
			result = append(result, user.Brief())
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
			OrgID    int64   `json:"org"`
			Username string  `json:"username" valid:"required"`
			Password string  `json:"password" valid:"required"`
			RoleIDs  []int64 `json:"roles"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		if _, err := govalidator.ValidateStruct(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		result := app2.TransactionDo(func(s store2.Store) interface{} {
			admin := s.MustGetUserFromContext(ctx)
			if exists, err := s.IsUserExists(form.Username); err != nil {
				return err
			} else if exists {
				return lang2.ErrUserExists
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
			role, err := s.CreateRole(app2.Config.DefaultOrganization(), form.Username, form.Username, lang2.Str(lang2.UserDefaultRoleDesc))
			if err != nil {
				return err
			}
			for _, res := range resource2.Guest {
				if res == resource2.Unknown {
					continue
				}
				res, err := s.GetApiResource(res)
				if err != nil {
					return err
				}
				_, err = role.SetPolicy(res, resource2.Invoke, resource2.Allow, nil)
				if err != nil {
					return err
				}
			}

			roles = append(roles, role)

			var org interface{}
			if app2.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app2.Config.DefaultOrganization()
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
			app2.Event.Publish(event.UserCreated, data.Get("adminID"), data.Get("userID"))
			return data.Pop("result")
		}

		return result
	})
}

func Detail(userID int64) hero.Result {
	return response2.Wrap(func() interface{} {
		user, err := app2.Store().GetUser(userID)
		if err != nil {
			return err
		}
		return user.Detail()
	})
}

func Update(userID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		result := app2.TransactionDo(func(s store2.Store) interface{} {
			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if app2.IsDefaultAdminUser(user) && !app2.IsDefaultAdminUser(admin) {
				return lang2.ErrNoPermission
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
				return lang2.ErrInvalidRequestData
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

			var data = model2.Map{}
			if form.Enable != nil {
				if false == *form.Enable {
					if app2.IsDefaultAdminUser(user) {
						return lang2.ErrFailedDisableDefaultUser
					}
					if user.Name() == admin.Name() {
						return lang2.ErrFailedDisableUserSelf
					}
				}
				data["enable"] = util.If(*form.Enable, status2.Enable, status2.Disable)
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
			app2.Event.Publish(event.UserUpdated, data.Get("adminID"), data.Get("userID"))
			return lang2.Ok
		}

		return result
	})
}

func Delete(userID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		result := app2.TransactionDo(func(s store2.Store) interface{} {
			user, err := s.GetUser(userID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if app2.IsDefaultAdminUser(user) {
				return lang2.ErrFailedRemoveDefaultUser
			}

			if user.Name() == admin.Name() {
				return lang2.ErrFailedRemoveUserSelf
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
			app2.Event.Publish(event.UserDeleted, data.Get("adminID"), data.Get("name"))
			return lang2.Ok
		}
		return result
	})
}

func UpdatePerm(userID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		key := fmt.Sprintf("UpdatePerm:%d", userID)
		if _, ok := global.Params.Get(key); ok {
			return lang2.ErrServerIsBusy
		}

		global.Params.Set(key, true)
		result := app2.TransactionDo(func(s store2.Store) interface{} {
			user, err := s.GetUser(userID)
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

			type P struct {
				ResourceClass int   `json:"class"`
				ResourceID    int64 `json:"id"`

				View *bool `json:"view"`   //是否可以观察资源
				Ctrl *bool `json:"ctrl"`   //是否可以控制资源
				Role *bool `json:"enable"` //角色是否启用
			}
			var form struct {
				Policies []P `json:"policies"`
			}
			if err = ctx.ReadJSON(&form); err != nil {
				return lang2.ErrInvalidRequestData
			}

			newRoles := hashset.New()
			for _, role := range roles {
				newRoles.Add(role.GetID())
			}

			admin := s.MustGetUserFromContext(ctx)

			//先处理角色设定
			for _, p := range form.Policies {
				if p.Role != nil {
					role, err := s.GetRole(p.ResourceID)
					if err != nil {
						return err
					}
					if app2.IsDefaultAdminUser(admin) || role.Name() != lang2.RoleSystemAdminName {
						if *p.Role {
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

			update := func(role model2.Role) interface{} {
				for _, p := range form.Policies {
					//角色设置，则跳过
					if p.Role != nil {
						continue
					}

					res, err := s.GetResource(resource2.Class(p.ResourceClass), p.ResourceID)
					if err != nil {
						return err
					}

					//Api权限不允许单独分配（只能通过角色分配）
					if res.ResourceClass() == resource2.Api {
						return lang2.ErrNoPermission
					}

					if p.View != nil {
						effect := util.If(*p.View, resource2.Allow, resource2.Deny).(resource2.Effect)
						_, err = role.SetPolicy(res, resource2.View, effect, map[model2.Resource]struct{}{})
						if err != nil {
							return err
						}
					}
					if p.Ctrl != nil {
						effect := util.If(*p.Ctrl, resource2.Allow, resource2.Deny).(resource2.Effect)
						_, err = role.SetPolicy(res, resource2.Ctrl, effect, map[model2.Resource]struct{}{})
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

			return lang2.ErrRoleNotFound
		})

		global.Params.Remove(key)
		if data, ok := result.(event.Data); ok {
			app2.Event.Publish(event.UserUpdated, data.Get("adminID"), data.Get("userID"))
			return lang2.Ok
		}

		return result
	})
}
