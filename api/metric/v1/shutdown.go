package v1

import (
	"context"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

func (s *Service) Shutdown(_ context.Context) error {
	s.grpcS.GracefulStop()
	logger.Info().Msg("Stopped gRPC service")
	return nil
}
