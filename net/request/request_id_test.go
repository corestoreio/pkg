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

package request_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/net/request"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/corestoreio/pkg/util/assert"
)

var _ mw.Middleware = (&request.ID{}).With() // test if function signature matches

func TestDefaultRequestPrefix(t *testing.T) {
	t.Parallel()

	var idGen = &request.ID{}

	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		id := w.Header().Get(request.HeaderIDKeyName)
		assert.Contains(t, id, "/")

	}, idGen.With())
	const regex = ".+/[A-Za-z0-9]+-[0-9]+"
	matchr := regexp.MustCompile(regex)

	bgwork.Wait(10, func(idx int) {

		req := httptest.NewRequest("GET", "/", nil)

		hpu := cstesting.NewHTTPParallelUsers(5, 10, 500, time.Millisecond)
		hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
			id := rec.Header().Get(request.HeaderIDKeyName)
			if !matchr.MatchString(id) {
				panic(fmt.Sprintf("ID %q does not match %q", id, regex))
			}
		}
		hpu.ServeHTTP(req, finalCH)
	})

	assert.Exactly(t, 500, int(*idGen.Count), "ID counts do not match")

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
