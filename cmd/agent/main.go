package main

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/collector"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"github.com/maybecoding/go-metrics.git/internal/agent/grpcsender"
	"github.com/maybecoding/go-metrics.git/internal/agent/sender"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
	"os/signal"
)

func main() {
	printInfo()
	// Config
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("failed to parse config: %s\n", err.Error())
		return
	}
	logger.Init(cfg.Log.Level)
	cfg.LogDebug()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var collect app.Collector = collector.New(ctx)
	var snd app.Sender
	if cfg.UseGRPC() {
		snd = grpcsender.New(ctx, cfg.Sender)
	} else {
		snd = sender.New(ctx, cfg.Sender)
	}
	a := app.New(collect, snd, cfg.App.CollectInterval, cfg.App.SendInterval)

	a.Run()

}
