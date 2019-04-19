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
		"text": this.format(msg),
	}
	if channel != "" {
		params["channel"] = channel
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	resp := httpclient.PostJson(url, string(body), timeout)
	if !resp.IsOK() {
		tmpfslog.Error("params: %v err: %v", params, resp.Error())
		return resp.Error()
	}
	return nil
}

func (this *BearyChat) format(msg event.Message) string {
	// return msg.ToJson(4)
	return msg.String()
}
