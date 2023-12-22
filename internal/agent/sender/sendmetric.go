package sender

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
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
	logger.Debug().Str("metric", string(rd)).Str("endpoint", j.endpoint).Msg("trying to send metric")

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

	// Получаем хеш
	hsHex := ""
	if j.hashKey != "" {
		h := hmac.New(sha256.New, []byte(j.hashKey))
		h.Write(rdGz)
		hs := h.Sum(nil)
		hsHex = hex.EncodeToString(hs)
	}
	// Создаем запрос
	cl := resty.New()
	var resp *resty.Response
	req := cl.R().
		SetBody(rdGz).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	if hsHex != "" {
		req.SetHeader("HashSHA256", hsHex)
	}

	for _, ri := range j.retryIntervals {
		resp, err = req.Post(j.endpoint)
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
		logger.Error().Str("status code", resp.Status()).Str("endpoint", j.endpoint).Msg("status code is")
		return
	}
}
