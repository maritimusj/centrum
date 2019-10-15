package main

import (
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

	flag.Parse()

	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	edge := json_rpc.New(devices.New())
	err := server.RegisterService(edge, "")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/rpc", server)

	log.Println("JSON RPC service listen and serving on port ", *port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *addr, *port), r); err != nil {
		log.Fatalf("Error serving: %s", err)
	}
}
