package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/response"
	"time"

	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
)

func RequireToken(p iris.Party) {
	jwtHandler := jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return app.Config.JwtTokenKey(), nil
		},
		Extractor: func(ctx iris.Context) (string, error) {
			return ctx.GetHeader("token"), nil
		},
		SigningMethod: jwt.SigningMethodHS512,
	})

	p.Use(jwtHandler.Serve)
	p.Use(hero.Handler(CheckUser))
}

func CheckUser(ctx iris.Context) {
	data := ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	user, err := app.Store().GetUser(data["sub"].(float64))
	if err == nil && user.IsEnabled() {
		ctx.Values().Set("__userID__", user.GetID())
		ctx.Next()
	} else {
		ctx.StatusCode(iris.StatusForbidden)
	}
}

func Login(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Username string `form:"username" validate:"required"`
			Password string `form:"password" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		user, err := app.Store().GetUser(form.Username)
		if err != nil {
			return err
		}
		if !user.IsEnabled() {
			return lang.ErrUserDisabled
		}
		if !user.CheckPassword(form.Password) {
			return lang.ErrPasswordWrong
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"sub": user.GetID(),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(app.Config.DefaultTokenExpiration()).Unix(),
		})

		token, err := claims.SignedString(app.Config.JwtTokenKey)
		if err != nil {
			return lang.InternalError(err)
		}
		return iris.Map{
			"token": token,
		}
	})
}
