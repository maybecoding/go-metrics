package memcollector

import (
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const (
	metricsCount     = 29
	collectTestCount = 50_000
)

func TestMemCollector(t *testing.T) {

	t.Run("#1 Metrics are collected and PoolCount set working", func(t *testing.T) {
		memColl := New()
		memColl.CollectMetrics()
		metrics := memColl.GetMetrics()
		assert.Equal(t, metricsCount, len(metrics))

		poolCountMetric := findMetricByName(metrics, "PoolCount")
		require.NotNil(t, poolCountMetric)
		assert.Equal(t, "1", poolCountMetric.Value)

		// Вызовем сбор метик разок другой
		for i := 0; i < collectTestCount; i += 1 {
			memColl.CollectMetrics()
		}
		metrics = memColl.GetMetrics()
		poolCountMetric = findMetricByName(metrics, "PoolCount")
		require.NotNil(t, poolCountMetric)
		expectedCount := strconv.FormatInt(collectTestCount+1, 10)
		assert.Equal(t, expectedCount, poolCountMetric.Value)

	})
}

func findMetricByName(metrics []*app.Metric, name string) *app.Metric {
	for _, metric := range metrics {
		if metric.Name == name {
			return metric
		}
	}
	return nil
}
