package Getui

var (
	defaultClient = NewClient()
)

func SetAppKey(appKey string) *Client {
	defaultClient.SetAppKey(appKey)
	return defaultClient
}

func SetAppId(appId string) *Client {
	defaultClient.SetAppId(appId)
	return defaultClient
}

func SetAppSecret(appSecret string) *Client {
	defaultClient.SetAppSecret(appSecret)
	return defaultClient
}

func SetMasterSecret(masterSecret string) *Client {
	defaultClient.SetMasterSecret(masterSecret)
	return defaultClient
}

func toSingle(cid string, title, body string) error {
	return defaultClient.toSingle(cid, title, body)
}
