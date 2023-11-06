package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
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
	err := http.ListenAndServe(addr, c.GetRouter())

	if err != nil {
		panic(fmt.Errorf("server on %v failed: %v", addr, err))
	}
}
