package main

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/maritimusj/centrum/edge/devices/InverseServer"
	"github.com/maritimusj/centrum/edge/devices/event"
	_ "github.com/maritimusj/centrum/edge/lang/zhCN"

	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/maritimusj/centrum/edge/devices"
	"github.com/maritimusj/centrum/json_rpc"
	log "github.com/sirupsen/logrus"
)

func main() {
	addr := flag.String("addr", "", "service addr")
	port := flag.Int("port", 1234, "service port")
	inversePort := flag.Int("i", 10502, "inverse server port")
	level := flag.String("l", "error", "log level")

	flag.Parse()

	l, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(l)

	//初始化event管理
	event.Init(context.Background())

	//初始化inverse Server
	err = InverseServer.Start(context.Background(), *addr, *inversePort)
	if err != nil {
		log.Fatal(err)
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
		log.Println("edge service listen on port ", *port)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *addr, *port), r); err != nil {
			log.Fatalf("error serving: %s", err)
		}
	}()

	const pidFile = "./edge.pid"
	pid := fmt.Sprintf("%d", os.Getpid())
	err = ioutil.WriteFile(pidFile, []byte(pid), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	println("exiting...")

	_ = os.Remove(pidFile)

	InverseServer.Close()
	runner.Close()
}
