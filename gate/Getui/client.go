package Getui

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/util/grand"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gtime"
)

const (
	BaseUrl   = "https://restapi.getui.com/v2/"
	IntentStr = "intent:#Intent;action=android.intent.action.oppopush;launchFlags=0x14000000;component=tech.cyhk.gwater;S.UP-OL-SU=true;S.title=%s;S.content=%s;S.payload=%s;end"
)

type Client struct {
	AppId        string
	AppKey       string
	AppSecret    string
	MasterSecret string
	BaseUrl      string
}

type Notification struct {
	Title     string  `json:"title"`
	Body      string  `json:"body"`
	ClickType string  `json:"click_type"`
	Intent    *string `json:"intent,omitempty"`
	Url       *string `json:"url,omitempty"`
}

type PushMessage struct {
	Duration     string        `json:"duration,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
	Transmission *string       `json:"transmission,omitempty"`
}

type Audience struct {
	Cid []string `json:"cid"`
}

type PushChannel struct {
	Android *Android `json:"android"`
	IOS     *IOS     `json:"ios"`
}

type IOS struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	APS     struct {
		Alert struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		} `json:"alert"`
		ContentAvailable int `json:"content_available"`
	} `json:"aps"`
}

func NewIOS(title, body, payload string) *IOS {
	ios := &IOS{
		Type:    "notify",
		Payload: payload,
	}
	ios.APS.Alert.Title = title
	ios.APS.Alert.Body = body
	return ios
}

func NewAndroid(title, body, payload string) *Android {
	android := &Android{}
	intent := fmt.Sprintf(IntentStr, title, body, payload)
	android.UPS.Notification = &Notification{
		Title:     title,
		Body:      body,
		ClickType: "intent",
		Intent:    &intent,
	}
	return android
}

type Android struct {
	UPS struct {
		Notification *Notification `json:"notification"`
	} `json:"ups"`
}

type PushRequest struct {
	ID          string                 `json:"request_id"`
	Settings    map[string]interface{} `json:"settings"`
	Audience    *Audience              `json:"audience"`
	PushMessage *PushMessage           `json:"push_message"`
	PushChannel *PushChannel           `json:"push_channel,omitempty"`
}

func (request *PushRequest) SetNotification(title, body string) *PushRequest {
	request.PushMessage.Notification = &Notification{
		Title:     title,
		Body:      body,
		ClickType: "none",
	}
	request.PushMessage.Transmission = nil
	return request
}

func (request *PushRequest) SetTransmission(title, body, payload string) *PushRequest {
	request.PushChannel = &PushChannel{
		Android: NewAndroid(title, body, payload),
		IOS:     NewIOS(title, body, payload),
	}
	request.PushMessage.Transmission = &payload
	request.PushMessage.Notification = nil
	return request
}

func (request *PushRequest) SetSettings(name string, settings interface{}) *PushRequest {
	request.Settings[name] = settings
	return request
}

func (request *PushRequest) AddCid(cid string) *PushRequest {
	request.Audience.Cid = append(request.Audience.Cid, cid)
	return request
}

func NewPushRequest() *PushRequest {
	request := &PushRequest{
		ID:          grand.Digits(16),
		Settings:    map[string]interface{}{},
		Audience:    &Audience{},
		PushMessage: &PushMessage{},
	}
	request.SetSettings("strategy", map[string]interface{}{
		"default": 1,
	})
	return request
}

type Response struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func NewClient() *Client {
	return &Client{}
}

func NewHttpClient() *ghttp.Client {
	return ghttp.NewClient().SetHeader("Content-Type", "application/json;charset=utf-8")
}

func (client *Client) SetAppId(appId string) *Client {
	client.AppId = appId
	client.BaseUrl = BaseUrl + appId
	return client
}

func (client *Client) SetAppKey(appKey string) *Client {
	client.AppKey = appKey
	return client
}

func (client *Client) SetAppSecret(appSecret string) *Client {
	client.AppSecret = appSecret
	return client
}

func (client *Client) SetMasterSecret(masterSecret string) *Client {
	client.MasterSecret = masterSecret
	return client
}

func (client *Client) getToken() (string, error) {
	ts := gtime.Now().TimestampMilliStr()
	signByte := sha256.Sum256([]byte(client.AppKey + ts + client.MasterSecret))
	sign := fmt.Sprintf("%x", signByte)

	//构建鉴权参数
	params := map[string]interface{}{
		"sign":      sign,
		"timestamp": ts,
		"appkey":    client.AppKey,
	}

	//参数JSON序列化
	jsonParams, _ := gjson.New(params).ToJson()
	r, err := NewHttpClient().Post(client.BaseUrl+"/auth", jsonParams)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = r.Close()
	}()

	//成功结果{"result":"ok","expire_time":"1568875773358","auth_token":"98b74e881a76dd5bbc98ec6cab8e650dfa33f24b1313507ba19c8947a320d7f5"}
	res := gjson.New(r.ReadAll())
	if res.GetInt("code") == 0 { //存入缓存存到过期时间前1秒
		token := res.GetString("data.token")
		if token != "" {
			expired := res.GetTime("data.expire_time").Sub(time.Now()) - time.Second
			gcache.Set("token", token, expired)
			return token, nil
		}
	}

	return "", errors.New(res.GetString("msg"))
}

func (client *Client) toSingle(cid string, title, body string) error {
	if client == nil || client.AppId == "" || client.AppKey == "" {
		return errors.New("GeTui service not properly initialized")
	}

	token, err := gcache.Get("token")
	if token == nil { //过期重新获取
		token, err = client.getToken()
		if err != nil {
			return err
		}
	}

	request := NewPushRequest()
	request.AddCid(cid)

	payload, _ := gjson.New(map[string]string{
		"title":   title,
		"content": body,
	}).ToJson()

	request.SetTransmission(title, body, string(payload))

	jsonParams, _ := gjson.New(request).ToJson()

	httpClient := NewHttpClient()
	httpClient.SetHeader("token", token.(string))

	r, err := httpClient.Post(client.BaseUrl+"/push/single/cid", jsonParams)
	if err != nil {
		return err
	}

	defer func() {
		_ = r.Close()
	}()

	res := gjson.New(r.ReadAll())
	if res.GetInt("code") == 0 {
		return nil
	}
	return errors.New(res.GetString("msg"))
}
