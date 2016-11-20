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

package loghttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/loghttp"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/util/cstesting"
)

func TestHTTPRequest_Race(t *testing.T) {

	logBuf := new(log.MutexBuffer)
	lg := logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))

	req := httptest.NewRequest("GET", "http://corestore.io/example", nil)
	req.RemoteAddr = "192.168.0.1"

	hpu := cstesting.NewHTTPParallelUsers(4, 10, 100, time.Microsecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		if have, want := rec.Code, http.StatusAlreadyReported; have != want {
			t.Errorf("Invalid Status Code! Have: %v Want: %v", have, want)
		}
	}
	hpu.ServeHTTP(
		req,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAlreadyReported)
			rec := w.(*httptest.ResponseRecorder)
			if lg.IsDebug() {
				lg.Debug("TestHTTPRequest_Race",
					log.Int("code", rec.Code),
					log.String("user_id", rec.Header().Get(cstesting.HeaderUserID)),
					log.String("loop_cnt", rec.Header().Get(cstesting.HeaderLoopID)),
					log.String("sleep_dur", rec.Header().Get(cstesting.HeaderSleep)),
					loghttp.Request("myReq", r))
			}
		}),
	)

	const search = `myReq: "GET http://corestore.io/example HTTP/1.1\r\n\r\n`
	if have, want := strings.Count(logBuf.String(), search), 40; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}
