package main

import (
	"github.com/maybecoding/go-metrics.git/cmd/agent/config"
	aApp "github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/httpSender"
	"github.com/maybecoding/go-metrics.git/internal/agent/memCollector"
)

func main() {
	// Config
	cfg := config.New()

	var memCollect aApp.Collector = memCollector.New()
	var httpSend aApp.Sender = httpSender.New(cfg.Sender.Address, cfg.Sender.Method, cfg.Sender.Template)

	app := aApp.New(memCollect, httpSend, cfg.App.SendIntervalSec, cfg.App.CollectIntervalSec)

	app.Start()

}
