package handlers

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/health"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net"
	"net/http"
	"sync"
)

const (
	FmtFloat = 'f'
)

type Handler struct {
	metric        *metricservice.MetricService
	pprofAddress  string
	server        *http.Server
	health        *health.Health
	pprofServer   *http.Server
	cfg           config.Server
	trustedSubNet *net.IPNet
}

func New(app *metricservice.MetricService, cfg config.Server, hl *health.Health) *Handler {
	_, ipNet, err := net.ParseCIDR(cfg.TrustedSubnet)
	if err != nil {
		logger.Error().Err(err).Msg("can't parse trusted subnet")
		ipNet = nil
	}
	return &Handler{metric: app, health: hl, cfg: cfg, pprofAddress: cfg.PprofAddress, trustedSubNet: ipNet}
}

func (c *Handler) Start(_ context.Context) error {
	// Инициализируем сервер
	c.server = &http.Server{Addr: c.cfg.Address, Handler: c.GetRouter()}

	// Инициализируем pprof
	c.pprofServer = &http.Server{Addr: c.pprofAddress, Handler: pprofRouter()}

	logger.Info().Str("address", c.pprofAddress).Msg("Start pprof server")
	go func() {
		err := c.pprofServer.ListenAndServe()
		if err != nil {
			logger.Error().Err(err).Msg("handlers - start - c.pprofServer.ListenAndServe")
		}
	}()

	logger.Info().Str("address", c.cfg.Address).Msg("Starting server")
	if c.cfg.CryptoKey != "" {
		return fmt.Errorf("server error %w, or server just stoped", c.server.ListenAndServeTLS(c.cfg.CryptoKey, c.cfg.CryptoKey))
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
