package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/perm"
	"github.com/maritimusj/centrum/web/response"
	"time"

	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
)

func RequireToken(p iris.Party, cfg config.Config) {
	jwtHandler := jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return cfg.JwtTokenKey(), nil
		},
		Extractor: func(c iris.Context) (string, error) {
			return c.GetHeader("token"), nil
		},
		SigningMethod: jwt.SigningMethodHS512,
	})

	p.Use(jwtHandler.Serve)
	p.Use(hero.Handler(CheckUser))
}

func CheckUser(c iris.Context, store store.Store, cfg config.Config) {
	data := c.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	user, err := store.GetUser(data["sub"].(float64))
	if err != nil {
		c.StatusCode(iris.StatusForbidden)
	} else {
		c.Values().Set(perm.DefaultEffect, cfg.DefaultEffect())
		c.Values().Set(perm.AdminUserKey, user)
		if user.Name() == cfg.DefaultUserName() {
			c.Values().Set(perm.DefaultAdminUserKey, true)
		}
		c.Next()
	}
}

func Login(ctx iris.Context, store store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Username string `form:"username" validate:"required"`
			Password string `form:"password" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		user, err := store.GetUser(form.Username)
		if err != nil {
			return err
		}
		if !user.CheckPassword(form.Password) {
			return lang.ErrPasswordWrong
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"sub": user.GetID(),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(cfg.MaxTokenExpiration()).Unix(),
		})

		token, err := claims.SignedString(cfg.JwtTokenKey())
		if err != nil {
			return lang.InternalError(err)
		}
		return iris.Map{
			"token": token,
		}
	})
}
