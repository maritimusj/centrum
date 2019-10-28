package api

import (
	"context"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/event"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/web/api/alarm"
	"github.com/maritimusj/centrum/web/api/edge"
	logStore "github.com/maritimusj/centrum/web/api/log"
	"github.com/maritimusj/centrum/web/api/my"
	"github.com/maritimusj/centrum/web/api/organization"
	"github.com/maritimusj/centrum/web/api/role"
	"github.com/maritimusj/centrum/web/api/statistics"
	"github.com/maritimusj/centrum/web/api/web"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/perm"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/maritimusj/centrum/web/api/device"
	"github.com/maritimusj/centrum/web/api/equipment"
	"github.com/maritimusj/centrum/web/api/group"
	"github.com/maritimusj/centrum/web/api/resource"
	"github.com/maritimusj/centrum/web/api/user"

	configAPI "github.com/maritimusj/centrum/web/api/config"
	ResourceDef "github.com/maritimusj/centrum/web/resource"
)

type Server interface {
	Register(values ...interface{})
	Start(ctx context.Context, cfg *config.Config)
	Wait()
}

type server struct {
	app *iris.Application
	wg  sync.WaitGroup
}

func New() Server {
	return &server{
		app: iris.Default(),
	}
}

func (server *server) Register(values ...interface{}) {
	hero.Register(values...)
}

func (server *server) Wait() {
	server.wg.Wait()
}

func (server *server) Start(ctx context.Context, cfg *config.Config) {
	server.app.Logger().SetLevel(cfg.LogLevel())

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD"},
		AllowCredentials: true,
	})

	//后台
	server.app.StaticWeb("/", "./public")

	//v1
	v1 := server.app.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	v1.PartyFunc("/web", func(p router.Party) {
		p.Post("/login", hero.Handler(web.Login))

		p.PartyFunc("/edge", func(p router.Party) {
			global.Params.Set("callbackURL", fmt.Sprintf("http://localhost:%d%s", app.Config.APIPort(), p.GetRelPath()))
			p.Post("/{id:int64}", hero.Handler(edge.Feedback))
		})

		p.PartyFunc("/", func(p router.Party) {
			web.RequireToken(p)
			p.Use(hero.Handler(perm.CheckApiPerm))

			p.PartyFunc("/config", func(p router.Party) {
				p.Get("/base", hero.Handler(configAPI.Base)).Name = ResourceDef.ConfigBaseDetail
				p.Put("/base", hero.Handler(configAPI.UpdateBase)).Name = ResourceDef.ConfigBaseUpdate
			})

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
				p.Get("/{groupID:int}/", hero.Handler(resource.List)).Name = ResourceDef.ResourceDetail
				p.Get("/", hero.Handler(resource.GroupList)).Name = ResourceDef.ResourceList
				p.Get("/search/", hero.Handler(resource.GetList)).Name = ResourceDef.ResourceDetail
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

				//实时状态
				p.Get("/{id:int64}/reset", hero.Handler(device.Reset)).Name = ResourceDef.DeviceStatus
				p.Get("/{id:int64}/status", hero.Handler(device.Status)).Name = ResourceDef.DeviceStatus
				p.Get("/{id:int64}/data", hero.Handler(device.Data)).Name = ResourceDef.DeviceData
				p.Put("/{id:int64}/{tagName:string}", hero.Handler(device.Ctrl)).Name = ResourceDef.DeviceCtrl
				p.Get("/{id:int64}/{tagName:string}", hero.Handler(device.GetCHValue)).Name = ResourceDef.DeviceCHValue

				//导出报表
				p.Post("/export", hero.Handler(statistics.Export))
			})
			//物理点位
			p.PartyFunc("/measure", func(p router.Party) {
				p.Delete("/{id:int64}", hero.Handler(device.DeleteMeasure)).Name = ResourceDef.MeasureDelete
				p.Get("/{id:int64}", hero.Handler(device.MeasureDetail)).Name = ResourceDef.MeasureDetail

				//历史趋势
				p.Post("/{id:int64}/statistics", hero.Handler(statistics.Measure))
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

				//实时状态
				p.Get("/{id:int64}/status", hero.Handler(equipment.Status)).Name = ResourceDef.EquipmentStatus
				p.Get("/{id:int64}/data", hero.Handler(equipment.Data)).Name = ResourceDef.EquipmentData
				p.Put("/{id:int64}/{stateID:int64}", hero.Handler(equipment.Ctrl)).Name = ResourceDef.EquipmentCtrl
				p.Get("/{id:int64}/{stateID:int64}", hero.Handler(equipment.GetCHValue)).Name = ResourceDef.EquipmentCHValue
			})

			//自定义点位
			p.PartyFunc("/state", func(p router.Party) {
				p.Get("/{id:int64}", hero.Handler(equipment.StateDetail)).Name = ResourceDef.StateDetail
				p.Put("/{id:int64}", hero.Handler(equipment.UpdateState)).Name = ResourceDef.StateUpdate
				p.Delete("/{id:int64}", hero.Handler(equipment.DeleteState)).Name = ResourceDef.StateDelete

				//历史趋势
				p.Post("/{id:int64}/statistics", hero.Handler(statistics.State))
			})

			//警报
			p.PartyFunc("/alarm", func(p router.Party) {
				p.Get("/", hero.Handler(alarm.List)).Name = ResourceDef.AlarmList
				p.Put("/{id:int64}", hero.Handler(alarm.Confirm)).Name = ResourceDef.AlarmConfirm
				p.Get("/{id:int64}", hero.Handler(alarm.Detail)).Name = ResourceDef.AlarmDetail
				p.Delete("/{id:int64}", hero.Handler(alarm.Delete)).Name = ResourceDef.AlarmDelete
			})

			//日志等级
			p.Get("/log/level", hero.Handler(logStore.Level)).Name = ResourceDef.LogLevelList
			//系统日志
			p.PartyFunc("/syslog", func(p router.Party) {
				p.Get("/", hero.Handler(logStore.List)).Name = ResourceDef.LogList
				p.Delete("/", hero.Handler(logStore.Delete)).Name = ResourceDef.LogDelete
			})
		})
	})

	addr := fmt.Sprintf("%s:%d", app.Config.APIAddr(), app.Config.APIPort())
	server.wg.Add(2)
	go func() {
		defer server.wg.Done()

		log.Trace("api server start at: ", addr)
		err := server.app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed))
		if err != nil {
			log.Error("listen: %s\n", err)
		}
	}()

	go func() {
		defer server.wg.Done()

		select {
		case <-ctx.Done():
			timeout, _ := context.WithTimeout(ctx, 6*time.Second)
			err := server.app.Shutdown(timeout)
			if err != nil {
				log.Tracef("shutdown http server: ", err)
			} else {
				log.Tracef("http server shutdown.")
			}
		}
	}()

	time.AfterFunc(3*time.Second, func() {
		app.Event.Publish(event.ApiServerStarted)
	})
}
