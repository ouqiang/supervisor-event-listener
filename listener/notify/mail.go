package notify

import (
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/ouqiang/supervisor-event-listener/event"
	"strings"
)

type Mail struct{}

func (mail *Mail) Send(message event.Message) error {
	body := message.String()
	body = strings.Replace(body, "\n", "<br>", -1)
	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", Conf.MailServer.User)
	gomailMessage.SetHeader("To", Conf.MailUser.Email...)
	gomailMessage.SetHeader("Subject", "Supervisor事件通知")
	gomailMessage.SetBody("text/html", body)

	var dialer *gomail.Dialer

	// 465 port should use SSL
	if Conf.MailServer.Port == 465 {
		dialer = &gomail.Dialer{
			Host:Conf.MailServer.Host,
			Port:Conf.MailServer.Port,
			Username:Conf.MailServer.User,
			Password:Conf.MailServer.Password,
			SSL: true,
		}
	} else {
		dialer = gomail.NewPlainDialer(
			Conf.MailServer.Host,
			Conf.MailServer.Port,
			Conf.MailServer.User,
			Conf.MailServer.Password,
		)
	}

	err := dialer.DialAndSend(gomailMessage)
	if err == nil {
		return nil
	}
	errorMessage := fmt.Sprintf("邮件发送失败#%s", err.Error())

	return errors.New(errorMessage)
}
