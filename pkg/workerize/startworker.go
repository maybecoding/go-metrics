package workerize

import "context"

// startWorker starts worker with checking context actuality
func (p *pool[T]) startWorker(ctx context.Context, ch chan T, num int) {
	if p.onStartWorker != nil {
		p.onStartWorker(ctx, num)
	}
	defer func() {
		if p.onFinishWorker != nil {
			p.onFinishWorker(ctx, num)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			for range ch {
			}
			return
		case m, ok := <-ch:
			if !ok {
				return
			}
			err := p.do(ctx, m)
			if err != nil && p.onDoErr != nil {
				p.onDoErr(ctx, err)
			}
		}
	}
}
