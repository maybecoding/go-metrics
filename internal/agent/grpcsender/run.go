package grpcsender

import (
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"sync"
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
	wg := &sync.WaitGroup{}
	for i := 0; i < s.cfg.NumWorkers; i += 1 {
		ii := i

		wg.Add(1)
		go func() {
			defer wg.Done()
			s.startWorker(cMt, ii)
		}()
	}
	wg.Wait()
}
