package notify

import (
	"github.com/ouqiang/supervisor-event-listener/config"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/tmpfslog"

	"fmt"
	"os"
	"time"
)

var (
	Conf  *config.Config
	queue chan event.Message
)

func init() {
	Conf = config.ParseConfig()

	tmpfslog.Info("loading config: %+v\n", Conf)
	queue = make(chan event.Message, 10)
	go start()
}

type Notifiable interface {
	Send(event.Message) error
}

func Push(header *event.Header, payload *event.Payload) {
	queue <- event.Message{header, payload}
}

func start() {
	var message event.Message
	var notifyHandler Notifiable
	for {
		message = <-queue
		tmpfslog.Info("%+v\n", message)
		switch Conf.NotifyType {
		case "mail":
			notifyHandler = &Mail{}
		case "slack":
			notifyHandler = &Slack{}
		case "webhook":
			notifyHandler = &WebHook{}
		case "bearychat":
			notifyHandler = &BearyChat{}
		}
		if notifyHandler == nil {
			continue
		}
		go send(notifyHandler, message)
		time.Sleep(1 * time.Second)
	}
}

func send(notifyHandler Notifiable, message event.Message) {
	// 最多重试3次
	tryTimes := 3
	i := 0
	for i < tryTimes {
		err := notifyHandler.Send(message)
		if err == nil {
			break
		}
		fmt.Fprintln(os.Stderr, err)
		time.Sleep(30 * time.Second)
		i++
	}
}
