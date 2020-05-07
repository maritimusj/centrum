package helper

import (
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
)

type Option struct {
	Limit       int64
	Offset      int64
	Kind        resource2.MeasureKind
	Class       resource2.Class
	OrgID       int64
	ParentID    *int64
	RoleID      *int64
	UserID      *int64
	GroupID     *int64
	DeviceID    int64
	MeasureID   int64
	EquipmentID int64
	StateID     int64

	Status *int64

	Name          string
	Keyword       string
	DefaultEffect resource2.Effect

	OrderBy string
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

func Kind(kind resource2.MeasureKind) OptionFN {
	return func(i *Option) {
		i.Kind = kind
	}
}

func Class(class resource2.Class) OptionFN {
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

func DefaultEffect(effect resource2.Effect) OptionFN {
	return func(i *Option) {
		i.DefaultEffect = effect
	}
}

func Device(deviceID int64) OptionFN {
	return func(i *Option) {
		i.DeviceID = deviceID
	}
}

func Measure(measureID int64) OptionFN {
	return func(i *Option) {
		i.MeasureID = measureID
	}
}

func Equipment(equipmentID int64) OptionFN {
	return func(i *Option) {
		i.EquipmentID = equipmentID
	}
}

func State(stateID int64) OptionFN {
	return func(i *Option) {
		i.StateID = stateID
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

func Status(status int64) OptionFN {
	return func(i *Option) {
		i.Status = &status
	}
}

func OrderBy(orderBy string) OptionFN {
	return func(i *Option) {
		i.OrderBy = orderBy
	}
}
