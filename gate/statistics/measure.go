package statistics

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/influxdata/influxdb1-client/models"
	"github.com/maritimusj/centrum/gate/lang"
)

func (client *Client) GetMeasureStats(dbName string, deviceID int64, tagName string, start, end *time.Time, interval interface{}) (*models.Row, error) {
	var i string
	switch v := interval.(type) {
	case time.Duration:
		i = v.String()
	case string:
		i = v
	case *string:
		if v != nil {
			i = *v
		}
	case int64:
		i = time.Duration(v).String()
	default:
	}

	var SQL strings.Builder
	if i != "" {
		SQL.WriteString(`SELECT max("val") FROM `)
	} else {
		SQL.WriteString(`SELECT "val","alarm" FROM `)
	}
	SQL.WriteString(fmt.Sprintf(`"%s"`, tagName))
	SQL.WriteString(fmt.Sprintf(` WHERE "uid"='%d' AND "time">='%s'`, deviceID, start.UTC().Format(time.RFC3339)))
	if end != nil {
		SQL.WriteString(fmt.Sprintf(` AND "time"<'%s'`, end.UTC().Format(time.RFC3339)))
	}

	if i != "" {
		SQL.WriteString(fmt.Sprintf(` GROUP BY time(%s)  fill(previous)`, i))
	}

	res, err := client.queryData(dbName, SQL.String())
	if err != nil {
		return nil, err
	}

	if res[0].Err != "" {
		return nil, lang.InternalError(errors.New(res[0].Err))
	}

	if len(res[0].Series) > 0 {
		return &res[0].Series[0], nil
	}

	return nil, lang.Error(lang.ErrNoStatisticsData)
}

func (client *Client) GetAlarmStats(dbName string, deviceID int64, start, end *time.Time) (*models.Row, error) {
	return nil, errors.New("not implement")
}
