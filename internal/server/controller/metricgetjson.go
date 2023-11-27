package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/model"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
	"strconv"
)

func (c *Controller) metricGetJSON(w http.ResponseWriter, r *http.Request) {
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

	m := model.Metrics{}
	err = decoder.Decode(&m)
	if err != nil {
		status, logMessage = http.StatusBadRequest, "failed to parse request JSON"
		return
	}
	if m.ID == "" || m.MType == "" {
		status, logMessage = http.StatusBadRequest, "isn't set metric name or metric type"
		return
	}

	value, err := c.app.GetMetric(m.MType, m.ID)
	logger.Log.Debug().Str("type", m.MType).Str("name", m.ID).Str("value", value).Msg("metric get json")

	if err != nil && errors.Is(err, app.ErrNoMetricValue) {
		status, logMessage = http.StatusNotFound, "metric isn't found"
		return
	}
	if err != nil {
		status, logMessage = http.StatusBadRequest, "metric can't be get"
		return
	}

	// Отправляем ответ
	// Сохраняем в структуру полученное значение
	if m.MType == "gauge" {
		gValue, err := strconv.ParseFloat(value, 64)
		m.Value = &gValue
		if err != nil {
			status, logMessage = http.StatusInternalServerError, "can't identify value due response prepare"
			return
		}
	} else if m.MType == "counter" {
		cValue, err := strconv.ParseInt(value, 10, 64)
		m.Delta = &cValue
		if err != nil {
			status, logMessage = http.StatusInternalServerError, "can't identify value due response prepare"
			return
		}
	} else {
		err = fmt.Errorf("get incorrect metric type from app")
		if err != nil {
			status, logMessage = http.StatusInternalServerError, "can't identify value due response prepare"
			return
		}
	}

	response, err = json.Marshal(m)
	if err != nil {
		status, logMessage = http.StatusInternalServerError, "can't prepare response"
		return
	}

	status = http.StatusOK

}
