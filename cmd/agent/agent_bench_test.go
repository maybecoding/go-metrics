package main

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/collector"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"github.com/maybecoding/go-metrics.git/internal/agent/sender"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"testing"
	"time"
)

// TestMainBench Тест используется для замера потребления памяти
func BenchmarkMain(b *testing.B) {

	cfg := &config.Config{
		App: config.App{
			CollectInterval: 0,
			SendInterval:    0,
		},
		Sender: config.Sender{
			Server:           "localhost:9099",
			EndpointTemplate: "%s://%s/update/",
			RetryIntervals:   []time.Duration{time.Millisecond, 3 * time.Millisecond, 5 * time.Millisecond},
			HashKey:          "key",
			NumWorkers:       10,
		},
		Log: config.Log{
			Level: "panic",
		},
	}
	logger.Init(cfg.Log.Level)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var collect app.Collector = collector.New(ctx)
	var snd app.Sender = sender.New(ctx, cfg.Sender)

	a := app.New(collect, snd, cfg.App.CollectInterval, cfg.App.SendInterval)

	a.Run()
}
