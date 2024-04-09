package grpcsender

import (
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	pb "github.com/maybecoding/go-metrics.git/pkg/metric_v1"
	"time"
)

func (s *Sender) sendMetric(m *app.Metrics) {
	pbM := pb.Metric{Id: m.ID, Type: m.MType, Value: m.Value, Delta: m.Delta}
	// Создаем попытки на отправку
	intervals := make([]time.Duration, 0, len(s.cfg.RetryIntervals)+1)
	intervals = append(intervals, 0)
	intervals = append(intervals, s.cfg.RetryIntervals...)
	for i, interval := range intervals {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(interval):
			_, err := s.clintGRPC.Set(s.ctx, &pbM)
			if err == nil {
				return
			} else {
				logger.Error().Err(err).Int("Attempt", i).Msg("grpcsender - sendMetric")
			}
		}
	}
}
