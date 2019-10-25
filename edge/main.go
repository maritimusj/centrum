package main

import (
	_ "github.com/maritimusj/centrum/edge/lang/zhCN"

	"context"
	"github.com/maritimusj/centrum/edge/devices/InverseServer"

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

	err = InverseServer.Start(context.Background(), *addr, *inversePort)
	if err != nil {
		log.Fatal(err)
	}

	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	edge := json_rpc.New(devices.New())
	err = server.RegisterService(edge, "")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/rpc", server)

	log.Println("edge service listen on port ", *port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *addr, *port), r); err != nil {
		log.Fatalf("error serving: %s", err)
	}
}
