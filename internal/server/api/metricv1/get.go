package metricv1

import (
	"context"
	"errors"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
	pb "github.com/maybecoding/go-metrics.git/pkg/metricv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Get(_ context.Context, mh *pb.MetricHeader) (*pb.Metric, error) {
	if mh == nil {
		return nil, status.Error(codes.InvalidArgument, "please provide correct request body")
	}
	var mt entity.Metrics
	mt.ID = mh.Id
	mt.MType = mh.Type
	err := s.metric.Get(&mt)
	if err != nil {
		if errors.Is(err, metricservice.ErrNoMetricValue) {
			return nil, status.Error(codes.NotFound, "metric not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error: %s", err.Error())
	}
	outMt := pb.Metric{Id: mh.Id, Type: mh.Type, Delta: mt.Delta, Value: mt.Value}
	return &outMt, nil
}
