package handlers

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/health"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
	"sync"
)

const (
	FmtFloat = 'f'
)

type Handler struct {
	metric        *metricservice.MetricService
	serverAddress string
	pprofAddress  string
	server        *http.Server
	health        *health.Health
	pprofServer   *http.Server
	hashKey       string
	cryptoKey     string
}

func New(app *metricservice.MetricService, cfg config.Server, hl *health.Health, hk string) *Handler {
	return &Handler{metric: app, health: hl, hashKey: hk, serverAddress: cfg.Address, pprofAddress: cfg.PprofAddress, cryptoKey: cfg.CryptoKey}
}

func (c *Handler) Start(_ context.Context) error {
	// Инициализируем сервер
	c.server = &http.Server{Addr: c.serverAddress, Handler: c.GetRouter()}

	// Инициализируем pprof
	c.pprofServer = &http.Server{Addr: c.pprofAddress, Handler: pprofRouter()}

	logger.Info().Str("address", c.pprofAddress).Msg("Start pprof server")
	go func() {
		err := c.pprofServer.ListenAndServe()
		if err != nil {
			logger.Error().Err(err).Msg("handlers - start - c.pprofServer.ListenAndServe")
		}
	}()

	logger.Info().Str("address", c.serverAddress).Msg("Starting server")
	if c.cryptoKey != "" {
		return fmt.Errorf("server error %w, or server just stoped", c.server.ListenAndServeTLS(c.cryptoKey, c.cryptoKey))
	}
	return fmt.Errorf("server error %w, or server just stoped", c.server.ListenAndServe())
}

var mtsPool = sync.Pool{
	New: func() any {
		ml := make(entity.MetricsList, 0, 50)
		return &ml
	},
}

func (c *Handler) Shutdown(_ context.Context) error {
	logger.Info().Msg("Ctrl + C command got, shutting down server")
	_ = c.pprofServer.Shutdown(context.Background())
	return c.server.Shutdown(context.Background())
}
