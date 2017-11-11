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

package request_test

import (
	"net/http/httptest"
	"testing"

	"net/http"
	"regexp"
	"time"

	"github.com/corestoreio/cspkg/net/mw"
	"github.com/corestoreio/cspkg/net/request"
	"github.com/corestoreio/cspkg/util/cstesting"
	"github.com/stretchr/testify/assert"
)

var idGen = &request.ID{}

func TestDefaultRequestPrefix(t *testing.T) {
	t.Parallel()
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		id := w.Header().Get(request.HeaderIDKeyName)
		assert.Contains(t, id, "/")

	}, idGen.With())

	req := httptest.NewRequest("GET", "/", nil)

	const regex = ".+/[A-Za-z0-9]+-[0-9]+"
	matchr := regexp.MustCompile(regex)
	hpu := cstesting.NewHTTPParallelUsers(5, 10, 500, time.Millisecond)
	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
		id := rec.Header().Get(request.HeaderIDKeyName)
		assert.True(t, matchr.MatchString(id), "ID %q does not match %q", id, regex)
	}
	hpu.ServeHTTP(req, finalCH)
	assert.Exactly(t, 50, int(*idGen.Count))
}

func BenchmarkWithRequestID(b *testing.B) {
	id := &request.ID{}
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		id := w.Header().Get(request.HeaderIDKeyName)
		if id == "" {
			b.Fatal("id is empty")
		}
	}, id.With())

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		finalCH.ServeHTTP(w, r)
	}
}
