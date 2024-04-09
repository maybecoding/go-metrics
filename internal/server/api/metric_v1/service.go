package metric_v1

import (
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	srv "github.com/maybecoding/go-metrics.git/pkg/metric_v1"
	"google.golang.org/grpc"
	"net"
)

type Service struct {
	srv.UnimplementedMetricsV1Server
	grpcS         *grpc.Server
	metric        *metricservice.MetricService
	cfg           config.Server
	trustedSubNet *net.IPNet
}

func New(metric *metricservice.MetricService, cfg config.Server) *Service {
	s := &Service{metric: metric, cfg: cfg}
	if cfg.TrustedSubnet != "" {
		_, ipNet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			logger.Error().Err(err).Msg("can't parse trusted subnet")
		} else {
			s.trustedSubNet = ipNet
		}
	}
	return s
}
