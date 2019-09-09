package api

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/hero"
	"gopkg.in/go-playground/validator.v9"

	"github.com/maritimusj/centrum/web/api/device"
	"github.com/maritimusj/centrum/web/api/equipment"
	"github.com/maritimusj/centrum/web/api/group"
	"github.com/maritimusj/centrum/web/api/resource"
	"github.com/maritimusj/centrum/web/api/role"
	"github.com/maritimusj/centrum/web/api/user"
	"github.com/maritimusj/centrum/web/perm"
)

type Server interface {
	Register(values ...interface{})
	Start(ctx context.Context) error
	Close()
}

type server struct {
	app *iris.Application
}

func (server *server) Register(values ...interface{}) {
	hero.Register(values...)
}

func (server *server) Start(ctx context.Context) error {
	hero.Register(validator.New())

	server.app.PartyFunc("/v1/web", func(p router.Party) {
		p.Use(hero.Handler(perm.Check))
		//资源
		p.PartyFunc("/resource", func(p router.Party) {
			p.Get("/", hero.Handler(resource.GroupList))
			p.Get("/{groupID:int}/", hero.Handler(resource.List))
		})

		//角色
		p.PartyFunc("/role", func(p router.Party) {
			p.Get("/", hero.Handler(role.List))
			p.Post("/", hero.Handler(role.Create))
			p.Get("/{id:int64}", hero.Handler(role.Detail))
			p.Put("/{id:int64}", hero.Handler(role.Update))
			p.Delete("/{id:int64}", hero.Handler(role.Delete))
		})

		//用户
		p.PartyFunc("/user", func(p router.Party) {
			p.Get("/", hero.Handler(user.List))
			p.Post("/", hero.Handler(user.Create))
			p.Get("/{id:int64}", hero.Handler(user.Detail))
			p.Put("/{id:int64}", hero.Handler(user.Update))
			p.Delete("/{id:int64}", hero.Handler(user.Delete))
		})

		//设备分组
		p.PartyFunc("/group", func(p router.Party) {
			p.Get("/", hero.Handler(group.List))
			p.Post("/", hero.Handler(group.Create))
			p.Get("/{id:int64}", hero.Handler(group.Detail))
			p.Put("/{id:int64}", hero.Handler(group.Update))
			p.Delete("/{id:int64}", hero.Handler(group.Delete))
		})

		//物理设备
		p.PartyFunc("/device", func(p router.Party) {
			p.Get("/", hero.Handler(device.List))
			p.Post("/", hero.Handler(device.Create))
			p.Get("/{id:int64}", hero.Handler(device.Detail))
			p.Put("/{id:int64}", hero.Handler(device.Update))
			p.Delete("/{id:int64}", hero.Handler(device.Delete))

			//物理点位
			p.Get("/{id:int64}/measure", hero.Handler(device.MeasureList))
		})
		//物理点位
		p.PartyFunc("/measure", func(p router.Party) {
			p.Get("/{id:int64}", hero.Handler(device.MeasureDetail))
		})

		//自定义设备
		p.PartyFunc("/equipment", func(p router.Party) {
			p.Get("/", hero.Handler(equipment.List))
			p.Post("/", hero.Handler(equipment.Create))
			p.Get("/{id:int64}", hero.Handler(equipment.Detail))
			p.Put("/{id:int64}", hero.Handler(equipment.Update))
			p.Delete("/{id:int64}", hero.Handler(equipment.Delete))

			//自定义点位
			p.Get("/{id:int64}/state", hero.Handler(equipment.StateList))
			p.Post("/{id:int64}/state", hero.Handler(equipment.CreateState))
		})
		//自定义点位
		p.PartyFunc("/state", func(p router.Party) {
			p.Get("/{id:int64}", hero.Handler(equipment.StateDetail))
			p.Put("/{id:int64}", hero.Handler(equipment.UpdateState))
			p.Delete("/{id:int64}", hero.Handler(equipment.DeleteState))
		})
	})
	return server.app.Run(iris.Addr(":9090"))
}

func (server *server) Close() {
}

func New() Server {
	return &server{}
}
