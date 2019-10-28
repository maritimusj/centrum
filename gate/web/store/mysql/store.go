package mysqlStore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	cache2 "github.com/maritimusj/centrum/gate/web/cache"
	memCache2 "github.com/maritimusj/centrum/gate/web/cache/memCache"
	db2 "github.com/maritimusj/centrum/gate/web/db"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	status2 "github.com/maritimusj/centrum/gate/web/status"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	"time"

	"github.com/maritimusj/centrum/synchronized"
	"github.com/maritimusj/centrum/util"

	log "github.com/sirupsen/logrus"
)

const (
	TbConfig          = "`config`"
	TbOrganization    = "`organizations`"
	TbUsers           = "`users`"
	TbRoles           = "`roles`"
	TbUserRoles       = "`user_roles`"
	TbPolicies        = "`policies`"
	TbGroups          = "`groups`"
	TbDevices         = "`devices`"
	TbMeasures        = "`measures`"
	TbDeviceGroups    = "`device_groups`"
	TbEquipments      = "`equipments`"
	TbStates          = "`states`"
	TbEquipmentGroups = "`equipment_groups`"
	TbApiResources    = "`api_resources`"
	TbAlarms          = "`alarms`"
)

type mysqlStore struct {
	db       db2.DB
	cache    cache2.Cache
	ctx      context.Context
	cleaners []func(string, interface{})
}

func New() store2.Store {
	return &mysqlStore{
		cache: memCache2.New(),
	}
}

func Attach(ctx context.Context, db db2.DB, cleaners ...func(key string, obj interface{})) store2.Store {
	s := storePool.Get().(*mysqlStore)
	s.ctx = ctx
	s.db = db
	s.cleaners = append(s.cleaners, cleaners...)
	return s
}

func parseOption(options ...helper2.OptionFN) *helper2.Option {
	option := helper2.Option{}
	for _, opt := range options {
		if opt != nil {
			opt(&option)
		}
	}
	return &option
}

