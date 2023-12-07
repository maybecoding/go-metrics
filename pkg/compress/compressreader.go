package compress

import (
	"io"
)

type compressReader struct {
	r  io.ReadCloser
	cr io.ReadCloser
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.cr.Read(p)
}

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
