package config

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/config"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/response"
)

type SysConfig struct {
	Title    string `json:"title"`
	RegOwner string `json:"reg_owner"`
	RegCode  string `json:"reg_code"`
}

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
	Sys     *SysConfig     `json:"sys"`
	Api     *APIConfig     `json:"api"`
	Def     *DefaultConfig `json:"default"`
	Log     *LogConfig     `json:"log"`
	Inverse *InverseConfig `json:"inverse"`
}

func Base() hero.Result {
	return response.Wrap(func() interface{} {
		return &Form{
			Sys: &SysConfig{
				Title:    app.Config.SysTitle(),
				RegOwner: app.Config.RegOwner(),
				RegCode:  app.Config.RegCode(),
			},
			Api: &APIConfig{
				Addr: app.Config.APIAddr(),
				Port: app.Config.APIPort(),
			},
			Def: &DefaultConfig{
				Username:        app.Config.DefaultUserName(),
				Organization:    app.Config.DefaultOrganization(),
				Effect:          int(app.Config.DefaultEffect()),
				PageSize:        app.Config.DefaultPageSize(),
				TokenExpiration: int64(app.Config.DefaultTokenExpiration().Seconds()),
			},
			Log: &LogConfig{
				Level:    app.Config.LogLevel(),
				FileName: app.Config.LogFileName(),
			},
		}
	})
}

func UpdateBase(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form Form
		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.Sys != nil {
			_ = app.Config.BaseConfig.SetOption(config.SysTitlePath, form.Sys.Title)
			_ = app.Config.BaseConfig.SetOption(config.SysRegCodePath, form.Sys.RegCode)
			_ = app.Config.BaseConfig.SetOption(config.SysRegOwnerPath, form.Sys.RegOwner)
		}

		if form.Api != nil {
			_ = app.Config.BaseConfig.SetOption(config.ApiAddrPath, form.Api.Addr)
			_ = app.Config.BaseConfig.SetOption(config.ApiPortPath, form.Api.Port)
		}

		if form.Def != nil {
			_ = app.Config.BaseConfig.SetOption(config.DefaultUserNamePath, form.Def.Username)
			_ = app.Config.BaseConfig.SetOption(config.DefaultOrganizationPath, form.Def.Organization)
			_ = app.Config.BaseConfig.SetOption(config.DefaultEffectPath, form.Def.Effect)
			_ = app.Config.BaseConfig.SetOption(config.DefaultPageSizePath, form.Def.PageSize)
			_ = app.Config.BaseConfig.SetOption(config.DefaultTokenExpirationPath, form.Def.TokenExpiration)
		}

		if form.Log != nil {
			_ = app.Config.BaseConfig.SetOption(config.LogLevelPath, form.Log.Level)
			_ = app.Config.BaseConfig.SetOption(config.LogFileNamePath, form.Log.FileName)
		}

		if form.Inverse != nil {
			_ = app.Config.BaseConfig.SetOption(config.InversePortPath, form.Inverse.Port)
		}

		return app.Config.BaseConfig.Save()
	})
}
