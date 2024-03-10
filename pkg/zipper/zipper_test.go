package zipper

import (
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testBytes       = []byte("Test string must be zipped very well.")
	testZippedBytes = []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 10, 73, 45, 46, 81, 40, 46, 41, 202, 204, 75, 87, 200, 45, 45, 46, 81, 72, 74, 85, 168, 202, 44, 40, 72, 77, 81, 40, 75, 45, 170, 84, 40, 79, 205, 201, 209, 3, 4, 0, 0, 255, 255, 155, 224, 55, 142, 37, 0, 0, 0}
)

func BenchmarkCompareZip(b *testing.B) {

	b.Run("Zip without optimization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ZippedBytesSimple(testBytes)
			require.NoError(b, err)
		}
	})

	b.Run("Zip with optimization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ZippedBytes(testBytes)
			require.NoError(b, err)
		}
	})
}

func TestZipper(t *testing.T) {
	t.Run("#1 ZippedBytesSimple smoke", func(t *testing.T) {
		zb, err := ZippedBytesSimple(testBytes)
		require.NoError(t, err)
		require.Equal(t, testZippedBytes, zb)
	})

	t.Run("#2 ZippedBytes smoke", func(t *testing.T) {
		zb, err := ZippedBytes(testBytes)
		require.NoError(t, err)
		require.Equal(t, testZippedBytes, zb)
	})

	t.Run("#3 ZippedBytesSimple produces same results as ZippedBytes", func(t *testing.T) {
		numBytes := 1000
		numIterations := 1000
		for i := 0; i < numIterations; i += 1 {
			bytes := make([]byte, numBytes)
			_, err := rand.Read(bytes)
			require.NoError(t, err)

			zbs, err := ZippedBytesSimple(bytes)
			require.NoError(t, err)

			zb, err := ZippedBytes(bytes)
			require.NoError(t, err)

			require.Equal(t, zbs, zb)
		}
		zb, err := ZippedBytesSimple(testBytes)
		require.NoError(t, err)
		require.Equal(t, testZippedBytes, zb)
	})
}
