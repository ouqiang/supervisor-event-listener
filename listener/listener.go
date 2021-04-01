package listener

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ouqiang/supervisor-event-listener/event"
	"github.com/ouqiang/supervisor-event-listener/notify"
)

var (
	// ErrPayloadLength ...
	ErrPayloadLength = errors.New("Header中len长度与实际读取长度不一致")
)

// Start ...
func Start() {
	go run()
}

func run() {
	for {
		defer func() {
			if err := recover(); err != nil {
				log.Print("panic", err)
			}
		}()
		listen()
		time.Sleep(time.Second)
	}
}

// 监听事件, 从标准输入获取事件内容
func listen() {
	reader := bufio.NewReader(os.Stdin)
	for {
		ready()
		msg, err := event.ReadMessage(reader)
		if err != nil {
			failure(err)
			continue
		}
		success()
		notify.Push(&msg)
	}
}

func ready() {
	fmt.Fprint(os.Stdout, "READY\n")
}

func success() {
	fmt.Fprint(os.Stdout, "RESULT 2\nOK")
}

func failure(err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Fprint(os.Stdout, "RESULT 2\nFAIL")
}
