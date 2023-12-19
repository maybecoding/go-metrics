package health

import "github.com/maybecoding/go-metrics.git/pkg/logger"

type Health struct {
	checks []func() error
}

func (h *Health) Check() bool {
	var err error
	for _, c := range h.checks {
		err = c()
		if err != nil {
			logger.Debug().Err(err).Msg("heath checker identified unhealthy service")
			return false
		}
	}
	return true
}

func (h *Health) Watch(hf func() error) {
	h.checks = append(h.checks, hf)
}

func New() *Health {
	return &Health{
		checks: make([]func() error, 0),
	}
}
