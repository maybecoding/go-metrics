// Package compress - package for middleware-functions for zip and unzip using gzip for http-handlers
package compress

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// HandlerFuncReader - middleware for read gzip by type http.HandlerFunc
func HandlerFuncReader(handlerFn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = newCompressReader(r.Body, cr)
			defer func() {
				_ = cr.Close()
			}()
		}

		// передаём управление хэндлеру
		handlerFn(w, r)
	}
}

// HandlerReader - middleware for read gzip by type  http.Handler
func HandlerReader(h http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
	return HandlerFuncReader(handlerFn)
}
