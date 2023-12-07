package compress

import (
	"compress/gzip"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
	"strings"
)

func HandlerFuncReader(handlerFn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		logger.Log.Debug().Str("encoding", contentEncoding).Msg("encoding")
		if strings.Contains(contentEncoding, "gzip") {
			logger.Log.Debug().Msg("body is compressed with gzip")
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

		// передаём управление хендлеру
		handlerFn(w, r)
	}
}

func HandlerReader(h http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
	return HandlerFuncReader(handlerFn)
}
