package compress

import (
	"io"
	"net/http"
)

type compressWriter struct {
	w  http.ResponseWriter
	cw io.WriteCloser
}

func (cw compressWriter) Header() http.Header {
	return cw.w.Header()
}
func (cw compressWriter) Write(b []byte) (int, error) {
	return cw.cw.Write(b)
}

func (cw compressWriter) WriteHeader(statusCode int) {
	cw.w.WriteHeader(statusCode)
}

func (cw compressWriter) Close() error {
	return cw.cw.Close()
}

func newCompressWriter(w http.ResponseWriter, cw io.WriteCloser) *compressWriter {
	return &compressWriter{
		w:  w,
		cw: cw,
	}
}
