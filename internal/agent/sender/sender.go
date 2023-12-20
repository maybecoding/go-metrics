package sender

import (
	"context"
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"sync"
	"time"
)

type Sender struct {
	endpoint       string
	retryIntervals []time.Duration
	hashKey        string
	ctx            context.Context
	numWorkers     int
}

func (j *Sender) SendWorker(inpM chan *app.Metrics, id int) {
	logger.Debug().Int("number", id).Msg("Started worker")
	defer func() {
		logger.Debug().Int("number", id).Msg("Stopped worker")
	}()
	for {
		select {
		case <-j.ctx.Done():
			return
		case m, ok := <-inpM:
			if !ok {
				return
			}
			j.sendMetric(m)
		}
	}
}

func (j *Sender) SendStart(inpM chan *app.Metrics) {
	wg := &sync.WaitGroup{}

	for i := 0; i < j.numWorkers; i += 1 {
		ii := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			j.SendWorker(inpM, ii)
		}()
	}
	wg.Wait()
}

func New(template, serverAddress string, retryIntervals []time.Duration, hashKey string, ctx context.Context, numWorkers int) *Sender {
	return &Sender{
		endpoint:       fmt.Sprintf(template, serverAddress),
		retryIntervals: retryIntervals,
		hashKey:        hashKey,
		ctx:            ctx,
		numWorkers:     numWorkers,
	}
}
