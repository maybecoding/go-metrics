package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	"github.com/maybecoding/go-metrics.git/internal/server/handlers"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

func Example_metricUpdateAllJSON() {
	// get handler
	contr := initAppGetHandler()
	// Prepare test server for testing handler
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	// Подготавливаем данные к отправке по каждой метрики по metricCnt
	const metricCnt = 2
	metrics := make([]entity.Metrics, 0, metricCnt*2)
	for i := 0; i <= metricCnt; i += 1 {
		// Метрика gauge
		value := float64(i)
		metrics = append(metrics, entity.Metrics{
			ID:    fmt.Sprintf("gauge_%d", i),
			MType: sapp.Gauge,
			Value: &value,
		})
		// Метрика counter
		delta := int64(i)
		metrics = append(metrics, entity.Metrics{
			ID:    fmt.Sprintf("counter_%d", i),
			MType: sapp.Counter,
			Delta: &delta,
		})
	}

	// Sending batch of messages
	fmt.Println("send batch of messages /updates/")
	// Prepare json with all messages for send
	jsonM, err := json.Marshal(metrics)
	fmt.Println("err1", err)
	// Prepare buffer for send
	buf := bytes.NewBuffer(jsonM)
	// Prepare new request with method, endpoint and data
	reqSendAll, err := http.NewRequest("POST", ts.URL+"/updates/", buf)
	fmt.Println("err2", err)
	// Set header
	reqSendAll.Header.Set("Content-Type", "application/json")
	// Send using test client
	respSendAll, err := ts.Client().Do(reqSendAll)
	fmt.Println("err3", err)
	_ = respSendAll.Body.Close()
	// Check status code
	fmt.Println("status", respSendAll.StatusCode)
	fmt.Println("-------------------------------")

	// Output:
	//send batch of messages /updates/
	//err1 <nil>
	//err2 <nil>
	//err3 <nil>
	//status 200
	//-------------------------------
}
func Example_metricUpdate() {
	// get handler
	contr := initAppGetHandler()
	// Prepare test server for testing handler
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	// Sending single gauge message
	fmt.Println("send single gauge message")
	// Prepare new request with method, endpoint
	reqSendG, err := http.NewRequest("POST", ts.URL+"/update/gauge/single/3232.23", nil)
	fmt.Println("err1", err)
	// Send using test client
	respSendG, err := ts.Client().Do(reqSendG)
	fmt.Println("err2", err)
	_ = respSendG.Body.Close()
	// Check status code
	fmt.Println("status", respSendG.StatusCode)
	fmt.Println("-------------------------------")

	// Output:
	//send single gauge message
	//err1 <nil>
	//err2 <nil>
	//status 200
	//-------------------------------
}

func Example_metricGetJSON() {
	// get handler
	contr := initAppGetHandler()
	// Prepare test server for testing handler
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	const metricCnt = 2
	// Show messages using post
	fmt.Println("show messages using post")
	// define metric types
	metricTypes := []string{"gauge", "counter"}
	for i := 0; i <= metricCnt; i += 1 { // metricCnt times
		for _, tp := range metricTypes { // for every metric type
			fmt.Println("#", i, tp)
			// prepare metric
			m := entity.Metrics{MType: tp, ID: fmt.Sprintf("%s_%d", tp, i)}
			// Prepare json with all messages for send
			bts, err := json.Marshal(m)
			fmt.Println("err1", err)
			// Prepare new request with method, endpoint and data
			req, err := http.NewRequest("POST", ts.URL+"/value/", bytes.NewBuffer(bts))
			fmt.Println("err2", err)
			// Send using test client
			resp, err := ts.Client().Do(req)
			fmt.Println("err3", err)
			// Check body
			data, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			fmt.Println("err4", err)
			if tp == "gauge" {
				fmt.Println("result", string(data))
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
	// get handler
	contr := initAppGetHandler()
	// Prepare test server for testing handler
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	// Show messages using get
	metricCnt := 2
	fmt.Println("show messages using get")
	// define metric types
	metricTypes := []string{"gauge", "counter"}
	for i := 0; i <= metricCnt; i += 1 { // metricCnt times
		for _, tp := range metricTypes {
			fmt.Println("#", i, tp)
			// Prepare new request with method, endpoint and data
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/value/%s/%s_%d", ts.URL, tp, tp, i), nil)
			fmt.Println("err1", err)
			// Send using test client
			resp, err := ts.Client().Do(req)
			fmt.Println("err2", err)
			// check body
			data, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			fmt.Println("err3", err)
			if tp == "gauge" {
				fmt.Println("result", data)
			}
			fmt.Println("-------------------------------")
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

func initAppGetHandler() *handlers.Handler {
	cfg := &config.Config{
		Log: config.Log{Level: "info"},
		Database: config.Database{
			ConnStr:        "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable",
			RetryIntervals: []time.Duration{time.Second},
		},
	}
	// Create app
	a := app.New(cfg).
		Init().       // Init app
		InitHandler() // Init handler
	// get handler
	return a.GetHandler()
}
