package sender

import (
	"context"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/internal/agent/config"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net"
	"sync"
)

type Sender struct {
	ctx context.Context
	cfg config.Sender
	ip  net.IP
}

func New(ctx context.Context, cfg config.Sender) *Sender {
	s := &Sender{
		ctx: ctx,
		cfg: cfg,
	}
	s.identifyIP()
	return s
}

func (j *Sender) Worker(inpM chan *app.Metrics, id int) {
	logger.Debug().Int("number", id).Msg("Started worker")
	defer func() {
		logger.Debug().Int("number", id).Msg("Stopped worker")
	}()
	for {
		select {
		case <-j.ctx.Done():
			for range inpM {
			}
			return
		case m, ok := <-inpM:
			if !ok {
				return
			}
			j.sendMetric(m)
		}
	}
}

func (j *Sender) Run(inpM chan *app.Metrics) {
	// Инициализируем массив с gzip-writer

	wg := &sync.WaitGroup{}

	for i := 0; i < j.cfg.NumWorkers; i += 1 {
		ii := i

		wg.Add(1)
		go func() {
			defer wg.Done()
			j.Worker(inpM, ii)
		}()
	}
	wg.Wait()
}
