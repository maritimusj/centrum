package mysql

import (
	"context"
	"database/sql"
	"github.com/maritimusj/centrum/db"
	"github.com/maritimusj/centrum/lang"
	"time"
)

type mysqlDB struct {
	db  *sql.DB
	ctx context.Context
}

func New() db.DB {
	return &mysqlDB{}
}

func (m *mysqlDB) TransactionDo(fn func(db db.DB) interface{}) interface{} {
	tx, err := m.db.Begin()
	if err != nil {
		return lang.InternalError(err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res := fn(tx)
	if res != nil {
		if err, ok := res.(error); ok {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return lang.InternalError(err)
	}

	return res
}

func Open(ctx context.Context, option map[string]interface{}) (db.WithTransaction, error) {
	if connStr, ok := option["connStr"].(string); ok {
		conn, err := sql.Open("mysql", connStr)
		if err != nil {
			return nil, lang.InternalError(err)
		}

		ctxTimeout, _ := context.WithTimeout(ctx, time.Second*3)
		err = conn.PingContext(ctxTimeout)
		if err != nil {
			return nil, lang.InternalError(err)
		}

		return &mysqlDB{
			db:  conn,
			ctx: ctx,
		}, nil
	}
	return nil, lang.Error(lang.ErrInvalidConnStr)
}

func (m *mysqlDB) Close() {
	if m != nil && m.db != nil {
		_ = m.db.Close()
	}
}

func (m *mysqlDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

func (m *mysqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.Query(query, args...)
}

func (m *mysqlDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.db.QueryRow(query, args...)
}
