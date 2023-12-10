package metricmemstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"time"
)

type MetricMemStorage struct {
	dataGauge      map[string]float64
	dataCounter    map[string]int64
	dumper         *Dumper
	backupInterval int64
}

func (mem *MetricMemStorage) Set(m *metric.Metrics) error {
	if m.MType == metric.Gauge {
		mem.dataGauge[m.ID] = *m.Value
	} else if m.MType == metric.Counter {
		mem.dataCounter[m.ID] += *m.Delta
	}

	return nil
}

func (mem *MetricMemStorage) Get(m *metric.Metrics) error {
	if m.MType == metric.Gauge {
		v, ok := mem.dataGauge[m.ID]
		if !ok {
			return metric.ErrNoMetricValue
		}
		m.Value = &v
	} else if m.MType == metric.Counter {
		d, ok := mem.dataCounter[m.ID]
		if !ok {
			return metric.ErrNoMetricValue
		}
		m.Delta = &d
	}

	return nil
}

func (mem *MetricMemStorage) GetAll() ([]*metric.Metrics, error) {
	mtr := make([]*metric.Metrics, 0, len(mem.dataGauge)+len(mem.dataCounter))
	for name, value := range mem.dataGauge {
		v := value
		mtr = append(mtr, &metric.Metrics{ID: name, MType: metric.Gauge, Value: &v})
	}
	for name, value := range mem.dataCounter {
		v := value
		mtr = append(mtr, &metric.Metrics{ID: name, MType: metric.Counter, Delta: &v})
	}
	return mtr, nil
}

func NewMemStorage(d *Dumper, backupInterval int64, restoreOnUp bool) *MetricMemStorage {
	mem := &MetricMemStorage{
		dataGauge:      make(map[string]float64),
		dataCounter:    make(map[string]int64),
		dumper:         d,
		backupInterval: backupInterval,
	}
	if restoreOnUp {
		mem.restoreMetrics()
	}
	return mem
}

func (mem *MetricMemStorage) StartBackupTimer(ctx context.Context) error {
	if mem.backupInterval == 0 {
		return nil
	}
	for {
		select {
		case <-time.After(time.Second * time.Duration(mem.backupInterval)):
			ms, err := mem.GetAll()
			if err != nil {
				logger.Log.Error().Err(err).Msg("error due get metrics for save")
			}
			err = mem.dumper.Save(ms)
			// Эту ошибку не выкидываем, она не критична
			if err != nil {
				logger.Log.Error().Err(err).Msg("error due saving metric")
			}
		case <-ctx.Done():
			logger.Log.Info().Msg("start saving metric on shutdown")
			ms, err := mem.GetAll()
			if err != nil {
				logger.Log.Error().Err(err).Msg("error due get metrics for save")
			}
			err = mem.dumper.Save(ms)
			if err != nil {
				return fmt.Errorf("error due saving metric %w", err)
			}
			logger.Log.Info().Msg("metric saved")
			return nil
		}
	}
}

func (mem *MetricMemStorage) restoreMetrics() {
	metrics, err := mem.dumper.Restore()
	if err != nil {
		logger.Log.Error().Err(err).Msg("error due metric restore")
	}
	for _, m := range metrics {
		err = mem.Set(m)
		if err != nil {
			logger.Log.Error().Err(err).Msg("error due restore metric")
		}
	}
}

func (mem *MetricMemStorage) Ping() error {
	return errors.New("incorrect ping call for storage in memory")
}

func (mem *MetricMemStorage) SetAll(mts []*metric.Metrics) error {
	// Чисто для обеспечениия обратной совместимости
	for _, mt := range mts {
		err := mem.Set(mt)
		if err != nil {
			return err
		}
	}
	return nil
}
