package app

import (
	"strconv"
)

type App struct {
	store Storage
}

type Storage interface {
	SetMetricGauge(gauge *MetricGauge)
	SetMetricCounter(counter *MetricCounter)

	GetMetricGauge(name string) (TypeGauge, error)
	GetMetricCounter(name string) (TypeCounter, error)

	GetMetricGaugeAll() []*MetricGauge
	GetMetricCounterAll() []*MetricCounter
}

const (
	FmtFloat = 'f'
)

func (a *App) UpdateMetric(mType string, name string, value string) error {
	switch mType {
	case Gauge:
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		a.store.SetMetricGauge(&MetricGauge{Name: name, Value: TypeGauge(value)})
	case Counter:
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		a.store.SetMetricCounter(&MetricCounter{Name: name, Value: TypeCounter(value)})
	default:
		return ErrMetricTypeIncorrect
	}
	return nil
}

func (a *App) GetMetric(mType string, name string) (string, error) {
	switch mType {
	case "gauge":
		value, err := a.store.GetMetricGauge(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(float64(value), FmtFloat, -1, 64), nil
	case "counter":
		value, err := a.store.GetMetricCounter(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(value), 10), nil
	default:
		return "", ErrNoMetricValue
	}
}

func (a *App) GetMetricsAll() []*Metric {
	mtrGauge := a.store.GetMetricGaugeAll()
	mtrCounter := a.store.GetMetricCounterAll()

	metrics := make([]*Metric, 0, len(mtrGauge)+len(mtrCounter))

	for _, m := range mtrGauge {
		name := m.Name
		value := strconv.FormatInt(int64(m.Value), 10)
		metrics = append(metrics, &Metric{Type: Gauge, Name: name, Value: value})
	}

	for _, m := range mtrCounter {
		name := m.Name
		value := strconv.FormatFloat(float64(m.Value), FmtFloat, -1, 64)
		metrics = append(metrics, &Metric{Type: Gauge, Name: name, Value: value})
	}

	return metrics
}

func New(store Storage) *App {
	return &App{store: store}
}
