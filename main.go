package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/maritimusj/centrum/app"
	"github.com/maritimusj/centrum/db"
	"github.com/maritimusj/centrum/helper"
	log "github.com/sirupsen/logrus"

	"github.com/maritimusj/centrum/lang"
	_ "github.com/maritimusj/centrum/lang/zhCN"

	"github.com/maritimusj/centrum/config"
	mysqlDB "github.com/maritimusj/centrum/db/mysql"
	"github.com/maritimusj/centrum/logStore"
	mysqlStore "github.com/maritimusj/centrum/store/mysql"
	"github.com/maritimusj/centrum/web/api"
)

func main() {
	//初始化配置
	cfg := config.New()

	//命令行参数
	logLevel := flag.String("l", cfg.LogLevel(), "log level, [trace,debug,info,warn,error,fatal]")
	resetDefaultUserPassword := flag.Bool("r", false, "reset default user password")
	flag.Parse()

	cfg.SetLogLevel(*logLevel)

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}

	//日志仓库
	logDBStore := logStore.New()
	err = logDBStore.Open(app.Ctx, cfg.LogFileName())
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.AddHook(logDBStore)
	log.SetLevel(level)

	defer func() {
		app.Cancel()
		logDBStore.Wait()
	}()

	//数据库连接
	conn, err := mysqlDB.Open(app.Ctx, map[string]interface{}{
		"connStr": cfg.DBConnStr(),
	})
	if err != nil {
		log.Fatal(err)
	}

	result := conn.TransactionDo(func(db db.DB) interface{} {
		store := mysqlStore.Attach(db)

		_, total, err := store.GetApiResourceList(helper.Limit(1))
		if err != nil {
			return err
		}
		if total == 0 {
			//初始化api资源
			err = store.InitApiResource()
			if err != nil {
				return err
			}
		}

		//创建默认组织
		defaultOrg := cfg.DefaultOrganization()
		org, err := store.GetOrganization(defaultOrg)
		if err != nil {
			if err != lang.Error(lang.ErrOrganizationNotFound) {
				return err
			}
			_, err := store.CreateOrganization(defaultOrg, defaultOrg)
			if err != nil {
				return err
			}
		} else {
			org.Enable()
			if err = org.Save(); err != nil {
				return err
			}
		}

		//初始化系统角色
		_, err = store.GetRole(lang.RoleSystemAdminName)
		if err != nil {
			if err != lang.Error(lang.ErrRoleNotFound) {
				return err
			}
			err = store.InitDefaultRoles(defaultOrg)
			if err != nil {
				return err
			}
		}

		//创建默认用户
		defaultUserName := cfg.DefaultUserName()
		user, err := store.GetUser(defaultUserName)
		if err != nil {
			if err != lang.Error(lang.ErrUserNotFound) {
				return err
			}
			user, err = store.CreateUser(defaultOrg, defaultUserName, []byte(defaultUserName), lang.RoleSystemAdminName)
			if err != nil {
				return err
			}
		} else if *resetDefaultUserPassword {
			user.Enable()
			user.ResetPassword(defaultUserName)
			if err = user.Save(); err != nil {
				return err
			}
			if err = user.SetRoles(lang.RoleSystemAdminName); err != nil {
				return err
			}
			log.Warnln(lang.Str(lang.DefaultUserPasswordResetOk))
		}
		return nil
	})

	if result != nil {
		log.Fatal(result.(error))
	}

	//API服务
	server := api.New()
	server.Register(conn, mysqlStore.Attach(conn), cfg, logDBStore)

	err = server.Start(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
