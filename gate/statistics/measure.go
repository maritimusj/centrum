package statistics

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	"github.com/maritimusj/centrum/lang"
	"strings"
	"time"
)

func (client *Client) GetMeasureStats(dbName string, deviceID int64, tagName string, start, end *time.Time, interval time.Duration) (*models.Row, error) {
	var SQL strings.Builder
	if interval > 0 {
		SQL.WriteString(`SELECT max("val") FROM `)
	} else {
		SQL.WriteString(`SELECT "val","alarm" FROM `)
	}
	SQL.WriteString(fmt.Sprintf(`"%s"`, tagName))
	SQL.WriteString(fmt.Sprintf(` WHERE "uid"='%d' AND "time">='%s'`, deviceID, start.UTC().Format(time.RFC3339)))
	if end != nil {
		SQL.WriteString(fmt.Sprintf(` AND "time"<'%s'`, end.UTC().Format(time.RFC3339)))
	}

	if interval > 0 {
		SQL.WriteString(fmt.Sprintf(` GROUP BY time(%s)  fill(0)`, interval.String()))
	}

	println(SQL.String())
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

	return nil, lang.Error(lang.ErrNotStatisticsData)
}
