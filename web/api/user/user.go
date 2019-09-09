package user

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
	"gopkg.in/go-playground/validator.v9"
)

func List(ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")

		users, total, err := s.GetUserList(store.Keyword(keyword), store.Page(page, pageSize))
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
			Name     string `json:"name" validate:"required"`
			Password []byte `json:"password" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		user, err := s.CreateUser(form.Name, form.Password)
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

func Update(userID int64, ctx iris.Context, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := s.GetUser(userID)
		if err != nil {
			return err
		}

		var form struct {
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
		if form.Title != nil {
			data["title"] = form.Title
		}
		if form.Mobile != nil {
			data["mobile"] = form.Mobile
		}
		if form.Email != nil {
			data["email"] = form.Email
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

func Delete(userID int64, s store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		user, err := s.GetUser(userID)
		if err != nil {
			return err
		}
		err = user.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
