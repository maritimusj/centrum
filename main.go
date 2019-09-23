package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/maritimusj/centrum/app"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/util"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/maritimusj/centrum/lang"
	_ "github.com/maritimusj/centrum/lang/zhCN"

	"github.com/maritimusj/centrum/web/api"
)

func main() {
	//命令行参数
	logLevel := flag.String("l", app.Config.LogLevel(), "log level, [trace,debug,info,warn,error,fatal]")
	resetDefaultUserPassword := flag.Bool("r", false, "reset default user password")
	resetDB := flag.Bool("flush", false, "erase all data in database")

	flag.Parse()

	if err := app.Init(*logLevel); err != nil {
		log.Fatal(err)
	}

	if *resetDB {
		code := util.RandStr(4, util.RandNum)
		fmt.Print(lang.Str(lang.ConfirmAdminPassword, code))

		var confirm string
		_, _ = fmt.Scanln(&confirm)
		if confirm != code {
			log.Fatal(lang.Error(lang.ErrConfirmCodeWrong))
		} else {
			res := app.TransactionDo(func(store store.Store) interface{} {
				return store.EraseAllData()
			})
			if res != nil {
				log.Fatal(res.(error))
			} else {
				fmt.Printf(lang.Str(lang.FlushDBOk))
				os.Exit(0)
			}
		}
	}

	//重置系统默认用户的密码和角色信息
	if *resetDefaultUserPassword {
		user, err := app.Store().GetUser(app.Config.DefaultUserName())
		if err != nil {
			log.Fatal(err)
		}

		user.Enable()
		user.ResetPassword(app.Config.DefaultUserName())
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
	err := server.Start(app.Config)
	if err != nil {
		log.Fatal(err)
	}
}
