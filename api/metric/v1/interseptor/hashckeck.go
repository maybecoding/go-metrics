package interseptor

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"hash"
	"strings"
)

func HashCheck(hashFunc func() hash.Hash, key, headerName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if key == "" {
			return handler(ctx, req)
		}
		// Получаем хэш из заголовка (положительный результат только если все будет хорошо
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			hashMd := md.Get(strings.ToLower(headerName))
			if len(hashMd) > 0 {
				protoMsg, ok := req.(proto.Message)
				if ok {
					data, err := proto.Marshal(protoMsg)
					if err == nil {
						hs := hmac.New(hashFunc, []byte(key))
						hs.Write(data)
						hsSum := hs.Sum(nil)
						hsHex := hex.EncodeToString(hsSum)
						if hsHex == hashMd[0] {
							return handler(ctx, req)
						}
					}
				}
			}
		}
		return nil, status.Errorf(codes.FailedPrecondition, "failed check hash")
	}
}
