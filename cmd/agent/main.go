package main

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/collector"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"github.com/maybecoding/go-metrics.git/internal/agent/sender"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
	"os/signal"
)

func main() {
	printInfo()
	// Config
	cfg := config.New()
	logger.Init(cfg.Log.Level)
	cfg.LogDebug()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	var collect app.Collector = collector.New(ctx)
	var snd app.Sender = sender.New(ctx, cfg.Sender)

	a := app.New(collect, snd, cfg.App.CollectInterval, cfg.App.SendInterval)

	a.Run()

}
