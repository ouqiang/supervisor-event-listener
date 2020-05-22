package notify

import (
	"syscall"

	"github.com/ouqiang/supervisor-event-listener/config"
	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/utils/errlog"

	"fmt"
	"os"
	"os/signal"
	"time"
)

var (
	confFilePath string
	Conf         *config.Config
	chanMsg      chan event.Message
	chanSig      chan os.Signal = make(chan os.Signal, 100)
	listerners   map[string]Notifiable
)

func Init(fpath string) error {
	errlog.Info("loading config: %s", fpath)
	if Conf != nil {
		return fmt.Errorf("init twice!!!")
	}
	Conf = config.ParseConfig(fpath)
	chanMsg = make(chan event.Message, 10)
	confFilePath = fpath
	signal.Notify(chanSig, syscall.SIGHUP)
	return nil
}

func Reload() error {
	fpath := confFilePath
	errlog.Info("loading config: %s", fpath)
	Conf = config.ParseConfig(fpath)
	return nil
}
func Push(msg event.Message) {
	chanMsg <- msg
}

func Start() {
	go start()
}

func start() {
	for {
		select {
		case msg := <-chanMsg:
			errlog.Debug("msg=%s", msg.ToJson(2))
			handleMessage(msg)
			time.Sleep(100 * time.Millisecond)
		case sig := <-chanSig:
			handleSignal(sig)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func handleSignal(sig os.Signal) error {
	if sig != syscall.SIGHUP {
		return fmt.Errorf("invalid signal %v", sig)
	}
	return Reload()
}

func handleMessage(msg event.Message) error {
	errlog.Debug("message: %+v\n", msg)
	var notifyHandler Notifiable
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
		return fmt.Errorf("invalid notify type %s", Conf.NotifyType)
	}
	go send(notifyHandler, msg)
	return nil
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
