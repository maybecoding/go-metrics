package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"strings"
	"time"
)

type Sender struct {
	endpoint       string
	retryIntervals []time.Duration
}

func (j *Sender) Send(metrics []*app.Metrics) {
	// Получаем json для отправки
	payload, err := json.Marshal(metrics)
	if err != nil {
		logger.Error().Err(err).Msg("error due marshal all metrics before send")
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
	cl := resty.New()
	var resp *resty.Response
	for _, ri := range j.retryIntervals {
		resp, err = cl.R().
			SetBody(buf).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			Post(j.endpoint)

		if err == nil || !strings.Contains(err.Error(), "connect") {
			break
		}
		logger.Debug().Err(err).Dur("duration", ri).Msg("connection error, trying after")
		time.Sleep(ri)
	}

	if err != nil {
		logger.Error().Err(err).Msg("can't do request")
		return
	}

	if resp.Status() != "200" {
		logger.Error().Str("status code", resp.Status()).Str("endpoint", j.endpoint).Msg("status code is")
		return
	}
}

func New(template, serverAddress string, retryIntervals []time.Duration) *Sender {
	return &Sender{
		endpoint:       fmt.Sprintf(template, serverAddress),
		retryIntervals: retryIntervals,
	}
}
