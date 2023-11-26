package controller

import (
	"encoding/json"
	"github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/logger"
	"github.com/maybecoding/go-metrics.git/internal/server/model"
	"net/http"
	"strconv"
)

func (c *Controller) metricUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response []byte
	)
	status := http.StatusOK
	logMessage := ""
	defer func() {

		if err != nil {
			logger.Log.Debug().Err(err).Msg(logMessage)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(response)

	}()

	// Получаем JSON и парсим
	decoder := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	m := model.Metrics{}
	err = decoder.Decode(&m)
	if err != nil {
		status, logMessage = http.StatusBadRequest, "failed to parse request JSON"
		return
	}

	var value string
	if m.Delta != nil {
		value = strconv.FormatInt(*m.Delta, 10)
	} else if m.Value != nil {
		value = strconv.FormatFloat(*m.Value, app.FmtFloat, -1, 64)
	} else {
		status, logMessage = http.StatusBadRequest, "metric value isn't provided"
		return
	}
	if err := c.app.UpdateMetric(m.MType, m.ID, value); err != nil {
		status, logMessage = http.StatusBadRequest, "failed to update metric in app"
		return
	}
	mOut := model.Metrics{
		MType: m.MType,
		ID:    m.ID,
	}
	procValStr, err := c.app.GetMetric(m.MType, m.ID)
	if err != nil {
		status, logMessage = http.StatusBadRequest, "failed to get metric value after update"
		return
	}
	if m.MType == app.Gauge {
		procValGauge, err := strconv.ParseFloat(procValStr, 64)
		if err != nil {
			status, logMessage = http.StatusBadRequest, "failed to get gauge from value after update"
			return
		}
		mOut.Value = &procValGauge
	} else if m.MType == app.Counter {
		procValCounter, err := strconv.ParseInt(procValStr, 10, 64)
		if err != nil {
			status, logMessage = http.StatusBadRequest, "failed to get gauge from value after update"
			return
		}
		mOut.Delta = &procValCounter
	}

	// Не будем доставать из хранилища и отдавать результат, отдадим полученную структуру
	response, err = json.Marshal(mOut)
	if err != nil {
		status, logMessage = http.StatusBadRequest, "failed to marshal json"
		return
	}

	status = http.StatusOK
}
