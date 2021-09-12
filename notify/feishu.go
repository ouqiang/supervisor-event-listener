package notify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/errlog"
	"github.com/ouqiang/supervisor-event-listener/utils/httpclient"
)

type Feishu conf.Feishu

func (this *Feishu) Send(msg *event.Message) error {
	url := this.URL
	timeout := this.Timeout
	return send2feishu(url, msg.String(), timeout)
}

func send2feishu(_url string, text string, timeout int) error {
	parse := func(_url string) (string, string) {
		tmpArr := strings.Split(_url, "?")
		if len(tmpArr) == 1 {
			return _url, ""
		} else if len(tmpArr) == 2 {
			_url = tmpArr[0]
			q, err := url.ParseQuery(tmpArr[1])
			if err != nil {
				panic(err)
			}
			return _url, q.Get("signKey")
		} else {
			panic(fmt.Errorf("invalid url: %s", _url))
		}
	}

	sign := func(secret string, timestamp int64) string {
		//timestamp + key  do sha256, then base64 encode
		stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)

		var data []byte
		h := hmac.New(sha256.New, []byte(stringToSign))
		_, err := h.Write(data)
		if err != nil {
			panic(err)
		}
		return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	params := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": text,
		},
	}

	_url, signKey := parse(_url)
	if signKey != "" {
		ts := time.Now().Unix()
		params["timestamp"] = ts
		params["sign"] = sign(signKey, ts)
	}
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	resp := httpclient.PostJson(_url, string(body), timeout)
	if !resp.IsOK() {
		errlog.Error("params: %v err: %v", params, resp.Error())
		return resp
	}
	return nil
}
