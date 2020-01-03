package edge

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/rpc/json"
	. "github.com/maritimusj/centrum/json_rpc"
)

type Balance struct {
	url   string
	total int
}
type EdgesMap struct {
	edges   []*Balance
	devices map[string]*Balance
	mu      sync.RWMutex
}

var (
	defaultEdgesMap = &EdgesMap{
		edges:   []*Balance{},
		devices: map[string]*Balance{},
	}
)

func Add(url string) {
	defaultEdgesMap.mu.Lock()
	defer defaultEdgesMap.mu.Unlock()

	defaultEdgesMap.edges = append(defaultEdgesMap.edges, &Balance{
		url:   url,
		total: 0,
	})
}

func Restart(url string) {
	_, _ = Invoke(url, "Edge.Restart", nil)
}

func Invoke(url, cmd string, request interface{}) (*Result, error) {
	message, err := json.EncodeClientRequest(cmd, request)
	if err != nil {
		return nil, err
	}

	log.Traceln("invoke: ", url, string(message))

	resp, err := http.Post(url, "application/json", bytes.NewReader(message))
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var reply Result

	if log.IsLevelEnabled(log.TraceLevel) {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("[invoke] %s, result: %s", err, string(data))
			return nil, err
		}

		log.Traceln("[invoke] ", string(data))

		err = json.DecodeClientResponse(bytes.NewReader(data), &reply)
		if err != nil {
			return nil, err
		}

	} else {
		err = json.DecodeClientResponse(resp.Body, &reply)
		if err != nil {
			return nil, err
		}
	}

	return &reply, nil
}

func GetBaseInfo(uid string) (map[string]interface{}, error) {
	defaultEdgesMap.mu.RLock()

	var url string
	if b, ok := defaultEdgesMap.devices[uid]; ok {
		url = b.url
	}

	defaultEdgesMap.mu.RUnlock()

	if url != "" {
		result, err := Invoke(url, "Edge.GetBaseInfo", uid)
		if err != nil {
			return map[string]interface{}{}, err
		}
		return result.Data.(map[string]interface{}), nil
	}

	return map[string]interface{}{}, errors.New("device not found")
}

func Reset(uid string) {
	defaultEdgesMap.mu.RLock()
	defer defaultEdgesMap.mu.RUnlock()

	if b, ok := defaultEdgesMap.devices[uid]; ok {
		go func() {
			_, _ = Invoke(b.url, "Edge.Reset", uid)
		}()
	}
}

func Remove(uid string) {
	defaultEdgesMap.mu.Lock()
	defer defaultEdgesMap.mu.Unlock()

	if b, ok := defaultEdgesMap.devices[uid]; ok {
		go func() {
			_, _ = Invoke(b.url, "Edge.Remove", uid)
		}()

		b.total -= 1
	}
}

func Active(conf *Conf) error {
	defaultEdgesMap.mu.Lock()

	var balance *Balance
	if v, ok := defaultEdgesMap.devices[conf.UID]; ok {
		balance = v
	} else {
		for _, b := range defaultEdgesMap.edges {
			if balance == nil || balance.total > b.total {
				balance = b
			}
		}
	}

	if balance != nil {
		if _, ok := defaultEdgesMap.devices[conf.UID]; !ok {
			balance.total += 1
			defaultEdgesMap.devices[conf.UID] = balance
		}
	}

	defaultEdgesMap.mu.Unlock()

	if balance != nil {
		_, err := Invoke(balance.url, "Edge.Active", conf)
		return err
	}

	return errors.New("no edge")
}

func SetValue(uid string, tag string, val interface{}) error {
	defaultEdgesMap.mu.RLock()

	var url string
	if b, ok := defaultEdgesMap.devices[uid]; ok {
		url = b.url
	}

	defaultEdgesMap.mu.RUnlock()

	if url != "" {
		_, err := Invoke(url, "Edge.SetValue", &Value{
			CH: CH{
				UID: uid,
				Tag: tag,
			},
			V: val,
		})
		return err
	}

	return errors.New("device not found")
}

func GetValue(uid string, tag string) (map[string]interface{}, error) {
	defaultEdgesMap.mu.RLock()

	var url string
	if b, ok := defaultEdgesMap.devices[uid]; ok {
		url = b.url
	}

	defaultEdgesMap.mu.RUnlock()

	if url != "" {
		result, err := Invoke(url, "Edge.GetValue", &CH{
			UID: uid,
			Tag: tag,
		})
		if err != nil {
			return nil, err
		}
		return result.Data.(map[string]interface{}), nil
	}

	return map[string]interface{}{}, errors.New("device not found")
}

func GetRealtimeData(uid string) ([]interface{}, error) {
	defaultEdgesMap.mu.RLock()

	var url string
	if b, ok := defaultEdgesMap.devices[uid]; ok {
		url = b.url
	}

	defaultEdgesMap.mu.RUnlock()

	if url != "" {
		result, err := Invoke(url, "Edge.GetRealtimeData", uid)
		if err != nil {
			return nil, err
		}
		return result.Data.([]interface{}), nil
	}

	return []interface{}{}, errors.New("device not found")
}
