package handlers

import (
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/pkg/health"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maybecoding/go-metrics.git/internal/server/metricmemstorage"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
	getResult string
	code      int
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestController(t *testing.T) {

	logger.Init("debug")
	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{name: "#1 Unknown metric", url: "/update/gauge2/Test/1", method: "POST", want: want{code: 400}},

		{name: "#2 Gauge success", url: "/update/gauge/Test/913372.185", method: "POST", want: want{code: 200}},
		{name: "#2. Gauge success", url: "/update/gauge/Test2/330095.942", method: "POST", want: want{code: 200}},
		{name: "#3 Counter success", url: "/update/counter/TestTestTestTestTest/312323", method: "POST", want: want{code: 200}},

		{name: "#4 Gauge bad value", url: "/update/gauge/Test/7.77x", method: "POST", want: want{code: 400}},
		{name: "#5 Counter bad value", url: "/update/counter/TestTestTestTestTest/1.1", method: "POST", want: want{code: 400}},

		{name: "#6 Gauge get", url: "/value/gauge/Test", method: "GET", want: want{code: 200, getResult: "913372.185"}},
		{name: "#6.1 Gauge get", url: "/value/gauge/Test2", method: "GET", want: want{code: 200, getResult: "330095.942"}},
		{name: "#7 Counter get", url: "/value/counter/TestTestTestTestTest", method: "GET", want: want{code: 200, getResult: "312323"}},
		{name: "#8 Undefined counter get", url: "/value/counter/Undefined", method: "GET", want: want{code: 404}},
	}

	dumper := metricmemstorage.NewDumper("")
	store := metricmemstorage.NewMemStorage(dumper, 10, false)
	a := metricservice.New(store)
	hl := health.New()
	contr := New(a, config.Server{}, hl)
	ts := httptest.NewServer(contr.GetRouter())
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, get := testRequest(t, ts, tt.method, tt.url)
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tt.want.code, res.StatusCode)
			if tt.method == "GET" {
				assert.Equal(t, tt.want.getResult, get)
			}
		})
	}
}
