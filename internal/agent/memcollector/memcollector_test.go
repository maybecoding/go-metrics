package memcollector

import (
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	metricsCount     = 29
	collectTestCount = 5_000
)

func TestMemCollector(t *testing.T) {

	t.Run("#1 Metrics are collected and PoolCount set working", func(t *testing.T) {
		memColl := New()
		memColl.CollectMetrics()
		metrics := memColl.GetMetrics()
		assert.Equal(t, metricsCount, len(metrics))

		poolCountMetric := findMetricByID(metrics, "PoolCount")
		require.NotNil(t, poolCountMetric)
		require.NotNil(t, poolCountMetric.Delta)
		assert.Equal(t, 1, *poolCountMetric.Delta)

		// Вызовем сбор метик разок другой
		for i := 0; i < collectTestCount; i += 1 {
			memColl.CollectMetrics()
		}
		metrics = memColl.GetMetrics()
		poolCountMetric = findMetricByID(metrics, "PoolCount")
		require.NotNil(t, poolCountMetric)
		require.NotNil(t, poolCountMetric.Delta)
		assert.Equal(t, collectTestCount+1, *poolCountMetric.Delta)

	})
}

func findMetricByID(metrics []*app.Metrics, name string) *app.Metrics {
	for _, metric := range metrics {
		if metric.ID == name {
			return metric
		}
	}
	return nil
}
