// Package health - Package exposes struct for health checks publication
package health

type Health struct {
	checks []func() error
}

// Check - Checks if all health checks is succeed
func (h *Health) Check() bool {
	var err error
	for _, c := range h.checks {
		err = c()
		if err != nil {
			return false
		}
	}
	return true
}

// Watch - add health check publication
func (h *Health) Watch(hf func() error) {
	h.checks = append(h.checks, hf)
}

// New - creates new struct for health checks publication
func New() *Health {
	return &Health{
		checks: make([]func() error, 0),
	}
}
