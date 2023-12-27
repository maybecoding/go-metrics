package hashcheck

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"hash"
	"io"
	"net/http"
)

type HashCheck struct {
	hashFunc   func() hash.Hash
	key        string
	headerName string
}

func New(hashFunc func() hash.Hash, key, headerName string) *HashCheck {
	return &HashCheck{
		hashFunc:   hashFunc,
		key:        key,
		headerName: headerName,
	}
}

func (hc *HashCheck) HandlerFunc(handlerFn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if hc.key != "" {
			hs := hmac.New(hc.hashFunc, []byte(hc.key))
			body, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()

			hs.Write(body)
			hsSum := hs.Sum(nil)
			hsHex := hex.EncodeToString(hsSum)
			hsHexReq := r.Header.Get(hc.headerName)
			if hsHex != hsHexReq {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				logger.Debug().Str("actual hash", hsHex).Str("declared hash", hsHexReq).Msg("error due hash check")
				_, _ = w.Write([]byte(""))
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		handlerFn(w, r)
	}
}

func (hc *HashCheck) Handler(h http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
	return hc.HandlerFunc(handlerFn)
}
