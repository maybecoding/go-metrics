package hasher

import (
	"crypto/hmac"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"hash"
)

func New(key string, hashFn func() hash.Hash) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {
		if key == "" {
			return nil
		}
		if body, ok := r.Body.([]byte); ok {
			hs := hmac.New(hashFn, []byte(key))
			hs.Write(body)
			hsSum := hs.Sum(nil)
			hsHex := hex.EncodeToString(hsSum)
			r.SetHeader("HashSHA256", hsHex)
		}
		return nil
	}
}
