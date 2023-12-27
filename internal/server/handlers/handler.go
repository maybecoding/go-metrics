package handlers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/maybecoding/go-metrics.git/pkg/hashcheck"
	"github.com/maybecoding/go-metrics.git/pkg/health"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/compress"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

const (
	FmtFloat = 'f'
)

type Handler struct {
	metric        *metric.Metric
	serverAddress string
	server        *http.Server
	health        *health.Health
	hashKey       string
}

func New(app *metric.Metric, serverAddress string, hl *health.Health, hk string) *Handler {
	return &Handler{metric: app, serverAddress: serverAddress, health: hl, hashKey: hk}
}

func (c *Handler) GetRouter() chi.Router {

	r := chi.NewRouter()

	// Подключаем логер
	r.Use(logger.Handler)

	// Подключаем проверку хэшей
	hashCheck := hashcheck.New(sha256.New, c.hashKey, "HashSHA256")
	//r.Use(hashCheck.Handler)

	// Установка значений
	r.Post("/update/{type}/{name}/{value}", c.metricUpdate)
	r.Post("/update/", compress.HandlerFuncReader(compress.HandlerFuncWriter(c.metricUpdateJSON, compress.BestSpeed)))

	r.Post("/updates/", hashCheck.HandlerFunc(compress.HandlerFuncReader(c.metricUpdateAllJSON)))

	// Получение значений
	r.Get("/value/{type}/{name}", c.metricGet)
	r.Post("/value/", compress.HandlerFuncReader(compress.HandlerFuncWriter(c.metricGetJSON, compress.BestSpeed)))

	// Отчет по метрикам
	r.Get("/", compress.HandlerFuncWriter(c.metricGetAll, compress.BestSpeed))

	// ping
	r.Get("/ping", c.ping)

	return r
}

func (c *Handler) Start() error {
	addr := c.serverAddress
	c.server = &http.Server{Addr: addr, Handler: c.GetRouter()}

	logger.Info().Str("address", addr).Msg("Starting server")
	return fmt.Errorf("server error %w, or server just stoped", c.server.ListenAndServe())
}

func (c *Handler) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	logger.Info().Msg("Ctrl + C command got, shutting down server")
	return c.server.Shutdown(context.Background())
}
