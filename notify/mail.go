package notify

import (
	"fmt"
	"strings"

	"github.com/go-gomail/gomail"
	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
)

type Mail conf.Mail

func (mail *Mail) Send(msg *event.Message) error {
	body := msg.String()
	body = strings.Replace(body, "\n", "<br>", -1)
	m := gomail.NewMessage()
	m.SetHeader("From", mail.ServerUser)
	m.SetHeader("To", mail.Receivers...)
	m.SetHeader("Subject", "Supervisor事件通知")
	m.SetBody("text/html", body)
	mailer := gomail.NewPlainDialer(
		mail.ServerHost,
		mail.ServerPort,
		mail.ServerUser,
		mail.ServerPassword,
	)
	err := mailer.DialAndSend(m)
	if err == nil {
		return nil
	}
	return fmt.Errorf("邮件发送失败#%v", err)
}
