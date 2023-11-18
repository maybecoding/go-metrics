package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

type (
	responseData struct {
		statusCode int
		contentLen int
	}

	proxyResponseWriter struct {
		wInter http.ResponseWriter
		resp   *responseData
	}
)

func newProxyResponseWriter(w http.ResponseWriter, resp *responseData) http.ResponseWriter {
	return &proxyResponseWriter{wInter: w, resp: resp}
}

func (w *proxyResponseWriter) Header() http.Header {
	return w.wInter.Header()
}
func (w *proxyResponseWriter) Write(b []byte) (int, error) {
	contentLen, err := w.wInter.Write(b)
	w.resp.contentLen = contentLen
	return contentLen, err
}
func (w *proxyResponseWriter) WriteHeader(statusCode int) {
	w.resp.statusCode = statusCode
	w.wInter.WriteHeader(statusCode)
}

// Log Общая переменная для логирования будет доступна всему коду
// Не лучшее решение, но самое простое
var Log *zerolog.Logger

func Init(level string) {

	// Пока используем консольный вывод
	//zl := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zl := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	switch level {
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zl.Debug().Msg("passed wrong error level")
	}
	zl.Debug().Str("log level", level).Msg("log initialized")

	Log = &zl
	return
}

func Handler(h http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()
		respData := responseData{}
		wproxy := newProxyResponseWriter(w, &respData)
		h.ServeHTTP(wproxy, r)
		Log.Debug().
			Str("URI", r.RequestURI).
			Dur("duration", time.Since(timeStart)).
			Str("method", r.Method).
			Int("content length", respData.contentLen).
			Int("status code", respData.statusCode).
			Msg("HTTP request handled")
	}

	return http.HandlerFunc(handlerFn)
}
