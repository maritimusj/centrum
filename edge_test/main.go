package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/rpc/json"
	"github.com/kr/pretty"
	. "github.com/maritimusj/centrum/json_rpc"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	url = "http://localhost:1234/rpc"
)

func invoke(cmd string, request interface{}) (*Result, error) {
	message, err := json.EncodeClientRequest(cmd, request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(message))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var reply Result
	err = json.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func main() {
	addr := flag.String("addr", "", "ip addr and port")
	uid := flag.String("uid", "", "uid of device")
	r := flag.Bool("r", false, "get realtime data")
	c := flag.String("c", "", "get channel data by tag name")
	u := flag.Int64("u", 0, "update di channel")
	s := flag.Bool("s", false, "stop device")

	flag.Parse()

	if *addr != "" {
		result, err := invoke("Edge.Active", Conf{
			UID:              *uid,
			Inverse:          false,
			Address:          *addr,
			Interval:         6 * time.Second,
			DB:               "gsd",
			InfluxDBAddress:  "http://localhost:8086",
			InfluxDBUserName: "",
			InfluxDBPassword: "",
		})

		if err != nil {
			log.Fatal("active failed:", err)
		}
		fmt.Printf("%# v", pretty.Formatter(result))
	}

	if *s {
		result, err := invoke("Edge.Remove", *uid)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%# v", pretty.Formatter(result))
		os.Exit(0)
	}


	if *c != "" {
		result, err := invoke("Edge.GetValue", &CH{
			UID: *uid,
			Tag: *c,
		})
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%# v", pretty.Formatter(result))
		os.Exit(0)
	}

	if *u != 0 {
		result, err := invoke("Edge.GetValue", &CH{
			UID: *uid,
			Tag: "DO-" + strconv.FormatInt(*u, 10),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%# v", pretty.Formatter(result))
		data := result.Data.(map[string]interface{})
		v := data["value"].(bool)
		tag := data["tag"].(string)
		result, err = invoke("Edge.SetValue", &Value{
			CH: CH{
				UID: *uid,
				Tag: tag,
			},
			V: !v,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%# v", pretty.Formatter(result))
		os.Exit(0)
	}

	if *r {
		result, err := invoke("Edge.GetRealtimeData", *uid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%# v", pretty.Formatter(result))
		os.Exit(0)
	}
}
