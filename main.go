package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/maritimusj/centrum/app"
	log "github.com/sirupsen/logrus"

	"github.com/maritimusj/centrum/lang"
	_ "github.com/maritimusj/centrum/lang/zhCN"

	"github.com/maritimusj/centrum/web/api"
)

func main() {
	//命令行参数
	logLevel := flag.String("l", app.Cfg.LogLevel(), "log level, [trace,debug,info,warn,error,fatal]")
	resetDefaultUserPassword := flag.Bool("r", false, "reset default user password")
	flag.Parse()

	if err := app.InitLog(*logLevel); err != nil {
		log.Fatal(err)
	}

	//数据库连接
	if err := app.InitDB(map[string]interface{}{
		"connStr": app.Cfg.DBConnStr(),
	}); err != nil {
		log.Fatal(err)
	}

	if err := app.Init(); err != nil {
		log.Fatal(err)
	}

	if *resetDefaultUserPassword {
		user, err := app.Store().GetUser(app.Cfg.DefaultUserName())
		if err != nil {
			log.Fatal(err)
		}

		user.Enable()
		user.ResetPassword(app.Cfg.DefaultUserName())
		if err = user.Save(); err != nil {
			log.Fatal(err)
		}
		if err = user.SetRoles(lang.RoleSystemAdminName); err != nil {
			log.Fatal()
		}
		log.Warnln(lang.Str(lang.DefaultUserPasswordResetOk))
	}

	//API服务
	server := api.New()
	err := server.Start(app.Cfg)
	if err != nil {
		log.Fatal(err)
	}
}