func (s *mysqlStore) EraseAllData() error {
	var statements = []string{
		"DELETE FROM " + TbOrganization,
		"DELETE FROM " + TbUsers,
		"DELETE FROM " + TbRoles,
		"DELETE FROM " + TbUserRoles,
		"DELETE FROM " + TbPolicies,
		"DELETE FROM " + TbGroups,
		"DELETE FROM " + TbDevices,
		"DELETE FROM " + TbMeasures,
		"DELETE FROM " + TbDeviceGroups,
		"DELETE FROM " + TbEquipments,
		"DELETE FROM " + TbStates,
		"DELETE FROM " + TbEquipmentGroups,
		"DELETE FROM " + TbApiResources,
		"DELETE FROM " + TbAlarms,
		"UPDATE `sqlite_sequence` SET seq = 0",
	}
	for _, st := range statements {
		_, err := s.db.Exec(st)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *mysqlStore) Close() {
	for _, fn := range s.cleaners {
		s.cache.Foreach(fn)
	}
	s.cache.Flush()
	storePool.Put(s)
}

func (s *mysqlStore) Cache() cache2.Cache {
	return s.cache
}

func (s *mysqlStore) loadConfig(id int64) (*Config, error) {
	var cfg = NewConfig(s, id)
	err := LoadData(s.db, TbConfig, map[string]interface{}{
		"name":       &cfg.name,
		"extra?":     &cfg.extra,
		"created_at": &cfg.createdAt,
		"update_at":  &cfg.updateAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrConfigNotFound)
	}

	return cfg, nil
}

func (s *mysqlStore) GetConfig(cfg interface{}) (model2.Config, error) {
	result := <-synchronized.Do(TbConfig, func() interface{} {
		var cfgID int64
		cfgID, err := s.getConfigID(cfg)
		if err != nil {
			return err
		}

		cfg, err := s.cache.LoadConfig(cfgID)
		if err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return cfg
		}

		cfg, err = s.loadConfig(cfgID)
		if err != nil {
			return err
		}

		err = s.cache.Save(cfg)
		if err != nil {
			return err
		}
		return cfg
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Config), nil
}

func (s *mysqlStore) CreateConfig(name string, data interface{}) (model2.Config, error) {
	result := <-synchronized.Do(TbConfig, func() interface{} {
		now := time.Now()
		data, err := json.Marshal(util.If(data != nil, data, "{}"))
		if err != nil {
			return err
		}
		cfgID, err := CreateData(s.db, TbConfig, map[string]interface{}{
			"name":       name,
			"extra":      data,
			"created_at": now,
			"update_at":  now,
		})
		if err != nil {
			return lang2.InternalError(err)
		}

		cfg, err := s.loadConfig(cfgID)
		if err != nil {
			return err
		}

		err = s.cache.Save(cfg)
		if err != nil {
			return err
		}

		return cfg
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Config), nil
}

func (s *mysqlStore) RemoveConfig(id interface{}) error {
	cfgID, err := s.getConfigID(id)
	if err != nil {
		return err
	}

	err = RemoveData(s.db, TbConfig, "id=?", cfgID)
	if err != nil {
		return lang2.InternalError(err)
	}

	s.cache.Remove(&Config{id: cfgID})
	return nil
}

func (s *mysqlStore) GetConfigList(options ...helper2.OptionFN) ([]model2.Config, int64, error) {
	option := parseOption(options...)
	var (
		fromSQL = "FROM " + TbConfig + " c  WHERE 1"
	)

	var params []interface{}

	if option.Keyword != "" {
		fromSQL += " AND c.name REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if option.GetTotal == nil || *option.GetTotal {
		if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
			return nil, 0, lang2.InternalError(err)
		}

		if total == 0 {
			return []model2.Config{}, 0, nil
		}
	}

	fromSQL += " ORDER BY c.id ASC"

	if option.Limit > 0 {
		fromSQL += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		fromSQL += " OFFSET ?"
		params = append(params, option.Offset)
	}
	rows, err := s.db.Query("SELECT r.id "+fromSQL, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var cfgID int64

	for rows.Next() {
		err = rows.Scan(&cfgID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Config{}, total, nil
		}
		ids = append(ids, cfgID)
	}

	var result []model2.Config
	for _, id := range ids {
		cfg, err := s.GetConfig(id)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, cfg)
	}

	return result, total, nil
}

func (s *mysqlStore) MustGetUserFromContext(ctx iris.Context) model2.User {
	userID := ctx.Values().GetInt64Default("__userID__", 0)
	if userID > 0 {
		user, err := s.GetUser(userID)
		if err != nil {
			panic(err)
		}
		return user
	}
	panic(lang2.Error(lang2.ErrInvalidUser))
}

func (s *mysqlStore) IsOrganizationExists(org interface{}) (bool, error) {
	if _, err := s.getOrganizationID(org); err != nil {
		if err != lang2.Error(lang2.ErrOrganizationNotFound) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (s *mysqlStore) IsUserExists(user interface{}) (bool, error) {
	if _, err := s.getUserID(user); err != nil {
		if err != lang2.Error(lang2.ErrUserNotFound) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (s *mysqlStore) IsRoleExists(role interface{}) (bool, error) {
	if _, err := s.getRoleID(role); err != nil {
		if err != lang2.Error(lang2.ErrRoleNotFound) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (s *mysqlStore) GetResourceGroupList() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":    resource2.Api,
			"title": lang2.ResourceClassTitle(resource2.Api),
		},
		map[string]interface{}{
			"id":    resource2.Group,
			"title": lang2.ResourceClassTitle(resource2.Group),
		},
		map[string]interface{}{
			"id":    resource2.Device,
			"title": lang2.ResourceClassTitle(resource2.Device),
		},
		map[string]interface{}{
			"id":    resource2.Measure,
			"title": lang2.ResourceClassTitle(resource2.Measure),
		},
		map[string]interface{}{
			"id":    resource2.Equipment,
			"title": lang2.ResourceClassTitle(resource2.Equipment),
		},
		map[string]interface{}{
			"id":    resource2.State,
			"title": lang2.ResourceClassTitle(resource2.State),
		},
	}
}

func (s *mysqlStore) loadOrganization(id int64) (*Organization, error) {
	var org = NewOrganization(s, id)
	err := LoadData(s.db, TbOrganization, map[string]interface{}{
		"enable":     &org.enable,
		"name":       &org.name,
		"title":      &org.title,
		"extra?":     &org.extra,
		"created_at": &org.createdAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrOrganizationNotFound)
	}
	return org, nil
}

func (s *mysqlStore) GetOrganization(id interface{}) (model2.Organization, error) {
	result := <-synchronized.Do(TbOrganization, func() interface{} {
		var orgID int64
		orgID, err := s.getOrganizationID(id)
		if err != nil {
			return err
		}

		org, err := s.cache.LoadOrganization(orgID)
		if err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return org
		}

		org, err = s.loadOrganization(orgID)
		if err != nil {
			return err
		}

		err = s.cache.Save(org)
		if err != nil {
			return err
		}
		return org
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Organization), nil
}

func (s *mysqlStore) CreateOrganization(name string, title string) (model2.Organization, error) {
	result := <-synchronized.Do(TbOrganization, func() interface{} {
		orgID, err := CreateData(s.db, TbOrganization, map[string]interface{}{
			"enable":     status2.Enable,
			"name":       name,
			"title":      title,
			"extra":      `{}`,
			"created_at": time.Now(),
		})
		if err != nil {
			return lang2.InternalError(err)
		}

		org, err := s.loadOrganization(orgID)
		if err != nil {
			return err
		}

		err = s.cache.Save(org)
		if err != nil {
			return err
		}
		return org
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Organization), nil
}

func (s *mysqlStore) RemoveOrganization(id interface{}) error {
	orgID, err := s.getOrganizationID(id)
	if err != nil {
		return err
	}

	err = RemoveData(s.db, TbOrganization, "id=?", orgID)
	if err != nil {
		return lang2.InternalError(err)
	}

	s.cache.Remove(&Organization{id: orgID})
	return nil
}

func (s *mysqlStore) GetOrganizationList(options ...helper2.OptionFN) ([]model2.Organization, int64, error) {
	option := parseOption(options...)

	var (
		from  = "FROM " + TbOrganization + " o"
		where = " WHERE 1"
	)

	var params []interface{}
	if option.Keyword != "" {
		where += " AND (o.name REGEXP ? OR o.title REGEXP ?)"
		params = append(params, option.Keyword, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(DISTINCT o.id) "+from+where, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.Organization{}, 0, nil
	}

	where += " ORDER BY o.id ASC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	rows, err := s.db.Query("SELECT DISTINCT o.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var userID int64

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Organization{}, total, nil
		}
		ids = append(ids, userID)
	}

	var result []model2.Organization
	for _, id := range ids {
		org, err := s.GetOrganization(id)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, org)
	}

	return result, total, nil
}

func (s *mysqlStore) loadUser(id int64) (*User, error) {
	var user = NewUser(s, id)
	err := LoadData(s.db, TbUsers, map[string]interface{}{
		"org_id":     &user.orgID,
		"enable":     &user.enable,
		"name":       &user.name,
		"title":      &user.title,
		"password":   &user.password,
		"mobile":     &user.mobile,
		"email":      &user.email,
		"created_at": &user.createdAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrUserNotFound)
	}

	return user, nil
}

func (s *mysqlStore) GetUser(user interface{}) (model2.User, error) {
	result := <-synchronized.Do(TbUsers, func() interface{} {
		userID, err := s.getUserID(user)
		if err != nil {
			return err
		}

		if user, err := s.cache.LoadUser(userID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return user
		}

		user, err := s.loadUser(userID)
		if err != nil {
			return err
		}

		err = s.cache.Save(user)
		if err != nil {
			return err
		}
		return user
	})

	if err, ok := result.(error); ok {
		return nil, err
	}

	return result.(model2.User), nil
}

func (s *mysqlStore) CreateUser(org interface{}, name string, password []byte, roles ...interface{}) (model2.User, error) {
	result := <-synchronized.Do(TbUsers, func() interface{} {
		orgID, err := s.getOrganizationID(org)
		if err != nil {
			return err
		}

		passwordData, err := util.HashPassword(password)
		if err != nil {
			return lang2.InternalError(err)
		}

		userID, err := CreateData(s.db, TbUsers, map[string]interface{}{
			"org_id":     orgID,
			"enable":     status2.Enable,
			"name":       name,
			"password":   passwordData,
			"title":      name,
			"mobile":     "",
			"email":      "",
			"created_at": time.Now(),
		})

		if err != nil {
			return lang2.InternalError(err)
		}

		user, err := s.loadUser(userID)
		if err != nil {
			return err
		}

		if len(roles) > 0 {
			err = user.SetRoles(roles...)
			if err != nil {
				return err
			}
		}

		err = s.cache.Save(user)
		if err != nil {
			return err
		}
		return user
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.User), nil
}

func (s *mysqlStore) RemoveUser(u interface{}) error {
	userID, err := s.getUserID(u)
	if err != nil {
		return err
	}

	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}

	err = user.SetRoles()
	if err != nil {
		return err
	}

	err = RemoveData(s.db, TbUsers, "id=?", userID)
	if err != nil {
		return lang2.InternalError(err)
	}

	s.cache.Remove(user)
	return nil
}

func (s *mysqlStore) GetUserList(options ...helper2.OptionFN) ([]model2.User, int64, error) {
	option := parseOption(options...)

	var (
		from  = "FROM " + TbUsers + " u"
		where = " WHERE 1"
	)

	var params []interface{}

	if option.OrgID > 0 {
		where += " AND u.org_id=?"
		params = append(params, option.OrgID)
	}

	if option.RoleID != nil && *option.RoleID > 0 {
		from += " LEFT JOIN " + TbUserRoles + " r ON u.id=r.user_id"
		where += " AND r.role_id=?"
		params = append(params, *option.RoleID)
	}

	if option.Keyword != "" {
		where += " AND (u.name REGEXP ? OR u.title REGEXP ? OR u.mobile REGEXP ? OR u.email REGEXP ?)"
		params = append(params, option.Keyword, option.Keyword, option.Keyword, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(DISTINCT u.id) "+from+where, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.User{}, 0, nil
	}

	where += " ORDER BY u.id ASC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	rows, err := s.db.Query("SELECT DISTINCT u.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var userID int64

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.User{}, total, nil
		}
		ids = append(ids, userID)
	}

	var result []model2.User
	for _, id := range ids {
		user, err := s.GetUser(id)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, user)
	}

	return result, total, nil
}

func (s *mysqlStore) loadRole(id int64) (*Role, error) {
	var role = NewRole(s, id)
	err := LoadData(s.db, TbRoles, map[string]interface{}{
		"org_id":     &role.orgID,
		"enable":     &role.enable,
		"name":       &role.name,
		"title":      &role.title,
		"desc":       &role.desc,
		"created_at": &role.createdAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrRoleNotFound)
	}
	return role, nil
}

func (s *mysqlStore) GetRole(role interface{}) (model2.Role, error) {
	result := <-synchronized.Do(TbRoles, func() interface{} {
		roleID, err := s.getRoleID(role)
		if err != nil {
			return err
		}
		if role, err := s.cache.LoadRole(roleID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return lang2.InternalError(err)
			}
		} else {
			return role
		}

		role, err := s.loadRole(roleID)
		if err != nil {
			return err
		}

		err = s.cache.Save(role)
		if err != nil {
			return err
		}
		return role
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Role), nil
}

func (s *mysqlStore) createRole(db db2.DB, org interface{}, name, title, desc string) (model2.Role, error) {
	result := <-synchronized.Do(TbRoles, func() interface{} {
		orgID, err := s.getOrganizationID(org)
		if err != nil {
			return err
		}
		roleID, err := CreateData(s.db, TbRoles, map[string]interface{}{
			"org_id":     orgID,
			"enable":     status2.Enable,
			"name":       name,
			"title":      title,
			"desc":       desc,
			"created_at": time.Now(),
		})
		if err != nil {
			return err
		}

		role, err := s.loadRole(roleID)
		if err != nil {
			return err
		}
		err = s.cache.Save(role)
		if err != nil {
			return err
		}
		return role
	})

	if err, ok := result.(error); ok {
		return nil, err
	}

	return result.(model2.Role), nil
}

func (s *mysqlStore) CreateRole(org interface{}, name, title, desc string) (model2.Role, error) {
	return s.createRole(s.db, org, name, title, desc)
}

func (s *mysqlStore) RemoveRole(role interface{}) error {
	roleID, err := s.getRoleID(role)
	if err != nil {
		return err
	}

	err = RemoveData(s.db, TbRoles, "id=?", roleID)
	if err != nil {
		return lang2.InternalError(err)
	}
	s.cache.Remove(&Role{id: roleID})
	return nil
}

func (s *mysqlStore) GetRoleList(options ...helper2.OptionFN) ([]model2.Role, int64, error) {
	option := parseOption(options...)
	var (
		fromSQL = "FROM " + TbRoles + " r "
	)

	var params []interface{}

	if option.UserID != nil {
		fromSQL += " INNER JOIN " + TbUserRoles + " u ON r.id=u.role_id WHERE u.user_id=?"
		params = append(params, *option.UserID)
	} else {
		fromSQL += " WHERE 1"
	}

	if option.OrgID > 0 {
		fromSQL += " AND r.org_id=?"
		params = append(params, option.OrgID)
	}

	if option.Keyword != "" {
		fromSQL += " AND (r.name LIKE ? OR r.title LIKE ?)"
		keyword := fmt.Sprintf("%%%s%%", option.Keyword)
		params = append(params, keyword, keyword)
	}

	var total int64
	if option.GetTotal == nil || *option.GetTotal {
		if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
			return nil, 0, lang2.InternalError(err)
		}

		if total == 0 {
			return []model2.Role{}, 0, nil
		}
	}

	fromSQL += " ORDER BY r.id ASC"

	if option.Limit > 0 {
		fromSQL += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		fromSQL += " OFFSET ?"
		params = append(params, option.Offset)
	}
	rows, err := s.db.Query("SELECT r.id "+fromSQL, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var roleID int64

	for rows.Next() {
		err = rows.Scan(&roleID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Role{}, total, nil
		}
		ids = append(ids, roleID)
	}

	var result []model2.Role
	for _, id := range ids {
		role, err := s.GetRole(id)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, role)
	}

	return result, total, nil
}

