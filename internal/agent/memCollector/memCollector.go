package memCollector

import (
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"math/rand"
	"runtime"
	"strconv"
)

const (
	MetricGauge   = "gauge"
	MetricCounter = "counter"
)

type MemCollector struct {
	gaugeMetrics map[string]float64
	poolCount    int64
	memStats     *runtime.MemStats
}

func New() *MemCollector {
	return &MemCollector{
		gaugeMetrics: make(map[string]float64),
		memStats:     &runtime.MemStats{},
	}
}

func (c *MemCollector) CollectMetrics() {
	// Собираем тетрики gauge
	runtime.ReadMemStats(c.memStats)
	c.gaugeMetrics["Alloc"] = float64(c.memStats.Alloc)
	c.gaugeMetrics["BuckHashSys"] = float64(c.memStats.BuckHashSys)
	c.gaugeMetrics["Frees"] = float64(c.memStats.Frees)
	c.gaugeMetrics["GCCPUFraction"] = c.memStats.GCCPUFraction
	c.gaugeMetrics["GCSys"] = float64(c.memStats.GCSys)
	c.gaugeMetrics["HeapAlloc"] = float64(c.memStats.HeapAlloc)
	c.gaugeMetrics["HeapIdle"] = float64(c.memStats.HeapIdle)
	c.gaugeMetrics["HeapInuse"] = float64(c.memStats.HeapInuse)
	c.gaugeMetrics["HeapObjects"] = float64(c.memStats.HeapObjects)
	c.gaugeMetrics["HeapReleased"] = float64(c.memStats.HeapReleased)
	c.gaugeMetrics["HeapSys"] = float64(c.memStats.HeapSys)
	c.gaugeMetrics["LastGC"] = float64(c.memStats.LastGC)
	c.gaugeMetrics["Lookups"] = float64(c.memStats.Lookups)
	c.gaugeMetrics["MCacheInuse"] = float64(c.memStats.MCacheInuse)
	c.gaugeMetrics["MCacheSys"] = float64(c.memStats.MCacheSys)
	c.gaugeMetrics["MSpanInuse"] = float64(c.memStats.MSpanInuse)
	c.gaugeMetrics["MSpanSys"] = float64(c.memStats.MSpanSys)
	c.gaugeMetrics["Mallocs"] = float64(c.memStats.Mallocs)
	c.gaugeMetrics["NextGC"] = float64(c.memStats.NextGC)
	c.gaugeMetrics["NumForcedGC"] = float64(c.memStats.NumForcedGC)
	c.gaugeMetrics["NumGC"] = float64(c.memStats.NumGC)
	c.gaugeMetrics["OtherSys"] = float64(c.memStats.OtherSys)
	c.gaugeMetrics["PauseTotalNs"] = float64(c.memStats.PauseTotalNs)
	c.gaugeMetrics["StackInuse"] = float64(c.memStats.StackInuse)
	c.gaugeMetrics["StackSys"] = float64(c.memStats.StackSys)
	c.gaugeMetrics["Sys"] = float64(c.memStats.Sys)
	c.gaugeMetrics["TotalAlloc"] = float64(c.memStats.TotalAlloc)
	c.gaugeMetrics["RandomValue"] = rand.Float64()

	// Собираем метики counter
	c.poolCount += 1

}

func (c *MemCollector) GetMetrics() []*app.Metric {
	metrics := make([]*app.Metric, len(c.gaugeMetrics)+1)
	i := 0
	for mType, m := range c.gaugeMetrics {
		metrics[i] = &app.Metric{Type: MetricGauge, Name: mType, Value: strconv.FormatFloat(m, 'E', -1, 64)}
		i += 1
	}
	metrics[i] = &app.Metric{Type: MetricCounter, Name: "PoolCount", Value: strconv.FormatInt(c.poolCount, 10)}
	return metrics
}
