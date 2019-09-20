package web

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/app"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/response"
	"time"

	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
	log "github.com/sirupsen/logrus"
)

func RequireToken(p iris.Party) {
	jwtHandler := jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return app.Cfg.JwtTokenKey(), nil
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
	s := app.Store()
	defer s.Close()

	data := ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)
	user, err := s.GetUser(data["sub"].(float64))
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

		s := app.Store()
		defer s.Close()

		user, err := s.GetUser(form.Username)
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
			"exp": time.Now().Add(app.Cfg.MaxTokenExpiration()).Unix(),
		})

		token, err := claims.SignedString(app.Cfg.JwtTokenKey())
		if err != nil {
			return lang.InternalError(err)
		}
		return iris.Map{
			"token": token,
		}
	})
}

func GetLogList(ctx iris.Context, src string) interface{} {
	level := ctx.URLParam("level")
	start := ctx.URLParamInt64Default("start", 0)
	page := ctx.URLParamInt64Default("page", 1)
	pageSize := ctx.URLParamInt64Default("pagesize", app.Cfg.DefaultPageSize())

	x := uint64(start)
	logs, total, err := app.LogDBStore.Get(src, level, &x, uint64((page-1)*pageSize), uint64(pageSize))
	if err != nil {
		return err
	}

	var result = make([]iris.Map, 0, len(logs))
	for _, v := range logs {
		r := map[string]interface{}{}
		err := json.Unmarshal(v.Content, &r)

		if err != nil {
			result = append(result, iris.Map{
				"id":  v.ID,
				"raw": string(v.Content),
			})
		} else {
			result = append(result, iris.Map{
				"id":      v.ID,
				"content": r,
			})
		}
	}
	return iris.Map{
		"stats": app.LogDBStore.Stats(),
		"start": x,
		"total": total,
		"list":  result,
	}
}

func DeleteLog(ctx iris.Context, src string) interface{} {
	s := app.Store()
	defer s.Close()

	admin := s.MustGetUserFromContext(ctx)

	err := app.LogDBStore.Delete(src)
	if err != nil {
		return err
	}

	log.Info(lang.Str(lang.LogDeletedByUser, admin.Name()))
	return lang.Ok
}
