package mysqlStore

import (
	"database/sql"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/db"
	"github.com/maritimusj/centrum/web/model"
)

func getOrganizationID(db db.DB, org interface{}) (int64, error) {
	switch v := org.(type) {
	case int64:
		if exists, err := IsDataExists(db, TbOrganization, "id=?", v); err != nil {
			return 0, lang.InternalError(err)
		} else if exists {
			return v, nil
		}
	case string:
		return getOrganizationIDByName(db, v)
	case model.Organization:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang.Error(lang.ErrOrganizationNotFound)
}

func getOrganizationIDByName(db db.DB, name string) (int64, error) {
	var orgID int64
	err := LoadData(db, TbOrganization, map[string]interface{}{
		"id": &orgID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.Error(lang.ErrOrganizationNotFound)
	}
	return orgID, nil
}

func getRoleID(db db.DB, role interface{}) (int64, error) {
	switch v := role.(type) {
	case int64:
		if exists, err := IsDataExists(db, TbRoles, "id=?", v); err != nil {
			return 0, lang.InternalError(err)
		} else if exists {
			return v, nil
		}
	case string:
		return getRoleIDByName(db, v)
	case model.Role:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang.Error(lang.ErrRoleNotFound)
}

func getRoleIDByName(db db.DB, name string) (int64, error) {
	var roleID int64
	err := LoadData(db, TbRoles, map[string]interface{}{
		"id": &roleID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.Error(lang.ErrRoleNotFound)
	}
	return roleID, nil
}

func getUserID(db db.DB, user interface{}) (int64, error) {
	checkExists := func(v interface{}) error {
		if exists, err := IsDataExists(db, TbUsers, "id=?", v); err != nil {
			return lang.InternalError(err)
		} else if !exists {
			return lang.Error(lang.ErrUserNotFound)
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
		return getUserIDByName(db, v)
	case model.User:
		if v != nil {
			return v.GetID(), nil
		}
	}
	return 0, lang.Error(lang.ErrUserNotFound)
}

func getUserIDByName(db db.DB, name string) (int64, error) {
	var userID int64
	err := LoadData(db, TbUsers, map[string]interface{}{
		"id": &userID,
	}, "name=?", name)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, lang.InternalError(err)
		}
		return 0, lang.Error(lang.ErrUserNotFound)
	}
	return userID, nil
}
