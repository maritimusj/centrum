package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/maritimusj/durafmt"

	"github.com/maritimusj/centrum/gate/web/edge"

	"github.com/spf13/viper"

	"github.com/maritimusj/centrum/gate/logStore"

	_ "github.com/mattn/go-sqlite3"

	"github.com/maritimusj/centrum/gate/lang"
	_ "github.com/maritimusj/centrum/gate/lang/enUS"
	_ "github.com/maritimusj/centrum/gate/lang/zhCN"

	webAPI "github.com/maritimusj/centrum/gate/web/api"
	webApp "github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.WithField("src", logStore.SystemLog).Errorln(err)
		}
	}()

	//命令行参数
	config := flag.String("config", "gate.yaml", "config file name")
	logLevel := flag.String("l", "", "log level, [trace,debug,info,warn,error,fatal]")
	resetDefaultUserPassword := flag.Bool("reset", false, "reset default user password")
	flushDB := flag.Bool("flush", false, "erase all data in database")
	langID := flag.Int("lang", lang.EnUS, "lang ID")

	flag.Parse()

	if *langID == lang.ZhCN || *langID == lang.EnUS {
		lang.Active(*langID)
	}

	durafmt.SetAlias("years", lang.Str(lang.Years))
	durafmt.SetAlias("weeks", lang.Str(lang.Weeks))
	durafmt.SetAlias("days", lang.Str(lang.Days))
	durafmt.SetAlias("hours", lang.Str(lang.Hours))
	durafmt.SetAlias("minutes", lang.Str(lang.Minutes))
	durafmt.SetAlias("seconds", lang.Str(lang.Seconds))
	durafmt.SetAlias("milliseconds", lang.Str(lang.Milliseconds))
	durafmt.SetAlias("microseconds", "'")

	viper.SetConfigFile(*config)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("influxdb.url", "http://localhost:8086")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}

	var edges []string
	if viper.IsSet("edges") {
		edges = viper.GetStringSlice("edges")
	} else {
		edges = []string{"http://localhost:1234/rpc"}
	}

	for _, url := range edges {
		edge.Add(url)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := webApp.Start(ctx, *logLevel); err != nil {
		log.Fatal(err)
	}

	if *flushDB {
		code := util.RandStr(4, util.RandNum)
		fmt.Print(lang.Str(lang.ConfirmAdminPassword, code))

		var confirm string
		_, _ = fmt.Scanln(&confirm)
		if confirm != code {
			log.Fatal(lang.Error(lang.ErrConfirmCodeWrong))
		} else {
			err := webApp.FlushDB()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf(lang.Str(lang.FlushDBOk))
				log.WithField("src", logStore.SystemLog).Warningln(lang.Str(lang.FlushDBOk))
				os.Exit(0)
			}
		}
	}

	if *resetDefaultUserPassword {
		err := webApp.ResetDefaultAdminUser()
		if err != nil {
			log.Fatal(err)
		}
		log.WithField("src", logStore.SystemLog).Warnln(lang.Str(lang.DefaultUserPasswordResetOk))
	}

	//API服务

	webAPI.Start(ctx, webApp.Config)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("exit...")

	cancel()

	webAPI.Wait()
	webApp.Close()
}
