package app

import (
	"context"
	"github.com/maritimusj/centrum/gate/config"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/logStore/bolt"
	"github.com/maritimusj/centrum/gate/statistics"
	db2 "github.com/maritimusj/centrum/gate/web/db"
	"github.com/maritimusj/centrum/gate/web/db/mysql"
	edge2 "github.com/maritimusj/centrum/gate/web/edge"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	mysqlStore2 "github.com/maritimusj/centrum/gate/web/store/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/asaskevich/EventBus"
)

var (
	Config *config.Config

	Ctx    context.Context
	cancel context.CancelFunc
	DB     db2.WithTransaction

	Event      = EventBus.New()
	LogDBStore = bolt.New()
	s          store2.Store
	StatsDB    = statistics.New()
)

func IsDefaultAdminUser(user model2.User) bool {
	return user.Name() == Config.DefaultUserName()
}

func Allow(user model2.User, res model2.Resource, action resource2.Action) bool {
	if IsDefaultAdminUser(user) {
		return true
	}

	allow, err := user.IsAllow(res, action)
	if err != nil && err == lang2.Error(lang2.ErrPolicyNotFound) {
		return Config.DefaultEffect() == resource2.Allow
	}
	return allow
}

func Store() store2.Store {
	return s
}

func NewStore(db db2.DB) store2.Store {
	return mysqlStore2.Attach(Ctx, db, func(key string, _ interface{}) {
		s.Cache().Remove(key)
	})
}

func TransactionDo(fn func(store2.Store) interface{}) interface{} {
	return DB.TransactionDo(func(db db2.DB) interface{} {
		s := NewStore(db)
		defer s.Close()

		return fn(s)
	})
}

func InitDB(params map[string]interface{}) error {
	conn, err := mysql.Open(Ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	DB = conn
	s = mysqlStore2.Attach(Ctx, DB)
	return nil
}

func InitLog(levelStr string) error {
	if levelStr == "" {
		levelStr = Config.LogLevel()
	}

	level, err := log.ParseLevel(levelStr)
	if err != nil {
		return err
	}

	Config.SetLogLevel(levelStr)

	//日志仓库
	err = LogDBStore.Open(map[string]interface{}{
		"filename": Config.LogFileName(),
	})
	if err != nil {
		return err
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.AddHook(LogDBStore)
	log.SetLevel(level)

	return nil
}

func Init(ctx context.Context, logLevel string) error {
	Ctx, cancel = context.WithCancel(ctx)

	//数据库连接
	if err := InitDB(map[string]interface{}{
		//"connStr": "root:12345678@/chuanyan?charset=utf8mb4&parseTime=true&loc=Local",
		"connStr": "./chuanyan.db",
	}); err != nil {
		return err
	}

	Config = config.New(Store())
	err := Config.Load()
	if err != nil {
		log.Error(err)
		return err
	}

	if err := InitLog(logLevel); err != nil {
		return err
	}

	if err := initEvent(); err != nil {
		return err
	}

	result := TransactionDo(func(s store2.Store) interface{} {
		_, total, err := s.GetApiResourceList(helper2.Limit(1))
		if err != nil {
			return err
		}

		//初始化api资源
		if total == 0 {
			err = s.InitApiResource()
			if err != nil {
				return err
			}
		}

		//创建默认组织
		defaultOrg := Config.DefaultOrganization()
		org, err := s.GetOrganization(defaultOrg)
		if err != nil {
			if err != lang2.Error(lang2.ErrOrganizationNotFound) {
				return err
			}
			org, err = s.CreateOrganization(defaultOrg, defaultOrg)
			if err != nil {
				return err
			}

			_, err = s.CreateGroup(org, lang2.Str(lang2.DefaultGroupTitle), lang2.Str(lang2.DefaultGroupDesc), 0)
			if err != nil {
				return err
			}
		} else {
			org.Enable()
			if err = org.Save(); err != nil {
				return err
			}
		}

		//初始化系统角色
		_, err = s.GetRole(lang2.RoleSystemAdminName)
		if err != nil {
			if err != lang2.Error(lang2.ErrRoleNotFound) {
				return err
			}
			err = s.InitDefaultRoles(defaultOrg)
			if err != nil {
				return err
			}
		}

		//创建默认用户
		defaultUserName := Config.DefaultUserName()
		_, err = s.GetUser(defaultUserName)
		if err != nil {
			if err != lang2.Error(lang2.ErrUserNotFound) {
				return err
			}
			_, err = s.CreateUser(defaultOrg, defaultUserName, []byte(defaultUserName), lang2.RoleSystemAdminName)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if result != nil {
		return result.(error)
	}

	if err := StatsDB.Open(map[string]interface{}{
		"connStr": "http://localhost:8086",
	}); err != nil {
		return err
	}

	return nil
}

func BootAllDevices() error {
	devices, _, err := Store().GetDeviceList()
	if err != nil {
		return err
	}
	for _, device := range devices {
		if err := edge2.ActiveDevice(device); err != nil {
			log.Error(err)
			device.Logger().Error(err)
		}
	}
	return nil
}

func Close() {
	cancel()
	Store().Close()
	LogDBStore.Close()
}

func FlushDB() error {
	res := TransactionDo(func(store store2.Store) interface{} {
		return store.EraseAllData()
	})
	if res != nil {
		return res.(error)
	}
	return nil
}

func ResetDefaultAdminUser() error {
	user, err := Store().GetUser(Config.DefaultUserName())
	if err != nil {
		return err
	}

	user.Enable()
	user.ResetPassword(Config.DefaultUserName())
	if err = user.Save(); err != nil {
		return err
	}
	if err = user.SetRoles(lang2.RoleSystemAdminName); err != nil {
		return err
	}
	return nil
}

func SetAllow(user model2.User, res model2.Resource, actions ...resource2.Action) error {
	if IsDefaultAdminUser(user) {
		return nil
	}
	return user.SetAllow(res, actions...)
}

func SetDeny(user model2.User, res model2.Resource, actions ...resource2.Action) error {
	if IsDefaultAdminUser(user) {
		return nil
	}
	return user.SetDeny(res, actions...)
}
