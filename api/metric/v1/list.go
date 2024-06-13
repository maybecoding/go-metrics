package v1

import (
	"context"
	"github.com/maybecoding/go-metrics.git/api/metric/v1/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) List(_ context.Context, _ *pb.Empty) (*pb.MetricList, error) {
	mts, err := s.metric.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %s", err.Error())
	}
	var outML pb.MetricList
	outML.Metrics = make([]*pb.Metric, 0, len(mts))
	for _, m := range mts {
		outM := pb.Metric{Id: m.ID, Type: m.MType, Value: m.Value, Delta: m.Delta}
		outML.Metrics = append(outML.Metrics, &outM)
	}
	return &outML, nil
}
