package metricmemstorage

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"sync"
	"time"
)

type MetricMemStorage struct {
	dataGauge      map[string]float64
	dataCounter    map[string]int64
	dumper         *Dumper
	backupInterval int64
	muGauge        sync.RWMutex
	muCounter      sync.RWMutex
}

func (mem *MetricMemStorage) Set(m metricservice.Metrics) error {
	if m.MType == metricservice.Gauge {
		mem.muGauge.Lock()
		mem.dataGauge[m.ID] = *m.Value
		mem.muGauge.Unlock()
	} else if m.MType == metricservice.Counter {
		mem.muCounter.Lock()
		mem.dataCounter[m.ID] += *m.Delta
		mem.muCounter.Unlock()
	}
	return nil
}

func (mem *MetricMemStorage) Get(m *metricservice.Metrics) error {
	if m.MType == metricservice.Gauge {
		mem.muGauge.RLock()
		v, ok := mem.dataGauge[m.ID]
		mem.muGauge.RUnlock()
		if !ok {
			return metricservice.ErrNoMetricValue
		}
		m.Value = &v
	} else if m.MType == metricservice.Counter {
		mem.muCounter.RLock()
		d, ok := mem.dataCounter[m.ID]
		mem.muCounter.RUnlock()
		if !ok {
			return metricservice.ErrNoMetricValue
		}
		m.Delta = &d
	}

	return nil
}

func (mem *MetricMemStorage) GetAll() ([]*metricservice.Metrics, error) {
	mtr := make([]*metricservice.Metrics, 0, len(mem.dataGauge)+len(mem.dataCounter))
	mem.muGauge.RLock()
	for name, value := range mem.dataGauge {
		v := value
		mtr = append(mtr, &metricservice.Metrics{ID: name, MType: metricservice.Gauge, Value: &v})
	}
	mem.muGauge.RUnlock()

	mem.muCounter.RLock()
	for name, value := range mem.dataCounter {
		v := value
		mtr = append(mtr, &metricservice.Metrics{ID: name, MType: metricservice.Counter, Delta: &v})
	}
	mem.muCounter.RUnlock()
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
				logger.Error().Err(err).Msg("error due get metrics for save")
			}
			err = mem.dumper.Save(ms)
			// Эту ошибку не выкидываем, она не критична
			if err != nil {
				logger.Error().Err(err).Msg("error due saving metric")
			}
		case <-ctx.Done():
			logger.Info().Msg("start saving metric on shutdown")
			ms, err := mem.GetAll()
			if err != nil {
				logger.Error().Err(err).Msg("error due get metrics for save")
			}
			err = mem.dumper.Save(ms)
			if err != nil {
				return fmt.Errorf("error due saving metric %w", err)
			}
			logger.Info().Msg("metric saved")
			return nil
		}
	}
}

func (mem *MetricMemStorage) restoreMetrics() {
	metrics, err := mem.dumper.Restore()
	if err != nil {
		logger.Error().Err(err).Msg("error due metric restore")
	}
	for _, m := range metrics {
		err = mem.Set(m)
		if err != nil {
			logger.Error().Err(err).Msg("error due restore metric")
		}
	}
}

func (mem *MetricMemStorage) SetAll(mts []metricservice.Metrics) error {
	// Чисто для обеспечения обратной совместимости
	for _, mt := range mts {
		err := mem.Set(mt)
		if err != nil {
			return err
		}
	}
	return nil
}
