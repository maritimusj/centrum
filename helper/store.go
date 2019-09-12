package helper

import (
	"database/sql"
)

type Option struct {
	Limit         int64
	Offset        int64
	Kind          int8
	Class         int8
	ParentID      *int64
	RoleID        *int64
	UserID        *int64
	GroupID       *int64
	DeviceID      int64
	EquipmentID   int64
	Name          string
	Keyword       string
	GetTotal      *bool
	DefaultEffect int8
}

type OptionFN func(*Option)

func Page(page, pageSize int64) OptionFN {
	return func(i *Option) {
		i.Offset = (page - 1) * pageSize
		i.Limit = pageSize
	}
}

func Limit(limit int64) OptionFN {
	return func(i *Option) {
		i.Limit = limit
	}
}

func Offset(offset int64) OptionFN {
	return func(i *Option) {
		i.Offset = offset
	}
}

func GetTotal(get bool) OptionFN {
	return func(i *Option) {
		var p = get
		i.GetTotal = &p
	}
}

func Kind(kind int8) OptionFN {
	return func(i *Option) {
		i.Kind = kind
	}
}

func Class(class int8) OptionFN {
	return func(i *Option) {
		i.Class = class
	}
}

func Group(groupID int64) OptionFN {
	return func(i *Option) {
		g := groupID
		i.GroupID = &g
	}
}

func Parent(parentID int64) OptionFN {
	return func(i *Option) {
		p := parentID
		i.ParentID = &p
	}
}

func Role(roleID int64) OptionFN {
	return func(i *Option) {
		p := roleID
		i.RoleID = &p
	}
}
func User(userID int64) OptionFN {
	return func(i *Option) {
		p := userID
		i.UserID = &p
	}
}

func DefaultEffect(effect int8) OptionFN {
	return func(i *Option) {
		i.DefaultEffect = effect
	}
}

func Device(deviceID int64) OptionFN {
	return func(i *Option) {
		i.DeviceID = deviceID
	}
}

func Equipment(equipmentID int64) OptionFN {
	return func(i *Option) {
		i.EquipmentID = equipmentID
	}
}

func Name(name string) OptionFN {
	return func(i *Option) {
		i.Name = name
	}
}

func Keyword(keyword string) OptionFN {
	return func(i *Option) {
		i.Keyword = keyword
	}
}

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
