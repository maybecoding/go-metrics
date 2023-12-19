package main

import (
	"context"
	"github.com/maybecoding/go-metrics.git/pkg/health"
	"os"
	"os/signal"
	"syscall"

	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/dbstorage"
	"github.com/maybecoding/go-metrics.git/internal/server/handlers"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/internal/server/metricmemstorage"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Получаем конфигурацию приложения
	cfg := config.NewConfig()
	logger.Init(cfg.Log.Level)
	logger.Debug().Str("backup file path", cfg.BackupStorage.Path).Msg("initialization")

	// Если задана база данных
	var store sapp.Store
	var memStore *metricmemstorage.MetricMemStorage
	var app *sapp.Metric

	// Контекст, который будет отменен при выходе из приложения Ctrl + C
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// HealthCheck
	hl := health.New()

	if cfg.Database.Use() {
		dbStore := dbstorage.New(cfg.Database.ConnStr, ctx, cfg.Database.RetryIntervals)
		// Просим HealthCheck присмотреть за БД
		hl.Watch(dbStore.Ping)
		store = dbStore
		defer dbStore.ConnectionClose()
	} else {
		dumper := metricmemstorage.NewDumper(cfg.BackupStorage.Path)
		memStore = metricmemstorage.NewMemStorage(dumper, cfg.BackupStorage.Interval, cfg.BackupStorage.IsRestoreOnUp)
		store = memStore
	}
	app = sapp.New(store)

	// Создаем контроллер и вверяем ему приложение
	contr := handlers.New(app, cfg.Server.Address, hl)

	g, gCtx := errgroup.WithContext(ctx)

	// Запускаем в приложении механизм бэкапирования
	if memStore != nil {
		g.Go(func() error {
			return memStore.StartBackupTimer(gCtx)
		})
	}

	// Запускаем сервер
	g.Go(func() error {
		return contr.Start()
	})

	// Запускаем выключатель для сервера
	g.Go(func() error {
		return contr.Shutdown(gCtx)
	})

	// Если вырубили приложение
	if err := g.Wait(); err != nil {
		logger.Info().Err(err).Msg("metric stopped")
	}

}