func (s *mysqlStore) loadPolicy(id int64) (model2.Policy, error) {
	var policy = NewPolicy(s, id)
	err := LoadData(s.db, TbPolicies, map[string]interface{}{
		"role_id":        &policy.roleID,
		"resource_class": &policy.resourceClass,
		"resource_id":    &policy.resourceID,
		"action?":        &policy.action,
		"effect?":        &policy.effect,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrPolicyNotFound)
	}
	return policy, nil
}

func (s *mysqlStore) GetPolicy(policyID int64) (model2.Policy, error) {
	result := <-synchronized.Do(TbPolicies, func() interface{} {
		if role, err := s.cache.LoadPolicy(policyID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return lang2.InternalError(err)
			}
		} else {
			return role
		}

		policy, err := s.loadPolicy(policyID)
		if err != nil {
			return err
		}

		err = s.cache.Save(policy)
		if err != nil {
			return err
		}
		return policy
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Policy), nil
}

func (s *mysqlStore) GetPolicyFrom(roleID int64, res model2.Resource, action resource2.Action) (model2.Policy, error) {
	var policyID int64
	err := LoadData(s.db, TbPolicies, map[string]interface{}{
		"id": &policyID,
	}, "role_id=? AND resource_class=? AND resource_id=? AND action=?", roleID, res.ResourceClass(), res.ResourceID(), action)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrPolicyNotFound)
	}

	return s.GetPolicy(policyID)

}

