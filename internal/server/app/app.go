package app

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"strconv"
	"time"
)

type App struct {
	store         Storage
	backupStorage BackupStorage
}

type Storage interface {
	SetMetricGauge(gauge *MetricGauge)
	SetMetricCounter(counter *MetricCounter)

	GetMetricGauge(name string) (TypeGauge, error)
	GetMetricCounter(name string) (TypeCounter, error)

	GetMetricGaugeAll() []*MetricGauge
	GetMetricCounterAll() []*MetricCounter
}

type BackupStorage interface {
	Save(metrics []*Metric) error
	Restore() ([]*Metric, error)
	GetBackupInterval() int64
	GetIsNeedRestore() bool
}

const (
	FmtFloat = 'f'
)

func (a *App) UpdateMetric(mType string, name string, value string) error {
	logger.Log.Debug().Str("type", mType).Str("ID", name).Str("value", value).Msg("UpdateMetric")
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
	if a.backupStorage.GetBackupInterval() == 0 {
		metrics := a.GetMetricsAll()
		err := a.backupStorage.Save(metrics)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error due save metrics in update metric")
		}
	}
	return nil
}

func (a *App) GetMetric(mType string, name string) (res string, e error) {
	defer func() {
		logger.Log.Debug().Str("type", mType).Str("ID", name).Str("value", res).Msg("GetMetric")
	}()

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
		value := strconv.FormatFloat(float64(m.Value), FmtFloat, -1, 64)
		metrics = append(metrics, &Metric{Type: Gauge, Name: name, Value: value})
	}

	for _, m := range mtrCounter {
		name := m.Name
		value := strconv.FormatInt(int64(m.Value), 10)
		metrics = append(metrics, &Metric{Type: Counter, Name: name, Value: value})
	}

	return metrics
}

func (a *App) StartBackupTimer(ctx context.Context) error {
	interval := a.backupStorage.GetBackupInterval()
	if interval == 0 {
		return nil
	}
	for {
		select {
		case <-time.After(time.Second * time.Duration(interval)):
			err := a.backupStorage.Save(a.GetMetricsAll())
			// Эту ошибку не выкидываем, она не критична
			if err != nil {
				logger.Log.Error().Err(err).Msg("error due saving metrics")
			}
		case <-ctx.Done():
			logger.Log.Info().Msg("start saving metrics on shutdown")
			err := a.backupStorage.Save(a.GetMetricsAll())
			if err != nil {
				return fmt.Errorf("error due saving metrics %w", err)
			}
			logger.Log.Info().Msg("metrics saved")
			return nil
		}
	}
}

func (a *App) restoreMetrics() {
	if !a.backupStorage.GetIsNeedRestore() {
		return
	}
	metrics, err := a.backupStorage.Restore()
	if err != nil {
		logger.Log.Error().Err(err).Msg("error due metrics restore")
	}
	for _, m := range metrics {
		err = a.UpdateMetric(m.Type, m.Name, m.Value)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error due restore metric")
		}
	}

}

func New(store Storage, backupStorage BackupStorage) *App {
	app := &App{store: store, backupStorage: backupStorage}
	app.restoreMetrics()
	return app
}
