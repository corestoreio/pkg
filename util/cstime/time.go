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

// RandTicker is similar to time.Ticker but ticks at random intervals between
// the min and max duration values.
type RandTicker struct {
	C     chan time.Time
	stopc chan struct{}
	min   time.Duration
	max   time.Duration
}

// NewRandTicker creates a new RandTicker. Min and max are durations of the
// shortest and longest allowed ticks. If max is < min, then max will be min.
// Ticker will run in a goroutine until explicitly stopped.
func NewRandTicker(min, max time.Duration) *RandTicker {
	if max < min {
		min, max = max, min
	}
	rt := &RandTicker{
		C:     make(chan time.Time),
		stopc: make(chan struct{}),
		min:   min,
		max:   max,
	}
	go rt.loop()
	return rt
}

// Stop terminates the ticker goroutine and closes the C channel.
func (rt *RandTicker) Stop() {
	close(rt.stopc)
	close(rt.C)
}

func (rt *RandTicker) loop() {
	t := time.NewTimer(rt.nextRandDur())
	defer t.Stop()
	for {
		// either a stop signal or a timeout
		select {
		case <-rt.stopc:
			return
		case <-t.C:
			select {
			case rt.C <- time.Now():
				t.Reset(rt.nextRandDur())
			default:
				// nothing to do
			}
		}
	}
}

func (rt *RandTicker) nextRandDur() time.Duration {
	interval := rand.Int63n(int64(rt.max-rt.min)) + int64(rt.min)
	return time.Duration(interval)
}
