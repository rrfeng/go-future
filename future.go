package future

import (
	"context"
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
	return NewWithContext(context.Background(), fn)
}

func NewWithContext[T any](ctx context.Context, fn func() (T, error)) *Future[T] {
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

	f.cancel()
	return result.dat, result.err
}

func (f *Future[T]) Cancel() {
	f.cancel()
}
