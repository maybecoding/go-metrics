package hashcheck

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	hashHeaderName = "HASH_HEADER"
	alg            = sha256.New
)

func TestHandler(t *testing.T) {

	type Req struct {
		body                string
		key                 string
		hash                string
		updateHashToCorrect bool
	}

	type Resp struct {
		code int
	}

	tests := []struct {
		name string
		req  Req
		resp Resp
	}{
		{"#1 No key no hash required", Req{"Some text", "", "", false}, Resp{http.StatusOK}},
		{"#2 Bad hash", Req{"Some text", "key", "incorrect", false}, Resp{http.StatusBadRequest}},
		{"#3 Good hash", Req{"Some text", "key", "fix me", true}, Resp{http.StatusOK}},
	}

	for _, tt := range tests {

		// Подготавливаем ожидаемые запросы
		buf := bytes.NewBuffer(nil)
		_, err := buf.Write([]byte(tt.req.body))
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/", buf)

		// Вычисляем хэш если требуется
		hash := tt.req.hash
		if tt.req.updateHashToCorrect {
			// Считаем hash
			hs := hmac.New(alg, []byte(tt.req.key))
			hs.Write([]byte(tt.req.body))
			hsSum := hs.Sum(nil)
			hash = hex.EncodeToString(hsSum)
		}

		if hash != "" {
			req.Header.Set(hashHeaderName, hash)
		}

		w := httptest.NewRecorder()

		hc := New(alg, tt.req.key, hashHeaderName)

		handler := hc.Handler(sendResponse(t, []byte("ok")))
		handler.ServeHTTP(w, req)

		resp := w.Result()
		_ = resp.Body.Close()
		require.Equal(t, tt.resp.code, resp.StatusCode)
	}
}

func sendResponse(t *testing.T, response []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(response)
		require.NoError(t, err)
	}
}
