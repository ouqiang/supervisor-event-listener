package notify

import (
	"errors"
	"fmt"
	"github.com/ouqiang/gocron/modules/httpclient"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils"
)

type Slack struct{}

func (slack *Slack) Send(message event.Message) error {
	body := slack.format(message.String(), Conf.Slack.Channel)
	timeout := 120
	response := httpclient.PostJson(Conf.Slack.WebHookUrl, body, timeout)
	if response.StatusCode == 200 {
		return nil
	}

	errorMessage := fmt.Sprintf("发送Slack消息失败#HTTP状态码-%d#HTTP-Body-%s",
		response.StatusCode, response.Body)

	return errors.New(errorMessage)
}

// 格式化消息内容
func (slack *Slack) format(content string, channel string) string {
	content = utils.EscapeJson(content)
	specialChars := []string{"&", "<", ">"}
	replaceChars := []string{"&amp;", "&lt;", "&gt;"}
	content = utils.ReplaceStrings(content, specialChars, replaceChars)

	return fmt.Sprintf(`{"text":"%s","username":"Supervisor事件通知", "channel":"%s"}`, content, channel)
}
