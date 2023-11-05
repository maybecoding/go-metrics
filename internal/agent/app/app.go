package app

import "time"

type Metric struct {
	Type  string
	Name  string
	Value string
}
type Collector interface {
	CollectMetrics()
	GetMetrics() []*Metric
}

type Sender interface {
	Send([]*Metric)
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

	// Всегда
	for {
		for sec := 1; sec <= period; sec += 1 {
			time.Sleep(time.Second)
			// Если самое время запускать сборку метик - запускаем
			if sec%a.CollectIntervalSec == 0 {
				a.Collector.CollectMetrics()
			}
			// Если наступило время отправлять
			if sec%a.SendIntervalSec == 0 {
				metrics := a.Collector.GetMetrics()
				a.Sender.Send(metrics)
			}
		}
	}

}
