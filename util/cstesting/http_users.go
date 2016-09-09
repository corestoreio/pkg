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

package cstesting

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// HTTPParallelUsers allows to run parallel and concurrent calls to a given
// http.Handler.
type HTTPParallelUsers struct {
	// Users or also known as number of threads
	Users int
	// Loops each user runs these loops
	Loops int
	// RampUpPeriod time to take to generate to the full request force. The
	// duration calculates: RampUpPeriod * Interval
	RampUpPeriod int
	// Interval an enum set of time.Nanosecond, time.Microsecond, time.Millisecond,
	// time.Second, time.Minute, time.Hour.
	Interval time.Duration
	// AssertResponse provides the possibility to check the written data after each
	// request.
	AssertResponse func(*httptest.ResponseRecorder)
}

// Header* got set within an user iteration to allow you to identify a request.
const (
	HeaderUserID = "X-Test-User-ID"
	HeaderLoopID = "X-Test-Loop-ID"
	HeaderSleep  = "X-Test-Sleep"
)

// NewHTTPParallelUsers initializes a new request producer. Users means the
// total amount of parallel users. Each user can execute a specific loopsPerUser
// count. The rampUpPeriod defines the total runtime of the test and the period
// it takes to produce the finally total amount of parallel requests. The
// interval applies to the exported constants of the time package:
// time.Nanosecond, time.Microsecond, time.Millisecond, time.Second,
// time.Minute, time.Hour. The total runtime calculates rampUpPeriod * interval.
// Every (rampUpPeriod / users) a new user starts with its requests. Each user
// request sleeps a specific equal time until the test ends. With the last
// started user the maximum amount of parallel requests will be reached.
func NewHTTPParallelUsers(users, loopsPerUser, rampUpPeriod int, interval time.Duration) HTTPParallelUsers {
	switch interval {
	case time.Nanosecond, time.Microsecond, time.Millisecond, time.Second, time.Minute, time.Hour:
	// ok
	default:
		panic(fmt.Sprintf("Unknown interval %s. Only allowed time.Nanosecond, time.Microsecond, etc", interval))
	}

	return HTTPParallelUsers{
		Users:        users,
		Loops:        loopsPerUser,
		RampUpPeriod: rampUpPeriod,
		Interval:     interval,
	}
}

func (hpu HTTPParallelUsers) sleepPerServeHTTP(userID int) time.Duration {
	d := (float64(hpu.RampUpPeriod) / float64(userID) / float64(hpu.Loops)) * float64(hpu.Interval)
	return time.Duration(d)
}

// serve runs the benchmark. r or rf can be nil, but not both.
func (hpu HTTPParallelUsers) serve(rf func() *http.Request, h http.Handler) {

	var user = func(wg *sync.WaitGroup, userID int) {
		for i := 1; i <= hpu.Loops; i++ {
			sl := hpu.sleepPerServeHTTP(userID)
			// go func(sl time.Duration) { // if go, then add WaitGroup
			w := httptest.NewRecorder()
			w.Header().Set(HeaderUserID, strconv.Itoa(userID))
			w.Header().Set(HeaderLoopID, strconv.Itoa(i))
			w.Header().Set(HeaderSleep, sl.String())

			h.ServeHTTP(w, rf())
			if hpu.AssertResponse != nil {
				hpu.AssertResponse(w)
			}
			// }(sl)
			time.Sleep(sl)
		}
		wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(hpu.Users)
	var startDelay = hpu.RampUpPeriod / hpu.Users
	var delay = new(int32)
	for j := 1; j <= hpu.Users; j++ {
		go func(userID int) {
			if startDelay > 0 && userID > 1 {
				atomic.AddInt32(delay, int32(startDelay))
				cd := atomic.LoadInt32(delay)
				if int(cd) <= hpu.RampUpPeriod {
					time.Sleep(time.Duration(cd) * hpu.Interval)
				}
			}
			go user(&wg, userID)
		}(j)
	}
	wg.Wait()
}

// ServeHTTP starts the testing and the request gets called with http.Handler.
// You might run into a race condition when trying to add a request body (an
// io.ReadCloser), because multiple reads and writes into the buffer. Use the
// function ServeHTTPNewRequest() if you need for each call to http.Handler a
// new request object.
func (hpu HTTPParallelUsers) ServeHTTP(r *http.Request, h http.Handler) {
	// should be refactored but for now quite ok
	// 10 threads, 20 seconds ramp-up - start with 1 user, each 2 seconds 1 user added
	hpu.serve(func() *http.Request {
		return r
	}, h)
}

// ServeHTTPNewRequest same as ServeHTTP() but creates for each iteration a new
// fresh request which will be passed to http.Handler. Does not trigger a race
// condition.
func (hpu HTTPParallelUsers) ServeHTTPNewRequest(rf func() *http.Request, h http.Handler) {
	hpu.serve(rf, h)
}
