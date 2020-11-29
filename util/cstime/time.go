// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cstime

import (
	"fmt"
	"math/rand"
	"time"
)

// ParseTimeStrict parses a formatted string and returns the time value it
// represents. The output is identical to time.Parse except it returns an
// error for strings that don't format to the input value.
//
// An example where the output differs from time.Parse would be:
// parseTimeStrict("1/2/06", "11/31/15")
//
// - time.Parse returns "2015-12-01 00:00:00 +0000 UTC"
//
// - ParseTimeStrict returns an error
func ParseTimeStrict(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return t, err
	}
	if t.Format(layout) != value {
		return t, fmt.Errorf("invalid time: %q", value)
	}
	return t, nil
}

// Source of the random ticker: github.com/fwojciec/clock

// RandomTicker is similar to time.Ticker but ticks at random intervals between
// the min and max duration values (stored internally as int64 nanosecond
// counts).
type RandomTicker struct {
	C            chan time.Time
	stopc        chan chan struct{}
	rand         *lockedSource
	min          int64
	max          int64
	durationKind time.Duration
}

// NewRandomTicker returns a pointer to an initialized instance of the
// RandomTicker. Min and max are durations of the shortest and longest allowed
// ticks. Ticker will run in a goroutine until explicitly stopped.
// durationKind must be one of the following constants: time.Nanosecond,
// time.Microsecond, time.Millisecond, time.Second, time.Minute, time.Hour or it
// panics.
// AFAIK: It uses time.NewTimer under the hood which increases memory usage in
// tight loops a lot as the timer cannot be garbage collected.
func NewRandomTicker(min, max, durationKind time.Duration) *RandomTicker {
	if min > max {
		min, max = max, min
	}

	rt := &RandomTicker{
		C:            make(chan time.Time),
		stopc:        make(chan chan struct{}),
		durationKind: durationKind,
		rand:         &lockedSource{src: rand.New(rand.NewSource(time.Now().UnixNano()))},
	}
	min = min.Truncate(durationKind)
	max = max.Truncate(durationKind)
	switch durationKind {
	case time.Nanosecond:
		rt.min = min.Nanoseconds()
		rt.max = max.Nanoseconds()
	case time.Microsecond:
		rt.min = min.Microseconds()
		rt.max = max.Microseconds()
	case time.Millisecond:
		rt.min = min.Milliseconds()
		rt.max = max.Milliseconds()
	case time.Second:
		rt.min = int64(min.Seconds())
		rt.max = int64(max.Seconds())
	case time.Minute:
		rt.min = int64(min.Minutes())
		rt.max = int64(max.Minutes())
	case time.Hour:
		rt.min = int64(min.Hours())
		rt.max = int64(max.Hours())
	default:
		panic(fmt.Sprintf("durationKind can only be time.Nanosecond, time.Microsecond, time.Millisecond, time.Second, time.Minute or time.Hour but you gave me: %s", durationKind))
	}
	go rt.loop()
	return rt
}

// Stop terminates the ticker goroutine and closes the C channel.
func (rt *RandomTicker) Stop() {
	c := make(chan struct{})
	rt.stopc <- c
	<-c
}

func (rt *RandomTicker) loop() {
	defer close(rt.C)
	t := time.NewTimer(rt.NextInterval())
	for {
		// either a stop signal or a timeout
		select {
		case c := <-rt.stopc:
			t.Stop()
			close(c)
			return
		case <-t.C:
			select {
			case rt.C <- time.Now():
				t.Stop()
				t = time.NewTimer(rt.NextInterval())
			default:
				// there could be noone receiving...
			}
		}
	}
}

// NextInterval returns a random duration between max and min in the provided
// durationKind.
func (rt *RandomTicker) NextInterval() time.Duration {
	interval := rand.Int63n(rt.max-rt.min) + rt.min
	return time.Duration(interval) * rt.durationKind
}
