package grpcsender

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/api/metric/v1/pb"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"time"
)

func (s *Sender) sendMetric(ctx context.Context, m *app.Metrics) error {
	pbM := pb.Metric{Id: m.ID, Type: m.MType, Value: m.Value, Delta: m.Delta}
	// Создаем попытки на отправку
	intervals := make([]time.Duration, 0, len(s.cfg.RetryIntervals)+1)
	intervals = append(intervals, 0)
	intervals = append(intervals, s.cfg.RetryIntervals...)
	for i, interval := range intervals {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(interval):
			_, err := s.clintGRPC.Set(ctx, &pbM)
			if err == nil {
				return nil
			} else {
				logger.Error().Err(err).Int("Attempt", i).Msg("grpcsender - sendMetric")
			}
		}
	}
	return fmt.Errorf("grpcsender - sendMetric - all %d attempts failed", len(intervals))
}
