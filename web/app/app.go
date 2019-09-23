package app

import (
	"context"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/db"
	mysqlDB "github.com/maritimusj/centrum/db/mysql"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/logStore"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	mysqlStore "github.com/maritimusj/centrum/store/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	Config = config.New()

	Ctx, cancel = context.WithCancel(context.Background())
	DB          db.WithTransaction

	LogDBStore = logStore.New()
	s          store.Store
)

func IsDefaultAdminUser(user model.User) bool {
	return user.Name() == Config.DefaultUserName()
}

func Allow(user model.User, res model.Resource, action resource.Action) bool {
	allow, err := user.IsAllow(res, action)
	if err != nil && err == lang.Error(lang.ErrPolicyNotFound) {
		return Config.DefaultEffect() == resource.Allow
	}
	return allow
}

func Store() store.Store {
	return s
}

func NewStore(db db.DB) store.Store {
	return mysqlStore.Attach(Ctx, db, func(key string, _ interface{}) {
		s.Cache().Remove(key)
	})
}

func TransactionDo(fn func(store.Store) interface{}) interface{} {
	return DB.TransactionDo(func(db db.DB) interface{} {
		s := NewStore(db)
		defer s.Close()

		return fn(s)
	})
}

func InitDB(params map[string]interface{}) error {
	conn, err := mysqlDB.Open(Ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	DB = conn
	s = mysqlStore.Attach(Ctx, DB)
	return nil
}

func InitLog(levelStr string) error {
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		return err
	}

	Config.SetLogLevel(levelStr)

	//日志仓库
	err = LogDBStore.Open(Ctx, Config.LogFileName())
	if err != nil {
		return err
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.AddHook(LogDBStore)
	log.SetLevel(level)

	return nil
}

func Init(logLevel string) error {
	if err := InitLog(logLevel); err != nil {
		return err
	}

	//数据库连接
	if err := InitDB(map[string]interface{}{
		"connStr": Config.DBConnStr(),
	}); err != nil {
		return err
	}

	result := TransactionDo(func(s store.Store) interface{} {
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
			_, err := s.CreateOrganization(defaultOrg, defaultOrg)
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
	return nil
}

func Close() {
	cancel()
	LogDBStore.Wait()
}
