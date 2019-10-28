package app

import (
	"context"
	"github.com/maritimusj/centrum/gate/config"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/logStore/bolt"
	"github.com/maritimusj/centrum/gate/statistics"
	"github.com/maritimusj/centrum/gate/web/db"
	"github.com/maritimusj/centrum/gate/web/db/mysql"
	"github.com/maritimusj/centrum/gate/web/edge"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	mysqlStore2 "github.com/maritimusj/centrum/gate/web/store/mysql"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/asaskevich/EventBus"
)

var (
	Config *config.Config

	Ctx    context.Context
	cancel context.CancelFunc
	DB     db.WithTransaction

	Event      = EventBus.New()
	LogDBStore = bolt.New()
	s          store2.Store
	StatsDB    = statistics.New()
)

func IsDefaultAdminUser(user model.User) bool {
	return user.Name() == Config.DefaultUserName()
}

func Allow(user model.User, res model.Resource, action resource.Action) bool {
	if IsDefaultAdminUser(user) {
		return true
	}

	allow, err := user.IsAllow(res, action)
	if err != nil && err == lang.Error(lang.ErrPolicyNotFound) {
		return Config.DefaultEffect() == resource.Allow
	}
	return allow
}

func Store() store2.Store {
	return s
}

func NewStore(db db.DB) store2.Store {
	return mysqlStore2.Attach(Ctx, db, func(key string, _ interface{}) {
		s.Cache().Remove(key)
	})
}

func TransactionDo(fn func(store2.Store) interface{}) interface{} {
	return DB.TransactionDo(func(db db.DB) interface{} {
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

	const dbFile =  "./chuanyan.db"
	var initDB bool
	if _, err := os.Stat(dbFile); err != nil  {
		if os.IsNotExist(err) {
			f, err := os.Create(dbFile)
			if err != nil {
				return lang.InternalError(err)
			}
			_ = f.Close()
			initDB = true
		}
	}
	//数据库连接
	if err := InitDB(map[string]interface{}{
		//"connStr": "root:12345678@/chuanyan?charset=utf8mb4&parseTime=true&loc=Local",
		"connStr": dbFile,
		"initDB": initDB,
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
		_, total, err := s.GetApiResourceList(helper.Limit(1))
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
			if err != lang.Error(lang.ErrOrganizationNotFound) {
				return err
			}
			org, err = s.CreateOrganization(defaultOrg, defaultOrg)
			if err != nil {
				return err
			}

			_, err = s.CreateGroup(org, lang.Str(lang.DefaultGroupTitle), lang.Str(lang.DefaultGroupDesc), 0)
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
		_, err = s.GetRole(lang.RoleSystemAdminName)
		if err != nil {
			if err != lang.Error(lang.ErrRoleNotFound) {
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
			if err != lang.Error(lang.ErrUserNotFound) {
				return err
			}
			_, err = s.CreateUser(defaultOrg, defaultUserName, []byte(defaultUserName), lang.RoleSystemAdminName)
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
		if err := edge.ActiveDevice(device); err != nil {
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
	if err = user.SetRoles(lang.RoleSystemAdminName); err != nil {
		return err
	}
	return nil
}

func SetAllow(user model.User, res model.Resource, actions ...resource.Action) error {
	if IsDefaultAdminUser(user) {
		return nil
	}
	return user.SetAllow(res, actions...)
}

func SetDeny(user model.User, res model.Resource, actions ...resource.Action) error {
	if IsDefaultAdminUser(user) {
		return nil
	}
	return user.SetDeny(res, actions...)
}
