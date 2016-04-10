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

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/store/storenet"
	"github.com/corestoreio/csfw/util/csjwt"
	"golang.org/x/net/context"
)

func bmServeHTTP(b *testing.B, opts ...ctxjwt.Option) {
	jwts, err := ctxjwt.NewService(opts...)
	if err != nil {
		b.Error(err)
	}
	cl := csjwt.MapClaims{
		"xfoo": "bar",
		"zfoo": 4711,
	}
	token, err := jwts.NewToken(scope.DefaultID, 0, cl)
	if err != nil {
		b.Error(err)
	}

	final := ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		return nil
	})
	jwtHandler := jwts.WithInitTokenAndStore()(final)

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

var keyBenchmarkHMACPW = ctxjwt.WithKey(scope.DefaultID, 0, csjwt.WithPassword([]byte(`Rump3lst!lzch3n`)))

// BenchmarkServeHTTPHMAC-4        	  100000	     15851 ns/op	    3808 B/op	      82 allocs/op Go 1.5.0
// BenchmarkServeHTTPHMAC-4        	  100000	     15550 ns/op	    4016 B/op	      72 allocs/op Go 1.6.0
func BenchmarkServeHTTPHMAC(b *testing.B) {
	bmServeHTTP(b, keyBenchmarkHMACPW)
}

// BenchmarkServeHTTPHMACSimpleBL-4	  100000	     16037 ns/op	    3808 B/op	      82 allocs/op Go 1.5.0
// BenchmarkServeHTTPHMACSimpleBL-4	  100000	     15765 ns/op	    4016 B/op	      72 allocs/op
func BenchmarkServeHTTPHMACSimpleBL(b *testing.B) {
	bl := ctxjwt.NewBlackListSimpleMap()
	bmServeHTTP(b,
		keyBenchmarkHMACPW,
		ctxjwt.WithBlacklist(bl),
	)
	// b.Logf("Blacklist Items %d", bl.Len())
}

// BenchmarkServeHTTPRSAGenerator-4	    5000	    328220 ns/op	   34544 B/op	     105 allocs/op Go 1.5.0
// BenchmarkServeHTTPRSAGenerator-4	    5000	    327690 ns/op	   34752 B/op	      95 allocs/op Go 1.6.0
func BenchmarkServeHTTPRSAGenerator(b *testing.B) {
	bmServeHTTP(b, ctxjwt.WithKey(scope.DefaultID, 0, csjwt.WithRSAGenerated()))
}

func getReq(b *testing.B, token []byte) *http.Request {
	req, err := http.NewRequest("GET", "http://abc.xyz", nil)
	if err != nil {
		b.Fatal(err)
	}
	ctxjwt.SetHeaderAuthorization(req, token)
	return req
}

// A nearly real world test as we're doing parallel requests to the
// middleware. Allocations are not that interesting because they include
// also NewRequest and ResponseRecorder.
// the number of allocs depends on the number of benchServeHTTPTokenCount.
// BenchmarkServeHTTP_DefaultConfig_BlackList_Parallel-4	  100000	     18993 ns/op	    9471 B/op	     127 allocs/op
// TODO(CS) lower to 80 allocs
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
	//jwts.Blacklist = ctxjwt.NewBlackListFreeCache(0)
	jwts.Blacklist = ctxjwt.NewBlackListSimpleMap()

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		//store.WithStorageConfig(cr), no configuration so config.ScopedGetter is nil
	)
	ctx := store.WithContextProvider(context.Background(), srv) // root context

	token, err := jwts.NewToken(scope.WebsiteID, 1, csjwt.MapClaims{ // 1 = website euro
		"someKey":          2.718281,
		storenet.ParamName: "at",
	})
	if err != nil {
		b.Fatal(err)
	}

	final := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
		_, err := ctxjwt.FromContext(ctx)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusUnavailableForLegalReasons)

		_, st, err := store.FromContextProvider(ctx)
		if err != nil {
			return err
		}
		if st.StoreCode() != "de" && st.StoreCode() != "at" {
			b.Fatalf("Unexpected Store: %s", st.StoreCode())
		}
		return nil
	})
	jwtHandler := jwts.WithInitTokenAndStore()(final)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder() // 3 allocs
			if err := jwtHandler.ServeHTTPContext(ctx, w, getReq(b, token)); err != nil {
				b.Fatal(err)
			}
			if w.Code != http.StatusUnavailableForLegalReasons {
				b.Fatalf("Response Code want %d; have %d", http.StatusUnavailableForLegalReasons, w.Code)
			}
		}
	})
	//b.Log("GC Pause:", gcPause())
}

//func gcPause() time.Duration {
//	runtime.GC()
//	var stats debug.GCStats
//	debug.ReadGCStats(&stats)
//	return stats.Pause[0]
//}

// BenchmarkServeHTTP_DefaultConfig_BlackList_Single-4  	   50000	     32826 ns/op	   11136 B/op	     127 allocs/op
// TODO(CS) lower to 80 allocs
func BenchmarkServeHTTP_DefaultConfig_BlackList_Single(b *testing.B) {

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
	//jwts.Blacklist = ctxjwt.NewBlackListFreeCache(0)
	jwts.Blacklist = ctxjwt.NewBlackListSimpleMap()

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.WebsiteID, "euro"),
		//store.WithStorageConfig(cr), no configuration so config.ScopedGetter is nil
	)
	ctx := store.WithContextProvider(context.Background(), srv) // root context

	claim := csjwt.MapClaims{
		"someKey":          3.14159,
		storenet.ParamName: "at",
	}

	token, err := jwts.NewToken(scope.WebsiteID, 1, claim) // 1 = website euro
	if err != nil {
		b.Fatal(err)
	}

	final := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
		_, err := ctxjwt.FromContext(ctx)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusUnavailableForLegalReasons)

		//if err := jwts.Logout(tk); err != nil {
		//	b.Fatal(err)
		//}

		return nil
	})
	jwtHandler := jwts.WithInitTokenAndStore()(final)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder() // 3 allocs
		if err := jwtHandler.ServeHTTPContext(ctx, w, getReq(b, token)); err != nil {
			b.Fatal(err)
		}
		if w.Code != http.StatusUnavailableForLegalReasons {
			b.Fatalf("Response Code want %d; have %d", http.StatusUnavailableForLegalReasons, w.Code)
		}
	}
}
