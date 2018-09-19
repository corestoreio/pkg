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

package mw_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/util/assert"
)

func TestErrorWithCode(t *testing.T) {
	eh := mw.ErrorWithStatusCode(http.StatusTeapot)
	rec := httptest.NewRecorder()
	eh(errors.New("Hello Error World")).ServeHTTP(rec, nil)
	assert.Exactly(t, http.StatusTeapot, rec.Code)
	assert.Contains(t, rec.Body.String(), `Hello Error World`)
	assert.Contains(t, rec.Body.String(), http.StatusText(http.StatusTeapot))
}

func TestErrorWithPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(string); ok {
				assert.Contains(t, s, `Oh dude, this handler should not be called, but it did.`)
			} else {
				t.Fatalf("recover should contain a string: %#v", r)
			}
		} else {
			t.Fatal("Expected a panic")
		}
	}()
	rec := httptest.NewRecorder()
	mw.ErrorWithPanic(errors.New("Oh dude, this handler should not be called, but it did.")).ServeHTTP(rec, nil)
	assert.Exactly(t, http.StatusInternalServerError, rec.Code)
}

func TestLogErrorWithStatusCode(t *testing.T) {

	var buf bytes.Buffer
	lg := logw.NewLog(logw.WithWriter(&buf), logw.WithLevel(logw.LevelDebug))

	eh := mw.LogErrorWithStatusCode(lg, http.StatusTeapot)
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	eh(errors.AlreadyRefunded.Newf("Invoice already refunded")).ServeHTTP(rec, r)
	assert.Exactly(t, http.StatusTeapot, rec.Code)
	assert.Contains(t, rec.Body.String(), "I'm a teapot\n")
	assert.Contains(t, rec.Body.String(), http.StatusText(http.StatusTeapot))
	assert.Contains(t, buf.String(), `mw.LogErrorWithStatusCode error: "Invoice already refunded" status_code: 418 request: "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"`)
}