func (s *mysqlStore) CreatePolicy(roleID int64, res model2.Resource, action resource2.Action, effect resource2.Effect) (model2.Policy, error) {
	result := <-synchronized.Do(TbPolicies, func() interface{} {
		policyID, err := CreateData(s.db, TbPolicies, map[string]interface{}{
			"role_id":        roleID,
			"resource_class": res.ResourceClass(),
			"resource_id":    res.ResourceID(),
			"action":         action,
			"effect":         effect,
			"created_at":     time.Now(),
		})

		if err != nil {
			return err
		}

		policy, err := s.loadPolicy(policyID)
		if err != nil {
			return err
		}
		err = s.cache.Save(policy)
		if err != nil {
			return err
		}
		return policy
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Policy), nil
}

func (s *mysqlStore) RemovePolicy(policyID int64) error {
	err := RemoveData(s.db, TbPolicies, "id=?", policyID)
	if err != nil {
		return lang2.InternalError(err)
	}
	s.cache.Remove(&Policy{id: policyID})
	return nil
}

func (s *mysqlStore) GetPolicyList(res model2.Resource, options ...helper2.OptionFN) ([]model2.Policy, int64, error) {
	option := parseOption(options...)

	var (
		fromSQL = "FROM " + TbPolicies + " WHERE 1"
	)

	var params []interface{}

	if option.RoleID != nil {
		fromSQL += " AND role_id=?"
		params = append(params, *option.RoleID)
	}

	if res != nil {
		fromSQL += " AND (resource_class=? AND resource_id=?)"
		params = append(params, res.ResourceClass(), res.ResourceID())
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.Policy{}, 0, nil
	}

	if option.Limit > 0 {
		fromSQL += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		fromSQL += " OFFSET ?"
		params = append(params, option.Offset)
	}

	rows, err := s.db.Query("SELECT id "+fromSQL, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var policyID int64

	for rows.Next() {
		err = rows.Scan(&policyID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Policy{}, total, nil
		}
		ids = append(ids, policyID)
	}

	var result []model2.Policy
	for _, id := range ids {
		policy, err := s.GetPolicy(id)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, policy)
	}

	return result, total, nil
}

func (s *mysqlStore) loadGroup(id int64) (model2.Group, error) {
	var group = NewGroup(s, id)

	err := LoadData(s.db, TbGroups, map[string]interface{}{
		"org_id":     &group.orgID,
		"parent_id":  &group.parentID,
		"title":      &group.title,
		"desc?":      &group.desc,
		"created_at": &group.createdAt,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrGroupNotFound)
	}
	return group, nil
}

func (s *mysqlStore) GetGroup(groupID int64) (model2.Group, error) {
	result := <-synchronized.Do(TbGroups, func() interface{} {
		if group, err := s.cache.LoadGroup(groupID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return group
		}

		group, err := s.loadGroup(groupID)
		if err != nil {
			return err
		}

		err = s.cache.Save(group)
		if err != nil {
			return err
		}
		return group
	})
	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Group), nil
}

func (s *mysqlStore) GetDeviceGroups(deviceID int64) ([]model2.Group, error) {
	const (
		SQL = "SELECT group_id FROM " + TbDeviceGroups + " WHERE device_id=?"
	)

	rows, err := s.db.Query(SQL, deviceID)
	if err != nil {
		return nil, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var groupID int64

	for rows.Next() {
		err = rows.Scan(&groupID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, lang2.InternalError(err)
			}
			return []model2.Group{}, nil
		}
		ids = append(ids, groupID)
	}

	var result []model2.Group
	for _, groupID := range ids {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return nil, err
		}
		result = append(result, group)
	}
	return result, nil
}

func (s *mysqlStore) GetEquipmentGroups(equipmentID int64) ([]model2.Group, error) {
	const (
		SQL = "SELECT group_id FROM " + TbEquipmentGroups + " WHERE equipment_id=?"
	)

	rows, err := s.db.Query(SQL, equipmentID)
	if err != nil {
		return nil, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var groupID int64

	for rows.Next() {
		err = rows.Scan(&groupID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, lang2.InternalError(err)
			}
			return []model2.Group{}, nil
		}
		ids = append(ids, groupID)
	}

	var result []model2.Group
	for _, groupID := range ids {
		group, err := s.GetGroup(groupID)
		if err != nil {
			return nil, err
		}
		result = append(result, group)
	}

	return result, nil
}

func (s *mysqlStore) CreateGroup(org interface{}, title string, desc string, parentID int64) (model2.Group, error) {
	result := <-synchronized.Do(TbGroups, func() interface{} {
		orgID, err := s.getOrganizationID(org)
		if err != nil {
			return err
		}
		data := map[string]interface{}{
			"org_id":     orgID,
			"parent_id":  parentID,
			"title":      title,
			"desc":       desc,
			"created_at": time.Now(),
		}

		groupID, err := CreateData(s.db, TbGroups, data)
		if err != nil {
			return err
		}
		group, err := s.loadGroup(groupID)
		if err != nil {
			return err
		}
		err = s.cache.Save(group)
		if err != nil {
			return err
		}
		return group
	})
	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Group), nil
}

func (s *mysqlStore) RemoveGroup(groupID int64) error {
	err := RemoveData(s.db, TbGroups, "id=?", groupID)
	if err != nil {
		return err
	}

	s.cache.Remove(&Group{id: groupID})
	return nil
}

func (s *mysqlStore) GetGroupList(options ...helper2.OptionFN) ([]model2.Group, int64, error) {
	option := parseOption(options...)
	var (
		from  = "FROM " + TbGroups + " g"
		where = " WHERE 1"
	)

	var params []interface{}

	if option.OrgID > 0 {
		where += " AND g.org_id=?"
		params = append(params, option.OrgID)
	}

	if option.UserID != nil {
		userID := *option.UserID
		if userID > 0 {
			from += fmt.Sprintf(` LEFT JOIN (
SELECT g.id,p.role_id,p.action,p.effect FROM %s g
INNER JOIN %s p ON p.resource_class=%d AND p.resource_id=g.id
INNER JOIN %s r ON p.role_id=r.id
WHERE p.role_id IN (SELECT role_id FROM %s WHERE user_id=%d)
) b ON g.id=b.id`, TbGroups, TbPolicies, resource2.Group, TbRoles, TbUserRoles, userID)

			if option.DefaultEffect == resource2.Allow {
				where += " AND ((b.action=0 AND b.effect=1) OR (ISNULL(b.action) AND ISNULL(b.effect)))"
			} else {
				where += " AND (b.action=0 AND b.effect=1)"
			}
		}
	}

	if option.ParentID != nil {
		where += " AND g.parent_id=?"
		params = append(params, *option.ParentID)
	}

	if option.Keyword != "" {
		where += " AND g.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(DISTINCT g.id) "+from+where, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.Group{}, 0, nil
	}

	where += " ORDER BY g.id ASC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	log.Trace("SELECT DISTINCT g.id " + from + where)
	rows, err := s.db.Query("SELECT DISTINCT g.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var groupID int64

	for rows.Next() {
		err = rows.Scan(&groupID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Group{}, total, nil
		}
		ids = append(ids, groupID)
	}

	var result []model2.Group
	for _, id := range ids {
		group, err := s.GetGroup(id)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, group)
	}

	return result, total, nil
}

