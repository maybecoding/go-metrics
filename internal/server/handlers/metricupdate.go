package handlers

import (
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
)

func (c *Handler) metricUpdate(w http.ResponseWriter, r *http.Request) {

	value := chi.URLParam(r, "value")

	m := metric.Metrics{
		ID:    chi.URLParam(r, "name"),
		MType: chi.URLParam(r, "type"),
	}

	logger.Debug().Str("ID", m.ID).Str("MType", m.MType).Str("value", value).Msg("UpdateMetric URL")

	if m.MType == metric.Gauge {
		gValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.Value = &gValue
	} else if m.MType == metric.Counter {
		cValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.Delta = &cValue
	}

	if err := c.metric.Set(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
