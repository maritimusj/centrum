package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"

	"github.com/maritimusj/centrum/gate/lang"
	_ "github.com/maritimusj/centrum/gate/lang/zhCN"

	webAPI "github.com/maritimusj/centrum/gate/web/api"
	webApp "github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/util"
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

	webAPI.Start(ctx, webApp.Config)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("exit...")

	cancel()

	webAPI.Wait()
	webApp.Close()
}
