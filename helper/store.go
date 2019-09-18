package helper

import (
	"github.com/maritimusj/centrum/resource"
)

type Option struct {
	Limit         int64
	Offset        int64
	Kind          resource.MeasureKind
	Class         resource.Class
	OrgID         int64
	ParentID      *int64
	RoleID        *int64
	UserID        *int64
	GroupID       *int64
	DeviceID      int64
	EquipmentID   int64
	Name          string
	Keyword       string
	GetTotal      *bool
	DefaultEffect resource.Effect
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

func Kind(kind resource.MeasureKind) OptionFN {
	return func(i *Option) {
		i.Kind = kind
	}
}

func Class(class resource.Class) OptionFN {
	return func(i *Option) {
		i.Class = class
	}
}

func Organization(orgID int64) OptionFN {
	return func(i *Option) {
		i.OrgID = orgID
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

func DefaultEffect(effect resource.Effect) OptionFN {
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
