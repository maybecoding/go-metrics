package metric

import (
	"errors"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

type (
	Metric struct {
		store Store
		//backupStorage BackupStorage
	}
	Metrics struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var (
	ErrMetricTypeIncorrect  = errors.New("metric type incorrect")
	ErrMetricValueIncorrect = errors.New("metric value incorrect")
	ErrNoMetricValue        = errors.New("no metric value")
)

type Store interface {
	Set(metric *Metrics) error
	Get(metric *Metrics) error

	GetAll() ([]*Metrics, error)
	Ping() error
}

func (ms *Metric) Set(m *Metrics) error {

	lg := logger.Log.Debug().Str("type", m.MType).Str("ID", m.ID)

	if m.MType != Gauge && m.MType != Counter {
		return ErrMetricTypeIncorrect
	}

	if m.MType == Gauge && m.Value == nil || m.MType == Counter && m.Delta == nil {
		return ErrNoMetricValue
	}

	if m.MType == Gauge {
		lg = lg.Float64("Value", *m.Value)
	} else {
		lg = lg.Int64("Value", *m.Delta)
	}
	lg.Msg("UpdateMetric")

	return ms.store.Set(m)

}

func (ms *Metric) Get(m *Metrics) (e error) {
	lg := logger.Log.Debug().Str("type", m.MType).Str("ID", m.ID) //.Str("value", res).Msg("GetMetric")

	if m.MType != Gauge && m.MType != Counter {
		return ErrNoMetricValue //ErrMetricTypeIncorrect
	}

	if m.Value != nil {
		m.Value = nil
	}
	if m.Delta != nil {
		m.Delta = nil
	}

	err := ms.store.Get(m)

	if m.MType == Gauge && m.Value != nil {
		lg = lg.Float64("Value", *m.Value)
	} else if m.MType == Counter && m.Delta != nil {
		lg = lg.Int64("Value", *m.Delta)
	}
	lg.Msg("GetMetric")

	return err
}

func (ms *Metric) GetAll() ([]*Metrics, error) {
	return ms.store.GetAll()
}

func (ms *Metric) Ping() error {
	return ms.store.Ping()
}

func New(store Store) *Metric {
	app := &Metric{store: store}
	//app.restoreMetrics()
	return app
}
