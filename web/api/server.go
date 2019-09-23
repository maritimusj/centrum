package api

import (
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/app"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/web/api/log"
	"github.com/maritimusj/centrum/web/api/my"
	"github.com/maritimusj/centrum/web/api/organization"
	"github.com/maritimusj/centrum/web/api/role"
	"github.com/maritimusj/centrum/web/api/web"
	"github.com/maritimusj/centrum/web/perm"
	"gopkg.in/go-playground/validator.v9"

	ResourceDef "github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/web/api/device"
	"github.com/maritimusj/centrum/web/api/equipment"
	"github.com/maritimusj/centrum/web/api/group"
	"github.com/maritimusj/centrum/web/api/resource"
	"github.com/maritimusj/centrum/web/api/user"
)

type Server interface {
	Register(values ...interface{})
	Start(cfg config.Config) error
	Close()
}

type server struct {
	app *iris.Application
}

func (server *server) Register(values ...interface{}) {
	hero.Register(values...)
}

func (server *server) Start(cfg config.Config) error {
	hero.Register(validator.New())
	server.app.Logger().SetLevel(cfg.LogLevel())

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD"},
		AllowCredentials: true,
	})

	v1 := server.app.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	v1.PartyFunc("/web", func(p router.Party) {
		p.Post("/login", hero.Handler(web.Login))

		p.PartyFunc("/", func(p router.Party) {
			web.RequireToken(p)
			p.Use(hero.Handler(perm.CheckApiPerm))

			//我的
			p.PartyFunc("/my", func(p router.Party) {
				p.Get("/profile", hero.Handler(my.Detail)).Name = ResourceDef.MyProfileDetail
				p.Put("/profile", hero.Handler(my.Update)).Name = ResourceDef.MyProfileUpdate

				//请求当前用户对于某个资源的权限情况
				p.Get("/perm/{class:string}", hero.Handler(my.Perm)).Name = ResourceDef.MyPerm
				p.Post("/perm/{class:string}", hero.Handler(my.MultiPerm)).Name = ResourceDef.MyPermMulti
			})

			//资源
			p.PartyFunc("/resource", func(p router.Party) {
				p.Get("/", hero.Handler(resource.GroupList)).Name = ResourceDef.ResourceList
				p.Get("/{groupID:int}/", hero.Handler(resource.List)).Name = ResourceDef.ResourceDetail
			})

			//组织
			p.PartyFunc("/org", func(p router.Party) {
				p.Get("/", hero.Handler(organization.List)).Name = ResourceDef.OrganizationList
				p.Post("/", hero.Handler(organization.Create)).Name = ResourceDef.OrganizationCreate
				p.Get("/{id:int64}", hero.Handler(organization.Detail)).Name = ResourceDef.OrganizationDetail
				p.Put("/{id:int64}", hero.Handler(organization.Update)).Name = ResourceDef.OrganizationUpdate
				p.Delete("/{id:int64}", hero.Handler(organization.Delete)).Name = ResourceDef.OrganizationDelete
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

				//给用户分配权限（通过用户同名角色）
				p.Put("/{id:int64}/perm", hero.Handler(user.UpdatePerm)).Name = ResourceDef.RoleUpdate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(user.LogList)).Name = ResourceDef.UserLogList
				p.Delete("/{id:int64}/log", hero.Handler(user.LogDelete)).Name = ResourceDef.UserLogDelete
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
				p.Post("/{id:int64}/measure", hero.Handler(device.CreateMeasure)).Name = ResourceDef.MeasureCreate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(device.LogList)).Name = ResourceDef.DeviceLogList
				p.Delete("/{id:int64}/log", hero.Handler(device.LogDelete)).Name = ResourceDef.DeviceLogDelete
			})
			//物理点位
			p.PartyFunc("/measure", func(p router.Party) {
				p.Delete("/{id:int64}", hero.Handler(device.DeleteMeasure)).Name = ResourceDef.MeasureDelete
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

				//日志
				p.Get("/{id:int64}/log", hero.Handler(equipment.LogList)).Name = ResourceDef.EquipmentLogList
				p.Delete("/{id:int64}/log", hero.Handler(equipment.LogDelete)).Name = ResourceDef.EquipmentLogDelete
			})

			//自定义点位
			p.PartyFunc("/state", func(p router.Party) {
				p.Get("/{id:int64}", hero.Handler(equipment.StateDetail)).Name = ResourceDef.StateDetail
				p.Put("/{id:int64}", hero.Handler(equipment.UpdateState)).Name = ResourceDef.StateUpdate
				p.Delete("/{id:int64}", hero.Handler(equipment.DeleteState)).Name = ResourceDef.StateDelete
			})

			//日志等级
			p.Get("/log/level", hero.Handler(log.Level)).Name = ResourceDef.LogLevelList
			//系统日志
			p.PartyFunc("/syslog", func(p router.Party) {
				p.Get("/", hero.Handler(log.List)).Name = ResourceDef.LogList
				p.Delete("/", hero.Handler(log.Delete)).Name = ResourceDef.LogDelete
			})
		})
	})

	addr := fmt.Sprintf("%s:%d", app.Config.APIAddr(), app.Config.APIPort())
	return server.app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed))
}

func (server *server) Close() {
}

func New() Server {
	return &server{
		app: iris.Default(),
	}
}
