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
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/blacklist"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"golang.org/x/net/context"
)

const benchTokenCount = 100

func benchBlackList(b *testing.B, bl ctxjwt.Blacklister) {
	jwts := ctxjwt.MustNewService()
	var tokens [benchTokenCount][]byte

	for i := 0; i < benchTokenCount; i++ {
		claim := jwtclaim.Map{
			"someKey": i,
		}
		tk, err := jwts.NewToken(scope.Default, 0, claim)
		if err != nil {
			b.Fatal(err)
		}
		tokens[i] = tk.Raw
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < benchTokenCount; i++ {
				if err := bl.Set(tokens[i], time.Minute); err != nil {
					b.Fatal(err)
				}
				if bl.Has(tokens[i]) == false {
					b.Fatalf("Cannot find token %s with index %d", tokens[i], i)
				}
			}
		}
	})
}

// BenchmarkBlackListMap_Parallel-4      	    2000	    586726 ns/op	   31686 B/op	     200 allocs/op
func BenchmarkBlackListMap_Parallel(b *testing.B) {
	bl := blacklist.NewBlackListSimpleMap()
	benchBlackList(b, bl)
}

// BenchmarkBlackListFreeCache_Parallel-4	   30000	     59542 ns/op	   31781 B/op	     300 allocs/op
func BenchmarkBlackListFreeCache_Parallel(b *testing.B) {
	bl := blacklist.NewBlackListFreeCache(0)
	benchBlackList(b, bl)
}

func bmServeHTTP(b *testing.B, opts ...ctxjwt.Option) {
	jwts, err := ctxjwt.NewService(opts...)
	if err != nil {
		b.Error(err)
	}
	cl := jwtclaim.Map{
		"xfoo": "bar",
		"zfoo": 4711,
	}
	token, err := jwts.NewToken(scope.Default, 0, cl)
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
	ctxjwt.SetHeaderAuthorization(req, token.Raw)
	w := httptest.NewRecorder()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
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

var keyBenchmarkHMACPW = ctxjwt.WithKey(scope.Default, 0, csjwt.WithPassword([]byte(`Rump3lst!lzch3n`)))

func BenchmarkServeHTTPHMAC(b *testing.B) {
	bmServeHTTP(b, keyBenchmarkHMACPW)
}

func BenchmarkServeHTTPHMACSimpleBL(b *testing.B) {
	bl := blacklist.NewBlackListSimpleMap()
	bmServeHTTP(b,
		keyBenchmarkHMACPW,
		ctxjwt.WithBlacklist(bl),
	)
	// b.Logf("Blacklist Items %d", bl.Len())
}

func BenchmarkServeHTTPRSAGenerator(b *testing.B) {
	bmServeHTTP(b, ctxjwt.WithKey(scope.Default, 0, csjwt.WithRSAGenerated()))
}

func BenchmarkServeHTTP_DefaultConfig_BlackList_Parallel(b *testing.B) {
	jwtHandler, ctx, token := benchmarkServeHTTPDefaultConfigBlackListSetup(b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchmarkServeHTTPDefaultConfigBlackListLoop(b, jwtHandler, ctx, token)
		}
	})
	//b.Log("GC Pause:", gcPause())
}

func BenchmarkServeHTTP_DefaultConfig_BlackList_Single(b *testing.B) {
	jwtHandler, ctx, token := benchmarkServeHTTPDefaultConfigBlackListSetup(b)
	for i := 0; i < b.N; i++ {
		benchmarkServeHTTPDefaultConfigBlackListLoop(b, jwtHandler, ctx, token)
	}
	//b.Log("GC Pause:", gcPause())
}

