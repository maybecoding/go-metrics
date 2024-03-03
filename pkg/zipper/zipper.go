// Package zipper - helper for zip bytes
package zipper

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"sync"
)

type zipper struct {
	bytesB *bytes.Buffer
	gzipW  *gzip.Writer
}

func newZipper() *zipper {
	z := &zipper{}
	z.bytesB = bytes.NewBuffer(nil)
	z.gzipW = gzip.NewWriter(z.bytesB)
	return z
}

func (z *zipper) zip(bs []byte) ([]byte, error) {
	z.bytesB.Reset()
	z.gzipW.Reset(z.bytesB)
	_, err := z.gzipW.Write(bs)
	if err != nil {
		return nil, fmt.Errorf("zipper - zip - z.gzipW.Write: %w", err)
	}
	err = z.gzipW.Close()
	if err != nil {
		return nil, fmt.Errorf("zipper - zip - z.gzipW.Close: %w", err)
	}
	return z.bytesB.Bytes(), nil
}

var zipPool = sync.Pool{New: func() any { return newZipper() }}
