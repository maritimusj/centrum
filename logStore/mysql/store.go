package mysql

import (
	"context"
	"github.com/maritimusj/centrum/logStore"
	"github.com/maritimusj/centrum/web/db"
	"github.com/sirupsen/logrus"
)

type mysqlStore struct {
	db db.DB
}

func New() logStore.Store {
	return &mysqlStore{}
}

func (m *mysqlStore) Open(ctx context.Context, option map[string]interface{}) error {
	panic("implement me")
}

func (m *mysqlStore) Wait() {
	panic("implement me")
}

func (m *mysqlStore) GetList(orgID int64, src, level string, start *uint64, offset, limit uint64) (result []*logStore.Data, total uint64, err error) {
	panic("implement me")
}

func (m *mysqlStore) Delete(orgID int64, src string) error {
	panic("implement me")
}

func (m *mysqlStore) Stats(orgID int64) map[string]interface{} {
	panic("implement me")
}

func (m *mysqlStore) Levels() []logrus.Level {
	panic("implement me")
}

func (m *mysqlStore) Fire(entry *logrus.Entry) error {
	panic("implement me")
}
