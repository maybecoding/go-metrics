package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerReader(t *testing.T) {

	type Req struct {
		body           string
		headerEncoding string
		mustZipped     bool
	}

	type Resp struct {
		body string
		code int
	}
	tests := []struct {
		name string
		req  Req
		resp Resp
	}{
		{
			name: "#1 Ordinary",
			req:  Req{body: "Some body", mustZipped: false, headerEncoding: "-"},
			resp: Resp{body: "", code: http.StatusOK},
		},
		{
			name: "#2 GZIPed",
			req:  Req{body: "Some body", mustZipped: true, headerEncoding: "gzip"},
			resp: Resp{body: "", code: http.StatusOK},
		},
		{
			name: "#3 no data",
			req:  Req{body: "", mustZipped: true, headerEncoding: "gzip"},
			resp: Resp{body: "", code: http.StatusInternalServerError},
		},
	}
	for _, test := range tests {

		var req *http.Request

		if test.req.body != "" {
			// Подготавливаем ожидаемые запросы
			buf := bytes.NewBuffer(nil)

			if test.req.mustZipped {
				zw := gzip.NewWriter(buf)
				_, err := zw.Write([]byte(test.req.body))
				require.NoError(t, err)
				err = zw.Close()
				require.NoError(t, err)
			} else {
				_, err := buf.Write([]byte(test.req.body))
				require.NoError(t, err)
			}

			req = httptest.NewRequest("POST", "/", buf)
		} else {
			req = httptest.NewRequest("POST", "/", nil)
		}
		req.Header.Set("Content-Encoding", test.req.headerEncoding)
		w := httptest.NewRecorder()

		handler := HandlerReader(checkRequestBody(t, test.req.body))
		handler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, test.resp.code, resp.StatusCode)
		_ = resp.Body.Close()
	}
}

func checkRequestBody(t *testing.T, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Проверяем тело запроса
		reqData, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer func() {
			_ = r.Body.Close()
		}()
		require.Equal(t, body, string(reqData))

		w.WriteHeader(http.StatusOK)
		require.NoError(t, err)
	}
}

func sendResponse(t *testing.T, response []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(response)
		require.NoError(t, err)
	}
}

func TestHandlerFuncWriter(t *testing.T) {
	type Req struct {
		body                 string
		headerAcceptEncoding string
		mustZipped           bool
		compressLevel        int
	}

	type Resp struct {
		body            string
		contentEncoding string
		code            int
	}
	tests := []struct {
		name string
		resp Resp
		req  Req
	}{
		{
			name: "#1 No zip required",
			req:  Req{body: "Some body", mustZipped: false, headerAcceptEncoding: "-", compressLevel: 0},
			resp: Resp{body: "", code: http.StatusOK, contentEncoding: ""},
		},
		{
			name: "#2 Zipping",
			req:  Req{body: "Some body", mustZipped: true, headerAcceptEncoding: "gzip", compressLevel: 0},
			resp: Resp{body: "", code: http.StatusOK, contentEncoding: "gzip"},
		},
		{
			name: "#3 with compress level",
			req:  Req{body: "Some body", mustZipped: true, headerAcceptEncoding: "gzip", compressLevel: 2},
			resp: Resp{body: "", code: http.StatusOK, contentEncoding: "gzip"},
		},
		{
			name: "#4 with wrong compress level",
			req:  Req{body: "Some body", mustZipped: true, headerAcceptEncoding: "gzip", compressLevel: 100},
			resp: Resp{body: "", code: http.StatusOK, contentEncoding: "gzip"},
		},
	}
	for _, test := range tests {

		// Подготавливаем ожидаемые запросы
		buf := bytes.NewBuffer(nil)
		_, err := buf.Write([]byte(test.req.body))
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/", buf)

		req.Header.Set("Accept-Encoding", test.req.headerAcceptEncoding)
		w := httptest.NewRecorder()

		var cmpLevels []int
		if test.req.compressLevel > 0 {
			cmpLevels = append(cmpLevels, test.req.compressLevel)
		}
		handler := HandlerFuncWriter(sendResponse(t, []byte(test.req.body)), cmpLevels...)
		handler(w, req)

		resp := w.Result()
		require.Equal(t, test.resp.code, resp.StatusCode)

		if test.resp.contentEncoding != "" {
			require.Equal(t, test.resp.contentEncoding, resp.Header.Get("Content-Encoding"))
		}
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		_ = resp.Body.Close()
		if test.req.mustZipped {
			respBody = unzip(t, respBody)
		}
		require.Equal(t, test.req.body, string(respBody))
	}
}

func unzip(t *testing.T, source []byte) []byte {
	zr, err := gzip.NewReader(bytes.NewReader(source))
	require.NoError(t, err)
	res, err := io.ReadAll(zr)
	_ = zr.Close()
	require.NoError(t, err)
	return res
}
