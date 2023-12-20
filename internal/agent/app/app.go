package app

import (
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
	SendStart(inpM chan *Metrics)
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

func (a App) Start() {

	// Запускаем go-рутину, которая собирает метрики с нужным интервалом
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.Collector.CollectMetrics(a.CollectInterval)
	}()

	// Запускаем регулярное получение метрик
	chMetrics := make(chan *Metrics)
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.Collector.FetchMetrics(chMetrics, a.SendInterval)
	}()

	// Запускаем воркеров по отправке
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.Sender.SendStart(chMetrics)
	}()

	wg.Wait()
	close(chMetrics)

}
