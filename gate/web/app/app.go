package app

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/kataras/iris"

	"github.com/maritimusj/centrum/gate/web/SysInfo"

	"github.com/shirou/gopsutil/host"

	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"

	"github.com/maritimusj/centrum/util"

	"github.com/spf13/viper"

	"github.com/asaskevich/EventBus"
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
	"github.com/maritimusj/centrum/gate/web/store"
	mysqlStore "github.com/maritimusj/centrum/gate/web/store/mysql"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/register"
	log "github.com/sirupsen/logrus"

	edgeLang "github.com/maritimusj/centrum/edge/lang"
)

var (
	Config *config.Config

	Ctx    context.Context
	cancel context.CancelFunc
	DB     db.WithTransaction

	Event      = EventBus.New()
	LogDBStore = bolt.New()
	s          store.Store
	StatsDB    = statistics.New()
)

var (
	__hostInfo SysInfo.Info

	__cpuInfo SysInfo.Info

	__cpuTimes SysInfo.Info

	__diskInfo SysInfo.Info

	__memInfo SysInfo.Info
)

func IsDefaultAdminUser(user model.User) bool {
	return user.Name() == Config.DefaultUserName()
}

func Allow(user model.User, res model.Resource, action resource.Action) bool {
	if IsDefaultAdminUser(user) {
		return true
	}

	allow, err := user.IsAllow(res, action)
	if err != nil && err == lang.ErrPolicyNotFound.Error() {
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
	conn, err := mysql.Open(Ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	if initDB, ok := params["initDB"].(bool); ok && initDB {
		_, err = conn.Exec(initDBSQL)
		if err != nil {
			return lang.InternalError(err)
		}
	}

	DB = conn
	s = mysqlStore.Attach(Ctx, DB)
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

func HostInfo() interface{} {
	if !__hostInfo.Have() {
		__hostInfo.Fetch(func() (interface{}, error) {
			data, err := host.Info()
			if err != nil {
				log.Warn(err)
				return map[string]interface{}{}, nil
			}
			bootTime := time.Unix(int64(data.BootTime), 0)
			return map[string]interface{}{
				"hostname":        data.Hostname,
				"uptime":          data.Uptime,
				"bootTime":        bootTime.Format(lang.DatetimeFormatterStr.Str()),
				"os":              data.OS,
				"platform":        data.Platform,
				"platformVersion": data.PlatformVersion,
				"kernelArch":      data.KernelArch,
			}, nil
		})
	}
	res := __hostInfo.Data()
	if res == nil {
		return iris.Map{}
	}
	return res
}

func CpuInfo() interface{} {
	if !__cpuInfo.Have() {
		__cpuInfo.Fetch(func() (interface{}, error) {
			var err error
			info, err := cpu.Info()
			if err != nil {
				log.Warningln(err)
			}

			xx := make([]interface{}, 0, len(info))
			for _, x := range info {
				xx = append(xx, map[string]interface{}{
					"cores":     x.Cores,
					"mhz":       x.Mhz,
					"modelName": x.ModelName,
				})
			}
			return xx, nil
		})
	}

	res := __cpuInfo.Data()
	if res == nil {
		return iris.Map{}
	}
	return res
}

func CpuTimes() interface{} {
	if __cpuTimes.Expired(1 * time.Second) {
		__cpuTimes.Fetch(func() (interface{}, error) {
			percent, err := cpu.Percent(1*time.Second, false)
			if err != nil {
				log.Warningln(err)
				return nil, err
			}
			return fmt.Sprintf("%.2f", percent[0]), nil
		})
	}

	res := __cpuTimes.Data()
	if res == nil {
		return ""
	}
	return res
}

func MemoryStatus() interface{} {
	if __memInfo.Expired(1 * time.Second) {
		__memInfo.Fetch(func() (interface{}, error) {
			m, err := mem.VirtualMemory()
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"total":       util.FormatFileSize(m.Total),
				"available":   util.FormatFileSize(m.Available),
				"used":        util.FormatFileSize(m.Used),
				"usedPercent": fmt.Sprintf("%.2f", m.UsedPercent),
			}, nil
		})
	}

	res := __memInfo.Data()
	if res == nil {
		return iris.Map{}
	}
	return res
}

func DiskStatus() interface{} {
	if __diskInfo.Expired(1 * time.Minute) {
		__diskInfo.Fetch(func() (interface{}, error) {
			ps, err := disk.Partitions(true)
			if err != nil {
				return nil, err
			}
			x := make([]map[string]interface{}, 0, len(ps))
			for _, p := range ps {
				v, _ := disk.Usage(p.Device)
				x = append(x, map[string]interface{}{
					"path":        v.Path,
					"total":       util.FormatFileSize(v.Total),
					"used":        util.FormatFileSize(v.Used),
					"usedPercent": fmt.Sprintf("%.2f", v.UsedPercent),
				})
			}
			return x, nil
		})
	}

	res := __diskInfo.Data()
	if res == nil {
		return iris.Map{}
	}
	return res
}

func SysStatus() interface{} {
	return map[string]interface{}{
		"host": HostInfo(),
		"mem":  MemoryStatus(),
		"disk": DiskStatus(),
		"cpu": map[string]interface{}{
			"cpus":        CpuInfo(),
			"usedPercent": CpuTimes(),
		},
	}
}

func Start(ctx context.Context, logLevel string) error {
	Ctx, cancel = context.WithCancel(ctx)

	const dbFile = "./chuanyan.db"
	var initDB bool

	if _, err := os.Stat(dbFile); err != nil {
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
		"connStr": dbFile,
		"initDB":  initDB,
	}); err != nil {
		return err
	}

	Config = config.New(Store())
	err := Config.Load()
	if err != nil {
		log.Error(err)
		return err
	}

	if viper.IsSet("gate") {
		_ = Config.BaseConfig.SetOption(config.ApiAddrPath, viper.GetString("gate.addr"))
		_ = Config.BaseConfig.SetOption(config.ApiPortPath, viper.GetString("gate.port"))
	}

	if viper.IsSet("influxdb") {
		_ = Config.BaseConfig.SetOption(config.InfluxDBUrl, viper.GetString("influxdb.url"))
		_ = Config.BaseConfig.SetOption(config.InfluxDBUserName, viper.GetString("influxdb.username"))
		_ = Config.BaseConfig.SetOption(config.InfluxDBPassword, viper.GetString("influxdb.password"))
	}

	if err := InitLog(logLevel); err != nil {
		return err
	}

	if err := initEvent(); err != nil {
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
			if err != lang.ErrOrganizationNotFound.Error() {
				return err
			}
			org, err = s.CreateOrganization(defaultOrg, defaultOrg)
			if err != nil {
				return err
			}

			_, err = s.CreateGroup(org, lang.DefaultGroupTitle.Str(), lang.DefaultGroupDesc.Str(), 0)
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
			if err != lang.ErrRoleNotFound.Error() {
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
			if err != lang.ErrUserNotFound.Error() {
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

func Fingerprints() string {
	return register.Fingerprints()
}

func IsRegistered() bool {
	return register.Verify(Config.RegOwner(), Config.RegCode())
}

func SaveRegisterInfo(owner, code string) error {
	if !register.Verify(owner, code) {
		return lang.ErrInvalidRegCode.Error()
	}

	_ = Config.BaseConfig.SetOption(config.SysRegOwnerPath, owner)
	_ = Config.BaseConfig.SetOption(config.SysRegCodePath, code)

	if err := Config.BaseConfig.Save(); err != nil {
		return err
	}
	return nil
}

func BootAllDevices() {
	if !IsRegistered() {
		return
	}

	select {
	case <-Ctx.Done():
		return
	default:
	}

	devices, _, err := Store().GetDeviceList()
	if err != nil {
		log.Error("[BootAllDevices] ", err)
		return
	}

	for _, device := range devices {
		if err := edge.ActiveDevice(device, Config); err != nil {
			log.Error("[BootAllDevices] ActiveDevice: ", err)

			device.Logger().Error(err)
			global.UpdateDeviceStatus(device, int(edgeLang.EdgeUnknownState), edgeLang.Str(edgeLang.EdgeUnknownState))
		}
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	}

	time.AfterFunc(15*time.Second, BootAllDevices)
}

func Close() {
	cancel()
	Store().Close()
	LogDBStore.Close()
}

func FlushDB() error {
	res := TransactionDo(func(store store.Store) interface{} {
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
