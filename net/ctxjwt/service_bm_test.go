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
	"sync"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/store/storenet"
	"golang.org/x/net/context"
)

func bmServeHTTP(b *testing.B, opts ...ctxjwt.Option) {
	jwts, err := ctxjwt.NewService(opts...)
	if err != nil {
		b.Error(err)
	}
	token, _, err := jwts.GenerateToken(scope.DefaultID, 0, map[string]interface{}{
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
	jwtHandler := jwts.WithParseAndValidate()(final)

	req, err := http.NewRequest("GET", "http://abc.xyz", nil)
	if err != nil {
		b.Error(err)
	}
	ctxjwt.SetHeaderAuthorization(req, token)
	w := httptest.NewRecorder()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv)

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
// BenchmarkServeHTTPHMAC-4        	  100000	     15550 ns/op	    4016 B/op	      72 allocs/op Go 1.6.0
func BenchmarkServeHTTPHMAC(b *testing.B) {
	password := []byte(`Rump3lst!lzch3n`)
	bmServeHTTP(b, ctxjwt.WithPassword(scope.DefaultID, 0, password))
}

// BenchmarkServeHTTPHMACSimpleBL-4	  100000	     16037 ns/op	    3808 B/op	      82 allocs/op Go 1.5.0
// BenchmarkServeHTTPHMACSimpleBL-4	  100000	     15765 ns/op	    4016 B/op	      72 allocs/op
func BenchmarkServeHTTPHMACSimpleBL(b *testing.B) {
	bl := ctxjwt.NewBlackListSimpleMap()
	password := []byte(`Rump3lst!lzch3n`)
	bmServeHTTP(b,
		ctxjwt.WithPassword(scope.DefaultID, 0, password),
		ctxjwt.WithBlacklist(bl),
	)
	// b.Logf("Blacklist Items %d", bl.Len())
}

// BenchmarkServeHTTPRSAGenerator-4	    5000	    328220 ns/op	   34544 B/op	     105 allocs/op Go 1.5.0
// BenchmarkServeHTTPRSAGenerator-4	    5000	    327690 ns/op	   34752 B/op	      95 allocs/op Go 1.6.0
func BenchmarkServeHTTPRSAGenerator(b *testing.B) {
	bmServeHTTP(b, ctxjwt.WithRSAGenerator(scope.DefaultID, 0))
}

const benchServeHTTPTokenCount = 100

// A nearly real world test as we're doing parallel requests to the
// middleware. Allocations are not that interesting because they include
// also NewRequest and ResponseRecorder.
// the number of allocs depends on the number of benchServeHTTPTokenCount.
// Map:  BenchmarkServeHTTP_DefaultConfig_BlackList_Parallel-4	    1000	   2333499 ns/op	  468948 B/op	    8077 allocs/op
// FC :  BenchmarkServeHTTP_DefaultConfig_BlackList_Parallel-4	    1000	   2239810 ns/op	  469178 B/op	    8080 allocs/op
// Null: BenchmarkServeHTTP_DefaultConfig_BlackList_Parallel-4	    2000	   2576187 ns/op	  452330 B/op	    7991 allocs/op
func BenchmarkServeHTTP_DefaultConfig_BlackList_Parallel(b *testing.B) {

	jwts := ctxjwt.MustNewService(
		ctxjwt.WithErrorHandler(scope.DefaultID, 0, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
			_, err := ctxjwt.FromContext(ctx)
			if err != nil {
				return err
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		})),
	)
	// below two lines comment out enables the null black list
	jwts.Blacklist = ctxjwt.NewBlackListFreeCache(0)
	//jwts.Blacklist = ctxjwt.NewBlackListSimpleMap()

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		//store.WithStorageConfig(cr), no configuration so config.ScopedGetter is nil
	)
	ctx := store.WithContextProvider(context.Background(), srv) // root context

	var tokens [benchServeHTTPTokenCount]string
	for i := 0; i < benchServeHTTPTokenCount; i++ {

		claim := map[string]interface{}{
			"someKey":          i,
			storenet.ParamName: "de",
		}

		var err error
		tokens[i], _, err = jwts.GenerateToken(scope.WebsiteID, 1, claim) // 1 = website euro
		if err != nil {
			b.Fatal(err)
		}
	}

	final := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
		_, err := ctxjwt.FromContext(ctx)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		return nil
	})
	jwtHandler := jwts.WithParseAndValidate()(final)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var wg sync.WaitGroup
		for pb.Next() {

			wg.Add(1)
			go func() {
				defer wg.Done()

				req, err := http.NewRequest("GET", "http://abc.xyz", nil)
				if err != nil {
					b.Fatal(err)
				}

				for i := 0; i < benchServeHTTPTokenCount; i++ {
					w := httptest.NewRecorder()
					ctxjwt.SetHeaderAuthorization(req, tokens[i])

					if err := jwtHandler.ServeHTTPContext(ctx, w, req); err != nil {
						b.Fatal(err)
					}
					if w.Code != http.StatusUnavailableForLegalReasons {
						b.Fatalf("Response Code want %d; have %d", http.StatusUnavailableForLegalReasons, w.Code)
					}
				}
			}()
		}
		wg.Wait()
	})
}
