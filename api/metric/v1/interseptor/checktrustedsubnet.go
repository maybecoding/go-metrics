package interseptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"strings"
)

func CheckTrustedSubnet(trustedSubNet *net.IPNet, ipAddrHeader string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if trustedSubNet == nil {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ipMd := md.Get(strings.ToLower(ipAddrHeader))
			if len(ipMd) > 0 {
				ipAddr := net.ParseIP(ipMd[0])
				if trustedSubNet.Contains(ipAddr) {
					return handler(ctx, req)
				}
			}
		}
		return nil, status.Errorf(codes.PermissionDenied, "failed check ip")

	}
}
