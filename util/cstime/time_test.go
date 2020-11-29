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

package cstime_test

import (
	"errors"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/cstime"
)

func TestParseTimeStrict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		layout  string
		value   string
		wantErr error
		want    string
	}{
		{"1/2/06", "11/31/15", errors.New("invalid time: \"11/31/15\""), ""},
		{"1/2/06", "11/30/15", nil, "2015-11-30 00:00:00 +0000 UTC"},
	}
	for _, test := range tests {

		tt, err := cstime.ParseTimeStrict(test.layout, test.value)
		if test.wantErr != nil {
			assert.Error(t, err, "Test %v", test)
			continue
		}
		assert.NoError(t, err, "Test %v", test)
		assert.Equal(t, test.want, tt.String(), "Test %v", test)
	}
}

func TestRandomTicker(t *testing.T) {
	t.Parallel()

	min := time.Duration(10)
	max := time.Duration(20)

	// tick can take a little longer since we're not adjusting it to account for
	// processing.
	precision := time.Duration(4)

	rt := cstime.NewRandomTicker(min*time.Millisecond, max*time.Millisecond, time.Millisecond)
	for i := 0; i < 5; i++ {
		t0 := time.Now()
		t1 := <-rt.C
		td := t1.Sub(t0)
		if td < min*time.Millisecond {
			t.Fatalf("tick was shorter than expected: %s", td)
		} else if td > (max+precision)*time.Millisecond {
			t.Fatalf("tick was longer than expected: %s", td)
		}
	}
	rt.Stop()
	time.Sleep((max + precision) * time.Millisecond)
	select {
	case v, ok := <-rt.C:
		if ok || !v.IsZero() {
			t.Fatal("ticker did not shut down")
		}
	default:
		t.Fatal("expected to receive close channel signal")
	}
}

func TestNewRandomTicker(t *testing.T) {
	rt := cstime.NewRandomTicker(10*time.Hour, 120*time.Minute, time.Millisecond)
	for i := 0; i < 10; i++ {
		ni := rt.NextInterval()
		t.Logf("%d: %s", i, ni)
	}
	rt.Stop()
}
