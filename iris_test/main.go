package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"log"
)

func main() {
	app := iris.New()
	app.Post("/", func(ctx iris.Context) {
		var form struct {
			Name *string `json:"name"`
		}
		err := ctx.ReadJSON(&form)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%#v", form)
		}
	})

	app.Get("/{id:int}/{id2:int64}", hero.Handler(func(first int, second int64, c iris.Context) {
		c.WriteString(fmt.Sprintf("first:%d, second:%d", first, second))
	}))

	log.Fatal(app.Run(iris.Addr(":8080")))
}
