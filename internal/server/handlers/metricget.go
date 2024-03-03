package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
)

func (c *Handler) metricGet(w http.ResponseWriter, r *http.Request) {

	m := &metricservice.Metrics{
		ID:    chi.URLParam(r, "name"),
		MType: chi.URLParam(r, "type"),
	}
	err := c.metric.Get(m)
	if err != nil {
		if errors.Is(err, metricservice.ErrNoMetricValue) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	res := ""
	if m.MType == metricservice.Gauge {
		res = strconv.FormatFloat(*m.Value, FmtFloat, -1, 64)
	} else if m.MType == metricservice.Counter {
		res = strconv.FormatInt(*m.Delta, 10)
	}
	_, _ = w.Write([]byte(res))

}
