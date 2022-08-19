package future

import (
	"context"
	"fmt"
	"time"
)

type Future[T any] struct {
	ret chan T
	ctx context.Context
	cancel context.CancelFunc
}

func New[T any](fn func() T) *Future[T] {
	c := make(chan T, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() { c <- fn() }()
	return &Future[T]{ret: c, ctx: ctx, cancel: cancel}
}

func WithContext[T any](fn func() T, ctx context.Context) *Future[T] {
	c := make(chan T, 1)
	fctx, cancel := context.WithCancel(ctx)

	go func() { c <- fn() }()
	return &Future[T]{ret: c, ctx: fctx, cancel: cancel}
}

func (f *Future[T]) Wait() (T, error) {
	var result T
	var err error
	select {
	case <-f.ctx.Done():
		err = f.ctx.Err()
	case tmp := <-f.ret:
		result = tmp
	}
	return result, err
}

func (f *Future[T]) Cancel() {
	f.cancel()
}
