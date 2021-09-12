package notify

import (
	"encoding/json"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/httpclient"
)

type WebHook conf.WebHook

func (hook *WebHook) Send(msg *event.Message) error {
	url := hook.URL
	timeout := hook.Timeout
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp := httpclient.PostJson(url, string(data), timeout)
	if !resp.IsOK() {
		return resp
	}
	return nil
}