func benchmarkServeHTTPDefaultConfigBlackListSetup(b *testing.B) (ctxhttp.Handler, context.Context, []byte) {

	jwts := ctxjwt.MustNewService(
		ctxjwt.WithErrorHandler(scope.Default, 0, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
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
	jwts.Blacklist = blacklist.NewBlackListSimpleMap()

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		//store.WithStorageConfig(cr), no configuration so config.ScopedGetter is nil
	)
	ctx := store.WithContextProvider(context.Background(), srv) // root context

	token, err := jwts.NewToken(scope.Website, 1, jwtclaim.Map{ // 1 = website euro
		"someKey":         2.718281,
		jwtclaim.KeyStore: "at",
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
	return jwtHandler, ctx, token.Raw
}

func getRequestWithToken(b *testing.B, token []byte) *http.Request {
	req, err := http.NewRequest("GET", "http://abc.xyz", nil)
	if err != nil {
		b.Fatal(err)
	}
	ctxjwt.SetHeaderAuthorization(req, token)
	return req
}

func benchmarkServeHTTPDefaultConfigBlackListLoop(b *testing.B, h ctxhttp.Handler, ctx context.Context, token []byte) {
	w := httptest.NewRecorder() // 3 allocs
	if err := h.ServeHTTPContext(ctx, w, getRequestWithToken(b, token)); err != nil {
		b.Fatal(err)
	}
	if w.Code != http.StatusUnavailableForLegalReasons {
		b.Fatalf("Response Code want %d; have %d", http.StatusUnavailableForLegalReasons, w.Code)
	}
}

// BenchmarkServeHTTP_MultiToken a bench mark which runs parallel and creates
// token for different store scopes. This means that the underlying map in the
// Service struct much performan many scope switches to return the correct scope.

// BenchmarkServeHTTP_MultiToken_MultiScope-4	  200000	     10332 ns/op	    3648 B/op	      64 allocs/op => null blacklist
// BenchmarkServeHTTP_MultiToken_MultiScope-4	  200000	     11583 ns/op	    3648 B/op	      64 allocs/op => map blacklist
// BenchmarkServeHTTP_MultiToken_MultiScope-4	  200000	      9800 ns/op	    3647 B/op	      64 allocs/op => freecache
// BenchmarkServeHTTP_MultiToken_MultiScope-4	  200000	      9580 ns/op	    3657 B/op	      63 allocs/op
// BenchmarkServeHTTP_MultiToken_MultiScope-4	  200000	      8366 ns/op	    3194 B/op	      43 allocs/op => no maps, freecache
func BenchmarkServeHTTP_MultiToken_MultiScope(b *testing.B) {

	jwts := ctxjwt.MustNewService(
		ctxjwt.WithErrorHandler(scope.Default, 0, ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
			_, err := ctxjwt.FromContext(ctx)
			if err != nil {
				return err
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		})),
		ctxjwt.WithExpiration(scope.Default, 0, time.Second*15),
		ctxjwt.WithExpiration(scope.Website, 1, time.Second*25),
		ctxjwt.WithKey(scope.Website, 1, csjwt.WithPasswordRandom()),
		ctxjwt.WithTemplateToken(scope.Website, 1, func() csjwt.Token {
			return csjwt.Token{
				Header: csjwt.NewHead(),
				Claims: jwtclaim.NewStore(),
			}
		}),
	)

	// below two lines comment out enables the null black list
	jwts.Blacklist = blacklist.NewBlackListFreeCache(0)
	//jwts.Blacklist = ctxjwt.NewBlackListSimpleMap()
	// for now it doesn't matter which blacklist version you use as the bottle neck
	// is somewhere else.

	var generateToken = func(storeCode string) []byte {
		s := jwtclaim.NewStore()
		s.Store = storeCode
		token, err := jwts.NewToken(scope.Website, 1, s)
		if err != nil {
			b.Fatal(err)
		}
		return token.Raw
	}

	// generate 9k tokens randomly distributed over those three scopes.
	const tokenCount = 9000
	var tokens [tokenCount][]byte
	var storeCodes = [...]string{"au", "de", "at", "uk", "nz"}
	for i := range tokens {
		tokens[i] = generateToken(storeCodes[rand.Intn(len(storeCodes))])

		// just add garbage to the blacklist
		tbl := generateToken(strconv.FormatInt(int64(i), 10))
		if err := jwts.Blacklist.Set(tbl, time.Millisecond*time.Microsecond*time.Duration(i)); err != nil {
			b.Fatal(err)
		}
	}

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		// scope store, that means you can switch to any store independent from its
		// website, no restricts apply. i need to refactor this behaviour to make
		// the store init process different to allow multiple websites but with
		// the same restrictions applied. maybe move the SetByCode to the client
		// side ...
		scope.MustSetByCode(scope.Store, "at"), // at default store for this context
		store.WithStorageConfig(cr),
	)
	ctx := store.WithContextProvider(context.Background(), srv) // root context

	final := ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
		tok, err := ctxjwt.FromContext(ctx)
		if err != nil {
			b.Logf("%#v", tok)
			return err
		}
		w.WriteHeader(http.StatusUnavailableForLegalReasons)

		_, st, err := store.FromContextProvider(ctx)
		if err != nil {
			return err
		}
		switch st.StoreCode() {
		case "de", "at", "uk", "nz", "au":
		default:
			b.Fatalf("Unexpected Store: %s", st.StoreCode())
		}
		return nil
	})
	jwtHandler := jwts.WithInitTokenAndStore()(final)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			w := httptest.NewRecorder() // 3 allocs
			if err := jwtHandler.ServeHTTPContext(ctx, w, getRequestWithToken(b, tokens[i%tokenCount])); err != nil {
				b.Fatal(cserr.NewMultiErr(err).VerboseErrors())
			}
			if w.Code != http.StatusUnavailableForLegalReasons {
				b.Fatalf("Response Code want %d; have %d", http.StatusUnavailableForLegalReasons, w.Code)
			}
			i++
		}
	})
	// b.Log("GC Pause:", cstesting.GCPause())
}
