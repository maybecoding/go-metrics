package main

import (
	aApp "github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"github.com/maybecoding/go-metrics.git/internal/agent/httpjsonsender"
	"github.com/maybecoding/go-metrics.git/internal/agent/logger"
	"github.com/maybecoding/go-metrics.git/internal/agent/memcollector"
)

func main() {
	// Config
	cfg := config.New()
	logger.Init(cfg.Log.Level)

	var memCollect aApp.Collector = memcollector.New()
	//var httpSend aApp.Sender = httpsender.New(cfg.Sender.Address, cfg.Sender.Method, cfg.Sender.Template)
	var jsonSender aApp.Sender = httpjsonsender.New(cfg.Sender.JSONEndpoint, cfg.Sender.Address)

	app := aApp.New(memCollect, jsonSender, cfg.App.SendIntervalSec, cfg.App.CollectIntervalSec)

	app.Start()

}
