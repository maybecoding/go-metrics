package controller

import (
	"github.com/go-chi/chi/v5"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	compress2 "github.com/maybecoding/go-metrics.git/pkg/compress"
	logger2 "github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
)

type Controller struct {
	app           *sapp.App
	serverAddress string
}

func New(app *sapp.App, serverAddress string) *Controller {
	return &Controller{app: app, serverAddress: serverAddress}
}

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()
	// подключаем логер
	r.Use(logger2.Handler)

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

func (c *Controller) Start() {

	addr := c.serverAddress
	logger2.Log.Info().Str("address", addr).Msg("Starting server")
	err := http.ListenAndServe(addr, c.GetRouter())

	if err != nil {
		logger2.Log.Fatal().Err(err).Msg("Failed to start server")
	}
}