func (s *mysqlStore) loadDevice(id int64) (model2.Device, error) {
	var device = NewDDevice(s, id)
	err := LoadData(s.db, TbDevices, map[string]interface{}{
		"org_id":     &device.orgID,
		"enable":     &device.enable,
		"title":      &device.title,
		"options?":   &device.options,
		"created_at": &device.createdAt,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrDeviceNotFound)
	}
	return device, nil
}

func (s *mysqlStore) GetDevice(deviceID int64) (model2.Device, error) {
	result := <-synchronized.Do(TbDevices, func() interface{} {
		if device, err := s.cache.LoadDevice(deviceID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return device
		}

		device, err := s.loadDevice(deviceID)
		if err != nil {
			return err
		}

		err = s.cache.Save(device)
		if err != nil {
			return err
		}
		return device
	})
	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Device), nil
}

func (s *mysqlStore) CreateDevice(org interface{}, title string, data map[string]interface{}) (model2.Device, error) {
	result := <-synchronized.Do(TbDevices, func() interface{} {
		orgID, err := s.getOrganizationID(org)
		if err != nil {
			return err
		}

		o, err := json.Marshal(data)
		if err != nil {
			return lang2.InternalError(err)
		}

		deviceID, err := CreateData(s.db, TbDevices, map[string]interface{}{
			"org_id":     orgID,
			"enable":     status2.Enable,
			"title":      title,
			"options":    o,
			"created_at": time.Now(),
		})
		if err != nil {
			return err
		}
		device, err := s.loadDevice(deviceID)
		if err != nil {
			return err
		}
		err = s.cache.Save(device)
		if err != nil {
			return err
		}
		return device
	})
	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Device), nil
}

func (s *mysqlStore) RemoveDevice(deviceID int64) error {
	err := RemoveData(s.db, TbDevices, "id=?", deviceID)
	if err != nil {
		return lang2.InternalError(err)
	}

	s.cache.Remove(&Device{id: deviceID})
	return nil
}

func (s *mysqlStore) GetDeviceList(options ...helper2.OptionFN) ([]model2.Device, int64, error) {
	option := parseOption(options...)
	var (
		from  = "FROM " + TbDevices + " d"
		where = " WHERE 1"
	)

	var params []interface{}

	if option.OrgID > 0 {
		where += " AND d.org_id=?"
		params = append(params, option.OrgID)
	}

	if option.UserID != nil {
		userID := *option.UserID
		if userID > 0 {
			from += fmt.Sprintf(` LEFT JOIN (
	SELECT d.id,p.role_id,p.action,p.effect FROM %s d 
	INNER JOIN %s p ON p.resource_class=%d AND p.resource_id=d.id 
	INNER JOIN %s r ON p.role_id=r.id
	WHERE p.role_id IN (SELECT role_id FROM %s WHERE user_id=%d)
) b ON d.id=b.id`, TbDevices, TbPolicies, resource2.Device, TbRoles, TbUserRoles, userID)

			if option.DefaultEffect == resource2.Allow {
				where += " AND ((b.action=0 AND b.effect=1) OR (ISNULL(b.action) AND ISNULL(b.effect)))"
			} else {
				where += " AND (b.action=0 AND b.effect=1)"
			}
		}
	}

	if option.GroupID != nil {
		from += " INNER JOIN " + TbDeviceGroups + " g ON d.id=g.device_id"
		where += " AND g.group_id=?"
		params = append(params, *option.GroupID)
	}

	if option.Keyword != "" {
		where += " AND d.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if option.GetTotal == nil || *option.GetTotal {
		if err := s.db.QueryRow("SELECT COUNT(DISTINCT d.id) "+from+where, params...).Scan(&total); err != nil {
			return nil, 0, lang2.InternalError(err)
		}

		if total == 0 {
			return []model2.Device{}, 0, nil
		}
	}

	where += " ORDER BY d.id ASC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	log.Trace("SELECT DISTINCT d.id " + from + where)
	rows, err := s.db.Query("SELECT DISTINCT d.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var deviceID int64

	for rows.Next() {
		err = rows.Scan(&deviceID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Device{}, total, nil
		}
		ids = append(ids, deviceID)
	}

	var result []model2.Device
	for _, id := range ids {
		device, err := s.GetDevice(id)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, device)
	}

	return result, total, nil
}

func (s *mysqlStore) loadMeasure(id int64) (model2.Measure, error) {
	var measure = NewMeasure(s, id)
	err := LoadData(s.db, TbMeasures, map[string]interface{}{
		"enable":     &measure.enable,
		"device_id":  &measure.deviceID,
		"title":      &measure.title,
		"tag":        &measure.tag,
		"kind":       &measure.kind,
		"created_at": &measure.createdAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrMeasureNotFound)
	}
	return measure, nil
}

func (s *mysqlStore) GetMeasureFromTagName(deviceID int64, tagName string) (model2.Measure, error) {
	result := <-synchronized.Do(TbMeasures, func() interface{} {
		measureID, err := s.getMeasureID(deviceID, tagName)
		if err != nil {
			return err
		}
		if measure, err := s.cache.LoadMeasure(measureID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return lang2.InternalError(err)
			}
		} else {
			return measure
		}

		measure, err := s.loadMeasure(measureID)
		if err != nil {
			return err
		}

		err = s.cache.Save(measure)
		if err != nil {
			return err
		}
		return measure
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Measure), nil
}

func (s *mysqlStore) GetMeasure(measureID int64) (model2.Measure, error) {
	result := <-synchronized.Do(TbMeasures, func() interface{} {
		if measure, err := s.cache.LoadMeasure(measureID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return measure
		}

		measure, err := s.loadMeasure(measureID)
		if err != nil {
			return err
		}

		err = s.cache.Save(measure)
		if err != nil {
			return err
		}
		return measure
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Measure), nil
}

func (s *mysqlStore) CreateMeasure(deviceID int64, title string, tag string, kind resource2.MeasureKind) (model2.Measure, error) {
	result := <-synchronized.Do(TbMeasures, func() interface{} {
		data := map[string]interface{}{
			"enable":     status2.Enable,
			"device_id":  deviceID,
			"title":      title,
			"tag":        tag,
			"kind":       kind,
			"created_at": time.Now(),
		}

		measureID, err := CreateData(s.db, TbMeasures, data)
		if err != nil {
			return err
		}

		measure, err := s.loadMeasure(measureID)
		if err != nil {
			return err
		}

		err = s.cache.Save(measure)
		if err != nil {
			return err
		}
		return measure
	})
	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Measure), nil
}

func (s *mysqlStore) RemoveMeasure(measureID int64) error {
	err := RemoveData(s.db, TbMeasures, "id=?", measureID)
	if err != nil {
		return err
	}

	s.cache.Remove(&Measure{id: measureID})
	return nil
}

