package mysqlStore

import (
	"database/sql"
	"github.com/maritimusj/centrum/db"
	"github.com/maritimusj/centrum/util"
	log "github.com/sirupsen/logrus"

	"errors"
	"fmt"
	"strings"
)

func LoadData(db db.DB, tbName string, data map[string]interface{}, cond string, params ...interface{}) error {
	if len(data) > 0 {
		var names = make([]string, 0, len(data))
		var values = make([]interface{}, 0, len(data))
		for k, v := range data {
			if strings.HasSuffix(k, "?") {
				name := strings.TrimSuffix(k, "?")
				names = append(names, fmt.Sprintf("IFNULL(`%s`,0) `%s`", name, name))
			} else {
				names = append(names, "`"+k+"`")
			}

			values = append(values, v)
		}

		var SQL strings.Builder
		SQL.WriteString("SELECT ")
		SQL.WriteString(strings.Join(names, ","))
		SQL.WriteString(" FROM ")
		SQL.WriteString(tbName)
		SQL.WriteString(" WHERE ")
		SQL.WriteString(cond)

		err := db.QueryRow(SQL.String(), params...).Scan(values...)
		log.Tracef("LoadData: %s => %s", SQL.String(), util.If(err != nil, err, "Ok"))
		if err != nil {
			return err
		}
		return nil
	}

	panic(errors.New("LoadData: empty data"))
}

func CreateData(db db.DB, tbName string, data map[string]interface{}) (int64, error) {
	if len(data) > 0 {
		var params = make([]string, 0, len(data))
		var values = make([]interface{}, 0, len(data))
		var placeHolders = make([]string, 0, len(data))

		for k, v := range data {
			params = append(params, "`"+k+"`")
			values = append(values, v)
			placeHolders = append(placeHolders, "?")
		}

		var SQL strings.Builder
		SQL.WriteString("INSERT INTO ")
		SQL.WriteString(tbName)
		SQL.WriteString("(")
		SQL.WriteString(strings.Join(params, ","))
		SQL.WriteString(") VALUES (")
		SQL.WriteString(strings.Join(placeHolders, ","))
		SQL.WriteString(")")

		result, err := db.Exec(SQL.String(), values...)
		log.Tracef("createData: %s => %s", SQL.String(), util.If(err != nil, err, "Ok"))

		if err != nil {
			return 0, err
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return lastInsertID, nil
	}

	panic(errors.New("CreateData: empty data"))
}

func SaveData(db db.DB, tbName string, data map[string]interface{}, cond string, params ...interface{}) error {
	if len(data) > 0 {
		var values = make([]interface{}, 0, len(data))
		var placeHolders = make([]string, 0, len(data))

		for k, v := range data {
			placeHolders = append(placeHolders, "`"+k+"`=?")
			values = append(values, v)
		}

		if len(params) > 0 {
			values = append(values, params...)
		}

		var SQL strings.Builder
		SQL.WriteString("UPDATE ")
		SQL.WriteString(tbName)
		SQL.WriteString(" SET ")
		SQL.WriteString(strings.Join(placeHolders, ","))
		SQL.WriteString(" WHERE ")
		SQL.WriteString(cond)

		_, err := db.Exec(SQL.String(), values...)
		log.Tracef("SaveData: %s => %s", SQL.String(), util.If(err != nil, err, "Ok"))
		if err != nil {
			return err
		}
		return nil
	}

	panic(errors.New("SaveData: empty data"))
}

func RemoveData(db db.DB, tbName string, cond string, params ...interface{}) error {
	SQL := "DELETE FROM " + tbName + " WHERE " + cond
	_, err := db.Exec(SQL, params...)
	log.Tracef("RemoveData: %s => %s", SQL, util.If(err != nil, err, "Ok"))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}

func IsDataExists(db db.DB, tbName string, cond string, params ...interface{}) (bool, error) {
	var total int64
	SQL := "SELECT COUNT(*) FROM " + tbName + " WHERE " + cond + " Limit 1"
	err := db.QueryRow(SQL, params...).Scan(&total)
	log.Tracef("IsDataExists: %s => %s", SQL, util.If(err != nil, err, "Ok"))
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return total > 0, nil
}
