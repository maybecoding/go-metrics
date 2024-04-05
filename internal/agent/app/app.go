package app

import (
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"sync"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

const (
	MetricGauge   = "gauge"
	MetricCounter = "counter"
)

type Collector interface {
	CollectMetrics(interval time.Duration)
	FetchMetrics(outM chan *Metrics, interval time.Duration)
}

type Sender interface {
	Run(inpM chan *Metrics)
}

type App struct {
	Collector
	Sender
	CollectInterval time.Duration
	SendInterval    time.Duration
}

func New(collector Collector, sender Sender, cInterval, sInterval time.Duration) *App {
	return &App{
		Collector:       collector,
		Sender:          sender,
		CollectInterval: cInterval,
		SendInterval:    sInterval,
	}
}

func (a App) Run() {

	// Запускаем go-рутину, которая собирает метрики с нужным интервалом
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			logger.Debug().Msg("stopped a.Collector.CollectMetrics")
		}()

		a.Collector.CollectMetrics(a.CollectInterval)
	}()

	// Запускаем регулярное получение метрик
	chMetrics := make(chan *Metrics)
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			logger.Debug().Msg("stopped a.Collector.FetchMetrics")
		}()
		a.Collector.FetchMetrics(chMetrics, a.SendInterval)
		close(chMetrics)
	}()

	// Запускаем воркеров по отправке
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			logger.Debug().Msg("stopped Sender.Run")
		}()
		a.Sender.Run(chMetrics)
	}()

	wg.Wait()
}
