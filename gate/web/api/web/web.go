package web

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"

	"github.com/maritimusj/centrum/gate/config"

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
	p.Use(hero.Handler(CheckReg))
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

/*
	注册码的形式为 0000-0000-0000
	以-分隔，共3段
	第一段为hmac的key
	计算方式：
	用上面的key，使用hmac-sha1算法取owner的hash值，比较hash值前4位是否与第二段相同，后4位是否与第三段相同
*/

func IsRegistered(owner, code string) bool {
	if code != "" {
		codes := strings.Split(code, "-")
		if len(codes) == 3 {
			hash := hmac.New(sha1.New, []byte(codes[0]))
			x := hex.EncodeToString(hash.Sum([]byte(owner)))
			return x[:4] == codes[1] && x[(len(x)-4):] == codes[2]
		}
	}

	return false
}

func CheckReg(ctx iris.Context) {
	if IsRegistered(app.Config.RegOwner(), app.Config.RegCode()) {
		ctx.Next()
	} else {
		ctx.StatusCode(iris.StatusUnauthorized)
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

		if !IsRegistered(owner, code) {
			return lang.ErrInvalidRegCode
		}
		_ = app.Config.BaseConfig.SetOption(config.SysRegOwnerPath, owner)
		_ = app.Config.BaseConfig.SetOption(config.SysRegCodePath, code)

		if err := app.Config.BaseConfig.Save(); err != nil {
			return err
		}

		return lang.Ok
	})
}

func GetReg() hero.Result {
	return response.Wrap(func() interface{} {
		return iris.Map{
			"registered": IsRegistered(app.Config.RegOwner(), app.Config.RegCode()),
			"owner":      app.Config.RegOwner(),
		}
	})
}

func Login(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		if !IsRegistered(app.Config.RegOwner(), app.Config.RegCode()) {
			return lang.ErrRegFirst
		}

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
			log.WithField("src", logStore.SystemLog).Infoln(lang.Str(lang.UserLoginFailedCauseDisabled, user.Name(), ctx.RemoteAddr()))
			return lang.ErrUserDisabled
		}
		if !user.CheckPassword(form.Password) {
			log.WithField("src", logStore.SystemLog).Infoln(lang.Str(lang.UserLoginFailedCausePasswordWrong, user.Name(), ctx.RemoteAddr()))
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
			log.WithField("src", logStore.SystemLog).Infoln(lang.Str(lang.UserLoginOk, user.Name(), ctx.RemoteAddr()))
		}

		return iris.Map{
			"token": token,
		}
	})
}
