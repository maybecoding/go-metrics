package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (c *Handler) metricUpdate(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")
	value := chi.URLParam(r, "value")

	if err := c.app.UpdateMetric(mType, name, value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
