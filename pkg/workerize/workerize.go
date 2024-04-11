// Package workerize for handling chan of some type with provided handlers and num of workers
package workerize

import (
	"context"
)

type pool[T any] struct {
	ch         chan T
	do         func(context.Context, T) error
	numWorkers int

	onStartWorker  func(context.Context, int)
	onFinishWorker func(context.Context, int)

	onDoErr func(context.Context, error)
}

// From returns workerized provided chan and setups n workers with handlers
func From[T any](ch chan T, do func(context.Context, T) error, numWorkers int) *pool[T] {
	return &pool[T]{ch: ch, do: do, numWorkers: numWorkers}
}

// OnStartWorker callback on start worker
func (p *pool[T]) OnStartWorker(fn func(context.Context, int)) *pool[T] {
	p.onStartWorker = fn
	return p
}

// OnFinishWorker callback on finish worker
func (p *pool[T]) OnFinishWorker(fn func(context.Context, int)) *pool[T] {
	p.onFinishWorker = fn
	return p
}

// OnDoError callback on handle func error
func (p *pool[T]) OnDoError(fn func(context.Context, error)) *pool[T] {
	p.onDoErr = fn
	return p
}
