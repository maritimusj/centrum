package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"github.com/maritimusj/centrum/lang"
	_ "github.com/maritimusj/centrum/lang/zhCN"
	"github.com/maritimusj/centrum/util"
	webAPI "github.com/maritimusj/centrum/web/api"
	webApp "github.com/maritimusj/centrum/web/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	//命令行参数
	logLevel := flag.String("l", "", "log level, [trace,debug,info,warn,error,fatal]")
	resetDefaultUserPassword := flag.Bool("reset", false, "reset default user password")
	flushDB := flag.Bool("flush", false, "erase all data in database")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := webApp.Init(ctx, *logLevel); err != nil {
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
				os.Exit(0)
			}
		}
	}

	if *resetDefaultUserPassword {
		err := webApp.ResetDefaultAdminUser()
		if err != nil {
			log.Fatal(err)
		}
		log.Warnln(lang.Str(lang.DefaultUserPasswordResetOk))
	}

	//API服务
	apiServer := webAPI.New()
	apiServer.Start(ctx, webApp.Config)

	WaitForExit()

	fmt.Println("exit...")

	cancel()

	apiServer.Wait()
	webApp.Close()
}

func WaitForExit() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
