package notify

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/errlog"
)

var (
	chanMsg    = make(chan *event.Message, 2048)
	chanReload = make(chan *conf.Config, 16)
	// notifiables  map[string]Notifiable
	notifiables atomic.Value
)

func Init(conf *conf.Config) error {
	if conf == nil {
		return fmt.Errorf("nil config")
	}

	m := map[string]Notifiable{}
	if conf.WebHook != nil {
		m["webhook"] = (*WebHook)(conf.WebHook)
	}
	if conf.Mail != nil {
		m["mail"] = (*Mail)(conf.Mail)
	}
	if conf.Slack != nil {
		m["slack"] = (*Slack)(conf.Slack)
	}
	if conf.BearyChat != nil {
		m["bearychat"] = (*BearyChat)(conf.BearyChat)
	}
	if conf.Feishu != nil {
		m["feishu"] = (*Feishu)(conf.Feishu)
	}
	notifiables.Store(m)
	return nil
}

func Reload(conf *conf.Config) {
	chanReload <- conf
}

func reload(conf *conf.Config) error {
	if err := Init(conf); err != nil {
		errlog.Info("loading config failed. %v", err)
		return err
	}
	names := []string{}
	for name := range Get() {
		names = append(names, name)
	}
	errlog.Info("reloaded config: %s", strings.Join(names, " "))
	return nil
}

func Push(msg *event.Message) {
	chanMsg <- msg
}

func Start() {
	go run()
}

func run() {
	for {
		select {
		case msg := <-chanMsg:
			errlog.Info("msg=%s", msg.ToJson(2))
			handleMessage(msg)
			time.Sleep(50 * time.Millisecond)
		case conf := <-chanReload:
			_ = reload(conf)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func handleMessage(msg *event.Message) error {
	errlog.Debug("message: %+v\n", msg)
	m := Get()
	for name, notifyHandler := range m {
		if notifyHandler == nil {
			errlog.Error("nil notify handler %s", name)
		}
		go send(notifyHandler, msg)
	}
	return nil
}

func send(notifyHandler Notifiable, msg *event.Message) {
	// 最多重试3次
	tryTimes := 3
	i := 0
	for i < tryTimes {
		err := notifyHandler.Send(msg)
		if err == nil {
			break
		}
		fmt.Fprintln(os.Stderr, err)
		time.Sleep(30 * time.Second)
		i++
	}
}

func Get() map[string]Notifiable {
	return notifiables.Load().(map[string]Notifiable)
}
