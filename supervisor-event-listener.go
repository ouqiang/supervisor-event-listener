package main

import (
	"flag"
	"log"
	"os"

	"github.com/ouqiang/supervisor-event-listener/listener"
	"github.com/ouqiang/supervisor-event-listener/listener/notify"
	"github.com/ouqiang/supervisor-event-listener/utils/tmpfslog"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Print("panic", err)
			os.Exit(127)
		}
	}()

	var configFile string
	var dryRun bool
	flag.StringVar(&configFile, "c", "/etc/supervisor-event-listener.ini", "config file")
	flag.BoolVar(&dryRun, "dryRun", false, "dry run, lint config file")
	flag.Parse()
	err := notify.Init(configFile)
	if err != nil {
		tmpfslog.Error("notify init failed. err: %+v", err)
		os.Exit(127)
	}
	if dryRun {
		return
	}

	notify.Start()
	listener.Start()
}
