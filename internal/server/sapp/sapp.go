package sapp

import (
	"fmt"
	"strconv"
)

type Storage interface {
	SetMetricGauge(metricName string, metricValue float64)
	SetMetricCounter(metricName string, metricValue int64)
}

type App struct {
	store Storage
}

func (a *App) UpdateMetric(mType string, name string, value string) error {
	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		a.store.SetMetricGauge(name, value)
	case "counter":
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		a.store.SetMetricCounter(name, value)
	default:
		return fmt.Errorf("incorrect metric type")
	}
	return nil
}

func New(store Storage) *App {
	return &App{store: store}
}
