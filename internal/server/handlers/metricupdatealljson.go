package handlers

import (
	"encoding/json"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
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

	var mts []metricservice.Metrics
	err := decoder.Decode(&mts)
	if err != nil {
		status = http.StatusBadRequest
		logger.Debug().Err(err).Msg("error due decode request")
		return
	}

	err = c.metric.SetAll(mts)
	if err != nil {
		status = http.StatusInternalServerError
		logger.Debug().Err(err).Msg("error due set all")
		return
	}
}
