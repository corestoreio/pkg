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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorWithCode(t *testing.T) {
	eh := mw.ErrorWithStatusCode(http.StatusTeapot)
	rec := httptest.NewRecorder()
	eh(errors.New("Hello Error World")).ServeHTTP(rec, nil)
	assert.Exactly(t, http.StatusTeapot, rec.Code)
	assert.Contains(t, rec.Body.String(), `Hello Error World`)
	assert.Contains(t, rec.Body.String(), http.StatusText(http.StatusTeapot))
}

func TestMustError(t *testing.T) {
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
	mw.MustError(errors.New("Oh dude, this handler should not be called, but it did.")).ServeHTTP(rec, nil)
	assert.Exactly(t, http.StatusInternalServerError, rec.Code)
}
