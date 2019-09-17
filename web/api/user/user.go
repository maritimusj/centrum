package user

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/kr/pretty"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/perm"
	"github.com/maritimusj/centrum/web/response"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")

		users, total, err := s.GetUserList(helper.Keyword(keyword), helper.Page(page, pageSize))
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

func Create(ctx iris.Context, s store.Store, validate *validator.Validate, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
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

		if _, err := s.GetUser(form.Username); err != lang.Error(lang.ErrUserNotFound) {
			return lang.ErrUserExists
		}

		var role model.Role
		var err error
		if form.RoleID != nil {
			role, err = s.GetRole(*form.RoleID)
			if err != nil {
				return err
			}
		}

		if cfg.IsRoleEnabled() && role == nil {
			return lang.ErrRoleNotFound
		}

		user, err := s.CreateUser(form.Username, []byte(form.Password), role)
		if err != nil {
			return err
		}

		return user.Simple()
	})
}

func Detail(userID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := s.GetUser(userID)
		if err != nil {
			return err
		}
		return user.Detail()
	})
}

func Update(userID int64, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := s.GetUser(userID)
		if err != nil {
			return err
		}

		if user.Name() == cfg.DefaultUserName() && perm.AdminUser(ctx).Name() != cfg.DefaultUserName() {
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

		if cfg.IsRoleEnabled() && form.Roles != nil {
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
				if user.Name() == cfg.DefaultUserName() {
					return lang.ErrFailedDisableDefaultUser
				}
				if user.Name() == perm.AdminUser(ctx).Name() {
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
}

func Delete(userID int64, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := s.GetUser(userID)
		if err != nil {
			return err
		}

		if user.Name() == cfg.DefaultUserName() {
			return lang.ErrFailedRemoveDefaultUser
		}

		if user.Name() == perm.AdminUser(ctx).Name() {
			return lang.ErrFailedRemoveUserSelf
		}

		err = user.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func UpdatePerm(userID int64, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := s.GetUser(userID)
		if err != nil {
			return err
		}

		if user.Name() == cfg.DefaultUserName() {
			return lang.ErrFailedEditDefaultUserPerm
		}

		roles, err := user.GetRoles()
		if err != nil {
			return err
		}

		update := func(role model.Role) interface{} {
			type P struct {
				ResourceClass int   `json:"class"`
				ResourceID    int64 `json:"id"`
				Action        int8  `json:"action"`
				Effect        int8  `json:"effect"`
			}
			var form struct {
				Title   *string `json:"title"`
				Polices []P     `json:"policies"`
			}
			if err = ctx.ReadJSON(&form); err != nil {
				return lang.ErrInvalidRequestData
			}

			fmt.Printf("%# v", pretty.Formatter(form))

			if len(form.Polices) > 0 {
				for _, p := range form.Polices {
					res, err := s.GetResource(resource.Class(p.ResourceClass), p.ResourceID)
					if err != nil {
						return err
					}
					_, err = role.SetPolicy(res, resource.Action(p.Action), resource.Effect(p.Effect))
					if err != nil {
						return err
					}
				}
			}

			if form.Title != nil && *form.Title != "" {
				role.SetTitle(*form.Title)
			}

			err = role.Save()
			if err != nil {
				return err
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
}
