package sender

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/maybecoding/go-metrics.git/pkg/workerize"
)

func (j *Sender) Run(inpM chan *app.Metrics) {
	workerize.From(inpM, j.sendMetric, j.cfg.NumWorkers).
		OnStartWorker(func(_ context.Context, num int) {
			logger.Debug().Int("number", num).Msg("Started worker")
		}).
		OnFinishWorker(func(_ context.Context, num int) {
			logger.Debug().Int("number", num).Msg("Finished worker")
		}).
		OnDoError(func(ctx context.Context, err error) {
			logger.Error().Err(err).Msg("worker error")
		}).
		Run(j.ctx)
}
