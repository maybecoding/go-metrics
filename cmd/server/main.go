package main

import (
	"context"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/backupstorage"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/handlers"
	"github.com/maybecoding/go-metrics.git/internal/server/memstorage"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Получаем конфигурацию приложения
	cfg := config.NewConfig()
	logger.Init(cfg.Log.Level)
	logger.Log.Debug().Str("backup file path", cfg.BackupStorage.Path).Msg("initialization")

	// Создаем хранилище и бэкапер
	var store sapp.Storage = smemstorage.NewMemStorage()
	var backupStorage sapp.BackupStorage = backupstorage.NewBackupStorage(
		cfg.BackupStorage.Interval,
		cfg.BackupStorage.Path,
		cfg.BackupStorage.IsRestoreOnUp,
	)
	app := sapp.New(store, backupStorage)

	// Создаем контроллер и вверяем ему приложение
	contr := handlers.New(app, cfg.Server.Address)

	// Контекст, который будет отменен при выходе из приложения Ctrl + C
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	g, gCtx := errgroup.WithContext(ctx)

	// Запускаем в приложении механизм бэкапирования
	g.Go(func() error {
		return app.StartBackupTimer(gCtx)
	})

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
		logger.Log.Info().Err(err).Msg("app stopped")
	}

}
