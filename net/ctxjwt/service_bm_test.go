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

package ctxjwt_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"golang.org/x/net/context"
)

// BenchmarkAuthorizationHMAC-4	  100000	     20215 ns/op	    5552 B/op	     105 allocs/op
func BenchmarkAuthorizationHMAC(b *testing.B) {

	/*
		that benchmark gives a false impression because we're also
		measuring the NewRequest/Response creation ...
	*/

	password := []byte(`Rump3lst!lzch3n`)
	service, err := ctxjwt.NewService(ctxjwt.WithPassword(password))
	if err != nil {
		b.Error(err)
	}
	tok, _, err := service.GenerateToken(map[string]interface{}{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	if err != nil {
		b.Error(err)
	}

	final := ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})
	jwtHandler := service.WithParseAndValidate()(final)

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// <15 allocs>
		req, err := http.NewRequest("GET", "http://auth.xyz", nil)
		if err != nil {
			b.Error(err)
		}
		req.Header.Set("Authorization", "Bearer "+tok)
		w := httptest.NewRecorder()
		//</>

		jwtHandler.ServeHTTPContext(ctx, w, req)

	}
}
