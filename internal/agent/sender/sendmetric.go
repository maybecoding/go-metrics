package sender

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/hasher"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
	"strings"
	"time"
)

func (j *Sender) sendMetric(mt *app.Metrics) {
	// Получаем json для отправки
	rd, err := json.Marshal(mt)
	if err != nil {
		logger.Error().Err(err).Msg("error due marshal all metrics before send")
		return
	}
	endpoint := fmt.Sprintf(j.cfg.EndpointTemplate, j.cfg.Server)
	logger.Debug().Str("endpoint", endpoint).Msg("Endpoint")
	logger.Debug().Str("metric", string(rd)).Str("endpoint", endpoint).Msg("trying to send metric")

	// Создаем сжатый поток
	buf := bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)

	// И записываем в него данные
	_, err = zw.Write(rd)
	if err != nil {
		logger.Error().Err(err).Msg("can't write into gzip writer")
		return
	}
	err = zw.Close()
	if err != nil {
		logger.Error().Err(err).Msg("can't close gzip writer")
		return
	}
	rdGz := buf.Bytes()

	// Создаем запрос
	cl := resty.New()
	cl.OnBeforeRequest(hasher.New(j.cfg.HashKey, sha256.New))
	var resp *resty.Response
	req := cl.R().
		SetBody(rdGz).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	for _, ri := range j.cfg.RetryIntervals {
		resp, err = req.Post(endpoint)
		if err != nil {
			logger.Error().Err(err).Msg("error due request")
		}
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

	if resp == nil {
		logger.Error().Msg("response is nil")
		return
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Error().Str("status code", resp.Status()).Str("endpoint", endpoint).Msg("status code is")
		return
	}
}
