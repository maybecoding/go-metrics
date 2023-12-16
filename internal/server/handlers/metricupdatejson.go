package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

func (c *Handler) metricUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response []byte
	)
	status := http.StatusOK
	logMessage := ""
	defer func() {

		if err != nil {
			logger.Debug().Err(err).Msg(logMessage)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(response)

	}()

	// Получаем JSON и парсим
	decoder := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	m := metric.Metrics{}
	err = decoder.Decode(&m)
	if err != nil {
		status, logMessage = http.StatusBadRequest, "failed to parse request JSON"
		return
	}

	if m.Delta != nil {
		logger.Debug().Str("ID", m.ID).Int64("Delta", *m.Delta).Msg("UpdateMetric JSON")
	}
	if m.Value != nil {
		logger.Debug().Str("ID", m.ID).Float64("Value", *m.Value).Msg("UpdateMetric JSON")
	}

	if err := c.metric.Set(&m); err != nil {
		status, logMessage = http.StatusBadRequest, "failed to update metric in metric"
		return
	}

	// Не будем доставать из хранилища и отдавать результат, отдадим полученную структуру
	response, err = json.Marshal(m)
	if err != nil {
		status, logMessage = http.StatusBadRequest, "failed to marshal json"
		return
	}

	status = http.StatusOK
}
