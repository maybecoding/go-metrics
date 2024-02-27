// Package hashcheck - package for middleware-functions for checking hash for http-handlers
package hashcheck

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
)

// HashCheck - struct for checking hash using key
type HashCheck struct {
	hashFunc   func() hash.Hash
	key        string
	headerName string
}

// New constructs new HashCheck struct with hash function, key and header for hash for checking
func New(hashFunc func() hash.Hash, key, headerName string) *HashCheck {
	return &HashCheck{
		hashFunc:   hashFunc,
		key:        key,
		headerName: headerName,
	}
}

// HandlerFunc - middleware function for checking hash
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
			fmt.Println("1", hsHex, "2", hsHexReq)
			if hsHex != hsHexReq {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(""))
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		handlerFn(w, r)
	}
}

// Handler - middleware function for checking hash
func (hc *HashCheck) Handler(h http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
	return hc.HandlerFunc(handlerFn)
}
