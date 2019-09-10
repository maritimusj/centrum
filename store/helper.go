package store

import (
	"database/sql"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
)

type Option struct {
	Limit       int64
	Offset      int64
	Kind        model.MeasureKind
	Class       resource.Class
	ParentID    *int64
	GroupID     *int64
	DeviceID    int64
	EquipmentID int64
	Keyword     string
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

func Kind(kind model.MeasureKind) OptionFN {
	return func(i *Option) {
		i.Kind = kind
	}
}

func Class(class resource.Class) OptionFN {
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
