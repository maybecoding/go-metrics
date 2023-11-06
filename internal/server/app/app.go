package app

import (
	"fmt"
	"strconv"
)

type App struct {
	store Storage
}

const (
	formatFloat = 'f'
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
		return strconv.FormatFloat(float64(value), formatFloat, -1, 64), nil
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

	metrics := make([]*Metric, len(mtrGauge)+len(mtrCounter))

	m, g, c := 0, 0, 0
	for g < len(mtrGauge) {
		fmt.Println("g:", g)
		name := mtrGauge[g].Name
		value := strconv.FormatInt(int64(mtrGauge[g].Value), 10)
		metrics[m] = &Metric{Type: Gauge, Name: name, Value: value}
		m += 1
		g += 1
	}
	for c < len(mtrCounter) {
		name := mtrGauge[c].Name
		value := strconv.FormatFloat(float64(mtrCounter[c].Value), formatFloat, -1, 64)
		metrics[m] = &Metric{Type: Counter, Name: name, Value: value}
		m += 1
		c += 1
	}

	return metrics
}

func New(store Storage) *App {
	return &App{store: store}
}
