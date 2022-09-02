package future

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestInt(t *testing.T) {
	type args struct {
		fn func() (int, error)
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test int",
			args: args{fn: func() (int, error) { time.Sleep(time.Second); return 1, nil }},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Async(tt.args.fn).Await(); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFuncArgs(t *testing.T) {
	type args struct {
		fn func(i int) (int, error)
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test func args",
			args: args{fn: func(i int) (int, error) { time.Sleep(time.Second); return i, nil }},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := 1
			if got, err := Async(func() (int, error) { return tt.args.fn(input) }).Await(); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCancel(t *testing.T) {
	type args struct {
		fn func(i int) (int, error)
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr string
	}{
		{
			name:    "test cancel",
			args:    args{fn: func(i int) (int, error) { time.Sleep(time.Second); return i, nil }},
			want:    0,
			wantErr: "context canceled",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := 1
			f := Async(func() (int, error) { return tt.args.fn(input) })
			f.Cancel()

			got, err := f.Await()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
			if err == nil || err.Error() != tt.wantErr {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeout(t *testing.T) {
	type args struct {
		fn func(i int) (int, error)
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr string
	}{
		{
			name:    "test timeout",
			args:    args{fn: func(i int) (int, error) { time.Sleep(time.Millisecond); return 2, nil }},
			want:    2,
			wantErr: "context deadline exceeded",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*2)
			f := AsyncWithContext(ctx, func() (int, error) { return tt.args.fn(1) })

			got, err := f.Await()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("result = %v, want %v", got, tt.want)
			}

			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
