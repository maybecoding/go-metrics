package controller

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	compress2 "github.com/maybecoding/go-metrics.git/pkg/compress"
	logger "github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
)

type Controller struct {
	app           *sapp.App
	serverAddress string
	server        *http.Server
}

func New(app *sapp.App, serverAddress string) *Controller {
	return &Controller{app: app, serverAddress: serverAddress}
}

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()
	// подключаем логер
	r.Use(logger.Handler)

	// Установка значениий
	r.Post("/update/{type}/{name}/{value}", c.metricUpdate)
	r.Post("/update/", compress2.HandlerFuncReader(compress2.HandlerFuncWriter(c.metricUpdateJSON, compress2.BestSpeed)))

	// Получение значениий
	r.Get("/value/{type}/{name}", c.metricGet)
	r.Post("/value/", compress2.HandlerFuncReader(compress2.HandlerFuncWriter(c.metricGetJSON, compress2.BestSpeed)))

	// Отчет по метрикам
	r.Get("/", compress2.HandlerFuncWriter(c.metricGetAll, compress2.BestSpeed))

	return r
}

func (c *Controller) Start() error {
	addr := c.serverAddress
	c.server = &http.Server{Addr: addr, Handler: c.GetRouter()}

	logger.Log.Info().Str("address", addr).Msg("Starting server")
	return fmt.Errorf("server error %w, or server just stoped", c.server.ListenAndServe())
}

func (c *Controller) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	logger.Log.Info().Msg("Ctrl + C command got, shutting down server")
	return c.server.Shutdown(context.Background())
}
