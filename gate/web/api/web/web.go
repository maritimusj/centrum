package web

import (
	"strings"
	"time"

	"github.com/maritimusj/centrum/global"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/logStore"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/response"
	log "github.com/sirupsen/logrus"

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

func Reg(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Owner string `json:"owner"`
			Code  string `json:"code"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		owner := strings.ToLower(strings.TrimSpace(form.Owner))
		code := strings.ToLower(strings.TrimSpace(form.Code))

		if err := app.SaveRegisterInfo(owner, code); err != nil {
			return err
		}

		return lang.Ok
	})
}

func GetReg() hero.Result {
	return response.Wrap(func() interface{} {
		return iris.Map{
			"registered":   app.IsRegistered(),
			"owner":        app.Config.RegOwner(),
			"fingerprints": app.Fingerprints(),
		}
	})
}

func Login(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Username string `json:"username" valid:"required"`
			Password string `json:"password" valid:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		user, err := app.Store().GetUser(form.Username)
		if err != nil {
			return err
		}
		if !user.IsEnabled() {
			log.WithField("src", logStore.SystemLog).Infoln(lang.UserLoginFailedCauseDisabled.Str(user.Name(), ctx.RemoteAddr()))
			return lang.ErrUserDisabled
		}
		if !user.CheckPassword(form.Password) {
			log.WithField("src", logStore.SystemLog).Infoln(lang.UserLoginFailedCausePasswordWrong.Str(user.Name(), ctx.RemoteAddr()))
			return lang.ErrPasswordWrong
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"sub": user.GetID(),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(app.Config.DefaultTokenExpiration()).Unix(),
		})

		token, err := claims.SignedString(app.Config.JwtTokenKey())
		if err != nil {
			return lang.InternalError(err)
		}

		if !ctx.URLParamExists("refresh") {
			log.WithField("src", logStore.SystemLog).Infoln(lang.UserLoginOk.Str(user.Name(), ctx.RemoteAddr()))
		}

		//注册用户，接收消息
		global.Create(token, user.GetID())

		return iris.Map{
			"token": token,
		}
	})
}
