package main

import (
	"flag"
	"log"
	"os"

	"github.com/ouqiang/supervisor-event-listener/conf"
	"github.com/ouqiang/supervisor-event-listener/listener"
	"github.com/ouqiang/supervisor-event-listener/notify"
	"github.com/ouqiang/supervisor-event-listener/utils/errlog"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Print("panic", err)
			os.Exit(127)
		}
	}()

	var configFile = flag.String("c", "/etc/supervisor-event-listener.ini", "config file")
	var dryRun = flag.Bool("dryRun", false, "dry run, lint config file")
	flag.Parse()

	if err := conf.Init(*configFile); err != nil {
		errlog.Error("config init failed. err: %v", err)
		os.Exit(127)
		return
	}
	if err := notify.Init(conf.Get()); err != nil {
		errlog.Error("notify init failed. err: %v", err)
		os.Exit(127)
		return
	}
	if *dryRun {
		return
	}
	notify.Start()
	listener.Start()
}
