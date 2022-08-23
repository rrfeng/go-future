package future

import (
	"context"
	"time"
)

type Future[T any] struct {
	ret    chan Result[T]
	ctx    context.Context
	cancel context.CancelFunc
}

type Result[T any] struct {
	dat T
	err error
}

func New[T any](fn func() (T, error)) *Future[T] {
	c := make(chan Result[T], 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ret, err := fn()
		c <- Result[T]{dat: ret, err: err}
	}()
	return &Future[T]{ret: c, ctx: ctx, cancel: cancel}
}

func WithContext[T any](fn func() (T, error), ctx context.Context) *Future[T] {
	c := make(chan Result[T], 1)
	fctx, cancel := context.WithCancel(ctx)

	go func() {
		ret, err := fn()
		c <- Result[T]{dat: ret, err: err}
	}()
	return &Future[T]{ret: c, ctx: fctx, cancel: cancel}
}

func (f *Future[T]) Wait() (T, error) {
	var result Result[T]

	select {
	case <-f.ctx.Done():
		result.err = f.ctx.Err()
	case ret := <-f.ret:
		result = ret
	}
	return result.dat, result.err
}

func (f *Future[T]) Cancel() {
	f.cancel()
}

