package future

import (
	"context"
)

type Future[T any] struct {
	ret    chan result[T]
	ctx    context.Context
	cancel context.CancelFunc
}

type result[T any] struct {
	dat T
	err error
}

func Async[T any](fn func() (T, error)) *Future[T] {
	return AsyncWithContext(context.Background(), fn)
}

func AsyncWithContext[T any](ctx context.Context, fn func() (T, error)) *Future[T] {
	c := make(chan result[T], 1)
	fctx, cancel := context.WithCancel(ctx)

	go func() {
		ret, err := fn()
		c <- result[T]{dat: ret, err: err}
		close(c)
	}()

	return &Future[T]{ret: c, ctx: fctx, cancel: cancel}
}

func (f *Future[T]) Await() (T, error) {
	defer f.cancel()

	var result result[T]

	select {
	case <-f.ctx.Done():
		select {
		case ret := <-f.ret:
			result = ret
		default:
			result.err = f.ctx.Err()
		}
	case ret := <-f.ret:
		result = ret
	}
	return result.dat, result.err
}

func (f *Future[T]) Cancel() {
	f.cancel()
}
