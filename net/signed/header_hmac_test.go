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

package signed_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

	"github.com/corestoreio/pkg/net/signed"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

var _ signed.HeaderParseWriter = (*signed.ContentHMAC)(nil)

func TestHMAC_Write(t *testing.T) {

	w := httptest.NewRecorder()
	sig := signed.ContentHMAC{
		Algorithm: "sha512",
	}
	sig.Write(w, []byte(`Hello Gophers`))

	const wantSig = `sha512 48656c6c6f20476f7068657273`
	if have, want := w.Header().Get(signed.HeaderContentHMAC), wantSig; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}

// 3000000	       494 ns/op	     144 B/op	       5 allocs/op
func BenchmarkHMAC_Write(b *testing.B) {
	const wantSig = `sha512 48656c6c6f20476f7068657273`
	w := httptest.NewRecorder()
	sig := signed.ContentHMAC{
		Algorithm: "sha512",
	}
	s := []byte(`Hello Gophers`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sig.Write(w, s)
	}
	if have, want := w.Header().Get(signed.HeaderContentHMAC), wantSig; have != want {
		b.Errorf("Have: %v Want: %v", have, want)
	}
}

func TestHMAC_Parse(t *testing.T) {
	var newReqHeader = func(value string) *http.Request {
		req := httptest.NewRequest("GET", "http://corestore.io", nil)
		req.Header.Set(signed.HeaderContentHMAC, value)
		return req
	}
	tests := []struct {
		req           *http.Request
		wantAlgorithm string
		wantHMAC      []byte
		wantErrBhf    errors.BehaviourFunc
	}{
		{
			newReqHeader(`sha1 48656c6c6f20476f7068657273`),
			"sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		{
			func() *http.Request {
				req := httptest.NewRequest("GET", "http://corestore.io", strings.NewReader("Hello\nWorld"))
				req.Header.Set("Trailer", signed.HeaderContentHMAC)
				req.Trailer = http.Header{}
				req.Trailer.Set(signed.HeaderContentHMAC, "sha1 48656c6c6f20476f7068657273")
				return req
			}(),
			"sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		{
			newReqHeader(`sha1	48656c6c6f20476f7068657273`),
			"sha1",
			nil,
			errors.IsNotValid, // because tab
		},
		{
			newReqHeader(`sha1 48656c6c6f20476f7068657273xx`),
			"sha1",
			nil,
			errors.IsNotValid, // because tab
		},
		{
			newReqHeader(`sha1 48656c6c6f20476f7068657273`),
			"sha2",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(`48656c6c6f20476f7068657273`),
			"sha2",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(``),
			"sha2",
			nil,
			errors.IsNotFound,
		},
	}
	for i, test := range tests {
		hm := signed.NewContentHMAC(test.wantAlgorithm)
		haveSig, haveErr := hm.Parse(test.req)
		if test.wantErrBhf != nil {
			assert.Nil(t, haveSig, "Index %d", i)
			assert.True(t, test.wantErrBhf(haveErr), "Error: %+v", haveErr)
			//t.Log(haveErr)
			continue
		}
		assert.Exactly(t, test.wantAlgorithm, hm.Algorithm, "Index %d", i)
		assert.Exactly(t, string(test.wantHMAC), string(haveSig), "Index %d", i)
		assert.NoError(t, haveErr, "Index %d: %+v", i, haveErr)
	}
}

// 10000000	       173 ns/op	      16 B/op	       1 allocs/op
func BenchmarkHMAC_Parse(b *testing.B) {

	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("Content-S1gnatur3", `sha1 48656c6c6f20476f7068657273`)

	sig := signed.ContentHMAC{
		Algorithm:  "sha1",
		HeaderName: "Content-S1gnatur3",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sig, err := sig.Parse(req)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if len(sig) < 3 {
			b.Fatal("Invalid length of signature")
		}
	}
}
