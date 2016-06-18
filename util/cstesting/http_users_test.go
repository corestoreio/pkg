// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package cstesting_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
)

func TestHTTPParallelUsers_WrongInterval(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok {
				if have, want := s, "Unknown interval 2s. Only allowed time.Nanosecond, time.Microsecond, etc"; have != want {
					t.Errorf("Have: %v Want: %v", have, want)
				}

			} else {
				t.Fatalf("Expecting a string in the panic; Got: %#v", r)
			}
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	_ = cstesting.NewHTTPParallelUsers(1, 1, 1, time.Second*2)
}

func TestHTTPParallelUsers_Single(t *testing.T) {
	tg := cstesting.NewHTTPParallelUsers(1, 1, 1, time.Nanosecond)
	req := httptest.NewRequest("GET", "http://corestore.io", nil)

	var reqCount int
	tg.ServeHTTP(req, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no race here because one single iteration
		reqCount++
	}))
	if have, want := reqCount, 1; have != want {
		t.Errorf("Request count mismatch! Have: %v Want: %v", have, want)
	}
}

func TestHTTPParallelUsers_Long(t *testing.T) {
	startTime := time.Now()
	const (
		users        = 4
		loops        = 10
		rampUpPeriod = 2
	)
	tg := cstesting.NewHTTPParallelUsers(users, loops, rampUpPeriod, time.Second)
	req := httptest.NewRequest("GET", "http://corestore.io", nil)

	tg.AssertResponse = func(rec *httptest.ResponseRecorder) {
		assert.NotEmpty(t, rec.Header().Get(cstesting.HeaderUserID))
		assert.NotEmpty(t, rec.Header().Get(cstesting.HeaderLoopID))
		assert.NotEmpty(t, rec.Header().Get(cstesting.HeaderSleep))
	}

	var reqCount = new(int32)
	tg.ServeHTTP(req, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//t.Logf("UserID %s LoopID %s Sleeping %s",
		//	rec.Header().Get(cstesting.HeaderUserID),
		//	rec.Header().Get(cstesting.HeaderLoopID),
		//	rec.Header().Get(cstesting.HeaderSleep),
		//)
		atomic.AddInt32(reqCount, 1)
	}))

	//t.Logf("Users %d Loops %d, RampUp %d", users, loops, rampUpPeriod)

	if have, want := *reqCount, int32(users*loops); have != want {
		t.Errorf("Request count mismatch! Have: %v Want: %v", have, want)
	}

	if have, want := int(time.Since(startTime).Seconds()), rampUpPeriod; have != want {
		t.Errorf("Test Running Time is weird! Have: %v Want: %v", have, want)
	}
}
