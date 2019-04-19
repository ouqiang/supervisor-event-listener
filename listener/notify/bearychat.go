package notify

import (
	"encoding/json"

	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/httpclient"
	"github.com/ouqiang/supervisor-event-listener/utils/tmpfslog"
)

type BearyChat struct{}

func (this *BearyChat) Send(msg event.Message) error {
	url := Conf.BearyChat.WebHookUrl
	channel := Conf.BearyChat.Channel
	timeout := Conf.BearyChat.Timeout

	params := map[string]interface{}{
		"text": msg.String(),
	}
	if channel != "" {
		params["channel"] = channel
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	tmpfslog.Info("url: %s", url)
	tmpfslog.Info("timeout: %d", timeout)
	tmpfslog.Info("params: %v", params)

	resp := httpclient.PostJson(url, string(body), timeout)
	if !resp.IsOK() {
		return resp.Error()
	}
	return nil
}

func (this *BearyChat) format(msg event.Message) string {
	return ""
}
