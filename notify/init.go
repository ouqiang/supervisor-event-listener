package notify

import (
	"sync/atomic"
	"syscall"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/errlog"

	"fmt"
	"os"
	"os/signal"
	"time"
)

var (
	confFilePath string
	chanMsg      = make(chan *event.Message, 2048)
	chanSig      = make(chan os.Signal, 128)
	// notifiables  map[string]Notifiable
	notifiables atomic.Value
)

func Init(conf *conf.Config) error {
	if conf == nil {
		return fmt.Errorf("nil config")
	}
	signal.Notify(chanSig, syscall.SIGHUP)

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

func Reload() error {
	fpath := confFilePath
	errlog.Info("loading config: %s", fpath)
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
		case sig := <-chanSig:
			handleSignal(sig)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func handleSignal(sig os.Signal) error {
	if sig != syscall.SIGHUP {
		return fmt.Errorf("invalid signal %v", sig)
	}
	return Reload()
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
