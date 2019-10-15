package edge

import (
	"bytes"
	"net/http"

	"github.com/gorilla/rpc/json"
	. "github.com/maritimusj/centrum/json_rpc"
)

const (
	url = "http://localhost:1234/rpc"
)

func Invoke(cmd string, request interface{}) (*Result, error) {
	message, err := json.EncodeClientRequest(cmd, request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(message))
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var reply Result
	err = json.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func GetBaseInfo(uid string) (map[string]interface{}, error) {
	result, err := Invoke("Edge.GetBaseInfo", uid)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return result.Data.(map[string]interface{}), nil
}

func Remove(uid string) {
	_, _ = Invoke("Edge.Remove", uid)
}

func Active(conf *Conf) error {
	_, err := Invoke("Edge.Active", conf)
	if err != nil {
		return err
	}
	return nil
}

func SetValue(uid string, tag string, val interface{}) error {
	_, err := Invoke("Edge.SetValue", &Value{
		CH: CH{
			UID: uid,
			Tag: tag,
		},
		V: val,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetValue(uid string, tag string) (map[string]interface{}, error) {
	result, err := Invoke("Edge.GetValue", &CH{
		UID: uid,
		Tag: tag,
	})
	if err != nil {
		return nil, err
	}
	return result.Data.(map[string]interface{}), nil
}

func GetRealtimeData(uid string) ([]interface{}, error) {
	result, err := Invoke("Edge.GetRealtimeData", uid)
	if err != nil {
		return nil, err
	}
	return result.Data.([]interface{}), nil
}
