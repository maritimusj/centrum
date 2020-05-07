package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/db"
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

	result := fn(tx)
	if result != nil {
		if errCode, ok := result.(lang.ErrIndex); ok && errCode != lang.Ok {
			return lang.Error(errCode)
		}
		if err, ok := result.(error); ok {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return lang.InternalError(err)
	}

	return result
}

func Open(ctx context.Context, option map[string]interface{}) (db.WithTransaction, error) {
	if connStr, ok := option["connStr"].(string); ok {
		conn, err := sql.Open("sqlite3", connStr)
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
	return nil, lang.ErrInvalidDBConnStr.Error()
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
