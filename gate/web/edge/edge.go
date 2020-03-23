package edge

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/rpc/v2/json"
	. "github.com/maritimusj/centrum/json_rpc"
)

type balanceCacheEntry struct {
	data interface{}
	exp  time.Time
}

func (e *balanceCacheEntry) IsExpired() bool {
	return e == nil || time.Now().Sub(e.exp) > time.Second
}

type Balance struct {
	url   string
	total int
	cache map[string]*balanceCacheEntry
	mu    sync.RWMutex
}

func (b *Balance) DeltaTotal(delta int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.total += delta
}

func (b *Balance) Remove(uid string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for path, _ := range b.cache {
		if strings.HasPrefix(path, uid) {
			delete(b.cache, path)
		}
	}
}

func (b *Balance) Set(uid string, key string, value interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	path := uid + "." + key
	if e, ok := b.cache[path]; ok {
		e.data = value
		e.exp = time.Now()
	} else {
		b.cache[path] = &balanceCacheEntry{
			data: value,
			exp:  time.Now(),
		}
	}
}

func (b *Balance) Get(uid string, key string) *balanceCacheEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	v, _ := b.cache[uid+"."+key]
	return v
}

type EdgesMap struct {
	edges   []*Balance
	devices map[string]*Balance
	mu      sync.RWMutex
}

func (e *EdgesMap) GetBalanceByDeviceUID(uid string) *Balance {
	e.mu.Lock()
	defer e.mu.Unlock()

	b, _ := e.devices[uid]
	return b
}

func (e *EdgesMap) AddDevice(uid string) *Balance {
	e.mu.Lock()
	defer e.mu.Unlock()

	var balance *Balance
	if v, ok := e.devices[uid]; ok {
		balance = v
	} else {
		for _, b := range e.edges {
			if balance == nil || balance.total > b.total {
				balance = b
			}
		}
	}

	if balance != nil {
		if _, ok := e.devices[uid]; !ok {
			balance.total += 1
			e.devices[uid] = balance
		}
	}

	return balance
}

var (
	defaultEdgesMap = &EdgesMap{
		edges:   []*Balance{},
		devices: map[string]*Balance{},
	}
)

//Add 增加一个edge URL
func Add(url string) {
	defaultEdgesMap.mu.Lock()
	defer defaultEdgesMap.mu.Unlock()

	defaultEdgesMap.edges = append(defaultEdgesMap.edges, &Balance{
		url:   url,
		total: 0,
		cache: map[string]*balanceCacheEntry{},
	})
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

//Restart 重启指定的edge
func Restart(url string) {
	_, _ = Invoke(url, "Edge.Restart", nil)
}

// GetBaseInfo 用于获取设备基本信息
func GetBaseInfo(uid string) (map[string]interface{}, error) {
	balance := defaultEdgesMap.GetBalanceByDeviceUID(uid)
	if balance == nil {
		return map[string]interface{}{}, errors.New("device not found")
	}

	if e := balance.Get(uid, "baseInfo"); !e.IsExpired() {
		if baseInfoData, ok := e.data.(map[string]interface{}); ok {
			return baseInfoData, nil
		}
	}

	result, err := Invoke(balance.url, "Edge.GetBaseInfo", uid)
	if err != nil {
		return map[string]interface{}{}, err
	}

	data, _ := result.Data.(map[string]interface{})
	balance.Set(uid, "baseInfo", data)

	return data, nil
}

// Reset 通知设备刷新配置
func Reset(uid string) {
	balance := defaultEdgesMap.GetBalanceByDeviceUID(uid)
	if balance != nil {
		if e := balance.Get(uid, "reset"); e.IsExpired() {
			_, _ = Invoke(balance.url, "Edge.Reset", uid)
			balance.Set(uid, "reset", true)
		}
	}
}

//Remove 移除一个设备
func Remove(uid string) {
	balance := defaultEdgesMap.GetBalanceByDeviceUID(uid)
	if balance != nil {
		_, _ = Invoke(balance.url, "Edge.Remove", uid)
		balance.DeltaTotal(-1)
		balance.Remove(uid)
	}
}

//Active 激活设备
func Active(conf *Conf) error {
	balance := defaultEdgesMap.AddDevice(conf.UID)
	if balance != nil {
		_, err := Invoke(balance.url, "Edge.Active", conf)
		return err
	}

	return errors.New("no edge")
}

//SetValue 设置设备指定点位的值
func SetValue(uid string, tag string, val interface{}) error {
	balance := defaultEdgesMap.GetBalanceByDeviceUID(uid)
	if balance != nil {
		_, err := Invoke(balance.url, "Edge.SetValue", &Value{
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

//GetValue 获取设备指定点位的值
func GetValue(uid string, tag string) (map[string]interface{}, error) {
	balance := defaultEdgesMap.GetBalanceByDeviceUID(uid)
	if balance != nil {
		if e := balance.Get(uid, tag); !e.IsExpired() {
			return e.data.(map[string]interface{}), nil
		}

		result, err := Invoke(balance.url, "Edge.GetValue", &CH{
			UID: uid,
			Tag: tag,
		})
		if err != nil {
			return nil, err
		}

		data, _ := result.Data.(map[string]interface{})

		balance.Set(uid, tag, data)

		return data, nil
	}

	return map[string]interface{}{}, errors.New("device not found")
}

//GetRealtimeData 获取指定设备的实时数据
func GetRealtimeData(uid string) ([]interface{}, error) {
	balance := defaultEdgesMap.GetBalanceByDeviceUID(uid)
	if balance != nil {
		if e := balance.Get(uid, "realtimeData"); !e.IsExpired() {
			return e.data.([]interface{}), nil
		}

		result, err := Invoke(balance.url, "Edge.GetRealtimeData", uid)
		if err != nil {
			return nil, err
		}

		data, _ := result.Data.([]interface{})
		balance.Set(uid, "realtimeData", data)

		return data, nil
	}

	return []interface{}{}, errors.New("device not found")
}
