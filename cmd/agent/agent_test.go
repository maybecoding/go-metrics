package main

import (
	"context"
	"encoding/json"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/collector"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"github.com/maybecoding/go-metrics.git/internal/agent/sender"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/compress"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type ServerRequest struct {
	Body   string
	Method string
	Err    error
}

func TestAgent(t *testing.T) {
	// Канал, по которому тестовый сервер вернет ответ
	chSrvReq := make(chan *ServerRequest, 30)
	ctx, cancel := context.WithCancel(context.Background())

	// Создаем тестовый сервер
	testHandler := compress.HandlerFuncReader(func(w http.ResponseWriter, r *http.Request) {
		srvReq := &ServerRequest{
			Method: r.Method,
		}
		defer func() {
			w.WriteHeader(http.StatusOK)
			chSrvReq <- srvReq
		}()
		defer func() {
			_ = r.Body.Close()
		}()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			srvReq.Err = err
			return
		}
		srvReq.Body = string(body)
	})

	ts := httptest.NewServer(testHandler)

	// Основная логика агента
	logger.Init("debug")
	logger.Debug().Str("test server endpoint", ts.URL).Msg("test server started")

	var collect = collector.New(ctx)
	cfgSender := config.Sender{
		EndpointTemplate: "%s",
		Server:           ts.URL,
		HashKey:          "",
		RetryIntervals:   []time.Duration{2000 * time.Millisecond, 2500 * time.Millisecond},
		NumWorkers:       2,
	}
	var snd app.Sender = sender.New(ctx, cfgSender)

	a := app.New(collect, snd, time.Duration(1000)*time.Millisecond, time.Duration(1600)*time.Millisecond)

	var wg sync.WaitGroup

	// Через 3 секунды отправим отмену
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(3500 * time.Millisecond)
		cancel()
	}()

	// Сама проверка
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-chSrvReq:
				logger.Debug().Interface("request", req).Msg("Message from server")
				var mReq app.Metrics
				assert.NoError(t, req.Err, "request must be without error")
				err := json.Unmarshal([]byte(req.Body), &mReq)
				assert.NoError(t, err, "request body must be a valid metric")
				mts := collect.GetMetrics()
				var mFound *app.Metrics
				for _, mt := range mts {
					if mt.ID == mReq.ID && mt.MType == mReq.MType {
						mFound = mt
						break
					}
				}
				assert.NotNil(t, mFound, "sent metric must be in collector")
				if mFound.MType == metricservice.Gauge {
					assert.NotNil(t, mFound.Value)
					assert.NotNil(t, mReq.Value)
					assert.Equal(t, *mReq.Value, *mFound.Value)
				} else {
					assert.NotNil(t, mFound.Delta)
					assert.NotNil(t, mReq.Delta)
					assert.Equal(t, *mReq.Delta, *mFound.Delta)
				}
			}
		}

	}()

	a.Run()
	wg.Wait()
	//time.Sleep(5 * time.Second)
}
