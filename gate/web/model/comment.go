package model

import "github.com/maritimusj/centrum/gate/web/helper"

type Comment interface {
	DBEntry
	OptionEntry
	Profile

	AlarmID() int64
	Alarm() (Alarm, error)

	UserID() int64
	User() (User, error)

	ParentID() int64
	Parent() (Comment, error)

	SetAlarmID(int64)
	SetUserID(int64)
	SetParentID(int64)

	GetReplyList(options ...helper.OptionFN) ([]Comment, int64, error)
}
