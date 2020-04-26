package mysqlStore

import (
	"database/sql"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/model"
)

func (s *mysqlStore) getConfigID(cfg interface{}) (int64, error) {
	if c, err := s.cache.LoadConfig(cfg); err == nil {
		return c.GetID(), nil
	}

	switch v := cfg.(type) {
	case int64:
		if exists, err := IsDataExists(s.db, TbConfig, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang.InternalError(err)
			}
			return 0, lang.ErrConfigNotFound.Error()
		} else if exists {
			return v, nil
		}
	case string:
		return s.getConfigIDByName(v)
	case model.Config:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang.ErrConfigNotFound.Error()
}

func (s *mysqlStore) getConfigIDByName(name string) (int64, error) {
	var cfgID int64
	err := LoadData(s.db, TbConfig, map[string]interface{}{
		"id": &cfgID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.ErrConfigNotFound.Error()
	}
	return cfgID, nil
}

func (s *mysqlStore) getOrganizationID(org interface{}) (int64, error) {
	if o, err := s.cache.LoadOrganization(org); err == nil {
		return o.GetID(), nil
	}

	switch v := org.(type) {
	case int64:
		if exists, err := IsDataExists(s.db, TbOrganization, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang.InternalError(err)
			}
			return 0, lang.ErrOrganizationNotFound.Error()
		} else if exists {
			return v, nil
		}
	case string:
		return s.getOrganizationIDByName(v)
	case model.Organization:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang.ErrOrganizationNotFound.Error()
}

func (s *mysqlStore) getOrganizationIDByName(name string) (int64, error) {
	var orgID int64
	err := LoadData(s.db, TbOrganization, map[string]interface{}{
		"id": &orgID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.ErrOrganizationNotFound.Error()
	}
	return orgID, nil
}

func (s *mysqlStore) getRoleID(role interface{}) (int64, error) {
	if r, err := s.cache.LoadRole(role); err == nil {
		return r.GetID(), nil
	}

	switch v := role.(type) {
	case int64:
		if exists, err := IsDataExists(s.db, TbRoles, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang.InternalError(err)
			}
			return 0, lang.ErrRoleNotFound.Error()
		} else if exists {
			return v, nil
		}
	case string:
		return s.getRoleIDByName(v)
	case model.Role:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang.ErrRoleNotFound.Error()
}

func (s *mysqlStore) getRoleIDByName(name string) (int64, error) {
	var roleID int64
	err := LoadData(s.db, TbRoles, map[string]interface{}{
		"id": &roleID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.ErrRoleNotFound.Error()
	}
	return roleID, nil
}

func (s *mysqlStore) getUserID(user interface{}) (int64, error) {
	if u, err := s.cache.LoadUser(user); err == nil {
		return u.GetID(), nil
	}

	checkExists := func(v interface{}) error {
		if exists, err := IsDataExists(s.db, TbUsers, "id=?", v); err != nil {
			return lang.InternalError(err)
		} else if !exists {
			return lang.ErrUserNotFound.Error()
		}
		return nil
	}
	switch v := user.(type) {
	case int64:
		if err := checkExists(v); err != nil {
			return 0, err
		}
		return v, nil
	case float64:
		if err := checkExists(v); err != nil {
			return 0, err
		}
		return int64(v), nil
	case string:
		return s.getUserIDByName(v)
	case model.User:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang.ErrUserNotFound.Error()
}

func (s *mysqlStore) getUserIDByName(name string) (int64, error) {
	var userID int64
	err := LoadData(s.db, TbUsers, map[string]interface{}{
		"id": &userID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.ErrUserNotFound.Error()
	}
	return userID, nil
}

func (s *mysqlStore) getMeasureID(deviceID int64, tag interface{}) (int64, error) {
	switch v := tag.(type) {
	case int64:
		if t, err := s.cache.LoadMeasure(v); err == nil {
			return t.GetID(), nil
		}
	case string:
		if t, err := s.cache.LoadMeasure(FormatMeasureName(deviceID, v)); err == nil {
			return t.GetID(), nil
		}
	}

	switch v := tag.(type) {
	case int64:
		if exists, err := IsDataExists(s.db, TbMeasures, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang.InternalError(err)
			}
			return 0, lang.ErrMeasureNotFound.Error()
		} else if exists {
			return v, nil
		}
	case string:
		return s.getMeasureIDByTagName(deviceID, v)
	case model.Measure:
		return v.GetID(), nil
	}
	return 0, lang.ErrMeasureNotFound.Error()
}

func (s *mysqlStore) getMeasureIDByTagName(deviceID int64, tagName string) (int64, error) {
	var measureID int64
	err := LoadData(s.db, TbMeasures, map[string]interface{}{
		"id": &measureID,
	}, "device_id=? AND tag=?", deviceID, tagName)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.ErrMeasureNotFound.Error()
	}
	return measureID, nil
}
