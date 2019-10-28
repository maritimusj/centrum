package mysqlStore

import (
	"database/sql"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	model2 "github.com/maritimusj/centrum/gate/web/model"
)

func (s *mysqlStore) getConfigID(cfg interface{}) (int64, error) {
	switch v := cfg.(type) {
	case int64:
		if _, err := s.cache.LoadConfig(v); err == nil {
			return v, nil
		}
		if exists, err := IsDataExists(s.db, TbConfig, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang2.InternalError(err)
			}
			return 0, lang2.Error(lang2.ErrConfigNotFound)
		} else if exists {
			return v, nil
		}
	case string:
		return s.getConfigIDByName(v)
	case model2.Config:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang2.Error(lang2.ErrConfigNotFound)
}

func (s *mysqlStore) getConfigIDByName(name string) (int64, error) {
	var cfgID int64
	err := LoadData(s.db, TbConfig, map[string]interface{}{
		"id": &cfgID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang2.InternalError(err)
		}
		return 0, lang2.Error(lang2.ErrConfigNotFound)
	}
	return cfgID, nil
}

func (s *mysqlStore) getOrganizationID(org interface{}) (int64, error) {
	switch v := org.(type) {
	case int64:
		if _, err := s.cache.LoadOrganization(v); err == nil {
			return v, nil
		}
		if exists, err := IsDataExists(s.db, TbOrganization, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang2.InternalError(err)
			}
			return 0, lang2.Error(lang2.ErrOrganizationNotFound)
		} else if exists {
			return v, nil
		}
	case string:
		return s.getOrganizationIDByName(v)
	case model2.Organization:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang2.Error(lang2.ErrOrganizationNotFound)
}

func (s *mysqlStore) getOrganizationIDByName(name string) (int64, error) {
	var orgID int64
	err := LoadData(s.db, TbOrganization, map[string]interface{}{
		"id": &orgID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang2.InternalError(err)
		}
		return 0, lang2.Error(lang2.ErrOrganizationNotFound)
	}
	return orgID, nil
}

func (s *mysqlStore) getRoleID(role interface{}) (int64, error) {
	switch v := role.(type) {
	case int64:
		if _, err := s.cache.LoadRole(v); err == nil {
			return v, nil
		}
		if exists, err := IsDataExists(s.db, TbRoles, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang2.InternalError(err)
			}
			return 0, lang2.Error(lang2.ErrRoleNotFound)
		} else if exists {
			return v, nil
		}
	case string:
		return s.getRoleIDByName(v)
	case model2.Role:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang2.Error(lang2.ErrRoleNotFound)
}

func (s *mysqlStore) getRoleIDByName(name string) (int64, error) {
	var roleID int64
	err := LoadData(s.db, TbRoles, map[string]interface{}{
		"id": &roleID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang2.InternalError(err)
		}
		return 0, lang2.Error(lang2.ErrRoleNotFound)
	}
	return roleID, nil
}

func (s *mysqlStore) getUserID(user interface{}) (int64, error) {
	checkExists := func(v interface{}) error {
		if exists, err := IsDataExists(s.db, TbUsers, "id=?", v); err != nil {
			return lang2.InternalError(err)
		} else if !exists {
			return lang2.Error(lang2.ErrUserNotFound)
		}
		return nil
	}
	switch v := user.(type) {
	case int64:
		if _, err := s.cache.LoadUser(v); err == nil {
			return v, nil
		}
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
	case model2.User:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang2.Error(lang2.ErrUserNotFound)
}

func (s *mysqlStore) getUserIDByName(name string) (int64, error) {
	var userID int64
	err := LoadData(s.db, TbUsers, map[string]interface{}{
		"id": &userID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang2.InternalError(err)
		}
		return 0, lang2.Error(lang2.ErrUserNotFound)
	}
	return userID, nil
}

func (s *mysqlStore) getMeasureID(deviceID int64, tag interface{}) (int64, error) {
	switch v := tag.(type) {
	case int64:
		if _, err := s.cache.LoadUser(v); err == nil {
			return v, nil
		}
		if exists, err := IsDataExists(s.db, TbMeasures, "id=?", v); err != nil {
			if err != sql.ErrNoRows {
				return 0, lang2.InternalError(err)
			}
			return 0, lang2.Error(lang2.ErrMeasureNotFound)
		} else if exists {
			return v, nil
		}
	case string:
		return s.getMeasureIDByTagName(deviceID, v)
	case model2.Measure:
		return v.GetID(), nil
	}
	return 0, lang2.Error(lang2.ErrMeasureNotFound)
}

func (s *mysqlStore) getMeasureIDByTagName(deviceID int64, tagName string) (int64, error) {
	var measureID int64
	err := LoadData(s.db, TbMeasures, map[string]interface{}{
		"id": &measureID,
	}, "device_id=? AND tag=?", deviceID, tagName)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang2.InternalError(err)
		}
		return 0, lang2.Error(lang2.ErrMeasureNotFound)
	}
	return measureID, nil
}
