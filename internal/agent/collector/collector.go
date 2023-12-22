package collector

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"runtime"
	"sync"
	"time"
)

type Collector struct {
	gaugeMetrics map[string]float64
	pollCount    int64

	memStats  *runtime.MemStats
	ctx       context.Context
	muGauge   sync.RWMutex
	muCounter sync.RWMutex
}

func New(ctx context.Context) *Collector {
	return &Collector{
		gaugeMetrics: make(map[string]float64),
		memStats:     &runtime.MemStats{},
		ctx:          ctx,
	}
}

func (c *Collector) CollectMetrics(interval time.Duration) {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(interval):
			c.collectAll()
		}
	}
}

func (c *Collector) FetchMetrics(outM chan *app.Metrics, interval time.Duration) {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(interval):
			c.muGauge.RLock()
			for mType, mt := range c.gaugeMetrics {
				m := mt
				outM <- &app.Metrics{MType: app.MetricGauge, ID: mType, Value: &m}
			}
			c.muGauge.RUnlock()

			c.muCounter.RLock()
			cnt := c.pollCount
			outM <- &app.Metrics{MType: app.MetricCounter, ID: "PollCount", Delta: &cnt}
			c.muCounter.RUnlock()
		}
	}

}

func (c *Collector) GetMetrics() []*app.Metrics {
	metrics := make([]*app.Metrics, 0, len(c.gaugeMetrics)+1)

	c.muGauge.RLock()
	for mType, mt := range c.gaugeMetrics {
		m := mt
		metrics = append(metrics, &app.Metrics{MType: app.MetricGauge, ID: mType, Value: &m})
	}
	c.muGauge.RUnlock()

	c.muCounter.RLock()
	cnt := c.pollCount
	c.muCounter.RUnlock()
	metrics = append(metrics, &app.Metrics{MType: app.MetricCounter, ID: "PollCount", Delta: &cnt})
	return metrics
}
