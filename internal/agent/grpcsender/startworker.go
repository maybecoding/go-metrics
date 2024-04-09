package grpcsender

import (
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

func (s *Sender) startWorker(cMt chan *app.Metrics, num int) {
	logger.Debug().Int("number", num).Msg("Started worker")
	defer func() {
		logger.Debug().Int("number", num).Msg("Stopped worker")
	}()
	for {
		select {
		case <-s.ctx.Done():
			for range cMt {
			}
			return
		case m, ok := <-cMt:
			if !ok {
				return
			}
			s.sendMetric(m)
		}
	}
}
