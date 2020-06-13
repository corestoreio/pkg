package bgwork

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestScheduleWorkAt(t *testing.T) {
	tests := []struct {
		name      string
		ctxD      time.Duration
		opts      ScheduleOptions
		err       error
		wantCount int
	}{
		{
			name: "err backoff",
			ctxD: time.Second,
			opts: ScheduleOptions{
				Sleep:   time.Millisecond * 20,
				Backoff: time.Millisecond,
			},
			err:       errors.New("err"),
			wantCount: 4,
		},
		{
			name: "single invocation",
			ctxD: time.Second,
			opts: ScheduleOptions{
				Sleep: time.Millisecond * 20,
			},
			wantCount: 1,
		},
		{
			name: "more invocation",
			ctxD: time.Second,
			opts: ScheduleOptions{
				Truncate: time.Nanosecond,
				Sleep:    time.Millisecond * 4,
			},
			wantCount: 4,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var c int
			cC := make(chan int)
			task := func() error {
				c++
				cC <- c
				return tt.err
			}

			ctx := ctxCancelAfter(tt.ctxD)()
			go ScheduleWorkAt(ctx, tt.opts, task)

			done := make(chan struct{})

			go func() {
				defer close(done)
				for c := range cC {
					if c == tt.wantCount {
						break
					}
				}
			}()

			select {
			case <-ctx.Done():
				t.Fatalf("Desired count %d not reached in time", tt.wantCount)
			case <-done:
				return
			}
		})
	}
}

func ctxCancelAfter(d time.Duration) func() context.Context {
	return func() context.Context {
		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(d, cancel)
		return ctx
	}
}
