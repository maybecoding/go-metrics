package metricservice

import (
	"errors"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

type (
	MetricService struct {
		store Store
		//backupStorage BackupStorage
	}
	Metrics struct {
		ID    string   `json:"id" db:"name"`               // Имя метрики
		MType string   `json:"type" db:"type"`             // Параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty" db:"delta"` // Значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty" db:"value"` // Значение метрики в случае передачи gauge
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

	Set(Metrics) error
	SetAll([]Metrics) error
}

func (ms *MetricService) Get(m *Metrics) (e error) {
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
		logger.Debug().Str("type", m.MType).Float64("Value", *m.Value).Msg("GetMetric")
	} else if m.MType == Counter && m.Delta != nil {
		logger.Debug().Str("type", m.MType).Int64("Value", *m.Delta).Msg("GetMetric")
	}

	return err
}

func (ms *MetricService) Set(m Metrics) error {

	if m.MType != Gauge && m.MType != Counter {
		return ErrMetricTypeIncorrect
	}

	if m.MType == Gauge && m.Value == nil || m.MType == Counter && m.Delta == nil {
		return ErrNoMetricValue
	}

	if m.MType == Gauge {
		m.Delta = nil
		logger.Debug().Str("type", m.MType).Str("ID", m.ID).Float64("Value", *m.Value).Msg("UpdateMetric")
	} else {
		m.Value = nil
		logger.Debug().Str("type", m.MType).Str("ID", m.ID).Int64("Value", *m.Delta).Msg("UpdateMetric")
	}

	return ms.store.Set(m)

}

func (ms *MetricService) GetAll() ([]*Metrics, error) {
	return ms.store.GetAll()
}

func (ms *MetricService) SetAll(mts []Metrics) error {
	err := ms.store.SetAll(mts)
	if err != nil {
		logger.Error().Err(err).Msg("error due SetAll metrics")
	}
	return err
}

func New(store Store) *MetricService {
	app := &MetricService{store: store}
	return app
}
