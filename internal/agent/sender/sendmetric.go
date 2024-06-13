package sender

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/hasher"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/maybecoding/go-metrics.git/pkg/zipper"
	"net/http"
	"os"
	"strings"
	"time"
)

func (j *Sender) sendMetric(_ context.Context, mt *app.Metrics) error {
	// Получаем json для отправки
	rd, err := json.Marshal(mt)
	if err != nil {
		return fmt.Errorf("sender - sendMetric - json.Marshal: %w", err)
	}
	protocol := "http"
	if j.cfg.CryptoKey != "" {
		protocol += "s"
	}
	endpoint := fmt.Sprintf(j.cfg.EndpointTemplate, protocol, j.cfg.Server)
	logger.Debug().Str("endpoint", endpoint).Msg("Endpoint")
	logger.Debug().Str("metric", string(rd)).Str("endpoint", endpoint).Msg("trying to send metric")

	// Создаем сжатый поток
	rdGz, err := zipper.ZippedBytes(rd)
	if err != nil {
		return fmt.Errorf("sender - sendMetric - zipper.ZippedBytes: %w", err)
	}
	// Создаем запрос
	cl := resty.New()
	if j.cfg.CryptoKey != "" {
		var crt []byte
		crt, err = os.ReadFile(j.cfg.CryptoKey)
		if err != nil {
			return fmt.Errorf("sender - sendMetric - error due read certificate: %w", err)
		}
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(crt)
		if !ok {
			return errors.New("error due append cert into roots")
		}
		cl.SetTLSClientConfig(&tls.Config{
			RootCAs: roots,
		})
	}
	cl.OnBeforeRequest(hasher.New(j.cfg.HashKey, sha256.New))
	var resp *resty.Response
	req := cl.R().
		SetBody(rdGz).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip")

	if j.ip != nil && j.cfg.IPAddrHeader != "" {
		req.SetHeader(j.cfg.IPAddrHeader, j.ip.String())
	}

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
		return fmt.Errorf("sender - sendMetric - can't do request: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("sender - sendMetric - response is nil: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("sender - sendMetric - endpoint(%s) status code is %s", endpoint, resp.Status())
	}
	return nil
}
