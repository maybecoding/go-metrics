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
	ErrMetricTypeIncorrect = errors.New("metric type incorrect")
	ErrNoMetricValue       = errors.New("no metric value")
)

type Store interface {
	Get(*Metrics) error
	GetAll() ([]*Metrics, error)

	Set(*Metrics) error
	SetAll([]*Metrics) error

	Ping() error
}

func (ms *Metric) Get(m *Metrics) (e error) {
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
		logger.Log.Debug().Str("type", m.MType).Float64("Value", *m.Value).Msg("GetMetric")
	} else if m.MType == Counter && m.Delta != nil {
		logger.Log.Debug().Str("type", m.MType).Int64("Value", *m.Delta).Msg("GetMetric")
	}

	return err
}

func (ms *Metric) Set(m *Metrics) error {

	if m.MType != Gauge && m.MType != Counter {
		return ErrMetricTypeIncorrect
	}

	if m.MType == Gauge && m.Value == nil || m.MType == Counter && m.Delta == nil {
		return ErrNoMetricValue
	}

	if m.MType == Gauge {
		m.Delta = nil
		logger.Log.Debug().Str("type", m.MType).Str("ID", m.ID).Float64("Value", *m.Value).Msg("UpdateMetric")
	} else {
		m.Value = nil
		logger.Log.Debug().Str("type", m.MType).Str("ID", m.ID).Int64("Value", *m.Delta).Msg("UpdateMetric")
	}

	return ms.store.Set(m)

}

func (ms *Metric) GetAll() ([]*Metrics, error) {
	return ms.store.GetAll()
}

func (ms *Metric) SetAll(mts []*Metrics) error {
	err := ms.store.SetAll(mts)
	logger.Log.Err(err).Msg("error due SetAll metrics")
	return err
}

func (ms *Metric) Ping() error {
	return ms.store.Ping()
}

func New(store Store) *Metric {
	app := &Metric{store: store}
	//app.restoreMetrics()
	return app
}
