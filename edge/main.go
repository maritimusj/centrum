package main

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/maritimusj/centrum/edge/devices/InverseServer"
	"github.com/maritimusj/centrum/edge/devices/event"

	"github.com/maritimusj/centrum/edge/lang"
	_ "github.com/maritimusj/centrum/edge/lang/enUS"
	_ "github.com/maritimusj/centrum/edge/lang/zhCN"

	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	"github.com/maritimusj/centrum/edge/devices"
	"github.com/maritimusj/centrum/json_rpc"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {
	logLevel := flag.String("l", "", "log error level")
	config := flag.String("config", "edge.yaml", "config file name")
	langID := flag.Int("lang", lang.EnUS, "language ID")

	flag.Parse()

	if *langID == lang.ZhCN || *langID == lang.EnUS {
		lang.Active(*langID)
	}

	viper.SetConfigFile(*config)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetDefault("edge.addr", "")
	viper.SetDefault("edge.port", 1234)

	//默认开启inverse server
	viper.SetDefault("inverse.enable", true)
	viper.SetDefault("inverse.addr", "")
	viper.SetDefault("inverse.port", 10502)

	viper.SetDefault("error.level", "error")

	var l log.Level
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}

	if *logLevel == "" {
		*logLevel = viper.GetString("error.level")
	}

	l, err = log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(l)

	//初始化event管理
	event.Init(context.Background())

	var (
		inverseEnable = viper.GetBool("inverse.enable")
		inverseAddr   = viper.GetString("inverse.addr")
		inversePort   = viper.GetInt("inverse.port")
	)
	if inverseEnable {
		//初始化inverse Server
		err = InverseServer.Start(context.Background(), inverseAddr, inversePort)
		if err != nil {
			log.Fatal(err)
		}
	}

	//初始化rpc服务
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	//初始化runner
	runner := devices.New()
	edge := json_rpc.New(runner)
	err = server.RegisterService(edge, "")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/rpc", server)

	go func() {
		var (
			addr = viper.GetString("edge.addr")
			port = viper.GetInt("edge.port")
		)

		log.Println("edge service listen on port: ", port)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), r); err != nil {
			log.Fatalf("error serving: %s", err)
		}
	}()

	pidFile := viper.GetString("pid.file")
	if pidFile != "" {
		pid := fmt.Sprintf("%d", os.Getpid())
		err = ioutil.WriteFile(pidFile, []byte(pid), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	var (
		pprofEnable = viper.GetBool("pprof.enable")
		addr        = viper.GetString("pprof.addr")
		port        = viper.GetInt("pprof.port")
	)
	if pprofEnable {
		go func() {
			_ = http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
		}()
	}

	quit := make(chan os.Signal)
	runner.RestartMainFN = func() {
		quit <- syscall.SIGINT
	}

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	println("exiting...")

	_ = os.Remove(pidFile)

	if inverseEnable {
		InverseServer.Close()
	}

	runner.Close()
}
