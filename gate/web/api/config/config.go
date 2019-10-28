package config

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/config"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	response2 "github.com/maritimusj/centrum/gate/web/response"
)

type APIConfig struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

type InverseConfig struct {
	Port int `json:"port"`
}

type DefaultConfig struct {
	Username        string `json:"username"`
	Organization    string `json:"organization"`
	Effect          int    `json:"effect"`
	PageSize        int64  `json:"pagesize"`
	TokenExpiration int64  `json:"expiration"`
}

type LogConfig struct {
	Level    string `json:"level"`
	FileName string `json:"filename"`
}

type Form struct {
	Api     *APIConfig     `json:"api"`
	Def     *DefaultConfig `json:"default"`
	Log     *LogConfig     `json:"log"`
	Inverse *InverseConfig `json:"inverse"`
}

func Base() hero.Result {
	return response2.Wrap(func() interface{} {
		return &Form{
			Api: &APIConfig{
				Addr: app2.Config.APIAddr(),
				Port: app2.Config.APIPort(),
			},
			Def: &DefaultConfig{
				Username:        app2.Config.DefaultUserName(),
				Organization:    app2.Config.DefaultOrganization(),
				Effect:          int(app2.Config.DefaultEffect()),
				PageSize:        app2.Config.DefaultPageSize(),
				TokenExpiration: int64(app2.Config.DefaultTokenExpiration().Seconds()),
			},
			Log: &LogConfig{
				Level:    app2.Config.LogLevel(),
				FileName: app2.Config.LogFileName(),
			},
		}
	})
}

func UpdateBase(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form Form
		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		if form.Api != nil {
			_ = app2.Config.BaseConfig.SetOption(config.ApiAddrPath, form.Api.Addr)
			_ = app2.Config.BaseConfig.SetOption(config.ApiPortPath, form.Api.Port)
		}

		if form.Def != nil {
			_ = app2.Config.BaseConfig.SetOption(config.DefaultUserNamePath, form.Def.Username)
			_ = app2.Config.BaseConfig.SetOption(config.DefaultOrganizationPath, form.Def.Organization)
			_ = app2.Config.BaseConfig.SetOption(config.DefaultEffectPath, form.Def.Effect)
			_ = app2.Config.BaseConfig.SetOption(config.DefaultPageSizePath, form.Def.PageSize)
			_ = app2.Config.BaseConfig.SetOption(config.DefaultTokenExpirationPath, form.Def.TokenExpiration)
		}

		if form.Log != nil {
			_ = app2.Config.BaseConfig.SetOption(config.LogLevelPath, form.Log.Level)
			_ = app2.Config.BaseConfig.SetOption(config.LogFileNamePath, form.Log.FileName)
		}

		if form.Inverse != nil {
			_ = app2.Config.BaseConfig.SetOption(config.InversePortPath, form.Inverse.Port)
		}

		return app2.Config.BaseConfig.Save()
	})
}
