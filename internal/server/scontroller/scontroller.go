package scontroller

import (
	"fmt"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/sapp"
	"net/http"
)

type Controller struct {
	app           *sapp.App
	serverAddress string
}

func New(app *sapp.App, serverAddress string) *Controller {
	return &Controller{app: app, serverAddress: serverAddress}
}

func (c *Controller) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", c.handleUpdate)
	mux.Handle("/", http.NotFoundHandler())

	addr := c.serverAddress
	err := http.ListenAndServe(addr, mux)

	if err != nil {
		panic(fmt.Errorf("server on %v failed: %v", addr, err))
	}
}
