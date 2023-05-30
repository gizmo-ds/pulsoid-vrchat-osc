package pulsoid

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rs/xid"
)

type WebSocketResult struct {
	Timestamp int64 `json:"timestamp"`
	Data      struct {
		HeartRate int `json:"heartRate"`
	} `json:"data"`
}

func GetRamielUrl(widgetId string) (string, error) {
	var result struct {
		Result struct {
			RamielUrl string `json:"ramielUrl"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	_, err := resty.New().R().
		SetResult(&result).
		SetBody(map[string]any{
			"method":  "getWidget",
			"jsonrpc": "2.0",
			"params":  map[string]any{"widgetId": widgetId},
			"id":      xid.New().String(),
		}).
		Post("https://pulsoid.net/v1/api/public/rpc")
	if err != nil {
		return "", err
	}
	if result.Error.Message != "" {
		return "", fmt.Errorf(result.Error.Message)
	}
	return result.Result.RamielUrl, err
}
