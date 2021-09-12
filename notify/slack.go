package notify

import (
	"fmt"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils"
	"github.com/ouqiang/supervisor-event-listener/utils/httpclient"
)

type Slack conf.Slack

func (s *Slack) Send(msg *event.Message) error {
	body := s.format(msg.String(), s.Channel)
	timeout := 120
	resp := httpclient.PostJson(s.URL, body, timeout)
	if !resp.IsOK() {
		return resp
	}
	return nil
}

// 格式化消息内容
func (s *Slack) format(content string, channel string) string {
	content = utils.EscapeJson(content)
	specialChars := []string{"&", "<", ">"}
	replaceChars := []string{"&amp;", "&lt;", "&gt;"}
	content = utils.ReplaceStrings(content, specialChars, replaceChars)

	return fmt.Sprintf(`{"text":"%s","username":"Supervisor事件通知", "channel":"%s"}`, content, channel)
}
