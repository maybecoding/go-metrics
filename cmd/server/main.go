package main

import (
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/backupstorage"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/controller"
	"github.com/maybecoding/go-metrics.git/internal/server/memstorage"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
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

	// Запускаем в приложении механихм бэкапирования
	go app.StartBackupTimer()

	// Создаем контроллер и вверяем ему приложение
	contr := controller.New(app, cfg.Server.Address)

	// Запускаем
	contr.Start()

}
