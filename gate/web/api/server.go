package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/hero"
	cfg "github.com/maritimusj/centrum/gate/config"
	"github.com/maritimusj/centrum/gate/event"
	"github.com/maritimusj/centrum/gate/web/api/alarm"
	"github.com/maritimusj/centrum/gate/web/api/brief"
	"github.com/maritimusj/centrum/gate/web/api/comment"
	"github.com/maritimusj/centrum/gate/web/api/config"
	"github.com/maritimusj/centrum/gate/web/api/device"
	"github.com/maritimusj/centrum/gate/web/api/edge"
	"github.com/maritimusj/centrum/gate/web/api/equipment"
	"github.com/maritimusj/centrum/gate/web/api/group"
	logStore "github.com/maritimusj/centrum/gate/web/api/log"
	"github.com/maritimusj/centrum/gate/web/api/my"
	"github.com/maritimusj/centrum/gate/web/api/organization"
	"github.com/maritimusj/centrum/gate/web/api/resource"
	"github.com/maritimusj/centrum/gate/web/api/role"
	"github.com/maritimusj/centrum/gate/web/api/statistics"
	"github.com/maritimusj/centrum/gate/web/api/user"
	"github.com/maritimusj/centrum/gate/web/api/web"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/perm"
	resourceDef "github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/global"

	log "github.com/sirupsen/logrus"
)

var (
	defaultAPIServer = New()
)

func Start(ctx context.Context, webDir string, cfg *cfg.Config) {
	defaultAPIServer.Start(ctx, webDir, cfg)
}

func Wait() {
	defaultAPIServer.Wait()
}

