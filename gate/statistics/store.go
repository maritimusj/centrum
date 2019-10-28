package statistics

import (
	_ "github.com/influxdata/influxdb1-client"
	db "github.com/influxdata/influxdb1-client/v2"
	lang2 "github.com/maritimusj/centrum/gate/lang"
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
			return lang2.InternalError(err)
		}

		client.db = c
		return nil
	}
	return lang2.Error(lang2.ErrInvalidDBConnStr)
}

func (client *Client) queryData(dbName string, cmd string) ([]db.Result, error) {
	q := db.NewQuery(cmd, dbName, "s")
	response, err := client.db.Query(q)
	if err != nil {
		return nil, lang2.InternalError(err)
	}
	if response.Error() != nil {
		return nil, lang2.InternalError(response.Error())
	}

	return response.Results, nil
}
