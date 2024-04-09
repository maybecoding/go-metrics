package metric_v1

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	pb "github.com/maybecoding/go-metrics.git/pkg/metric_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) SetAll(_ context.Context, pbMts *pb.MetricList) (*pb.Empty, error) {
	if pbMts == nil {
		return nil, status.Error(codes.InvalidArgument, "please provide correct request body")
	}
	for _, pbM := range pbMts.Metrics {
		if pbM == nil {
			continue
		}
		m := entity.Metrics{ID: pbM.Id, MType: pbM.Type, Value: pbM.Value, Delta: pbM.Delta}
		err := s.metric.Set(m)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal server error: %s", err.Error())
		}
	}
	return new(pb.Empty), nil
}
