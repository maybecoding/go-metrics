package metric_v1

import (
	"context"
	"errors"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	pb "github.com/maybecoding/go-metrics.git/pkg/metric_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Set(ctx context.Context, pbM *pb.Metric) (*pb.Empty, error) {
	if pbM == nil {
		return nil, status.Error(codes.InvalidArgument, "please provide correct request body")
	}
	m := entity.Metrics{ID: pbM.Id, MType: pbM.Type, Value: pbM.Value, Delta: pbM.Delta}
	err := s.metric.Set(m)
	if err != nil {
		if errors.Is(err, metricservice.ErrMetricTypeIncorrect) {
			return nil, status.Error(codes.InvalidArgument, "metric type incorrect, must me 'gauge' or 'counter'")
		} else if errors.Is(err, metricservice.ErrNoMetricValue) {
			return nil, status.Error(codes.InvalidArgument, "for type 'gauge' must be passed 'value', for 'counter' 'delta'")
		}
		return nil, status.Errorf(codes.Internal, "internal server error: %s", err.Error())
	}
	return new(pb.Empty), nil
}
