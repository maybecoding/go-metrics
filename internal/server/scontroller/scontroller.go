package controller

import "github.com/maybecoding/go-metrics.git/internal/server/core"

type Controller struct {
	core core.Core
}

func New(core core.Core, serverAddress string) Controller {

}
