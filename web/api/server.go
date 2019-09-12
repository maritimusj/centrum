package api

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/web/api/web"
	"github.com/maritimusj/centrum/web/perm"
	"gopkg.in/go-playground/validator.v9"

	ResourceDef "github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/web/api/device"
	"github.com/maritimusj/centrum/web/api/equipment"
	"github.com/maritimusj/centrum/web/api/group"
	"github.com/maritimusj/centrum/web/api/resource"
	"github.com/maritimusj/centrum/web/api/role"
	"github.com/maritimusj/centrum/web/api/user"
)

type Server interface {
	Register(values ...interface{})
	Start(ctx context.Context, cfg config.Config) error
	Close()
}

type server struct {
	app *iris.Application
}

func (server *server) Register(values ...interface{}) {
	hero.Register(values...)
}

func (server *server) Start(ctx context.Context, cfg config.Config) error {
	hero.Register(validator.New())

	server.app.Post("/v1/web/login/", hero.Handler(web.Login))
	server.app.PartyFunc("/v1/web", func(p router.Party) {
		web.RequireToken(p, cfg)
		p.Use(hero.Handler(perm.Check))
		//资源
		p.PartyFunc("/resource", func(p router.Party) {
			p.Get("/", hero.Handler(resource.GroupList))
			p.Get("/{groupID:int}/", hero.Handler(resource.List))
		})

		//角色
		p.PartyFunc("/role", func(p router.Party) {
			p.Get("/", hero.Handler(role.List)).Name = ResourceDef.RoleList
			p.Post("/", hero.Handler(role.Create)).Name = ResourceDef.RoleCreate
			p.Get("/{id:int64}", hero.Handler(role.Detail)).Name = ResourceDef.RoleDetail
			p.Put("/{id:int64}", hero.Handler(role.Update)).Name = ResourceDef.RoleUpdate
			p.Delete("/{id:int64}", hero.Handler(role.Delete)).Name = ResourceDef.RoleDelete
		})

		//用户
		p.PartyFunc("/user", func(p router.Party) {
			p.Get("/", hero.Handler(user.List)).Name = ResourceDef.UserList
			p.Post("/", hero.Handler(user.Create)).Name = ResourceDef.UserCreate
			p.Get("/{id:int64}", hero.Handler(user.Detail)).Name = ResourceDef.UserDetail
			p.Put("/{id:int64}", hero.Handler(user.Update)).Name = ResourceDef.UserUpdate
			p.Delete("/{id:int64}", hero.Handler(user.Delete)).Name = ResourceDef.UserDelete
		})

		//设备分组
		p.PartyFunc("/group", func(p router.Party) {
			p.Get("/", hero.Handler(group.List)).Name = ResourceDef.GroupList
			p.Post("/", hero.Handler(group.Create)).Name = ResourceDef.GroupCreate
			p.Get("/{id:int64}", hero.Handler(group.Detail)).Name = ResourceDef.GroupDetail
			p.Put("/{id:int64}", hero.Handler(group.Update)).Name = ResourceDef.GroupUpdate
			p.Delete("/{id:int64}", hero.Handler(group.Delete)).Name = ResourceDef.GroupDelete
		})

		//物理设备
		p.PartyFunc("/device", func(p router.Party) {
			p.Get("/", hero.Handler(device.List)).Name = ResourceDef.DeviceList
			p.Post("/", hero.Handler(device.Create)).Name = ResourceDef.DeviceCreate
			p.Get("/{id:int64}", hero.Handler(device.Detail)).Name = ResourceDef.DeviceDetail
			p.Put("/{id:int64}", hero.Handler(device.Update)).Name = ResourceDef.DeviceUpdate
			p.Delete("/{id:int64}", hero.Handler(device.Delete)).Name = ResourceDef.DeviceDelete

			//物理点位
			p.Get("/{id:int64}/measure", hero.Handler(device.MeasureList)).Name = ResourceDef.MeasureList
		})
		//物理点位
		p.PartyFunc("/measure", func(p router.Party) {
			p.Get("/{id:int64}", hero.Handler(device.MeasureDetail)).Name = ResourceDef.MeasureDetail
		})

		//自定义设备
		p.PartyFunc("/equipment", func(p router.Party) {
			p.Get("/", hero.Handler(equipment.List)).Name = ResourceDef.EquipmentList
			p.Post("/", hero.Handler(equipment.Create)).Name = ResourceDef.EquipmentCreate
			p.Get("/{id:int64}", hero.Handler(equipment.Detail)).Name = ResourceDef.EquipmentDetail
			p.Put("/{id:int64}", hero.Handler(equipment.Update)).Name = ResourceDef.EquipmentUpdate
			p.Delete("/{id:int64}", hero.Handler(equipment.Delete)).Name = ResourceDef.EquipmentDelete

			//自定义点位
			p.Get("/{id:int64}/state", hero.Handler(equipment.StateList)).Name = ResourceDef.StateList
			p.Post("/{id:int64}/state", hero.Handler(equipment.CreateState)).Name = ResourceDef.StateCreate
		})
		//自定义点位
		p.PartyFunc("/state", func(p router.Party) {
			p.Get("/{id:int64}", hero.Handler(equipment.StateDetail)).Name = ResourceDef.StateDetail
			p.Put("/{id:int64}", hero.Handler(equipment.UpdateState)).Name = ResourceDef.StateUpdate
			p.Delete("/{id:int64}", hero.Handler(equipment.DeleteState)).Name = ResourceDef.StateDelete
		})
	})
	return server.app.Run(iris.Addr(":9090"))
}

func (server *server) Close() {
}

func New() Server {
	return &server{
		app: iris.Default(),
	}
}
