package metric_v1

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/server/api/metric_v1/interseptor"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	pb "github.com/maybecoding/go-metrics.git/pkg/metric_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

func (s *Service) Run(ctx context.Context) error {
	listen, err := net.Listen("tcp", s.cfg.GRPCAddress)
	if err != nil {
		return fmt.Errorf("metric_v1 - Run - net.Listen: %w", err)
	}
	var opt []grpc.ServerOption
	// Если необходимо использовать TLS
	if s.cfg.CryptoKey != "" {
		cred, err := credentials.NewServerTLSFromFile(s.cfg.CryptoKey, s.cfg.CryptoKey)
		if err != nil {
			return fmt.Errorf("metric_v1 - Run -NewServerTLSFromFile: %w", err)
		}
		// Создание экземпляра gRPC-сервера с поддержкой SSL/TLS
		opt = append(opt, grpc.Creds(cred))
	}
	opt = append(opt, grpc.ChainUnaryInterceptor(
		interseptor.CheckTrustedSubnet(s.trustedSubNet, s.cfg.IPAddrHeader),
		interseptor.HashCheck(sha256.New, s.cfg.HashKey, "HashSHA256")))

	s.grpcS = grpc.NewServer(opt...)

	pb.RegisterMetricsV1Server(s.grpcS, s)
	logger.Info().Str("address", s.cfg.GRPCAddress).Msg("Start grpc server")
	if err = s.grpcS.Serve(listen); err != nil {
		return fmt.Errorf("metric_v1 - Run - grpcS.Serve: %w", err)
	}

	return nil
}
