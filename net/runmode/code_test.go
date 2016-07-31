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

package runmode_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/corestoreio/csfw/net/runmode"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

var _ runmode.StoreCodeProcesser = (*runmode.ProcessStoreCode)(nil)

func TestProcessStoreCode_FromRequest(t *testing.T) {

	var getCookieRequest = func(c *http.Cookie) *http.Request {
		rootRequest, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("Root request error: %s", err)
		}
		if c != nil {
			rootRequest.AddCookie(c)
		}
		return rootRequest
	}

	var getGETRequest = func(kv ...string) *http.Request {
		reqURL := "http://corestore.io/"
		var uv url.Values
		if len(kv)%2 == 0 {
			uv = make(url.Values)
			for i := 0; i < len(kv); i = i + 2 {
				uv.Set(kv[i], kv[i+1])
			}
			reqURL = reqURL + "?" + uv.Encode()
		}
		rootRequest, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			t.Fatalf("Root request error: %s", err)
		}
		return rootRequest
	}

	tests := []struct {
		req      *http.Request
		wantCode string
	}{
		{
			getCookieRequest(&http.Cookie{Name: runmode.FieldName, Value: "dede"}),
			"dede",
		},
		{
			getCookieRequest(&http.Cookie{Name: runmode.FieldName, Value: "ded'e"}),
			"",
		},
		{
			getCookieRequest(&http.Cookie{Name: "invalid", Value: "dede"}),
			"",
		},

		{
			getGETRequest(runmode.URLFieldName, "dede"),
			"dede",
		},
		{
			getGETRequest(runmode.URLFieldName, "ded¢e"),
			"",
		},
		{
			getGETRequest("invalid", "dede"),
			"",
		},
	}
	for i, test := range tests {
		c := &runmode.ProcessStoreCode{URLFieldName: runmode.URLFieldName, FieldName: runmode.FieldName}
		code := c.FromRequest(test.req)
		assert.Exactly(t, test.wantCode, code, "Index %d", i)
	}
}

var benchmarkProcessStoreCode_FromRequest_Cookie string

//BenchmarkProcessStoreCode_FromRequest_Cookie/Found-4         	  500000	      3047 ns/op	     296 B/op	       3 allocs/op
//BenchmarkProcessStoreCode_FromRequest_Cookie/NotFound-4      	10000000	       110 ns/op	       0 B/op	       0 allocs/op
func BenchmarkProcessStoreCode_FromRequest_Cookie(b *testing.B) {
	c := &runmode.ProcessStoreCode{URLFieldName: runmode.URLFieldName, FieldName: runmode.FieldName}
	const defaultCookie = `mage-translation-storage=%7B%7D; mage-translation-file-version=%7B%7D; mage-cache-storage=%7B%7D; mage-cache-storage-section-invalidation=%7B%7D; mage-cache-sessid=true; PHPSESSID=ogb786ncug3gunsnoevjem7n32; form_key=6DnQ2Xiy2oMpp7FB`

	b.Run("Found", func(b *testing.B) {

		req := httptest.NewRequest("GET", "https://corestoreio.io?a=b", nil)
		req.Header.Set("Cookie", defaultCookie)
		req.AddCookie(&http.Cookie{Name: runmode.FieldName, Value: "dede"})

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			benchmarkProcessStoreCode_FromRequest_Cookie = c.FromRequest(req)
			if benchmarkProcessStoreCode_FromRequest_Cookie == "" {
				b.Fatal("benchmarkProcessStoreCode_FromRequest_Cookie is empty")
			}
		}
	})

	b.Run("NotFound", func(b *testing.B) {
		req := httptest.NewRequest("GET", "https://corestoreio.io?c=d", nil)
		req.Header.Set("Cookie", defaultCookie)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			benchmarkProcessStoreCode_FromRequest_Cookie = c.FromRequest(req)
			if benchmarkProcessStoreCode_FromRequest_Cookie != "" {
				b.Fatal("benchmarkProcessStoreCode_FromRequest_Cookie is NOT empty")
			}
		}
	})

}
