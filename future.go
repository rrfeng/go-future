package future

import (
	"context"
)

type Future[T any] interface {
	Await() (T, error)
	Cancel()
}

type future[T any] struct {
	ret    chan result[T]
	ctx    context.Context
	cancel context.CancelFunc
}

type result[T any] struct {
	dat T
	err error
}

func Async[T any](fn func() (T, error)) Future[T] {
	return AsyncWithContext(context.Background(), fn)
}

func AsyncWithContext[T any](ctx context.Context, fn func() (T, error)) Future[T] {
	c := make(chan result[T], 1)
	fctx, cancel := context.WithCancel(ctx)

	go func() {
		ret, err := fn()
		c <- result[T]{dat: ret, err: err}
	}()
	return &future[T]{ret: c, ctx: fctx, cancel: cancel}
}

func (f *future[T]) Await() (T, error) {
	var result result[T]

	select {
	case <-f.ctx.Done():
		result.err = f.ctx.Err()
	case ret := <-f.ret:
		result = ret
	}

	f.cancel()
	return result.dat, result.err
}

func (f *future[T]) Cancel() {
	f.cancel()
}
