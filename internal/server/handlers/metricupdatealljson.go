package handlers

import (
	"encoding/json"
	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"net/http"
)

func (c *Handler) metricUpdateAllJSON(w http.ResponseWriter, r *http.Request) {

	status := http.StatusOK
	defer func() {
		w.WriteHeader(status)
	}()

	decoder := json.NewDecoder(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	var mts []*metric.Metrics
	err := decoder.Decode(&mts)
	if err != nil {
		status = http.StatusBadRequest
		return
	}

	err = c.metric.SetAll(mts)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}
}
