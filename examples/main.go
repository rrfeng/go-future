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
	f1 := future.New(MyFunc1)
	input := "test args"
	ctx, _ := context.WithCancel(context.Background())
	f2 := future.NewWithContext(ctx, func() (string, error) { return MyFunc2(input) })

	r1, e1 := f1.Wait()
	r2, e2 := f2.Wait()
	fmt.Printf("MyFunc1 result: %v, error: %v\n", r1, e1)
	fmt.Printf("MyFunc2 result: %v, error: %v\n", r2, e2)
}
