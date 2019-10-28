package api

import (
	"context"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/config"
	"github.com/maritimusj/centrum/gate/event"
	alarm2 "github.com/maritimusj/centrum/gate/web/api/alarm"
	config2 "github.com/maritimusj/centrum/gate/web/api/config"
	device2 "github.com/maritimusj/centrum/gate/web/api/device"
	edge2 "github.com/maritimusj/centrum/gate/web/api/edge"
	equipment2 "github.com/maritimusj/centrum/gate/web/api/equipment"
	group2 "github.com/maritimusj/centrum/gate/web/api/group"
	log2 "github.com/maritimusj/centrum/gate/web/api/log"
	my2 "github.com/maritimusj/centrum/gate/web/api/my"
	organization2 "github.com/maritimusj/centrum/gate/web/api/organization"
	resource3 "github.com/maritimusj/centrum/gate/web/api/resource"
	role2 "github.com/maritimusj/centrum/gate/web/api/role"
	statistics2 "github.com/maritimusj/centrum/gate/web/api/statistics"
	user2 "github.com/maritimusj/centrum/gate/web/api/user"
	web2 "github.com/maritimusj/centrum/gate/web/api/web"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	perm2 "github.com/maritimusj/centrum/gate/web/perm"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/global"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
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
		p.Post("/login", hero.Handler(web2.Login))

		p.PartyFunc("/edge", func(p router.Party) {
			global.Params.Set("callbackURL", fmt.Sprintf("http://localhost:%d%s", app2.Config.APIPort(), p.GetRelPath()))
			p.Post("/{id:int64}", hero.Handler(edge2.Feedback))
		})

		p.PartyFunc("/", func(p router.Party) {
			web2.RequireToken(p)
			p.Use(hero.Handler(perm2.CheckApiPerm))

			p.PartyFunc("/config", func(p router.Party) {
				p.Get("/base", hero.Handler(config2.Base)).Name = resource2.ConfigBaseDetail
				p.Put("/base", hero.Handler(config2.UpdateBase)).Name = resource2.ConfigBaseUpdate
			})

			//我的
			p.PartyFunc("/my", func(p router.Party) {
				p.Get("/profile", hero.Handler(my2.Detail)).Name = resource2.MyProfileDetail
				p.Put("/profile", hero.Handler(my2.Update)).Name = resource2.MyProfileUpdate

				//请求当前用户对于某个资源的权限情况
				p.Get("/perm/{class:string}", hero.Handler(my2.Perm)).Name = resource2.MyPerm
				p.Post("/perm/{class:string}", hero.Handler(my2.MultiPerm)).Name = resource2.MyPermMulti
			})

			//资源
			p.PartyFunc("/resource", func(p router.Party) {
				p.Get("/{groupID:int}/", hero.Handler(resource3.List)).Name = resource2.ResourceDetail
				p.Get("/", hero.Handler(resource3.GroupList)).Name = resource2.ResourceList
				p.Get("/search/", hero.Handler(resource3.GetList)).Name = resource2.ResourceDetail
			})

			//组织
			p.PartyFunc("/org", func(p router.Party) {
				p.Get("/", hero.Handler(organization2.List)).Name = resource2.OrganizationList
				p.Post("/", hero.Handler(organization2.Create)).Name = resource2.OrganizationCreate
				p.Get("/{id:int64}", hero.Handler(organization2.Detail)).Name = resource2.OrganizationDetail
				p.Put("/{id:int64}", hero.Handler(organization2.Update)).Name = resource2.OrganizationUpdate
				p.Delete("/{id:int64}", hero.Handler(organization2.Delete)).Name = resource2.OrganizationDelete
			})

			//角色
			p.PartyFunc("/role", func(p router.Party) {
				p.Get("/", hero.Handler(role2.List)).Name = resource2.RoleList
				p.Post("/", hero.Handler(role2.Create)).Name = resource2.RoleCreate
				p.Get("/{id:int64}", hero.Handler(role2.Detail)).Name = resource2.RoleDetail
				p.Put("/{id:int64}", hero.Handler(role2.Update)).Name = resource2.RoleUpdate
				p.Delete("/{id:int64}", hero.Handler(role2.Delete)).Name = resource2.RoleDelete
			})

			//用户
			p.PartyFunc("/user", func(p router.Party) {
				p.Get("/", hero.Handler(user2.List)).Name = resource2.UserList
				p.Post("/", hero.Handler(user2.Create)).Name = resource2.UserCreate
				p.Get("/{id:int64}", hero.Handler(user2.Detail)).Name = resource2.UserDetail
				p.Put("/{id:int64}", hero.Handler(user2.Update)).Name = resource2.UserUpdate
				p.Delete("/{id:int64}", hero.Handler(user2.Delete)).Name = resource2.UserDelete

				//给用户分配权限（通过用户同名角色）
				p.Put("/{id:int64}/perm", hero.Handler(user2.UpdatePerm)).Name = resource2.RoleUpdate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(user2.LogList)).Name = resource2.UserLogList
				p.Delete("/{id:int64}/log", hero.Handler(user2.LogDelete)).Name = resource2.UserLogDelete
			})

			//设备分组
			p.PartyFunc("/group", func(p router.Party) {
				p.Get("/", hero.Handler(group2.List)).Name = resource2.GroupList
				p.Post("/", hero.Handler(group2.Create)).Name = resource2.GroupCreate
				p.Get("/{id:int64}", hero.Handler(group2.Detail)).Name = resource2.GroupDetail
				p.Put("/{id:int64}", hero.Handler(group2.Update)).Name = resource2.GroupUpdate
				p.Delete("/{id:int64}", hero.Handler(group2.Delete)).Name = resource2.GroupDelete
			})

			//物理设备
			p.PartyFunc("/device", func(p router.Party) {
				p.Get("/", hero.Handler(device2.List)).Name = resource2.DeviceList
				p.Post("/", hero.Handler(device2.Create)).Name = resource2.DeviceCreate
				p.Get("/{id:int64}", hero.Handler(device2.Detail)).Name = resource2.DeviceDetail
				p.Put("/{id:int64}", hero.Handler(device2.Update)).Name = resource2.DeviceUpdate
				p.Delete("/{id:int64}", hero.Handler(device2.Delete)).Name = resource2.DeviceDelete

				//物理点位
				p.Get("/{id:int64}/measure", hero.Handler(device2.MeasureList)).Name = resource2.MeasureList
				p.Post("/{id:int64}/measure", hero.Handler(device2.CreateMeasure)).Name = resource2.MeasureCreate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(device2.LogList)).Name = resource2.DeviceLogList
				p.Delete("/{id:int64}/log", hero.Handler(device2.LogDelete)).Name = resource2.DeviceLogDelete

				//实时状态
				p.Get("/{id:int64}/reset", hero.Handler(device2.Reset)).Name = resource2.DeviceStatus
				p.Get("/{id:int64}/status", hero.Handler(device2.Status)).Name = resource2.DeviceStatus
				p.Get("/{id:int64}/data", hero.Handler(device2.Data)).Name = resource2.DeviceData
				p.Put("/{id:int64}/{tagName:string}", hero.Handler(device2.Ctrl)).Name = resource2.DeviceCtrl
				p.Get("/{id:int64}/{tagName:string}", hero.Handler(device2.GetCHValue)).Name = resource2.DeviceCHValue

				//导出报表
				p.Post("/export", hero.Handler(statistics2.Export))
			})
			//物理点位
			p.PartyFunc("/measure", func(p router.Party) {
				p.Delete("/{id:int64}", hero.Handler(device2.DeleteMeasure)).Name = resource2.MeasureDelete
				p.Get("/{id:int64}", hero.Handler(device2.MeasureDetail)).Name = resource2.MeasureDetail

				//历史趋势
				p.Post("/{id:int64}/statistics", hero.Handler(statistics2.Measure))
			})

			//自定义设备
			p.PartyFunc("/equipment", func(p router.Party) {
				p.Get("/", hero.Handler(equipment2.List)).Name = resource2.EquipmentList
				p.Post("/", hero.Handler(equipment2.Create)).Name = resource2.EquipmentCreate
				p.Get("/{id:int64}", hero.Handler(equipment2.Detail)).Name = resource2.EquipmentDetail
				p.Put("/{id:int64}", hero.Handler(equipment2.Update)).Name = resource2.EquipmentUpdate
				p.Delete("/{id:int64}", hero.Handler(equipment2.Delete)).Name = resource2.EquipmentDelete

				//自定义点位
				p.Get("/{id:int64}/state", hero.Handler(equipment2.StateList)).Name = resource2.StateList
				p.Post("/{id:int64}/state", hero.Handler(equipment2.CreateState)).Name = resource2.StateCreate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(equipment2.LogList)).Name = resource2.EquipmentLogList
				p.Delete("/{id:int64}/log", hero.Handler(equipment2.LogDelete)).Name = resource2.EquipmentLogDelete

				//实时状态
				p.Get("/{id:int64}/status", hero.Handler(equipment2.Status)).Name = resource2.EquipmentStatus
				p.Get("/{id:int64}/data", hero.Handler(equipment2.Data)).Name = resource2.EquipmentData
				p.Put("/{id:int64}/{stateID:int64}", hero.Handler(equipment2.Ctrl)).Name = resource2.EquipmentCtrl
				p.Get("/{id:int64}/{stateID:int64}", hero.Handler(equipment2.GetCHValue)).Name = resource2.EquipmentCHValue
			})

			//自定义点位
			p.PartyFunc("/state", func(p router.Party) {
				p.Get("/{id:int64}", hero.Handler(equipment2.StateDetail)).Name = resource2.StateDetail
				p.Put("/{id:int64}", hero.Handler(equipment2.UpdateState)).Name = resource2.StateUpdate
				p.Delete("/{id:int64}", hero.Handler(equipment2.DeleteState)).Name = resource2.StateDelete

				//历史趋势
				p.Post("/{id:int64}/statistics", hero.Handler(statistics2.State))
			})

			//警报
			p.PartyFunc("/alarm", func(p router.Party) {
				p.Get("/", hero.Handler(alarm2.List)).Name = resource2.AlarmList
				p.Put("/{id:int64}", hero.Handler(alarm2.Confirm)).Name = resource2.AlarmConfirm
				p.Get("/{id:int64}", hero.Handler(alarm2.Detail)).Name = resource2.AlarmDetail
				p.Delete("/{id:int64}", hero.Handler(alarm2.Delete)).Name = resource2.AlarmDelete
			})

			//日志等级
			p.Get("/log/level", hero.Handler(log2.Level)).Name = resource2.LogLevelList
			//系统日志
			p.PartyFunc("/syslog", func(p router.Party) {
				p.Get("/", hero.Handler(log2.List)).Name = resource2.LogList
				p.Delete("/", hero.Handler(log2.Delete)).Name = resource2.LogDelete
			})
		})
	})

	addr := fmt.Sprintf("%s:%d", app2.Config.APIAddr(), app2.Config.APIPort())
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
		app2.Event.Publish(event.ApiServerStarted)
	})
}
