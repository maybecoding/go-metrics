package grpcsender

import "github.com/maybecoding/go-metrics.git/pkg/logger"

func (s *Sender) Terminate() {
	if s.clientConn != nil {
		_ = s.clientConn.Close()
	}
	logger.Info().Msg("Stopped gRPC server")
}
