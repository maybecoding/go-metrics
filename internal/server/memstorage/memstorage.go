package smemstorage

import (
	"github.com/maybecoding/go-metrics.git/internal/server/app"
)

type MemStorage struct {
	dataGauge   map[string]app.TypeGauge
	dataCounter map[string]app.TypeCounter
}

func (m *MemStorage) SetMetricGauge(metric *app.MetricGauge) {
	m.dataGauge[metric.Name] = metric.Value
}

func (m *MemStorage) SetMetricCounter(metric *app.MetricCounter) {
	m.dataCounter[metric.Name] += metric.Value
}

func (m *MemStorage) GetMetricGauge(name string) (app.TypeGauge, error) {
	res, ok := m.dataGauge[name]
	if !ok {
		return 0, app.ErrNoMetricValue
	}
	return res, nil
}

func (m *MemStorage) GetMetricCounter(name string) (app.TypeCounter, error) {
	res, ok := m.dataCounter[name]
	if !ok {
		return 0, app.ErrNoMetricValue
	}
	return res, nil
}

func (m *MemStorage) GetMetricGaugeAll() []*app.MetricGauge {
	mtr := make([]*app.MetricGauge, len(m.dataGauge))
	i := 0
	for name, value := range m.dataGauge {
		mtr[i] = &app.MetricGauge{Name: name, Value: value}
		i += 1
	}
	return mtr
}

func (m *MemStorage) GetMetricCounterAll() []*app.MetricCounter {
	mtr := make([]*app.MetricCounter, len(m.dataCounter))
	i := 0
	for name, value := range m.dataCounter {
		mtr[i] = &app.MetricCounter{Name: name, Value: value}
		i += 1
	}
	return mtr
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		dataGauge:   make(map[string]app.TypeGauge),
		dataCounter: make(map[string]app.TypeCounter),
	}
}
