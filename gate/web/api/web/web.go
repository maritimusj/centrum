package web

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	"time"

	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
)

func RequireToken(p iris.Party) {
	jwtHandler := jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return app2.Config.JwtTokenKey(), nil
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
	user, err := app2.Store().GetUser(data["sub"].(float64))
	if err == nil && user.IsEnabled() {
		ctx.Values().Set("__userID__", user.GetID())
		ctx.Next()
	} else {
		ctx.StatusCode(iris.StatusForbidden)
	}
}

func Login(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Username string `form:"username" valid:"required"`
			Password string `form:"password" valid:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		user, err := app2.Store().GetUser(form.Username)
		if err != nil {
			return err
		}
		if !user.IsEnabled() {
			return lang2.ErrUserDisabled
		}
		if !user.CheckPassword(form.Password) {
			return lang2.ErrPasswordWrong
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"sub": user.GetID(),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(app2.Config.DefaultTokenExpiration()).Unix(),
		})

		token, err := claims.SignedString(app2.Config.JwtTokenKey())
		if err != nil {
			return lang2.InternalError(err)
		}
		return iris.Map{
			"token": token,
		}
	})
}
