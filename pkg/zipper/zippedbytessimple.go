package zipper

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

// ZippedBytesSimple - function for zip bytes without optimization
func ZippedBytesSimple(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)
	// И записываем в него данные
	_, err := zw.Write(b)
	if err != nil {
		return nil, fmt.Errorf("zipper - ZippedBytesSimple - z.gzipW.Write: %w", err)
	}
	err = zw.Close()
	if err != nil {
		return nil, fmt.Errorf("zipper - ZippedBytesSimple - z.gzipW.Close: %w", err)
	}
	return buf.Bytes(), nil
}
