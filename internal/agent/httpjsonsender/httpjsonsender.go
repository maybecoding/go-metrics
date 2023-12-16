package httpjsonsender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
)

type HTTPJSONSender struct {
	endpoint string
}

func (j *HTTPJSONSender) Send(metrics []*app.Metrics) {
	for _, metric := range metrics {
		j.sendMetric(metric)
	}
}

func (j *HTTPJSONSender) sendMetric(metric *app.Metrics) {
	// Получаем json для отправки
	payload, err := json.Marshal(metric)
	if err != nil {
		logger.Error().Err(err).Msg("error due marshal metric before send")
		return
	}
	logger.Debug().Str("json", string(payload)).Msg("trying to send json")

	// Создаем сжатый поток
	buf := bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)

	// И записываем в него данные
	_, err = zw.Write(payload)
	if err != nil {
		logger.Error().Err(err).Msg("can't write into gzip writer")
		return
	}
	err = zw.Close()
	if err != nil {
		logger.Error().Err(err).Msg("can't close gzip writer")
		return
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", j.endpoint, buf)
	if err != nil {
		logger.Error().Err(err).Msg("can't create request")
		return
	}
	// Устанавливаем заголовок
	req.Header.Set("Content-Type", "application/json")
	// Не забываем указать что это сжатые данные
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("can't do request")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		logger.Error().Int("status code", resp.StatusCode).Str("endpoint", j.endpoint).Msg("status code is")
		return
	}
}

func New(template, serverAddress string) *HTTPJSONSender {
	return &HTTPJSONSender{
		endpoint: fmt.Sprintf(template, serverAddress),
	}
}
