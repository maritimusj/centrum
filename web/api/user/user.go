package user

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/status"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/util"
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

func Create(ctx iris.Context, s store.Store, validate *validator.Validate) hero.Result {
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

		if form.Roles != nil {
			roles := make([]interface{}, 0, len(form.Roles))
			for _, role := range form.Roles {
				roles = append(roles, role)
			}
			err = user.SetRoles(roles...)
			if err != nil {
				return err
			}
		}

		if form.Password != nil {
			err = user.ResetPassword(*form.Password)
			if err != nil {
				return err
			}
		}
		var data = model.Map{}
		if form.Enable != nil {
			if !*form.Enable && user.Name() == cfg.DefaultUserName() {
				return lang.ErrFailedRemoveDefaultUser
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
			err = user.Update(data)
			if err != nil {
				return err
			}
		}
		err = user.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(userID int64, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := s.GetUser(userID)
		if err != nil {
			return err
		}

		if user.Name() == cfg.DefaultUserName() {
			return lang.ErrFailedRemoveDefaultUser
		}

		err = user.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
