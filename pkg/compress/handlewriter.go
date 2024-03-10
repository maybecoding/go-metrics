package compress

import (
	"compress/gzip"
	"net/http"
	"strings"
)

const (
	BestSpeed = iota + 1
	BestCompression
	DefaultCompression
)

var gzipLevels = map[int]int{
	BestSpeed:          gzip.BestSpeed,
	BestCompression:    gzip.BestCompression,
	DefaultCompression: gzip.DefaultCompression,
}

// HandlerFuncWriter - middleware for compression response body using gzip
// compLevel - optional parameter for set compression level, if passed more than one all after first is ignored
func HandlerFuncWriter(handlerFn http.HandlerFunc, compLevel ...int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// По умолчанию оригинальный ResponseWriter
		ow := w
		// Логика для gzip (можно подключить другие)
		ae := r.Header.Get("Accept-Encoding")
		if strings.Contains(ae, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			// Пытаемся понять какой уровень сжатия
			level := DefaultCompression
			if 1 <= len(compLevel) {
				level = compLevel[0]
				if level != BestSpeed && level != BestCompression && level != DefaultCompression {
					level = DefaultCompression
				}
			}

			gzipLevel := gzipLevels[level]
			cw, err := gzip.NewWriterLevel(w, gzipLevel)
			if err != nil {
				cw = gzip.NewWriter(w)
			}
			defer func() {
				_ = cw.Close()
			}()
			ow = newCompressWriter(w, cw)
		}
		handlerFn(ow, r)
	}
}
