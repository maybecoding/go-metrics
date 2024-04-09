package grpcsender

import (
	"crypto/sha256"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/hasher"
	pb "github.com/maybecoding/go-metrics.git/pkg/metricv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Sender) initClient() error {
	var err error
	var cred credentials.TransportCredentials
	if s.cfg.CryptoKey != "" {
		cred, err = credentials.NewClientTLSFromFile(s.cfg.CryptoKey, "")
		if err != nil {
			return fmt.Errorf("grpcsender - initClient - NewClientTLSFromFile '%s': %w", s.cfg.CryptoKey, err)
		}
	} else {
		cred = insecure.NewCredentials()
	}

	s.clientConn, err = grpc.NewClient(s.cfg.GRPCServer,
		grpc.WithTransportCredentials(cred),
		grpc.WithUnaryInterceptor(hasher.ClientInterceptor(s.cfg.HashKey, sha256.New)))
	if err != nil {
		return fmt.Errorf("grpcsender - initClient - grpc.NewClient(%s): %w", s.cfg.GRPCServer, err)
	}
	s.clintGRPC = pb.NewMetricsV1Client(s.clientConn)
	return nil
}
