package collector

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	metricsCount     = 31
	collectTestCount = 100
)

func TestMemCollector(t *testing.T) {

	t.Run("#1 Metrics are collected and PollCount set working", func(t *testing.T) {
		logger.Init("debug")
		ctx, cancel := context.WithCancel(context.Background())
		memColl := New(ctx)
		memColl.collectAll()
		metrics := memColl.GetMetrics()
		assert.Greater(t, len(metrics), metricsCount)

		pollCountMetric := findMetricByID(metrics, "PollCount")
		require.NotNil(t, pollCountMetric)
		require.NotNil(t, pollCountMetric.Delta)
		assert.Equal(t, int64(1), *pollCountMetric.Delta)

		// Вызовем сбор метик разок другой
		for i := 0; i < collectTestCount; i += 1 {
			memColl.collectAll()
		}
		metrics = memColl.GetMetrics()
		pollCountMetric = findMetricByID(metrics, "PollCount")
		require.NotNil(t, pollCountMetric)
		require.NotNil(t, pollCountMetric.Delta)
		assert.Equal(t, int64(collectTestCount+1), *pollCountMetric.Delta)
		cancel()
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
