package app

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/dbstorage"
	"github.com/maybecoding/go-metrics.git/internal/server/handlers"
	"github.com/maybecoding/go-metrics.git/internal/server/metricmemstorage"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/health"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/maybecoding/go-metrics.git/pkg/starter"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	store     metricservice.Store
	metricSrv *metricservice.MetricService
	handler   *handlers.Handler
	starter   *starter.Starter
	cfg       *config.Config
	hl        *health.Health
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Init() *App {
	logger.Init(a.cfg.Log.Level)
	a.cfg.LogDebug()
	logger.Debug().Str("backup file path", a.cfg.BackupStorage.Path).Msg("initialization")

	// Контекст, который будет отменен при выходе из приложения Ctrl + C
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// HealthCheck
	a.hl = health.New()

	// Starter - ИНИЦИАЛИЗАЦИЯ
	a.starter = starter.New(ctx)

	// Выбираем и инициализируем store в зависимости от настроек
	if a.cfg.Database.Use() {
		dbStore := dbstorage.New(a.cfg.Database.ConnStr, ctx, a.cfg.Database.RetryIntervals)
		if a.cfg.Database.RunMigrations {
			dbstorage.Migrate(a.cfg.Database.ConnStr)
		}
		// Просим HealthCheck присмотреть за БД
		a.hl.Watch(dbStore.Ping)
		a.store = dbStore
		a.starter.OnShutdown(func(_ context.Context) error {
			dbStore.ConnectionClose()
			return nil
		})
	} else {
		dumper := metricmemstorage.NewDumper(a.cfg.BackupStorage.Path)
		memStore := metricmemstorage.NewMemStorage(dumper, a.cfg.BackupStorage.Interval, a.cfg.BackupStorage.IsRestoreOnUp)
		// START mem store backup
		a.starter.OnRun(memStore.StartBackupTimer)
		a.store = memStore
	}
	// Создаем сервис
	a.metricSrv = metricservice.New(a.store)
	return a
}

func (a *App) InitHandler() *App {
	if a.metricSrv == nil {
		panic("first app.Start method must be called, service is not initialized")
	}
	// Create  handler
	a.handler = handlers.New(a.metricSrv, a.cfg.Server, a.hl)
	// ЗАПУСКАЕМ контроллер
	a.starter.OnRun(a.handler.Start)
	// После завершения приложения потушим сервер
	a.starter.OnShutdown(a.handler.Shutdown)
	return a
}

func (a *App) Run() {
	// Ожидаем завершения работы запущенных рутин через start.Run
	if err := a.starter.Run(); err != nil {
		logger.Info().Err(err).Msg("metric stopped")
	}
}

func (a *App) GetMetricService() *metricservice.MetricService {
	return a.metricSrv
}

func (a *App) GetHandler() *handlers.Handler {
	return a.handler
}
