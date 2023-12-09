package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
)

func (c *Handler) metricGet(w http.ResponseWriter, r *http.Request) {

	m := &metric.Metrics{
		ID:    chi.URLParam(r, "name"),
		MType: chi.URLParam(r, "type"),
	}
	err := c.metric.Get(m)
	if err != nil {
		if errors.Is(err, metric.ErrNoMetricValue) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	res := ""
	if m.MType == metric.Gauge {
		res = strconv.FormatFloat(*m.Value, FmtFloat, -1, 64)
	} else if m.MType == metric.Counter {
		res = strconv.FormatInt(*m.Delta, 10)
	}
	_, _ = w.Write([]byte(res))

}
