package main

import (
	aApp "github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"github.com/maybecoding/go-metrics.git/internal/agent/httpsender"
	"github.com/maybecoding/go-metrics.git/internal/agent/memcollector"
)

func main() {
	// Config
	cfg := config.New()

	var memCollect aApp.Collector = memcollector.New()
	var httpSend aApp.Sender = httpsender.New(cfg.Sender.Address, cfg.Sender.Method, cfg.Sender.Template)

	app := aApp.New(memCollect, httpSend, cfg.App.SendIntervalSec, cfg.App.CollectIntervalSec)

	app.Start()

}
