package zipper

// ZippedBytes - optimized function for zip bytes using pool
func ZippedBytes(b []byte) ([]byte, error) {
	z := zipPool.Get().(*zipper)
	defer zipPool.Put(z)
	return z.zip(b)
}
