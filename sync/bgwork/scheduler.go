package bgwork

import (
	"context"
	"time"

	"github.com/corestoreio/log"
)

type ScheduleOptions struct {
	Log      log.Logger
	Truncate time.Duration
	Sleep    time.Duration // duration to sleep for the next call of function t
	Backoff  time.Duration // delay the next try in duration x in case of an error, default 5s.
}

// ScheduleWorkAt runs the task t
func ScheduleWorkAt(ctx context.Context, opts ScheduleOptions, t func() error) {
	if opts.Truncate == 0 {
		opts.Truncate = 24 * time.Hour
	}
	if opts.Sleep == 0 {
		opts.Sleep = time.Minute
	}
	if opts.Backoff == 0 {
		opts.Backoff = 5 * time.Second
	}

	nextRun := make(chan struct{})
	var backoff time.Duration = 1

	go func() {
		nextRun <- struct{}{}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-nextRun:
			if err := t(); err != nil {
				nextRetry := opts.Backoff * backoff
				if opts.Log != nil && opts.Log.IsInfo() {
					opts.Log.With(log.Err(err), log.String("nextRetry", nextRetry.String())).Info("Error running scheduled task")
				}
				time.AfterFunc(nextRetry, func() {
					nextRun <- struct{}{}
				})
				backoff *= 2

			} else {
				backoff = 1
				now := time.Now()
				syncT := now.Truncate(opts.Truncate).Add(opts.Truncate + opts.Sleep)
				d := syncT.Sub(now)
				if opts.Log != nil && opts.Log.IsInfo() {
					opts.Log.With(log.String("duration", d.String())).Info("Scheduling next update")
				}
				time.AfterFunc(d, func() {
					nextRun <- struct{}{}
				})
			}
		}
	}
}
