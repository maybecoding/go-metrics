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
	"time"
)

func main() {
	// Config
	cfg := config.New()
	logger.Init(cfg.Log.Level)
	cfg.LogDebug()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	time.Sleep(time.Second)
	var collect app.Collector = collector.New(ctx)
	var snd app.Sender = sender.New(cfg.Sender.EndpointTemplate, cfg.Sender.Address, cfg.Sender.RetryIntervals, cfg.Sender.HashKey, ctx, cfg.Sender.NumWorkers)

	a := app.New(collect, snd, time.Duration(cfg.App.CollectIntervalSec)*time.Second, time.Duration(cfg.App.SendIntervalSec)*time.Second)

	a.Start()

}
