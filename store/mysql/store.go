package mysqlStore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maritimusj/centrum/cache"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/util"
	"sync"
	"time"
)

const (
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
)

type mysqlStore struct {
	db    *sql.DB
	cache cache.Cache

	lockerMap map[string]sync.Mutex
	mu        sync.Mutex
	wg        sync.WaitGroup
	ctx       context.Context
}

func New() store.Store {
	return &mysqlStore{
		lockerMap: make(map[string]sync.Mutex),
	}
}

func parseOption(options ...store.OptionFN) *store.Option {
	option := store.Option{}
	for _, opt := range options {
		if opt != nil {
			opt(&option)
		}
	}
	return &option
}

func (s *mysqlStore) TransactionDo(fn func(db store.DB) interface{}) interface{} {
	tx, err := s.db.Begin()
	if err != nil {
		return lang.InternalError(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res := fn(tx)
	if res != nil {
		if err, ok := res.(error); ok {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return lang.InternalError(err)
	}

	return res
}

func (s *mysqlStore) Synchronized(name string, fn func() interface{}) <-chan interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.lockerMap[name]
	if ok {
		v := sync.Mutex{}
		s.lockerMap[name] = v
	}

	resultChan := make(chan interface{})
	go func() {
		v.Lock()
		s.wg.Add(1)

		defer func() {
			close(resultChan)
			s.wg.Done()
			v.Unlock()
		}()

		select {
		case <-s.ctx.Done():
			resultChan <- s.ctx.Err()
			return
		default:
			resultChan <- fn()
			return
		}

	}()

	return resultChan
}

func (s *mysqlStore) Open(ctx context.Context, option map[string]interface{}) error {
	if c, ok := option["cache"].(cache.Cache); ok {
		s.cache = c
	} else {
		panic(errors.New("invalid cache"))
	}

	if connStr, ok := option["connStr"].(string); ok {
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			return lang.InternalError(err)
		}

		ctxTimeout, _ := context.WithTimeout(ctx, time.Second*3)
		err = db.PingContext(ctxTimeout)
		if err != nil {
			return lang.InternalError(err)
		}

		s.db = db
		s.ctx = ctx
		return nil
	}
	return lang.Error(lang.ErrInvalidConnStr)
}

func (s *mysqlStore) Close() {
	if s != nil && s.db != nil {
		_ = s.db.Close()
	}
}

func (s *mysqlStore) GetResourceGroupList() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":    resource.Api,
			"title": lang.ResourceClassTitle(resource.Api),
		},
		map[string]interface{}{
			"id":    resource.Group,
			"title": lang.ResourceClassTitle(resource.Group),
		},
		map[string]interface{}{
			"id":    resource.Device,
			"title": lang.ResourceClassTitle(resource.Device),
		},
		map[string]interface{}{
			"id":    resource.Measure,
			"title": lang.ResourceClassTitle(resource.Measure),
		},
		map[string]interface{}{
			"id":    resource.Equipment,
			"title": lang.ResourceClassTitle(resource.Equipment),
		},
		map[string]interface{}{
			"id":    resource.State,
			"title": lang.ResourceClassTitle(resource.State),
		},
	}
}

func (s *mysqlStore) getUser(db DB, id int64) (*User, error) {
	var user = NewUser(s, id)
	err := LoadData(db, TbUsers, map[string]interface{}{
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
			return nil, lang.InternalError(err)
		}
		return nil, lang.Error(lang.ErrUserNotFound)
	}

	return user, nil
}

