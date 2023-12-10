package httpjsonbatchsender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
	"strings"
	"time"
)

type HTTPJSONBatchSender struct {
	endpoint       string
	retryIntervals []time.Duration
}

func (j *HTTPJSONBatchSender) Send(metrics []*app.Metrics) {
	// Получаем json для отправки
	payload, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error due marshal all metrics before send")
		return
	}
	logger.Log.Debug().Str("json", string(payload)).Msg("trying to send json")

	// Создаем сжатый поток
	buf := bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)

	// И записываем в него данные
	_, err = zw.Write(payload)
	if err != nil {
		logger.Log.Error().Err(err).Msg("can't write into gzip writer")
		return
	}
	err = zw.Close()
	if err != nil {
		logger.Log.Error().Err(err).Msg("can't close gzip writer")
		return
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", j.endpoint, buf)
	if err != nil {
		logger.Log.Error().Err(err).Msg("can't create request")
		return
	}
	// Устанавливаем заголовок
	req.Header.Set("Content-Type", "application/json")
	// Не забываем указать что это сжатые данные
	req.Header.Set("Content-Encoding", "gzip")

	var resp *http.Response
	for _, ri := range j.retryIntervals {
		resp, err = http.DefaultClient.Do(req)
		// TODO найти более подходящий способ понять, что проблема с соединением
		if err == nil || !strings.Contains(err.Error(), "connect") {
			break
		}
		logger.Log.Debug().Err(err).Dur("duration", ri).Msg("connection error, trying after")
		time.Sleep(ri)
	}

	if err != nil {
		logger.Log.Error().Err(err).Msg("can't do request")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		logger.Log.Error().Int("status code", resp.StatusCode).Str("endpoint", j.endpoint).Msg("status code is")
		return
	}
}

func New(template, serverAddress string, retryIntervals []time.Duration) *HTTPJSONBatchSender {
	return &HTTPJSONBatchSender{
		endpoint:       fmt.Sprintf(template, serverAddress),
		retryIntervals: retryIntervals,
	}
}