type Server interface {
	Register(values ...interface{})
	Start(ctx context.Context, webDir string, cfg *cfg.Config)
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

func (server *server) Start(ctx context.Context, webDir string, cfg *cfg.Config) {
	//server.app.Logger().SetLevel(cfg.LogLevel())
	server.app.Logger().SetOutput(ioutil.Discard)

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD"},
		AllowCredentials: true,
	})

	//后台
	server.app.StaticWeb("/", webDir)

	//v1
	v1 := server.app.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	v1.PartyFunc("/reg", func(p router.Party) {
		p.Get("/", hero.Handler(web.GetReg))
		p.Post("/", hero.Handler(web.Reg))
	})

	v1.PartyFunc("/web", func(p router.Party) {
		p.Post("/login", hero.Handler(web.Login))

		//edge 回调
		p.PartyFunc("/edge", func(p router.Party) {
			_ = global.Params.Set("callbackURL", fmt.Sprintf("http://localhost:%d%s", app.Config.APIPort(), p.GetRelPath()))
			p.Post("/{id:int64}", hero.Handler(edge.Feedback))
		})

		p.PartyFunc("/", func(p router.Party) {
			web.RequireToken(p)
			p.Use(hero.Handler(perm.CheckApiPerm))

			p.PartyFunc("/config", func(p router.Party) {
				p.Get("/base", hero.Handler(config.Base)).Name = resourceDef.ConfigBaseDetail
				p.Put("/base", hero.Handler(config.UpdateBase)).Name = resourceDef.ConfigBaseUpdate
			})

			//我的
			p.PartyFunc("/my", func(p router.Party) {
				p.Get("/profile", hero.Handler(my.Detail)).Name = resourceDef.MyProfileDetail
				p.Put("/profile", hero.Handler(my.Update)).Name = resourceDef.MyProfileUpdate

				//请求当前用户对于某个资源的权限情况
				p.Get("/perm/{class:string}", hero.Handler(my.Perm)).Name = resourceDef.MyPerm
				p.Post("/perm/{class:string}", hero.Handler(my.MultiPerm)).Name = resourceDef.MyPermMulti
			})

			//系统简讯
			p.PartyFunc("/brief", func(p router.Party) {
				p.Get("/", hero.Handler(brief.Simple)).Name = resourceDef.SysBrief
			})

			//资源
			p.PartyFunc("/resource", func(p router.Party) {
				p.Get("/{groupID:int}/", hero.Handler(resource.List)).Name = resourceDef.ResourceDetail
				p.Get("/", hero.Handler(resource.GroupList)).Name = resourceDef.ResourceList
				p.Get("/search/", hero.Handler(resource.GetList)).Name = resourceDef.ResourceDetail
			})

			//组织
			p.PartyFunc("/org", func(p router.Party) {
				p.Get("/", hero.Handler(organization.List)).Name = resourceDef.OrganizationList
				p.Post("/", hero.Handler(organization.Create)).Name = resourceDef.OrganizationCreate
				p.Get("/{id:int64}", hero.Handler(organization.Detail)).Name = resourceDef.OrganizationDetail
				p.Put("/{id:int64}", hero.Handler(organization.Update)).Name = resourceDef.OrganizationUpdate
				p.Delete("/{id:int64}", hero.Handler(organization.Delete)).Name = resourceDef.OrganizationDelete
			})

			//角色
			p.PartyFunc("/role", func(p router.Party) {
				p.Get("/", hero.Handler(role.List)).Name = resourceDef.RoleList
				p.Post("/", hero.Handler(role.Create)).Name = resourceDef.RoleCreate
				p.Get("/{id:int64}", hero.Handler(role.Detail)).Name = resourceDef.RoleDetail
				p.Put("/{id:int64}", hero.Handler(role.Update)).Name = resourceDef.RoleUpdate
				p.Delete("/{id:int64}", hero.Handler(role.Delete)).Name = resourceDef.RoleDelete
			})

			//用户
			p.PartyFunc("/user", func(p router.Party) {
				p.Get("/", hero.Handler(user.List)).Name = resourceDef.UserList
				p.Post("/", hero.Handler(user.Create)).Name = resourceDef.UserCreate
				p.Get("/{id:int64}", hero.Handler(user.Detail)).Name = resourceDef.UserDetail
				p.Put("/{id:int64}", hero.Handler(user.Update)).Name = resourceDef.UserUpdate
				p.Delete("/{id:int64}", hero.Handler(user.Delete)).Name = resourceDef.UserDelete

				//给用户分配权限（通过用户同名角色）
				p.Put("/{id:int64}/perm", hero.Handler(user.UpdatePerm)).Name = resourceDef.RoleUpdate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(user.LogList)).Name = resourceDef.UserLogList
				p.Delete("/{id:int64}/log", hero.Handler(user.LogDelete)).Name = resourceDef.UserLogDelete
			})

			//设备分组
			p.PartyFunc("/group", func(p router.Party) {
				p.Get("/", hero.Handler(group.List)).Name = resourceDef.GroupList
				p.Post("/", hero.Handler(group.Create)).Name = resourceDef.GroupCreate
				p.Get("/{id:int64}", hero.Handler(group.Detail)).Name = resourceDef.GroupDetail
				p.Put("/{id:int64}", hero.Handler(group.Update)).Name = resourceDef.GroupUpdate
				p.Delete("/{id:int64}", hero.Handler(group.Delete)).Name = resourceDef.GroupDelete
			})

			//物理设备
			p.PartyFunc("/device", func(p router.Party) {
				p.Get("/", hero.Handler(device.List)).Name = resourceDef.DeviceList
				p.Post("/status", hero.Handler(device.MultiStatus)).Name = resourceDef.DeviceList
				p.Post("/", hero.Handler(device.Create)).Name = resourceDef.DeviceCreate

				p.Get("/{id:int64}", hero.Handler(device.Detail)).Name = resourceDef.DeviceDetail
				p.Put("/{id:int64}", hero.Handler(device.Update)).Name = resourceDef.DeviceUpdate
				p.Delete("/{id:int64}", hero.Handler(device.Delete)).Name = resourceDef.DeviceDelete

				//物理点位
				p.Get("/{id:int64}/measure", hero.Handler(device.MeasureList)).Name = resourceDef.MeasureList
				p.Post("/{id:int64}/measure", hero.Handler(device.CreateMeasure)).Name = resourceDef.MeasureCreate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(device.LogList)).Name = resourceDef.DeviceLogList
				p.Delete("/{id:int64}/log", hero.Handler(device.LogDelete)).Name = resourceDef.DeviceLogDelete

				//实时状态
				p.Get("/{id:int64}/reset", hero.Handler(device.Reset)).Name = resourceDef.DeviceStatus
				p.Get("/{id:int64}/status", hero.Handler(device.Status)).Name = resourceDef.DeviceStatus
				p.Get("/{id:int64}/data", hero.Handler(device.Data)).Name = resourceDef.DeviceData
				p.Put("/{id:int64}/{tagName:string}", hero.Handler(device.Ctrl)).Name = resourceDef.DeviceCtrl
				p.Get("/{id:int64}/{tagName:string}", hero.Handler(device.GetCHValue)).Name = resourceDef.DeviceCHValue

				//导出报表
				p.Get("/export/{uid:string}/stats", hero.Handler(statistics.ExportStats))
				p.Get("/export/{uid:string}/download", hero.Handler(statistics.ExportDownload))
				p.Post("/export", hero.Handler(statistics.Export)).Name = resourceDef.DataExport
			})
			//物理点位
			p.PartyFunc("/measure", func(p router.Party) {
				p.Delete("/{id:int64}", hero.Handler(device.DeleteMeasure)).Name = resourceDef.MeasureDelete
				p.Get("/{id:int64}", hero.Handler(device.MeasureDetail)).Name = resourceDef.MeasureDetail

				//历史趋势
				p.Post("/{id:int64}/statistics", hero.Handler(statistics.Measure)).Name = resourceDef.DeviceStatistics
			})

			//自定义设备
			p.PartyFunc("/equipment", func(p router.Party) {
				p.Get("/", hero.Handler(equipment.List)).Name = resourceDef.EquipmentList
				p.Post("/status", hero.Handler(equipment.MultiStatus)).Name = resourceDef.EquipmentList
				p.Post("/", hero.Handler(equipment.Create)).Name = resourceDef.EquipmentCreate
				p.Get("/{id:int64}", hero.Handler(equipment.Detail)).Name = resourceDef.EquipmentDetail
				p.Put("/{id:int64}", hero.Handler(equipment.Update)).Name = resourceDef.EquipmentUpdate
				p.Delete("/{id:int64}", hero.Handler(equipment.Delete)).Name = resourceDef.EquipmentDelete

				//自定义点位
				p.Get("/{id:int64}/state", hero.Handler(equipment.StateList)).Name = resourceDef.StateList
				p.Post("/{id:int64}/state", hero.Handler(equipment.CreateState)).Name = resourceDef.StateCreate

				//日志
				p.Get("/{id:int64}/log", hero.Handler(equipment.LogList)).Name = resourceDef.EquipmentLogList
				p.Delete("/{id:int64}/log", hero.Handler(equipment.LogDelete)).Name = resourceDef.EquipmentLogDelete

				//实时状态
				p.Get("/{id:int64}/status", hero.Handler(equipment.Status)).Name = resourceDef.EquipmentStatus
				p.Get("/{id:int64}/data", hero.Handler(equipment.Data)).Name = resourceDef.EquipmentData
				p.Put("/{id:int64}/{stateID:int64}", hero.Handler(equipment.Ctrl)).Name = resourceDef.EquipmentCtrl
				p.Get("/{id:int64}/{stateID:int64}", hero.Handler(equipment.GetCHValue)).Name = resourceDef.EquipmentCHValue
			})

			//自定义点位
			p.PartyFunc("/state", func(p router.Party) {
				p.Get("/{id:int64}", hero.Handler(equipment.StateDetail)).Name = resourceDef.StateDetail
				p.Put("/{id:int64}", hero.Handler(equipment.UpdateState)).Name = resourceDef.StateUpdate
				p.Delete("/{id:int64}", hero.Handler(equipment.DeleteState)).Name = resourceDef.StateDelete

				//历史趋势
				p.Post("/{id:int64}/statistics", hero.Handler(statistics.State)).Name = resourceDef.EquipmentStatistics
			})

			//警报
			p.PartyFunc("/alarm", func(p router.Party) {
				p.Get("/", hero.Handler(alarm.List)).Name = resourceDef.AlarmList
				p.Put("/{id:int64}", hero.Handler(alarm.Confirm)).Name = resourceDef.AlarmConfirm
				p.Get("/{id:int64}", hero.Handler(alarm.Detail)).Name = resourceDef.AlarmDetail
				p.Delete("/{id:int64}", hero.Handler(alarm.Delete)).Name = resourceDef.AlarmDelete

				p.Get("/{alarm:int64}/comments", hero.Handler(comment.List)).Name = resourceDef.CommentList

				//历史趋势
				p.Post("/statistics", hero.Handler(statistics.Alarm))

				//导出报表
				p.Get("/export", hero.Handler(alarm.Export))
			})

			//警报备注
			p.PartyFunc("/comment", func(p router.Party) {
				p.Get("/", hero.Handler(comment.List)).Name = resourceDef.CommentList
				p.Post("/", hero.Handler(comment.Create)).Name = resourceDef.CommentCreate
				p.Get("/{id:int64}", hero.Handler(comment.Detail)).Name = resourceDef.CommentDetail
				p.Delete("/{id:int64}", hero.Handler(comment.Delete)).Name = resourceDef.CommentDelete
			})

			//日志等级
			p.Get("/log/level", hero.Handler(logStore.Level)).Name = resourceDef.SysBrief
			//系统日志
			p.PartyFunc("/syslog", func(p router.Party) {
				p.Get("/", hero.Handler(logStore.List)).Name = resourceDef.LogList
				p.Delete("/", hero.Handler(logStore.Delete)).Name = resourceDef.LogDelete
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
