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

package ctxjwt_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"golang.org/x/net/context"
)

func bmServeHTTP(b *testing.B, opts ...ctxjwt.Option) {
	service, err := ctxjwt.NewService(opts...)
	if err != nil {
		b.Error(err)
	}
	token, _, err := service.GenerateToken(map[string]interface{}{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	if err != nil {
		b.Error(err)
	}

	final := ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		return nil
	})
	jwtHandler := service.WithParseAndValidate()(final)

	req, err := http.NewRequest("GET", "http://abc.xyz", nil)
	if err != nil {
		b.Error(err)
	}
	ctxjwt.SetHeaderAuthorization(req, token)
	w := httptest.NewRecorder()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jwtHandler.ServeHTTPContext(ctx, w, req); err != nil {
			b.Error(err)
		}
		if w.Code != http.StatusTeapot {
			b.Errorf("Response Code want %d; have %d", http.StatusTeapot, w.Code)
		}
	}
}

// BenchmarkServeHTTPHMAC-4        	  100000	     15851 ns/op	    3808 B/op	      82 allocs/op Go 1.5.0
func BenchmarkServeHTTPHMAC(b *testing.B) {
	password := []byte(`Rump3lst!lzch3n`)
	bmServeHTTP(b, ctxjwt.WithPassword(password))
}

// BenchmarkServeHTTPHMACSimpleBL-4	  100000	     16037 ns/op	    3808 B/op	      82 allocs/op Go 1.5.0
func BenchmarkServeHTTPHMACSimpleBL(b *testing.B) {
	bl := ctxjwt.NewSimpleMapBlackList()
	password := []byte(`Rump3lst!lzch3n`)
	bmServeHTTP(b,
		ctxjwt.WithPassword(password),
		ctxjwt.WithBlacklist(bl),
	)
	b.Logf("Blacklist Items %d", bl.Len())
}

// BenchmarkServeHTTPRSAGenerator-4	    5000	    328220 ns/op	   34544 B/op	     105 allocs/op Go 1.5.0
func BenchmarkServeHTTPRSAGenerator(b *testing.B) {
	bmServeHTTP(b, ctxjwt.WithRSAGenerator())
}
