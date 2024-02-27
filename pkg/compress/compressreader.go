package compress

import (
	"io"
)

type compressReader struct {
	r  io.ReadCloser
	cr io.ReadCloser
}

// Read - Реализация метода Read в оболочке над кастомным ридером с функцией закрытия исходного ридера
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.cr.Read(p)
}

// Close - Реализация метода Close в оболочке над кастомным ридером с функцией закрытия исходного ридера
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.cr.Close()
}

func newCompressReader(r io.ReadCloser, cr io.ReadCloser) *compressReader {
	return &compressReader{
		r:  r,
		cr: cr,
	}
}
