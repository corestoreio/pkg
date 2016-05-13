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

package jwtauth_test

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/jwtauth"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/blacklist"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
)

const benchTokenCount = 100

func benchBlackList(b *testing.B, bl jwtauth.Blacklister) {
	jwts := jwtauth.MustNewService()
	var tokens [benchTokenCount][]byte

	for i := 0; i < benchTokenCount; i++ {
		claim := jwtclaim.Map{
			"someKey": i,
		}
		tk, err := jwts.NewToken(scope.Default, 0, claim)
		if err != nil {
			b.Fatal(errors.PrintLoc(err))
		}
		tokens[i] = tk.Raw
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < benchTokenCount; i++ {
				if err := bl.Set(tokens[i], time.Minute); err != nil {
					b.Fatal(errors.PrintLoc(err))
				}
				if !bl.Has(tokens[i]) {
					b.Fatalf("Cannot find token %s with index %d", tokens[i], i)
				}
			}
		}
	})
}

func BenchmarkBlackListMap_Parallel(b *testing.B) {
	bl := blacklist.NewMap()
	benchBlackList(b, bl)
}

func BenchmarkBlackListFreeCache_Parallel(b *testing.B) {
	bl := blacklist.NewFreeCache(0)
	benchBlackList(b, bl)
}

func bmServeHTTP(b *testing.B, opts ...jwtauth.Option) {
	jwts, err := jwtauth.NewService(opts...)
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

	final := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	jwtHandler := jwts.WithInitTokenAndStore()(final)

	req, err := http.NewRequest("GET", "http://abc.xyz", nil)
	if err != nil {
		b.Error(err)
	}
	jwtauth.SetHeaderAuthorization(req, token.Raw)
	w := httptest.NewRecorder()

	cr := cfgmock.NewService()
	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		store.WithStorageConfig(cr),
	)
	dsv, err := srv.DefaultStoreView()
	ctx := store.WithContextRequestedStore(context.Background(), dsv, err)
	req = req.WithContext(ctx)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jwtHandler.ServeHTTP(w, req)
		if w.Code != http.StatusTeapot {
			b.Errorf("Response Code want %d; have %d", http.StatusTeapot, w.Code)
		}
	}
}

var keyBenchmarkHMACPW = jwtauth.WithKey(scope.Default, 0, csjwt.WithPassword([]byte(`Rump3lst!lzch3n`)))

func BenchmarkServeHTTPHMAC(b *testing.B) {
	bmServeHTTP(b, keyBenchmarkHMACPW)
}

func BenchmarkServeHTTPHMACSimpleBL(b *testing.B) {
	bl := blacklist.NewMap()
	bmServeHTTP(b,
		keyBenchmarkHMACPW,
		jwtauth.WithBlacklist(bl),
	)
	// b.Logf("Blacklist Items %d", bl.Len())
}

func BenchmarkServeHTTPRSAGenerator(b *testing.B) {
	bmServeHTTP(b, jwtauth.WithKey(scope.Default, 0, csjwt.WithRSAGenerated()))
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

func benchmarkServeHTTPDefaultConfigBlackListSetup(b *testing.B) (http.Handler, context.Context, []byte) {

	jwts := jwtauth.MustNewService(
		jwtauth.WithErrorHandler(scope.Default, 0, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := jwtauth.FromContext(r.Context())
			if err != nil {
				b.Fatal(err) // epic fail
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		})),
	)
	// below two lines comment out enables the null black list
	//jwts.Blacklist = jwtauth.NewBlackListFreeCache(0)
	jwts.Blacklist = blacklist.NewMap()

	srv := storemock.NewEurozzyService(
		scope.MustSetByCode(scope.Website, "euro"),
		//store.WithStorageConfig(cr), no configuration so config.ScopedGetter is nil
	)

	dsv, err := srv.DefaultStoreView()
	ctx := store.WithContextRequestedStore(context.Background(), dsv, err)

	token, err := jwts.NewToken(scope.Website, 1, jwtclaim.Map{ // 1 = website euro
		"someKey":         2.718281,
		jwtclaim.KeyStore: "at",
	})
	if err != nil {
		b.Fatal(errors.PrintLoc(err))
	}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := jwtauth.FromContext(r.Context())
		if err != nil {
			b.Fatal(err)
		}
		w.WriteHeader(http.StatusUnavailableForLegalReasons)

		st, err := store.FromContextRequestedStore(ctx)
		if err != nil {
			b.Fatal(err)
		}
		if st.StoreCode() != "de" && st.StoreCode() != "at" {
			b.Fatalf("Unexpected Store: %s", st.StoreCode())
		}
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
	jwtauth.SetHeaderAuthorization(req, token)
	return req
}

func benchmarkServeHTTPDefaultConfigBlackListLoop(b *testing.B, h http.Handler, ctx context.Context, token []byte) {
	w := httptest.NewRecorder() // 3 allocs
	h.ServeHTTP(w, getRequestWithToken(b, token).WithContext(ctx))

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

	jwts := jwtauth.MustNewService(
		jwtauth.WithErrorHandler(scope.Default, 0, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := jwtauth.FromContext(r.Context())
			if err != nil {
				b.Fatal(err)
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		})),
		jwtauth.WithExpiration(scope.Default, 0, time.Second*15),
		jwtauth.WithExpiration(scope.Website, 1, time.Second*25),
		jwtauth.WithKey(scope.Website, 1, csjwt.WithPasswordRandom()),
		jwtauth.WithTemplateToken(scope.Website, 1, func() csjwt.Token {
			return csjwt.Token{
				Header: csjwt.NewHead(),
				Claims: jwtclaim.NewStore(),
			}
		}),
	)

	// below two lines comment out enables the null black list
	jwts.Blacklist = blacklist.NewFreeCache(0)
	//jwts.Blacklist = jwtauth.NewBlackListSimpleMap()
	// for now it doesn't matter which blacklist version you use as the bottle neck
	// is somewhere else.

	var generateToken = func(storeCode string) []byte {
		s := jwtclaim.NewStore()
		s.Store = storeCode
		token, err := jwts.NewToken(scope.Website, 1, s)
		if err != nil {
			b.Fatal(errors.PrintLoc(err))
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
			b.Fatal(errors.PrintLoc(err))
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
	dsv, err := srv.DefaultStoreView()
	ctx := store.WithContextRequestedStore(context.Background(), dsv, err) // root context

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tok, err := jwtauth.FromContext(ctx)
		if err != nil {
			b.Fatalf("Error: %s\n%#v", err, tok)
		}
		w.WriteHeader(http.StatusUnavailableForLegalReasons)

		st, err := store.FromContextRequestedStore(ctx)
		if err != nil {
			b.Fatal(err)
		}
		switch st.StoreCode() {
		case "de", "at", "uk", "nz", "au":
			// do nothing all good
		default:
			b.Fatalf("Unexpected Store: %s", st.StoreCode())
		}
	})
	jwtHandler := jwts.WithInitTokenAndStore()(final)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			w := httptest.NewRecorder() // 3 allocs
			jwtHandler.ServeHTTP(w, getRequestWithToken(b, tokens[i%tokenCount]).WithContext(ctx))

			if w.Code != http.StatusUnavailableForLegalReasons {
				b.Fatalf("Response Code want %d; have %d", http.StatusUnavailableForLegalReasons, w.Code)
			}
			i++
		}
	})
	// b.Log("GC Pause:", cstesting.GCPause())
}