func (s *mysqlStore) GetMeasureList(options ...helper2.OptionFN) ([]model2.Measure, int64, error) {
	option := parseOption(options...)

	var (
		from  = "FROM " + TbMeasures + " m"
		where = " WHERE 1"
	)

	var params []interface{}

	if option.UserID != nil {
		userID := *option.UserID
		if userID > 0 {
			from += fmt.Sprintf(` LEFT JOIN (
SELECT m.id,p.role_id,p.action,p.effect FROM %s m
INNER JOIN %s p ON p.resource_class=%d AND p.resource_id=m.id
INNER JOIN %s r ON p.role_id=r.id
WHERE p.role_id IN (SELECT role_id FROM %s WHERE user_id=%d)
) b ON m.id=b.id`, TbMeasures, TbPolicies, resource2.Measure, TbRoles, TbUserRoles, userID)

			if option.DefaultEffect == resource2.Allow {
				where += " AND ((b.action=0 AND b.effect=1) OR (ISNULL(b.action) AND ISNULL(b.effect)))"
			} else {
				where += " AND (b.action=0 AND b.effect=1)"
			}
		}
	}

	if option.DeviceID > 0 {
		where += " AND m.device_id=?"
		params = append(params, option.DeviceID)
	}

	if option.Kind != resource2.AllKind {
		where += " AND m.kind=?"
		params = append(params, option.Kind)
	}

	if option.Keyword != "" {
		where += " AND m.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(DISTINCT m.id) "+from+where, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.Measure{}, 0, nil
	}

	where += " ORDER BY m.id ASC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	log.Trace("SELECT DISTINCT d.id " + from + where)
	rows, err := s.db.Query("SELECT DISTINCT m.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var measureID int64

	for rows.Next() {
		err = rows.Scan(&measureID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Measure{}, total, nil
		}
		ids = append(ids, measureID)
	}

	var result []model2.Measure
	for _, id := range ids {
		measure, err := s.GetMeasure(id)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, measure)
	}

	return result, total, nil
}

