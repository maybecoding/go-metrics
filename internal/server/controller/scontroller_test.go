package controller

import (
	"github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/memstorage"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type want struct {
	code int
}

func TestController(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{name: "#1 Bad url", url: "/updat/gauge/Alloc/1", method: "POST", want: want{code: 404}},
		{name: "#2 Bad url", url: "/u", method: "POST", want: want{code: 404}},
		{name: "#3 Bad url", url: "/update/gauge/Alloc", method: "POST", want: want{code: 404}},
		{name: "#4 Unknown metric", url: "/update/gauge2/Test/1", method: "POST", want: want{code: 400}},
		{name: "#5 Unknown metric", url: "/update/counter2/Test/1", method: "POST", want: want{code: 400}},

		{name: "#6 Gauge success", url: "/update/gauge/Test/7.77", method: "POST", want: want{code: 200}},
		{name: "#7 Counter success", url: "/update/counter/TestTestTestTestTest/312323", method: "POST", want: want{code: 200}},

		{name: "#8 Gauge bad value", url: "/update/gauge/Test/7.77x", method: "POST", want: want{code: 400}},
		{name: "#9 Counter bad value", url: "/update/counter/TestTestTestTestTest/1.1", method: "POST", want: want{code: 400}},
	}

	store := smemstorage.NewMemStorage()
	app := app.New(store)
	contr := New(app, "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			contr.handleUpdate(w, req)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
