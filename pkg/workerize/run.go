package workerize

import (
	"context"
	"sync"
)

// Run Creates configured num of routines and starts workers
func (p *pool[T]) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for i := 0; i < p.numWorkers; i += 1 {
		ii := i

		wg.Add(1)
		go func() {
			defer wg.Done()
			p.startWorker(ctx, p.ch, ii)
		}()
	}
	wg.Wait()
}
