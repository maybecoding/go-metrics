package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/app"
	"net/http"
)

func (c *Handler) metricGet(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	value, err := c.app.GetMetric(mType, name)
	if err != nil {
		if errors.Is(err, app.ErrNoMetricValue) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(value))
}