func (s *mysqlStore) loadEquipment(id int64) (model2.Equipment, error) {
	var equipment = NewEquipment(s, id)
	err := LoadData(s.db, TbEquipments, map[string]interface{}{
		"org_id":     &equipment.orgID,
		"enable":     &equipment.enable,
		"title":      &equipment.title,
		"desc":       &equipment.desc,
		"created_at": &equipment.createdAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, lang2.Error(lang2.ErrEquipmentNotFound)
	}
	return equipment, nil
}

func (s *mysqlStore) GetEquipment(equipmentID int64) (model2.Equipment, error) {
	result := <-synchronized.Do(TbEquipments, func() interface{} {
		if equipment, err := s.cache.LoadEquipment(equipmentID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return equipment
		}

		equipment, err := s.loadEquipment(equipmentID)
		if err != nil {
			return err
		}

		err = s.cache.Save(equipment)
		if err != nil {
			return err
		}

		return equipment
	})

	if err, ok := result.(error); ok {
		return nil, err
	}

	return result.(model2.Equipment), nil
}

func (s *mysqlStore) CreateEquipment(org interface{}, title, desc string) (model2.Equipment, error) {
	result := <-synchronized.Do(TbEquipments, func() interface{} {
		orgID, err := s.getOrganizationID(org)
		if err != nil {
			return err
		}

		equipmentID, err := CreateData(s.db, TbEquipments, map[string]interface{}{
			"enable":     status2.Enable,
			"org_id":     orgID,
			"title":      title,
			"desc":       desc,
			"created_at": time.Now(),
		})
		if err != nil {
			return lang2.InternalError(err)
		}

		equipment, err := s.loadEquipment(equipmentID)
		if err != nil {
			return err
		}
		err = s.cache.Save(equipment)
		if err != nil {
			return err
		}
		return equipment
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Equipment), nil
}

func (s *mysqlStore) RemoveEquipment(equipmentID int64) error {
	err := RemoveData(s.db, TbEquipments, "id=?", equipmentID)
	if err != nil {
		return err
	}
	s.cache.Remove(&Equipment{id: equipmentID})
	return nil
}

func (s *mysqlStore) GetEquipmentList(options ...helper2.OptionFN) ([]model2.Equipment, int64, error) {
	option := parseOption(options...)
	var (
		from  = "FROM " + TbEquipments + " e"
		where = " WHERE 1"
	)

	var params []interface{}

	if option.OrgID > 0 {
		where += " AND e.org_id=?"
		params = append(params, option.OrgID)
	}

	if option.UserID != nil {
		userID := *option.UserID
		if userID > 0 {
			from += fmt.Sprintf(` LEFT JOIN (
	SELECT e.id,p.role_id,p.action,p.effect FROM %s e 
	INNER JOIN %s p ON p.resource_class=%d AND p.resource_id=e.id 
	INNER JOIN %s r ON p.role_id=r.id
	WHERE p.role_id IN (SELECT role_id FROM %s WHERE user_id=%d)
) b ON e.id=b.id`, TbEquipments, TbPolicies, resource2.Equipment, TbRoles, TbUserRoles, userID)

			if option.DefaultEffect == resource2.Allow {
				where += " AND ((b.action=0 AND b.effect=1) OR (ISNULL(b.action) AND ISNULL(b.effect)))"
			} else {
				where += " AND (b.action=0 AND b.effect=1)"
			}
		}
	}

	if option.GroupID != nil {
		from += " INNER JOIN " + TbEquipmentGroups + " g ON e.id=g.equipment_id"
		where += " AND g.group_id=?"
		params = append(params, *option.GroupID)
	}

	if option.Keyword != "" {
		where += " AND e.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if option.GetTotal == nil || *option.GetTotal {
		if err := s.db.QueryRow("SELECT COUNT(DISTINCT e.id) "+from+where, params...).Scan(&total); err != nil {
			return nil, 0, lang2.InternalError(err)
		}

		if total == 0 {
			return []model2.Equipment{}, 0, nil
		}
	}

	where += " ORDER BY e.id ASC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	rows, err := s.db.Query("SELECT DISTINCT e.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var equipmentID int64

	for rows.Next() {
		err = rows.Scan(&equipmentID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Equipment{}, total, nil
		}

		ids = append(ids, equipmentID)
	}

	var result []model2.Equipment
	for _, id := range ids {
		equipment, err := s.GetEquipment(id)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, equipment)
	}

	return result, total, nil
}

func (s *mysqlStore) loadState(id int64) (model2.State, error) {
	var state = NewState(s, id)
	err := LoadData(s.db, TbStates, map[string]interface{}{
		"enable":       &state.enable,
		"title":        &state.title,
		"desc?":        &state.desc,
		"equipment_id": &state.equipmentID,
		"measure_id":   &state.measureID,
		"script":       &state.script,
		"created_at":   &state.createdAt,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrStateNotFound)
	}
	return state, nil
}

func (s *mysqlStore) GetState(stateID int64) (model2.State, error) {
	result := <-synchronized.Do(TbStates, func() interface{} {
		if state, err := s.cache.LoadState(stateID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return state
		}

		state, err := s.loadState(stateID)
		if err != nil {
			return err
		}

		err = s.cache.Save(state)
		if err != nil {
			return err
		}
		return state
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.State), nil
}

func (s *mysqlStore) CreateState(equipmentID, measureID int64, title, desc, script string) (model2.State, error) {
	result := <-synchronized.Do(TbStates, func() interface{} {
		data := map[string]interface{}{
			"enable":       status2.Enable,
			"title":        title,
			"desc":         desc,
			"equipment_id": equipmentID,
			"measure_id":   measureID,
			"script":       script,
			"created_at":   time.Now(),
		}

		stateID, err := CreateData(s.db, TbStates, data)
		if err != nil {
			return err
		}
		state, err := s.loadState(stateID)
		if err != nil {
			return err
		}

		err = s.cache.Save(state)
		if err != nil {
			return err
		}
		return state
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.State), nil
}

func (s *mysqlStore) RemoveState(stateID int64) error {
	err := RemoveData(s.db, TbStates, "id=?", stateID)
	if err != nil {
		return err
	}

	s.cache.Remove(&State{id: stateID})
	return nil
}

func (s *mysqlStore) GetStateList(options ...helper2.OptionFN) ([]model2.State, int64, error) {
	option := parseOption(options...)

	var (
		from  = "FROM " + TbStates + " s"
		where = " WHERE 1"
	)

	var params []interface{}

	if option.UserID != nil {
		userID := *option.UserID
		if userID > 0 {
			from += fmt.Sprintf(` LEFT JOIN (
SELECT s.id,p.role_id,p.action,p.effect FROM %s s
INNER JOIN %s p ON p.resource_class=%d AND p.resource_id=s.id
INNER JOIN %s r ON p.role_id=r.id
WHERE p.role_id IN (SELECT role_id FROM %s WHERE user_id=%d)
) b ON s.id=b.id`, TbStates, TbPolicies, resource2.State, TbRoles, TbUserRoles, userID)

			if option.DefaultEffect == resource2.Allow {
				where += " AND ((b.action=0 AND b.effect=1) OR (ISNULL(b.action) AND ISNULL(b.effect)))"
			} else {
				where += " AND (b.action=0 AND b.effect=1)"
			}
		}
	}

	if option.EquipmentID > 0 {
		where += " AND s.equipment_id=?"
		params = append(params, option.EquipmentID)
	}

	//if option.Kind != resource.AllKind {
	//	where += " AND s.kind=?"
	//	params = append(params, option.Kind)
	//}

	if option.Keyword != "" {
		where += " AND s.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(DISTINCT s.id) "+from+where, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.State{}, 0, nil
	}

	where += " ORDER BY s.id ASC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	log.Trace("SELECT DISTINCT s.id " + from + where)
	rows, err := s.db.Query("SELECT DISTINCT s.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var stateID int64

	for rows.Next() {
		err = rows.Scan(&stateID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.State{}, total, nil
		}

		ids = append(ids, stateID)
	}

	var result []model2.State
	for _, stateID := range ids {
		state, err := s.GetState(stateID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, state)
	}

	return result, total, nil
}

func (s *mysqlStore) loadAlarm(id int64) (model2.Alarm, error) {
	var alarm = NewAlarm(s, id)
	err := LoadData(s.db, TbAlarms, map[string]interface{}{
		"org_id":     &alarm.orgID,
		"status":     &alarm.status,
		"device_id":  &alarm.deviceID,
		"measure_id": &alarm.measureID,
		"extra":      &alarm.extra,
		"created_at": &alarm.createdAt,
		"updated_at": &alarm.updatedAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrAlarmNotFound)
	}
	return alarm, nil
}

func (s *mysqlStore) GetAlarm(alarmID int64) (model2.Alarm, error) {
	result := <-synchronized.Do(TbAlarms, func() interface{} {
		if alarm, err := s.cache.LoadAlarm(alarmID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return alarm
		}

		alarm, err := s.loadAlarm(alarmID)
		if err != nil {
			return err
		}

		err = s.cache.Save(alarm)
		if err != nil {
			return err
		}
		return alarm
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Alarm), nil
}

func (s *mysqlStore) CreateAlarm(device model2.Device, measureID int64, data map[string]interface{}) (model2.Alarm, error) {
	result := <-synchronized.Do(TbAlarms, func() interface{} {
		extra, err := json.Marshal(data)
		if err != nil {
			return lang2.InternalError(err)
		}

		now := time.Now()
		data := map[string]interface{}{
			"org_id":     device.OrganizationID(),
			"status":     status2.Unconfirmed,
			"device_id":  device.GetID(),
			"measure_id": measureID,
			"extra":      extra,
			"created_at": now,
			"updated_at": now,
		}

		alarmID, err := CreateData(s.db, TbAlarms, data)
		if err != nil {
			return err
		}
		alarm, err := s.loadAlarm(alarmID)
		if err != nil {
			return err
		}

		err = s.cache.Save(alarm)
		if err != nil {
			return err
		}
		return alarm
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model2.Alarm), nil
}

func (s *mysqlStore) RemoveAlarm(alarmID int64) error {
	err := RemoveData(s.db, TbAlarms, "id=?", alarmID)
	if err != nil {
		return lang2.InternalError(err)
	}

	s.cache.Remove(&Alarm{id: alarmID})
	return nil
}

func (s *mysqlStore) GetLastUnconfirmedAlarm(device model2.Device, measureID int64) (model2.Alarm, error) {
	const (
		SQL = "SELECT id FROM " + TbAlarms + " WHERE device_id=? AND measure_id=? AND status=? LIMIT 1"
	)

	var alarmID int64
	if err := s.db.QueryRow(SQL, device.GetID(), measureID, status2.Unconfirmed).Scan(&alarmID); err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrAlarmNotFound)
	}

	return s.GetAlarm(alarmID)
}

func (s *mysqlStore) GetAlarmList(start, end *time.Time, options ...helper2.OptionFN) ([]model2.Alarm, int64, error) {
	option := parseOption(options...)

	var (
		from  = "FROM " + TbAlarms + " a"
		where = " WHERE 1"
	)

	var params []interface{}
	if option.UserID != nil {
		userID := *option.UserID
		if userID > 0 {
			from += fmt.Sprintf(` LEFT JOIN (
SELECT m.id,p.role_id,p.action,p.effect FROM %s m
INNER JOIN %s p ON p.resource_class=%d AND p.resource_id=m.id
INNER JOIN %s r ON p.role_id=r.id
WHERE p.role_id IN (SELECT role_id FROM %s WHERE user_id=%d)
) b ON a.measure_id=b.id`, TbMeasures, TbPolicies, resource2.Measure, TbRoles, TbUserRoles, userID)

			if option.DefaultEffect == resource2.Allow {
				where += " AND ((b.action=0 AND b.effect=1) OR (ISNULL(b.action) AND ISNULL(b.effect)))"
			} else {
				where += " AND (b.action=0 AND b.effect=1)"
			}
		}
	}

	if option.DeviceID > 0 {
		where += " AND a.device_id=?"
		params = append(params, option.DeviceID)
	}

	if option.MeasureID > 0 {
		where += " AND a.measure_id=?"
		params = append(params, option.MeasureID)
	}

	if start != nil {
		where += " AND a.created_at>=?"
		params = append(params, *start)
	}

	if end != nil {
		where += " AND a.created_at<?"
		params = append(params, *end)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(DISTINCT a.id) "+from+where, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.Alarm{}, 0, nil
	}

	where += " ORDER BY a.updated_at DESC"

	if option.Limit > 0 {
		where += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		where += " OFFSET ?"
		params = append(params, option.Offset)
	}

	log.Trace("SELECT DISTINCT a.id " + from + where)
	rows, err := s.db.Query("SELECT DISTINCT a.id,a.updated_at "+from+where, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var (
		ids     []int64
		alarmID int64
		updated time.Time
	)

	for rows.Next() {
		err = rows.Scan(&alarmID, &updated)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.Alarm{}, total, nil
		}
		ids = append(ids, alarmID)
	}

	var result []model2.Alarm
	for _, id := range ids {
		alarm, err := s.GetAlarm(id)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, alarm)
	}

	return result, total, nil
}

func (s *mysqlStore) GetResourceList(class resource2.Class, options ...helper2.OptionFN) ([]model2.Resource, int64, error) {
	var result []model2.Resource
	switch class {
	case resource2.Api:
		res, total, err := s.GetApiResourceList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, r := range res {
			result = append(result, r)
		}
		return result, total, nil
	case resource2.Group:
		groups, total, err := s.GetGroupList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, group := range groups {
			result = append(result, group)
		}
		return result, total, nil
	case resource2.Device:
		devices, total, err := s.GetDeviceList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, device := range devices {
			result = append(result, device)
		}
		return result, total, nil
	case resource2.Measure:
		measures, total, err := s.GetMeasureList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, measure := range measures {
			result = append(result, measure)
		}
		return result, total, nil
	case resource2.Equipment:
		equipments, total, err := s.GetEquipmentList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, equipment := range equipments {
			result = append(result, equipment)
		}
		return result, total, nil
	case resource2.State:
		states, total, err := s.GetStateList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, state := range states {
			result = append(result, state)
		}
		return result, total, nil
	default:
		panic(errors.New("GetResourceList: unknown resource type"))
	}
}

func (s *mysqlStore) GetResource(class resource2.Class, resourceID int64) (model2.Resource, error) {
	switch class {
	case resource2.Api:
		res, err := s.GetApiResource(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource2.Group:
		res, err := s.GetGroup(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource2.Device:
		res, err := s.GetDevice(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource2.Measure:
		res, err := s.GetMeasure(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource2.Equipment:
		res, err := s.GetEquipment(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource2.State:
		res, err := s.GetState(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, lang2.Error(lang2.ErrInvalidResourceClassID)
	}
}

func (s *mysqlStore) loadApiResource(resID int64) (model2.ApiResource, error) {
	var apiRes = NewApiResource(s, resID)
	err := LoadData(s.db, TbApiResources, map[string]interface{}{
		"name":  &apiRes.name,
		"title": &apiRes.title,
		"desc":  &apiRes.desc,
	}, "id=?", resID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang2.InternalError(err)
		}
		return nil, lang2.Error(lang2.ErrApiResourceNotFound)
	}
	return apiRes, nil
}

func (s *mysqlStore) GetApiResource(res interface{}) (model2.ApiResource, error) {
	result := <-synchronized.Do(TbApiResources, func() interface{} {
		var resID int64
		switch v := res.(type) {
		case int64:
			resID = v
		case string:
			err := LoadData(s.db, TbApiResources, map[string]interface{}{
				"id": &resID,
			}, "name=?", v)
			if err != nil {
				if err != sql.ErrNoRows {
					return lang2.InternalError(err)
				}
				return lang2.Error(lang2.ErrApiResourceNotFound)
			}
		default:
			panic(errors.New("GetApiResource: unknown api resource"))
		}

		if res, err := s.cache.LoadApiResource(resID); err != nil {
			if err != lang2.Error(lang2.ErrCacheNotFound) {
				return err
			}
		} else {
			return res
		}

		res, err := s.loadApiResource(resID)
		if err != nil {
			return err
		}
		err = s.cache.Save(res)
		if err != nil {
			return err
		}

		return res
	})

	if err, ok := result.(error); ok {
		return nil, err
	}

	return result.(model2.ApiResource), nil
}

func (s *mysqlStore) GetApiResourceList(options ...helper2.OptionFN) ([]model2.ApiResource, int64, error) {
	option := parseOption(options...)

	var (
		fromSQL = "FROM " + TbApiResources + " WHERE title != ''"
	)

	var params []interface{}
	if option.Name != "" {
		fromSQL += " AND name REGEXP ?"
		params = append(params, option.Name)
	}

	if option.Keyword != "" {
		fromSQL += " AND title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
		return nil, 0, lang2.InternalError(err)
	}

	if total == 0 {
		return []model2.ApiResource{}, 0, nil
	}

	fromSQL += " ORDER BY id ASC"

	if option.Limit > 0 {
		fromSQL += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		fromSQL += " OFFSET ?"
		params = append(params, option.Offset)
	}

	rows, err := s.db.Query("SELECT id "+fromSQL, params...)
	if err != nil {
		return nil, 0, lang2.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var ids []int64
	var resID int64

	for rows.Next() {
		err = rows.Scan(&resID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang2.InternalError(err)
			}
			return []model2.ApiResource{}, total, nil
		}
		ids = append(ids, resID)
	}

	_ = rows.Close()
	var result []model2.ApiResource
	for _, id := range ids {
		res, err := s.GetApiResource(id)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, res)
	}

	return result, total, nil
}

func (s *mysqlStore) InitApiResource() error {
	result := <-synchronized.Do(TbApiResources, func() interface{} {
		err := RemoveData(s.db, TbApiResources, "1")
		if err != nil {
			return err
		}
		for _, entry := range lang2.ApiResourcesMap() {
			_, err := CreateData(s.db, TbApiResources, map[string]interface{}{
				"name":  entry[0],
				"title": entry[1],
				"desc":  entry[2],
			})
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

func (s *mysqlStore) InitDefaultRoles(org interface{}) error {
	for pair, apiRes := range lang2.DefaultRoles() {
		role, err := s.createRole(s.db, org, pair[0], pair[1], pair[2])
		if err != nil {
			return err
		}

		for _, api := range apiRes {
			if api == resource2.Unknown {
				continue
			}

			res, err := s.GetApiResource(api)
			if err != nil {
				return err
			}
			_, err = role.SetPolicy(res, resource2.Invoke, resource2.Allow, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
