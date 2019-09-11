package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/maritimusj/centrum/cache/memCache"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	mysqlStore "github.com/maritimusj/centrum/store/mysql"
	"github.com/maritimusj/centrum/web/api"
	"log"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	cfg := config.New()

	s := mysqlStore.New()
	err := s.Open(ctx, map[string]interface{}{
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
	_, err = s.GetUser(defaultUserName)
	if err != nil {
		if err != lang.Error(lang.ErrUserNotFound) {
			log.Fatal(err)
		}
		_, err := s.CreateUser(defaultUserName, []byte(defaultUserName), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	server := api.New()
	server.Register(s, cfg)

	err = server.Start(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
