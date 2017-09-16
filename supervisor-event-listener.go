package main

import (
	"github.com/ouqiang/supervisor-event-listener/listener"
)

func main() {
	for {
		listener.Start()
	}
}
