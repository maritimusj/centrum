package mysqlStore

import (
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Comment struct {
	id       int64
	parentID int64

	refID  int64
	userID int64

	extra     []byte
	createdAt time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewComment(s *mysqlStore, id int64) *Comment {
	return &Comment{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
}

func (c *Comment) GetID() int64 {
	return c.id
}

func (c *Comment) AlarmID() int64 {
	return c.refID
}

func (c *Comment) SetAlarmID(id int64) {
	if c.refID != id {
		c.refID = id
		c.dirty.Set("ref_id", func() interface{} {
			return c.refID
		})
	}
}

func (c *Comment) Alarm() (model.Alarm, error) {
	return c.store.GetAlarm(c.refID)
}

func (c *Comment) UserID() int64 {
	return c.userID
}

func (c *Comment) SetUserID(id int64) {
	if c.userID != id {
		c.userID = id
		c.dirty.Set("user_id", func() interface{} {
			return c.userID
		})
	}
}

func (c *Comment) User() (model.User, error) {
	return c.store.GetUser(c.userID)
}

func (c *Comment) ParentID() int64 {
	return c.parentID
}

func (c *Comment) SetParentID(id int64) {
	if c.parentID != id {
		c.parentID = id
		c.dirty.Set("parent_id", func() interface{} {
			return c.parentID
		})
	}
}

func (c *Comment) Parent() (model.Comment, error) {
	return c.store.GetComment(c.parentID)
}

func (c *Comment) GetReplyList(options ...helper.OptionFN) ([]model.Comment, int64, error) {
	options = append(options, helper.Parent(c.id))
	alarm, _ := c.Alarm()
	return c.store.GetCommentList(alarm, 0, options...)
}

func (c *Comment) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Comment) Save() error {
	if c.dirty.Any() {
		err := SaveData(c.store.db, TbComments, c.dirty.Data(true), "id=?", c.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (c *Comment) Destroy() error {
	if c == nil {
		return lang.Error(lang.ErrCommentNotFound)
	}
	comments, _, err := c.store.GetCommentList(nil, 0, helper.Parent(c.id))
	if err != nil {
		return err
	}
	for _, comment := range comments {
		comment.SetParentID(0)
		if err := comment.Save(); err != nil {
			return err
		}
	}

	return c.store.RemoveComment(c.id)
}

func (c *Comment) Option() map[string]interface{} {
	return gjson.ParseBytes(c.extra).Value().(map[string]interface{})
}

func (c *Comment) GetOption(path string) gjson.Result {
	if c != nil {
		return gjson.GetBytes(c.extra, path)
	}
	return gjson.Result{}
}

func (c *Comment) SetOption(path string, value interface{}) error {
	if c != nil {
		data, err := sjson.SetBytes(c.extra, path, value)
		if err != nil {
			return err
		}

		c.extra = data
		c.dirty.Set("extra", func() interface{} {
			return c.extra
		})

		return nil
	}

	return lang.Error(lang.ErrDeviceNotFound)
}

func (c *Comment) Simple() model.Map {
	if c == nil {
		return model.Map{}
	}

	return model.Map{
		"id":         c.id,
		"parent":     c.parentID,
		"user":       c.userID,
		"comment":    c.Option(),
		"created_at": c.createdAt.Format("2006-01-02 15:04:05"),
	}
}

func (c *Comment) Brief() model.Map {
	if c == nil {
		return model.Map{}
	}
	user, _ := c.User()
	brief := model.Map{
		"id":         c.id,
		"alarm":      c.refID,
		"parent":     c.parentID,
		"comment":    c.Option(),
		"created_at": c.createdAt.Format("2006-01-02 15:04:05"),
	}
	if user != nil {
		brief["user"] = user.Simple()
	}
	return brief
}

func (c *Comment) Detail() model.Map {
	if c == nil {
		return model.Map{}
	}
	user, _ := c.User()
	detail := model.Map{
		"id":         c.id,
		"alarm":      c.refID,
		"parent":     c.parentID,
		"comment":    c.Option(),
		"created_at": c.createdAt.Format("2006-01-02 15:04:05"),
	}
	if user != nil {
		detail["user"] = user.Simple()
	}
	return detail
}
