package main

import (
	"context"
	"github.com/maritimusj/centrum/cache/memCache"
	"github.com/maritimusj/centrum/config"
	mysqlStore "github.com/maritimusj/centrum/store/mysql"
	"github.com/maritimusj/centrum/web/api"
	"log"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	cfg := config.New()

	s := mysqlStore.New()
	err := s.Open(ctx, map[string]interface{}{
		"cache": memCache.New(),
	})
	if err != nil {
		log.Fatal(err)
	}

	server := api.New()
	server.Register(s, cfg)

	err = server.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
