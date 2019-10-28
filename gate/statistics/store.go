package statistics

import (
	_ "github.com/influxdata/influxdb1-client"
	db "github.com/influxdata/influxdb1-client/v2"
	"github.com/maritimusj/centrum/lang"
)

type Client struct {
	db db.Client
}

func New() *Client {
	return &Client{}
}

func (client *Client) Open(option map[string]interface{}) error {
	if v, ok := option["connStr"].(string); ok {
		username, _ := option["username"].(string)
		password, _ := option["password"].(string)
		c, err := db.NewHTTPClient(db.HTTPConfig{
			Addr:     v,
			Username: username,
			Password: password,
		})

		if err != nil {
			return lang.InternalError(err)
		}

		client.db = c
		return nil
	}
	return lang.Error(lang.ErrInvalidDBConnStr)
}

func (client *Client) queryData(dbName string, cmd string) ([]db.Result, error) {
	q := db.NewQuery(cmd, dbName, "s")
	response, err := client.db.Query(q)
	if err != nil {
		return nil, lang.InternalError(err)
	}
	if response.Error() != nil {
		return nil, lang.InternalError(response.Error())
	}

	return response.Results, nil
}
