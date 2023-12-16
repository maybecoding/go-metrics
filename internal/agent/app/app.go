package app

import (
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

const (
	MetricGauge   = "gauge"
	MetricCounter = "counter"
)

type Collector interface {
	CollectMetrics()
	GetMetrics() []*Metrics
}

type Sender interface {
	Send([]*Metrics)
}

type App struct {
	Collector
	Sender
	SendIntervalSec    int
	CollectIntervalSec int
}

func New(collector Collector, sender Sender, sendIntervalSec int, collectIntervalSec int) *App {
	return &App{
		Collector:          collector,
		Sender:             sender,
		CollectIntervalSec: collectIntervalSec,
		SendIntervalSec:    sendIntervalSec,
	}
}

func (a App) Start() {
	// Пока без горутин и мьютексов будем читерить))
	period := a.SendIntervalSec * a.CollectIntervalSec
	logger.Info().Int("period", period).Msg("starting collecting and sending metric")

	// Всегда
	for {
		for sec := 1; sec <= period; sec += 1 {
			time.Sleep(time.Second)
			// Если самое время запускать сборку метик - запускаем
			if sec%a.CollectIntervalSec == 0 {
				a.Collector.CollectMetrics()
				logger.Debug().Msg("metric collected (I hope)")
			}
			// Если наступило время отправлять
			if sec%a.SendIntervalSec == 0 {
				metrics := a.Collector.GetMetrics()
				a.Sender.Send(metrics)
				logger.Debug().Msg("metric sent (I hope)")
			}
		}
	}

}
