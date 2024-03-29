package sparkdesk

import (
	factory "chat/adapter/common"
	"chat/globals"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type ChatInstance struct {
	AppId     string
	ApiSecret string
	ApiKey    string
	Endpoint  string
}

func TransformAddr(model string) string {
	switch model {
	case globals.SparkDesk:
		return "v1.1"
	case globals.SparkDeskV2:
		return "v2.1"
	case globals.SparkDeskV3:
		return "v3.1"
	case globals.SparkDeskV35:
		return "v3.5"
	default:
		return "v1.1"
	}
}

func TransformModel(model string) string {
	switch model {
	case globals.SparkDesk:
		return "general"
	case globals.SparkDeskV2:
		return "generalv2"
	case globals.SparkDeskV3:
		return "generalv3"
	case globals.SparkDeskV35:
		return "generalv3.5"
	default:
		return "general"
	}
}

func NewChatInstanceFromConfig(conf globals.ChannelConfig) factory.Factory {
	params := conf.SplitRandomSecret(3)

	return &ChatInstance{
		AppId:     params[0],
		ApiSecret: params[1],
		ApiKey:    params[2],
		Endpoint:  conf.GetEndpoint(),
	}
}

func (c *ChatInstance) CreateUrl(endpoint, host, date, auth string) string {
	v := make(url.Values)
	v.Add("host", host)
	v.Add("date", date)
	v.Add("authorization", auth)
	return fmt.Sprintf("%s?%s", endpoint, v.Encode())
}

func (c *ChatInstance) Sign(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// GenerateUrl will generate the signed url for sparkdesk api
func (c *ChatInstance) GenerateUrl(endpoint string) string {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}

	date := time.Now().UTC().Format(time.RFC1123)
	data := strings.Join([]string{
		fmt.Sprintf("host: %s", uri.Host),
		fmt.Sprintf("date: %s", date),
		fmt.Sprintf("GET %s HTTP/1.1", uri.Path),
	}, "\n")

	signature := c.Sign(data, c.ApiSecret)
	authorization := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf(
			"hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"",
			c.ApiKey,
			"hmac-sha256",
			"host date request-line",
			signature,
		),
	))

	return c.CreateUrl(endpoint, uri.Host, date, authorization)
}
