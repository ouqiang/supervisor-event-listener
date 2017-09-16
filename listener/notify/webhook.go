package notify

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/httpclient"
)

type WebHook struct{}

func (hook *WebHook) Send(message event.Message) error {
	encodeMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	timeout := 60
	response := httpclient.PostJson(Conf.WebHook.Url, string(encodeMessage), timeout)

	if response.StatusCode == 200 {
		return nil
	}
	errorMessage := fmt.Sprintf("webhook执行失败#HTTP状态码-%d#HTTP-BODY-%s", response.StatusCode, response.Body)
	return errors.New(errorMessage)
}
