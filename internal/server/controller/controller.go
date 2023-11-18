package controller

import (
	"github.com/go-chi/chi/v5"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/logger"
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
	r.Post("/update/{type}/{name}/{value}", c.metricUpdate)
	r.Get("/value/{type}/{name}", c.metricGet)
	r.Get("/", c.metricGetAll)
	return r
}

func (c *Controller) Start() {

	addr := c.serverAddress
	logger.Log.Info().Str("address", addr).Msg("Starting server")
	err := http.ListenAndServe(addr, logger.Handler(c.GetRouter()))

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start server")
	}
}
