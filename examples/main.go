package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rrfeng/go-future"
)

func MyFunc1() (string, error) {
	time.Sleep(time.Second)
	return "this is MyFunc1", nil
}

func MyFunc2(arg string) (string, error) {
	time.Sleep(time.Second)
	return "this is MyFunc2, arg: " + arg, nil
}

func main() {
	f1 := future.Async(MyFunc1)
	// future.Cancel() should be called for prevent context memory leak,
	// though future.Wait() will call cancel() internally, but we cannot
	// assume what will happen before you call future.Wait()
	defer f1.Cancel()

	input := "test args"
	ctx, cancel := context.WithCancel(context.Background())
	f2 := future.AsyncWithContext(ctx, func() (string, error) { return MyFunc2(input) })
	defer cancel()

	r1, e1 := f1.Await()
	r2, e2 := f2.Await()
	fmt.Printf("MyFunc1 result: %v, error: %v\n", r1, e1)
	fmt.Printf("MyFunc2 result: %v, error: %v\n", r2, e2)
}
