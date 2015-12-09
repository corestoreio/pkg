// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ctxmw_test

// Idea: github.com/rs/xaccess Copyright (c) 2015 Olivier Poitrey <rs@dailymotion.com> MIT License

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"bytes"

	"github.com/corestoreio/csfw/net/ctxlog"
	"github.com/corestoreio/csfw/net/ctxmw"
	"github.com/corestoreio/csfw/utils/log"
)

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
	assert.Equal(t, "ok", ctxmw.ResponseStatus(fakeContext{err: nil}, http.StatusOK))
	assert.Equal(t, "canceled", ctxmw.ResponseStatus(fakeContext{err: context.Canceled}, http.StatusOK))
	assert.Equal(t, "timeout", ctxmw.ResponseStatus(fakeContext{err: context.DeadlineExceeded}, http.StatusOK))
	assert.Equal(t, "error", ctxmw.ResponseStatus(fakeContext{err: nil}, http.StatusFound))
}

func TestWithAccessLog(t *testing.T) {
	var buf bytes.Buffer
	defer buf.Reset()

	ctx := ctxlog.WithContext(context.Background(), log.NewStdLogger(log.SetStdWriter(&buf)))

	finalH := ctxhttp.Chain(
		ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusTeapot)
			_, err := w.Write([]byte{'1', '2', '3'})
			time.Sleep(time.Millisecond)
			return err
		}),
		ctxmw.WithAccessLog(),
	)

	r, _ := http.NewRequest("GET", "/gopherine", nil)
	r.RemoteAddr = "127.0.0.1"
	r.Header.Set("User-Agent", "Mozilla")
	r.Header.Set("Referer", "http://rustlang.org")

	w := httptest.NewRecorder()
	if err := finalH.ServeHTTPContext(ctx, w, r); err != nil {
		t.Fatal(err)
	}

	assert.Exactly(t, `123`, w.Body.String())
	assert.Exactly(t, http.StatusTeapot, w.Code)

	want1 := `request error: "" method: "GET" uri: "/gopherine" type: "access" status: "error" status_code: 418 duration:`
	want2 := `size: 3 remote_addr: "127.0.0.1" user_agent: "Mozilla" referer: "http://rustlang.org"`
	assert.Contains(t, buf.String(), want1)
	assert.Contains(t, buf.String(), want2)
}
