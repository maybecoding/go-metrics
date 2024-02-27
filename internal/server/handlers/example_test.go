package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/dbstorage"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/metric"
	"github.com/maybecoding/go-metrics.git/pkg/health"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

func Example_metricUpdateAllJSON() {
	logger.Init("info")

	// Если задана база данных
	var store sapp.Store
	var app *sapp.Metric

	ctx := context.Background()

	// HealthCheck
	hl := health.New()

	dbStore := dbstorage.New("postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", ctx, []time.Duration{time.Second})
	store = dbStore
	defer dbStore.ConnectionClose()

	app = sapp.New(store)

	// Создаем контроллер и вверяем ему приложение
	//contr := handlers.New(app, "", hl, hashKey)
	contr := New(app, "", hl, "")
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	// Подготавливаем данные к отправке
	const metricCnt = 2
	metrics := make([]sapp.Metrics, 0, metricCnt*2)

	for i := 0; i <= metricCnt; i += 1 {
		value := float64(i)
		metrics = append(metrics, sapp.Metrics{
			ID:    fmt.Sprintf("gauge_%d", i),
			MType: sapp.Gauge,
			Value: &value,
		})

		delta := int64(i)
		metrics = append(metrics, sapp.Metrics{
			ID:    fmt.Sprintf("counter_%d", i),
			MType: sapp.Counter,
			Delta: &delta,
		})
	}

	// Sending batch of messages
	{
		fmt.Println("send batch of messages /updates/")
		jsonM, err := json.Marshal(metrics)
		fmt.Println("err1", err)

		buf := bytes.NewBuffer(jsonM)
		reqSendAll, err := http.NewRequest("POST", ts.URL+"/updates/", buf)
		fmt.Println("err2", err)

		reqSendAll.Header.Set("Content-Type", "application/json")

		respSendAll, err := ts.Client().Do(reqSendAll)
		fmt.Println("err3", err)
		_ = respSendAll.Body.Close()
		fmt.Println("status", respSendAll.StatusCode)
		fmt.Println("-------------------------------")
	}

	// Output:
	//send batch of messages /updates/
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//status 200
	//-------------------------------
}
func Example_metricUpdate() {
	logger.Init("info")

	// Если задана база данных
	var store sapp.Store
	var app *sapp.Metric

	ctx := context.Background()

	// HealthCheck
	hl := health.New()

	dbStore := dbstorage.New("postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", ctx, []time.Duration{time.Second})
	store = dbStore
	defer dbStore.ConnectionClose()

	app = sapp.New(store)

	// Создаем контроллер и вверяем ему приложение
	//contr := handlers.New(app, "", hl, hashKey)
	contr := New(app, "", hl, "")
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	// Подготавливаем данные к отправке

	// Sending single gauge message
	{
		fmt.Println("send single gauge message /update/gauge/single/3232.23")
		reqSendG, err := http.NewRequest("POST", ts.URL+"/update/gauge/single/3232.23", nil)
		fmt.Println("err1", err)
		respSendG, err := ts.Client().Do(reqSendG)
		fmt.Println("err2", err)
		_ = respSendG.Body.Close()
		fmt.Println("status", respSendG.StatusCode)
		fmt.Println("-------------------------------")
	}

	// Output:
	//send single gauge message /update/gauge/single/3232.23
	//err1 <nil>
	//err2 <nil>
	//status 200
	//-------------------------------
}

func Example_metricGetJSON() {
	logger.Init("info")

	// Если задана база данных
	var store sapp.Store
	var app *sapp.Metric

	ctx := context.Background()

	// HealthCheck
	hl := health.New()

	dbStore := dbstorage.New("postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", ctx, []time.Duration{time.Second})
	store = dbStore
	defer dbStore.ConnectionClose()

	app = sapp.New(store)

	// Создаем контроллер и вверяем ему приложение
	//contr := handlers.New(app, "", hl, hashKey)
	contr := New(app, "", hl, "")
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	// Подготавливаем данные к отправке
	const metricCnt = 2
	// Show messages using post
	{
		fmt.Println("show messages using post")
		metricTypes := []string{"gauge", "counter"}
		for i := 0; i <= metricCnt; i += 1 {
			for _, tp := range metricTypes {
				fmt.Println("#", i, tp)
				m := sapp.Metrics{MType: tp, ID: fmt.Sprintf("%s_%d", tp, i)}
				bts, err := json.Marshal(m)
				fmt.Println("err1", err)
				req, err := http.NewRequest("POST", ts.URL+"/value/", bytes.NewBuffer(bts))
				fmt.Println("err2", err)
				resp, err := ts.Client().Do(req)
				fmt.Println("err3", err)
				data, err := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				fmt.Println("err4", err)
				if tp == "gauge" {
					fmt.Println("result", string(data))
				}
			}
		}
	}
	//Output:
	//show messages using post
	//# 0 gauge
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//err4 <nil>
	//result {"id":"gauge_0","type":"gauge","value":0}
	//# 0 counter
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//err4 <nil>
	//# 1 gauge
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//err4 <nil>
	//result {"id":"gauge_1","type":"gauge","value":1}
	//# 1 counter
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//err4 <nil>
	//# 2 gauge
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//err4 <nil>
	//result {"id":"gauge_2","type":"gauge","value":2}
	//# 2 counter
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//err4 <nil>
}

func Example_metricGet() {
	logger.Init("info")

	// Если задана база данных
	var store sapp.Store
	var app *sapp.Metric

	ctx := context.Background()

	// HealthCheck
	hl := health.New()

	dbStore := dbstorage.New("postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", ctx, []time.Duration{time.Second})
	store = dbStore
	defer dbStore.ConnectionClose()

	app = sapp.New(store)

	// Создаем контроллер и вверяем ему приложение
	//contr := handlers.New(app, "", hl, hashKey)
	contr := New(app, "", hl, "")
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	// Show messages using get
	metricCnt := 2
	{
		fmt.Println("show messages using get")
		metricTypes := []string{"gauge", "counter"}
		for i := 0; i <= metricCnt; i += 1 {
			for _, tp := range metricTypes {
				fmt.Println("#", i, tp)
				req, err := http.NewRequest("GET", fmt.Sprintf("%s/value/%s/%s_%d", ts.URL, tp, tp, i), nil)
				fmt.Println("err1", err)
				resp, err := ts.Client().Do(req)
				fmt.Println("err2", err)
				data, err := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				fmt.Println("err3", err)
				if tp == "gauge" {
					fmt.Println("result", data)
				}
				fmt.Println("-------------------------------")
			}
		}
	}
	//Output:
	//show messages using get
	//# 0 gauge
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//result [48]
	//-------------------------------
	//# 0 counter
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//-------------------------------
	//# 1 gauge
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//result [49]
	//-------------------------------
	//# 1 counter
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//-------------------------------
	//# 2 gauge
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//result [50]
	//-------------------------------
	//# 2 counter
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//-------------------------------
}
