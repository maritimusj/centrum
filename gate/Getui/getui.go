package Getui

import (
	"time"

	"github.com/maritimusj/centrum/gate/lang"

	"github.com/maritimusj/centrum/gate/properties"
	"github.com/maritimusj/centrum/gate/web/model"
)

/**
一个用户可以登录多个clientID,但每个clientID只能对应一个用户
*/
func Register(clientID string, user model.User) error {
	x := properties.LoadString("clientId", clientID)
	if x != "" {
		err := properties.Delete("user", x, clientID)
		if err != nil {
			return err
		}
	}

	err := properties.Write(time.Now().Format("2006-01-02 15:04:05"), "user", user.Name(), clientID)
	if err != nil {
		return err
	}

	return properties.Write(user.Name(), "clientId", clientID)
}

/**
向指定用户推送消息
*/
func SendTo(user model.User, title, content string) {
	result := properties.LoadAllString("user", user.Name())
	for clientID := range result {
		err := defaultClient.toSingle(clientID, title, content)
		if err != nil {
			user.Logger().Warnln(lang.GeTuiSendMessageFailed.Str(err))
		}
	}
}
