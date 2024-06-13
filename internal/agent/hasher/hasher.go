package hasher

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
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

func ClientInterceptor(key string, hashFn func() hash.Hash) grpc.UnaryClientInterceptor {
	if key == "" {
		return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		protoMsg, ok := req.(proto.Message)
		if !ok {
			logger.Error().Msg("can't calc hash: req.(proto.Message)")
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		data, err := proto.Marshal(protoMsg)
		if err != nil {
			logger.Error().Err(err).Msg("can't calc hash")
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		hs := hmac.New(hashFn, []byte(key))
		hs.Write(data)
		hsSum := hs.Sum(nil)
		hsHex := hex.EncodeToString(hsSum)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(map[string]string{})
		}
		md.Set("hashsha256", hsHex)

		newCtx := metadata.NewOutgoingContext(ctx, md)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}
