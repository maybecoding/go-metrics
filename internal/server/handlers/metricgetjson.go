package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

func (c *Handler) metricGetJSON(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response []byte
	)
	status := http.StatusOK
	logMessage := ""

	defer func() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if logMessage != "" {
			logger.Log.Debug().Err(err).Msg(logMessage)
		}
		_, _ = w.Write(response)
	}()

	// Получаем JSON и парсим
	decoder := json.NewDecoder(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()

	m := metric.Metrics{}
	err = decoder.Decode(&m)
	if err != nil {
		status, logMessage = http.StatusBadRequest, "failed to parse request JSON"
		return
	}

	err = c.metric.Get(&m)

	if err != nil && errors.Is(err, metric.ErrNoMetricValue) {
		status, logMessage = http.StatusNotFound, "metric isn't found"
		return
	}
	if err != nil {
		status, logMessage = http.StatusBadRequest, "metric can't be get"
		return
	}

	response, err = json.Marshal(m)
	if err != nil {
		status, logMessage = http.StatusInternalServerError, "can't prepare response"
		return
	}

	status = http.StatusOK

}
