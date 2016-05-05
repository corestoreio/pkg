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
	"github.com/stretchr/testify/assert"
)

func TestDefaultRequestPrefix(t *testing.T) {
	s := mw.RequestIDService{}
	s.Init()
	p := s.NewID()
	assert.Exactly(t, "-1", p[len(p)-2:])
	assert.Contains(t, p, "/")
}

func testWithRequestID(t *testing.T, gen mw.RequestIDGenerator) {
	var opt mw.Option
	if gen != nil {
		opt = mw.SetRequestIDGenerator(gen)
	}

	finalCH := mw.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := w.Header().Get(mw.RequestIDHeader)
		assert.Exactly(t, "-2", id[len(id)-2:])
		assert.Contains(t, id, "/")

	}, mw.WithRequestID(opt))

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	finalCH.ServeHTTP(w, r)
}

func TestWithRequestIDDefault(t *testing.T) {
	testWithRequestID(t, nil)
}

type testGenerator struct{}

func (testGenerator) Init() {

}
func (testGenerator) NewID() string {
	return "goph/er-2"
}

func TestWithRequestIDCustom(t *testing.T) {
	testWithRequestID(t, testGenerator{})
}

// BenchmarkWithRequestID-4	 3000000	       432 ns/op	      64 B/op	       3 allocs/op
func BenchmarkWithRequestID(b *testing.B) {

	finalCH := mw.Chain(func(w http.ResponseWriter, r *http.Request) {
		id := w.Header().Get(mw.RequestIDHeader)
		if id == "" {
			b.Fatal("id is empty")
		}
	}, mw.WithRequestID())

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