func (s *mysqlStore) GetUser(user interface{}) (model.User, error) {
	result := <-s.Synchronized(TbUsers, func() interface{} {
		var userID int64
		switch v := user.(type) {
		case int64:
			userID = v
		case float64:
			userID = int64(v)
		case string:
			id, err := getUserIDByName(s.db, v)
			if err != nil {
				return err
			}
			userID = id
		default:
			panic(errors.New("GetUser: unknown user"))
		}
		if user, err := s.cache.LoadUser(userID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
				return err
			}
		} else {
			return user
		}

		user, err := s.getUser(s.db, userID)
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

	return result.(model.User), nil
}

func (s *mysqlStore) CreateUser(name string, password []byte, role model.Role) (model.User, error) {
	result := <-s.Synchronized(TbUsers, func() interface{} {
		return s.TransactionDo(func(db store.DB) interface{} {
			passwordData, err := util.HashPassword(password)
			if err != nil {
				return lang.InternalError(err)
			}
			userID, err := CreateData(db, TbUsers, map[string]interface{}{
				"enable":     status.Enable,
				"name":       name,
				"password":   passwordData,
				"title":      name,
				"mobile":     "",
				"email":      "",
				"created_at": time.Now(),
			})

			if err != nil {
				return lang.InternalError(err)
			}

			user, err := s.getUser(db, userID)
			if err != nil {
				return err
			}

			if role == nil {
				role, err = s.createRole(db, name)
				if err != nil {
					return err
				}
			}

			err = user.SetRoles(role)
			if err != nil {
				return err
			}

			err = s.cache.Save(user)
			if err != nil {
				return err
			}
			return user
		})
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model.User), nil
}

func (s *mysqlStore) RemoveUser(userID int64) error {
	result := <-s.Synchronized(TbUsers, func() interface{} {
		err := RemoveData(s.db, TbUsers, "id=?", userID)
		if err != nil {
			return lang.InternalError(err)
		}
		s.cache.Remove(&User{id: userID})
		return nil
	})

	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetUserList(options ...store.OptionFN) ([]model.User, int64, error) {
	option := parseOption(options...)

	var (
		fromSQL = "FROM " + TbUsers + " WHERE 1"
	)

	var params []interface{}
	if option.Keyword != "" {
		fromSQL += " AND (name REGEXP ? OR title REGEXP ? OR mobile REGEXP ? OR email REGEXP ?)"
		params = append(params, option.Keyword, option.Keyword, option.Keyword, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
		return nil, 0, lang.InternalError(err)
	}

	if total == 0 {
		return []model.User{}, 0, nil
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
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.User
	var userID int64

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.User{}, total, nil
		}

		role, err := s.GetUser(userID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, role)
	}

	return result, total, nil
}

func (s *mysqlStore) loadRole(id int64) (*Role, error) {
	var role = NewRole(s, id)
	err := LoadData(s.db, TbRoles, map[string]interface{}{
		"enable":     &role.enable,
		"title":      &role.title,
		"created_at": &role.createdAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang.InternalError(err)
		}
		return nil, lang.Error(lang.ErrRoleNotFound)
	}
	return role, nil
}

func (s *mysqlStore) GetRole(roleID int64) (model.Role, error) {
	result := <-s.Synchronized(TbRoles, func() interface{} {
		if role, err := s.cache.LoadRole(roleID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
				return lang.InternalError(err)
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
	return result.(model.Role), nil
}

func (s *mysqlStore) createRole(db DB, title string) (model.Role, error) {
	result := <-s.Synchronized(TbRoles, func() interface{} {
		roleID, err := CreateData(s.db, TbRoles, map[string]interface{}{
			"enable":     status.Enable,
			"title":      title,
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

	return result.(model.Role), nil
}
func (s *mysqlStore) CreateRole(title string) (model.Role, error) {
	return s.createRole(s.db, title)
}

func (s *mysqlStore) RemoveRole(roleID int64) error {
	result := <-s.Synchronized(TbRoles, func() interface{} {
		err := RemoveData(s.db, TbRoles, "id=?", roleID)
		if err != nil {
			return lang.InternalError(err)
		}
		s.cache.Remove(&Role{id: roleID})
		return nil
	})

	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetRoleList(options ...store.OptionFN) ([]model.Role, int64, error) {
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

	if option.Keyword != "" {
		fromSQL += " AND r.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if option.GetTotal == nil || *option.GetTotal {
		if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
			return nil, 0, lang.InternalError(err)
		}

		if total == 0 {
			return []model.Role{}, 0, nil
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
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.Role
	var roleID int64

	for rows.Next() {
		err = rows.Scan(&roleID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.Role{}, total, nil
		}

		role, err := s.GetRole(roleID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, role)
	}

	return result, total, nil
}

func (s *mysqlStore) loadPolicy(id int64) (model.Policy, error) {
	var policy = NewPolicy(s, id)
	err := LoadData(s.db, TbPolicies, map[string]interface{}{
		"role_id":        &policy.roleID,
		"resource_class": &policy.resourceClass,
		"resource_id":    &policy.resourceID,
		"action":         &policy.action,
		"effect":         &policy.effect,
		"created_at":     &policy.createdAt,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang.InternalError(err)
		}
		return nil, lang.Error(lang.ErrPolicyNotFound)
	}
	return policy, nil
}

func (s *mysqlStore) GetPolicy(policyID int64) (model.Policy, error) {
	result := <-s.Synchronized(TbPolicies, func() interface{} {
		if role, err := s.cache.LoadPolicy(policyID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
				return lang.InternalError(err)
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
	return result.(model.Policy), nil
}

func (s *mysqlStore) CreatePolicyIsNotExists(roleID int64, res resource.Resource, action resource.Action, defaultEffect resource.Effect) (model.Policy, error) {
	result := <-s.Synchronized(TbPolicies, func() interface{} {
		var policyID int64
		err := LoadData(s.db, TbPolicies, map[string]interface{}{
			"id": &policyID,
		}, "role_id=? AND resource_class=? AND resource_id=? AND action=?", roleID, res.ResourceClass(), res.ResourceID(), action)

		if err != nil {
			if err != sql.ErrNoRows {
				return lang.InternalError(err)
			}
			policyID, err = CreateData(s.db, TbPolicies, map[string]interface{}{
				"role_id":        roleID,
				"resource_class": res.ResourceClass(),
				"resource_id":    res.ResourceID(),
				"action":         action,
				"effect":         defaultEffect,
				"created_at":     time.Now(),
			})
			if err != nil {
				return lang.InternalError(err)
			}
		}

		policy, err := s.GetPolicy(policyID)
		if err != nil {
			return err
		}
		return policy
	})

	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(model.Policy), nil
}

func (s *mysqlStore) CreatePolicy(roleID int64, res resource.Resource, action resource.Action, effect resource.Effect) (model.Policy, error) {
	result := <-s.Synchronized(TbPolicies, func() interface{} {
		policyID, err := CreateData(s.db, TbPolicies, map[string]interface{}{
			"enable":         status.Enable,
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
	return result.(model.Policy), nil
}

func (s *mysqlStore) RemovePolicy(policyID int64) error {
	result := <-s.Synchronized(TbPolicies, func() interface{} {
		err := RemoveData(s.db, TbPolicies, "id=?", policyID)
		if err != nil {
			return lang.InternalError(err)
		}
		s.cache.Remove(&Policy{id: policyID})
		return nil
	})

	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetPolicyList(res resource.Resource, options ...store.OptionFN) ([]model.Policy, int64, error) {
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
		return nil, 0, lang.InternalError(err)
	}

	if total == 0 {
		return []model.Policy{}, 0, nil
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
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.Policy
	var policyID int64

	for rows.Next() {
		err = rows.Scan(&policyID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.Policy{}, total, nil
		}

		role, err := s.GetPolicy(policyID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, role)
	}
	return result, total, nil
}

func (s *mysqlStore) loadGroup(id int64) (model.Group, error) {
	var group = NewGroup(s, id)

	err := LoadData(s.db, TbGroups, map[string]interface{}{
		"parent_id":  &group.parentID,
		"title":      &group.title,
		"created_at": &group.createdAt,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang.InternalError(err)
		}
		return nil, lang.Error(lang.ErrGroupNotFound)
	}
	return group, nil
}

func (s *mysqlStore) GetGroup(groupID int64) (model.Group, error) {
	result := <-s.Synchronized(TbGroups, func() interface{} {
		if group, err := s.cache.LoadGroup(groupID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
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
	return result.(model.Group), nil
}

func (s *mysqlStore) CreateGroup(title string, parentID int64) (model.Group, error) {
	result := <-s.Synchronized(TbGroups, func() interface{} {
		data := map[string]interface{}{
			"enable":     status.Enable,
			"parent_id":  parentID,
			"title":      title,
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
	return result.(model.Group), nil
}

func (s *mysqlStore) RemoveGroup(groupID int64) error {
	result := <-s.Synchronized(TbGroups, func() interface{} {
		err := RemoveData(s.db, TbGroups, "id=?", groupID)
		if err != nil {
			return err
		}

		s.cache.Remove(&Group{id: groupID})
		return nil
	})

	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetGroupList(options ...store.OptionFN) ([]model.Group, int64, error) {
	option := parseOption(options...)
	var (
		fromSQL = "FROM " + TbGroups + " WHERE 1"
	)

	var params []interface{}

	if option.ParentID != nil {
		fromSQL += " AND parent_id=?"
		params = append(params, *option.ParentID)
	}

	if option.Keyword != "" {
		fromSQL += " AND title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
		return nil, 0, lang.InternalError(err)
	}

	if total == 0 {
		return []model.Group{}, 0, nil
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
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.Group
	var groupID int64

	for rows.Next() {
		err = rows.Scan(&groupID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.Group{}, total, nil
		}

		group, err := s.GetGroup(groupID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, group)
	}

	return result, total, nil
}

func (s *mysqlStore) loadDevice(id int64) (model.Device, error) {
	var device = NewDDevice(s, id)
	err := LoadData(s.db, TbDevices, map[string]interface{}{
		"enable":     &device.enable,
		"title":      &device.title,
		"options":    &device.options,
		"created_at": &device.createdAt,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, lang.InternalError(err)
		}
		return nil, lang.Error(lang.ErrDeviceNotFound)
	}
	return device, nil
}

func (s *mysqlStore) GetDevice(deviceID int64) (model.Device, error) {
	result := <-s.Synchronized(TbDevices, func() interface{} {
		if device, err := s.cache.LoadDevice(deviceID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
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
	return result.(model.Device), nil
}

func (s *mysqlStore) CreateDevice(title string, data map[string]interface{}) (model.Device, error) {
	result := <-s.Synchronized(TbDevices, func() interface{} {
		o, err := json.Marshal(data)
		if err != nil {
			return lang.InternalError(err)
		}

		deviceID, err := CreateData(s.db, TbDevices, map[string]interface{}{
			"enable":     status.Enable,
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
	return result.(model.Device), nil
}

func (s *mysqlStore) RemoveDevice(deviceID int64) error {
	result := <-s.Synchronized(TbDevices, func() interface{} {
		err := RemoveData(s.db, TbDevices, "id=?", deviceID)
		if err != nil {
			return lang.InternalError(err)
		}

		s.cache.Remove(&Device{id: deviceID})
		return nil
	})

	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetDeviceList(options ...store.OptionFN) ([]model.Device, int64, error) {
	option := parseOption(options...)
	var (
		from  = "FROM " + TbDevices + " d"
		where = " WHERE 1"
	)

	var params []interface{}

	if option.UserID != nil {
		userID := *option.UserID
		if userID > 0 {
			var joinWay = util.If(option.DefaultEffect == resource.Allow, "LEFT", "INNER")
			from += fmt.Sprintf(" %s JOIN %s p ON p.resource_class=%d AND p.resource_id=d.id", joinWay, TbPolicies, resource.Device)
			where += " AND p.role_id IN (SELECT role_id FROM " + TbUserRoles + " WHERE user_id=?)"
			params = append(params, userID)
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
		if err := s.db.QueryRow("SELECT COUNT(*) "+from+where, params...).Scan(&total); err != nil {
			return nil, 0, lang.InternalError(err)
		}

		if total == 0 {
			return []model.Device{}, 0, nil
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

	println("SELECT d.id " + from + where)
	rows, err := s.db.Query("SELECT d.id "+from+where, params...)
	if err != nil {
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.Device
	var deviceID int64

	for rows.Next() {
		err = rows.Scan(&deviceID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.Device{}, total, nil
		}

		device, err := s.GetDevice(deviceID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, device)
	}

	return result, total, nil
}

func (s *mysqlStore) loadMeasure(id int64) (model.Measure, error) {
	var measure = Measure{id: id, store: s}
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
			return nil, lang.InternalError(err)
		}
		return nil, lang.Error(lang.ErrMeasureNotFound)
	}
	return &measure, nil
}

func (s *mysqlStore) GetMeasure(measureID int64) (model.Measure, error) {
	result := <-s.Synchronized(TbMeasures, func() interface{} {
		if measure, err := s.cache.LoadMeasure(measureID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
				return lang.InternalError(err)
			}
		} else {
			return measure
		}

		role, err := s.loadMeasure(measureID)
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
	return result.(model.Measure), nil
}

func (s *mysqlStore) CreateMeasure(deviceID int64, title string, tag string, kind model.MeasureKind) (model.Measure, error) {
	result := <-s.Synchronized(TbMeasures, func() interface{} {
		data := map[string]interface{}{
			"enable":    status.Enable,
			"device_id": deviceID,
			"title":     title,
			"tag":       tag,
			"kind":      kind,
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
	return result.(model.Measure), nil
}

func (s *mysqlStore) RemoveMeasure(measureID int64) error {
	result := <-s.Synchronized(TbMeasures, func() interface{} {
		err := RemoveData(s.db, TbMeasures, "id=?", measureID)
		if err != nil {
			return err
		}
		s.cache.Remove(&Measure{id: measureID})
		return nil
	})
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetMeasureList(options ...store.OptionFN) ([]model.Measure, int64, error) {
	option := parseOption(options...)

	var (
		fromSQL = "FROM " + TbMeasures + " WHERE 1"
	)

	var params []interface{}

	if option.DeviceID > 0 {
		fromSQL += " AND device_id=?"
		params = append(params, option.DeviceID)
	}

	if option.Kind != model.AllKind {
		fromSQL += " AND kind=?"
		params = append(params, option.Kind)
	}

	if option.Keyword != "" {
		fromSQL += " AND title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
		return nil, 0, lang.InternalError(err)
	}

	if total == 0 {
		return []model.Measure{}, 0, nil
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
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.Measure
	var measureID int64

	for rows.Next() {
		err = rows.Scan(&measureID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.Measure{}, total, nil
		}

		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, measure)
	}

	return result, total, nil
}

func (s *mysqlStore) loadEquipment(id int64) (model.Equipment, error) {
	var equipment = NewEquipment(s, id)
	err := LoadData(s.db, TbEquipments, map[string]interface{}{
		"enable":     &equipment.enable,
		"title":      &equipment.title,
		"desc":       &equipment.desc,
		"created_at": &equipment.createdAt,
	}, "id=?", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, lang.Error(lang.ErrEquipmentNotFound)
	}
	return equipment, nil
}

func (s *mysqlStore) GetEquipment(equipmentID int64) (model.Equipment, error) {
	result := <-s.Synchronized(TbEquipments, func() interface{} {
		if equipment, err := s.cache.LoadEquipment(equipmentID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
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

	return result.(model.Equipment), nil
}

func (s *mysqlStore) CreateEquipment(title, desc string) (model.Equipment, error) {
	result := <-s.Synchronized(TbEquipments, func() interface{} {
		equipmentID, err := CreateData(s.db, TbEquipments, map[string]interface{}{
			"enable":     status.Enable,
			"title":      title,
			"desc":       desc,
			"created_at": time.Now(),
		})
		if err != nil {
			return lang.InternalError(err)
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
	return result.(model.Equipment), nil
}

func (s *mysqlStore) RemoveEquipment(equipmentID int64) error {
	result := <-s.Synchronized(TbEquipments, func() interface{} {
		err := RemoveData(s.db, TbEquipments, "id=?", equipmentID)
		if err != nil {
			return err
		}
		s.cache.Remove(&Equipment{id: equipmentID})
		return nil
	})

	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetEquipmentList(options ...store.OptionFN) ([]model.Equipment, int64, error) {
	option := parseOption(options...)
	var (
		fromSQL = "FROM " + TbEquipments + " e"
	)

	var params []interface{}
	if option.GroupID != nil {
		fromSQL += " INNER JOIN " + TbEquipmentGroups + " g ON e.id=g.equip_id WHERE g.group_id=?"
		params = append(params, *option.GroupID)
	}

	if option.Keyword != "" {
		fromSQL += " AND e.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
		return nil, 0, lang.InternalError(err)
	}

	if total == 0 {
		return []model.Equipment{}, 0, nil
	}

	fromSQL += " ORDER BY e.id ASC"

	if option.Limit > 0 {
		fromSQL += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		fromSQL += " OFFSET ?"
		params = append(params, option.Offset)
	}
	rows, err := s.db.Query("SELECT e.id "+fromSQL, params...)
	if err != nil {
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.Equipment
	var roleID int64

	for rows.Next() {
		err = rows.Scan(&roleID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.Equipment{}, total, nil
		}

		role, err := s.GetEquipment(roleID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, role)
	}

	return result, total, nil
}

func (s *mysqlStore) loadState(id int64) (model.State, error) {
	var state = NewState(s, id)
	err := LoadData(s.db, TbStates, map[string]interface{}{
		"enable":       &state.enable,
		"title":        &state.title,
		"equipment_id": &state.equipmentID,
		"measure_id":   &state.measureID,
		"script":       &state.script,
		"createdAt":    &state.createdAt,
	}, "id=?", id)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, lang.Error(lang.ErrStateNotFound)
	}
	return state, nil
}

func (s *mysqlStore) GetState(stateID int64) (model.State, error) {
	result := <-s.Synchronized(TbStates, func() interface{} {
		if state, err := s.cache.LoadState(stateID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
				return lang.InternalError(err)
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
	return result.(model.State), nil
}

func (s *mysqlStore) CreateState(equipmentID int64, measureID int64, title string, script string) (model.State, error) {
	result := <-s.Synchronized(TbStates, func() interface{} {
		data := map[string]interface{}{
			"enable":       status.Enable,
			"title":        title,
			"equipment_id": equipmentID,
			"measure_id":   measureID,
			"script":       script,
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
	return result.(model.State), nil
}

func (s *mysqlStore) RemoveState(stateID int64) error {
	result := <-s.Synchronized(TbStates, func() interface{} {
		err := RemoveData(s.db, TbStates, "id=?", stateID)
		if err != nil {
			return err
		}
		s.cache.Remove(&State{id: stateID})
		return nil
	})

	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *mysqlStore) GetStateList(options ...store.OptionFN) ([]model.State, int64, error) {
	option := parseOption(options...)

	var (
		fromSQL = "FROM " + TbStates + " s INNER JOIN " + TbMeasures + " m ON s.measure_id=m.id WHERE 1"
	)

	var params []interface{}

	if option.EquipmentID > 0 {
		fromSQL += " AND s.equipment_id=?"
		params = append(params, option.EquipmentID)
	}

	if option.Kind != model.AllKind {
		fromSQL += " AND m.kind=?"
		params = append(params, option.Kind)
	}

	if option.Keyword != "" {
		fromSQL += " AND s.title REGEXP ?"
		params = append(params, option.Keyword)
	}

	var total int64
	if err := s.db.QueryRow("SELECT COUNT(*) "+fromSQL, params...).Scan(&total); err != nil {
		return nil, 0, lang.InternalError(err)
	}

	if total == 0 {
		return []model.State{}, 0, nil
	}

	fromSQL += " ORDER BY s.id ASC"

	if option.Limit > 0 {
		fromSQL += " LIMIT ?"
		params = append(params, option.Limit)
	}

	if option.Offset > 0 {
		fromSQL += " OFFSET ?"
		params = append(params, option.Offset)
	}

	rows, err := s.db.Query("SELECT s.id "+fromSQL, params...)
	if err != nil {
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.State
	var stateID int64

	for rows.Next() {
		err = rows.Scan(&stateID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.State{}, total, nil
		}

		state, err := s.GetState(stateID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, state)
	}

	return result, total, nil
}

func (s *mysqlStore) GetResourceList(class resource.Class, options ...store.OptionFN) ([]resource.Resource, int64, error) {
	var result []resource.Resource
	switch class {
	case resource.Api:
		res, total, err := s.GetApiResourceList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, r := range res {
			result = append(result, r)
		}
		return result, total, nil
	case resource.Group:
		groups, total, err := s.GetGroupList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, group := range groups {
			result = append(result, group)
		}
		return result, total, nil
	case resource.Device:
		devices, total, err := s.GetDeviceList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, device := range devices {
			result = append(result, device)
		}
		return result, total, nil
	case resource.Measure:
		measures, total, err := s.GetMeasureList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, measure := range measures {
			result = append(result, measure)
		}
		return result, total, nil
	case resource.Equipment:
		equipments, total, err := s.GetEquipmentList(options...)
		if err != nil {
			return nil, 0, err
		}
		for _, equipment := range equipments {
			result = append(result, equipment)
		}
		return result, total, nil
	case resource.State:
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

func (s *mysqlStore) GetResource(class resource.Class, resourceID int64) (resource.Resource, error) {
	switch class {
	case resource.Api:
		res, err := s.GetApiResource(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource.Group:
		res, err := s.GetGroup(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource.Device:
		res, err := s.GetDevice(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource.Measure:
		res, err := s.GetMeasure(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource.Equipment:
		res, err := s.GetEquipment(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	case resource.State:
		res, err := s.GetState(resourceID)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		panic(errors.New("GetResource: unknown resource class"))
	}
}

func (s *mysqlStore) loadApiResource(resID int64) (model.ApiResource, error) {
	var apiRes = NewApiResource(s, resID)
	err := LoadData(s.db, TbApiResources, map[string]interface{}{
		"`name`":  &apiRes.name,
		"`title`": &apiRes.title,
		"`desc`":  &apiRes.desc,
	}, "id=?", resID)
	if err != nil {
		return nil, err
	}
	return apiRes, nil
}

func (s *mysqlStore) GetApiResource(res interface{}) (model.ApiResource, error) {
	result := <-s.Synchronized(TbApiResources, func() interface{} {
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
					return lang.InternalError(err)
				}
				return lang.Error(lang.ErrApiResourceNotFound)
			}
		default:
			panic(errors.New("GetApiResource: unknown api resource"))
		}

		if res, err := s.cache.LoadApiResource(resID); err != nil {
			if err != lang.Error(lang.ErrCacheNotFound) {
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

	return result.(model.ApiResource), nil
}

func (s *mysqlStore) GetApiResourceList(options ...store.OptionFN) ([]model.ApiResource, int64, error) {
	option := parseOption(options...)

	var (
		fromSQL = "FROM " + TbApiResources + " WHERE 1"
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
		return nil, 0, lang.InternalError(err)
	}

	if total == 0 {
		return []model.ApiResource{}, 0, nil
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
		return nil, 0, lang.InternalError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []model.ApiResource
	var stateID int64

	for rows.Next() {
		err = rows.Scan(&stateID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, 0, lang.InternalError(err)
			}
			return []model.ApiResource{}, total, nil
		}

		res, err := s.GetApiResource(stateID)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, res)
	}

	return result, total, nil
}

func (s *mysqlStore) InitApiResource() error {
	result := <-s.Synchronized(TbApiResources, func() interface{} {
		return s.TransactionDo(func(db store.DB) interface{} {
			err := RemoveData(db, TbApiResources, "1")
			if err != nil {
				return err
			}
			for _, entry := range lang.ApiResourcesMap {
				_, err := CreateData(db, TbApiResources, map[string]interface{}{
					"`name`":  entry[0],
					"`title`": entry[1],
					"`desc`":  entry[2],
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
	})

	if result != nil {
		return result.(error)
	}
	return nil
}
