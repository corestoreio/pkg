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

package accesslog_test

// Idea: github.com/rs/xaccess Copyright (c) 2015 Olivier Poitrey <rs@dailymotion.com> MIT License

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corestoreio/cspkg/net/accesslog"
	"github.com/corestoreio/cspkg/net/mw"
	"github.com/corestoreio/log/logw"
	"github.com/rs/xstats"
	"github.com/stretchr/testify/assert"
)

var _ xstats.XStater = (*accesslog.BlackholeXStat)(nil)

type fakeContext struct {
	err error
}

func (c fakeContext) Err() error {
	return c.err
}

func (c fakeContext) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), true
}

func (c fakeContext) Done() <-chan struct{} {
	return make(chan struct{})
}

func (c fakeContext) Value(key interface{}) interface{} {
	return nil
}

func TestResponseStatus(t *testing.T) {
	assert.Equal(t, "ok", accesslog.ResponseStatus(fakeContext{err: nil}, http.StatusOK))
	assert.Equal(t, "canceled", accesslog.ResponseStatus(fakeContext{err: context.Canceled}, http.StatusOK))
	assert.Equal(t, "timeout", accesslog.ResponseStatus(fakeContext{err: context.DeadlineExceeded}, http.StatusOK))
	assert.Equal(t, "error", accesslog.ResponseStatus(fakeContext{err: nil}, http.StatusFound))
}

func TestWithAccessLog(t *testing.T) {
	var buf bytes.Buffer
	defer buf.Reset()

	testLog := logw.NewLog(logw.WithWriter(&buf))

	finalH := mw.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte{'1', '2', '3'})
			time.Sleep(time.Millisecond)
		}),
		accesslog.WithAccessLog(accesslog.BlackholeXStat{}, testLog),
	)

	r, _ := http.NewRequest("GET", "/gopherine", nil)
	r.RemoteAddr = "127.0.0.1"
	r.Header.Set("User-Agent", "Mozilla")
	r.Header.Set("Referer", "http://rustlang.org")

	w := httptest.NewRecorder()
	finalH.ServeHTTP(w, r)

	assert.Exactly(t, `123`, w.Body.String())
	assert.Exactly(t, http.StatusTeapot, w.Code)

	want1 := `method: "GET" uri: "/gopherine" type: "access" status: "error" status_code: 418 duration:`
	want2 := `size: 3 remote_addr: "127.0.0.1" user_agent: "Mozilla" referer: "http://rustlang.org"`
	assert.Contains(t, buf.String(), want1)
	assert.Contains(t, buf.String(), want2)
}
