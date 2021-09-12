package notify

import (
	"encoding/json"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/errlog"
	"github.com/ouqiang/supervisor-event-listener/utils/httpclient"
)

type BearyChat conf.BearyChat

func (this *BearyChat) Send(msg *event.Message) error {
	url := this.URL
	channel := this.Channel
	timeout := this.Timeout

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
		errlog.Error("params: %v err: %v", params, resp.Error())
		return resp
	}
	return nil
}

func (this *BearyChat) format(msg *event.Message) string {
	// return msg.ToJson(4)
	return msg.String()
}
