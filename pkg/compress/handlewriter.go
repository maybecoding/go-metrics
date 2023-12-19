package compress

import (
	"compress/gzip"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
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

func HandlerFuncWriter(handlerFn http.HandlerFunc, compLevel ...int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// По умолчанию оригинальный ResponseWriter
		ow := w
		// Логика для gzip (можно подключить другие)
		ae := r.Header.Get("Accept-Encoding")
		logger.Debug().Str("Accept-Encoding", ae).Msg("on compression")
		if strings.Contains(ae, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			// Пытаемся понять какой уровень сжатия
			level := DefaultCompression
			if 1 <= len(compLevel) {
				level = compLevel[0]
				if level != BestSpeed && level != BestCompression && level != DefaultCompression {
					logger.Error().Int("level", level).Msg("compression write level is wrong, using default")
					level = DefaultCompression
				}
			}

			gzipLevel, ok := gzipLevels[level]
			if !ok {
				logger.Error().Int("metric compression level", level).Msg("gzip compression level can't be identify by mapping, using default")
				level = gzipLevels[DefaultCompression]
			}
			cw, err := gzip.NewWriterLevel(w, gzipLevel)
			if err != nil {
				logger.Error().Str("compression", "gzip").Int("level", level).Msg("can't set compression level, using default")
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
