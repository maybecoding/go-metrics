package grpcsender

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/maybecoding/go-metrics.git/pkg/workerize"
)

func (s *Sender) Run(cMt chan *app.Metrics) {
	logger.Info().Msg("Run gRPC sender")
	err := s.initClient()
	if err != nil {
		logger.Error().Err(err).Msg("grpcsender - Run")
		return
	}
	s.setupContext()

	go func() {
		<-s.ctx.Done()
		s.Terminate()
	}()
	workerize.From(cMt, s.sendMetric, s.cfg.NumWorkers).
		OnStartWorker(func(_ context.Context, num int) {
			logger.Debug().Int("number", num).Msg("Started worker")
		}).
		OnFinishWorker(func(_ context.Context, num int) {
			logger.Debug().Int("number", num).Msg("Finished worker")
		}).
		OnDoError(func(ctx context.Context, err error) {
			logger.Error().Err(err).Msg("worker error")
		}).
		Run(s.ctx)
}
