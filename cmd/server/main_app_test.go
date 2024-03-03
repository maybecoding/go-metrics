package main

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/dbstorage"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkSetAll(b *testing.B) {
	logger.Init("fatal")

	// Если задана база данных
	var store sapp.Store
	var app *sapp.MetricService

	ctx := context.Background()

	dbStore := dbstorage.New("postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", ctx, []time.Duration{time.Second})
	store = dbStore
	defer dbStore.ConnectionClose()

	app = sapp.New(store)

	b.ResetTimer()
	times := 10_000
	// Создадим массив из элементов к отправке
	for i := 0; i < times; i++ {
		b.StopTimer()
		value := rand.Float64()
		mGauge := sapp.Metrics{
			ID:    fmt.Sprintf("Gauge_%d", i),
			MType: sapp.Gauge,
			Value: &value,
		}

		delta := int64(100)
		mCounter := sapp.Metrics{
			ID:    fmt.Sprintf("Counter_%d", i),
			MType: sapp.Counter,
			Delta: &delta,
		}

		b.StartTimer()
		err := app.Set(mGauge)
		require.NoError(b, err)

		err = app.Set(mCounter)
		require.NoError(b, err)
	}

}
