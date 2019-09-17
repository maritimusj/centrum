package main

import (
	"context"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/maritimusj/centrum/lang"
	_ "github.com/maritimusj/centrum/lang/zhCN"

	"github.com/maritimusj/centrum/cache/memCache"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/logStore"
	mysqlStore "github.com/maritimusj/centrum/store/mysql"
	"github.com/maritimusj/centrum/web/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

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
	err = logDBStore.Open(ctx, cfg.LogFileName())
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.AddHook(logDBStore)
	log.SetLevel(level)

	defer func() {
		cancel()
		logDBStore.Wait()
	}()

	//数据库连接
	s := mysqlStore.New()
	err = s.Open(ctx, map[string]interface{}{
		"cache":   memCache.New(),
		"connStr": cfg.DBConnStr(),
	})
	if err != nil {
		log.Fatal(err)
	}

	//初始化api资源
	err = s.InitApiResource()
	if err != nil {
		log.Fatal(err)
	}

	//创建默认用户
	defaultUserName := cfg.DefaultUserName()
	user, err := s.GetUser(defaultUserName)
	if err != nil {
		if err != lang.Error(lang.ErrUserNotFound) {
			log.Fatal(err)
		}
		_, err := s.CreateUser(defaultUserName, []byte(defaultUserName), nil)
		if err != nil {
			log.Fatal(err)
		}
	} else if *resetDefaultUserPassword {
		user.Enable()
		user.ResetPassword(defaultUserName)
		if err = user.Save(); err != nil {
			log.Fatal(err)
		}
		log.Warnln(lang.Str(lang.DefaultUserPasswordResetOk))
	}

	//API服务
	server := api.New()
	server.Register(s, cfg, logDBStore)

	err = server.Start(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
