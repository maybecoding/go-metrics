package httpjsonsender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/logger"
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
	payload, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error due marshal metric before send")
		return
	}
	req, err := http.NewRequest("POST", j.endpoint, bytes.NewReader(payload))
	if err != nil {
		logger.Log.Error().Err(err).Msg("can't create request")
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("can't do request")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		logger.Log.Error().Int("status code", resp.StatusCode).Msg("status code is")
		return
	}
}

func New(template, serverAddress string) *HTTPJSONSender {
	return &HTTPJSONSender{
		endpoint: fmt.Sprintf(template, serverAddress),
	}
}
