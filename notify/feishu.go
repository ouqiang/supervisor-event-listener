package notify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/errlog"
	"github.com/ouqiang/supervisor-event-listener/utils/httpclient"
)

type Feishu conf.Feishu

func (this *Feishu) Send(msg *event.Message) error {
	url := this.URL
	timeout := 6
	calcSign := func(secret string, timestamp int64) string {
		//timestamp + key 做sha256, 再进行base64 encode
		stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)

		var data []byte
		h := hmac.New(sha256.New, []byte(stringToSign))
		_, err := h.Write(data)
		if err != nil {
			panic(err)
		}
		return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}
	_ = calcSign
	params := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": msg.String(),
		},
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

// sleep and rety
