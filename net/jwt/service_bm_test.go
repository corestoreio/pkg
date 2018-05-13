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

package jwt_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/net/jwt"
	"github.com/corestoreio/pkg/storage/containable"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/store/storemock"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
)

func bmWithToken(b *testing.B, opts ...jwt.Option) {
	jwts := jwt.MustNew(opts...)

	if err := jwts.Options(jwt.WithRootConfig(cfgmock.NewService())); err != nil {
		b.Fatal(err)
	}
	cl := jwtclaim.Map{
		"xfoo": "bar",
		"zfoo": 4711,
	}
	token, err := jwts.NewToken(scope.Default.WithID(0), cl)
	if err != nil {
		b.Error(err)
	}

	jwtHandler := jwts.WithToken(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	req, err := http.NewRequest("GET", "http://abc.xyz", nil)
	if err != nil {
		b.Error(err)
	}
	jwt.SetHeaderAuthorization(req, token.Raw)
	req = req.WithContext(scope.WithContext(req.Context(), 0, 0)) // Default Scope

	bf := func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			jwtHandler.ServeHTTP(w, req)
			if w.Code != http.StatusTeapot {
				b.Errorf("Response Code want %d; have %d", http.StatusTeapot, w.Code)
			}
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(bf)
}

var keyBenchmarkHMACPW = jwt.WithKey(csjwt.WithPassword([]byte(`Rump3lst!lzch3n`))) // for scope.DefaultTypeID

// 200000	      8474 ns/op	    2698 B/op	      63 allocs/op <= Go 1.7
func BenchmarkWithToken_HMAC_InMemoryBL(b *testing.B) {
	bl := containable.NewInMemory()
	bmWithToken(b, keyBenchmarkHMACPW, jwt.WithBlacklist(bl))
	// b.Logf("Blacklist Items %d", bl.Len())
}

// 30000	     55376 ns/op	    9180 B/op	      92 allocs/op <= Go 1.7
func BenchmarkWithToken_RSAGenerator_2048(b *testing.B) {
	bmWithToken(b, jwt.WithKey(csjwt.WithRSAGenerated()))
}

func getRequestWithToken(b *testing.B, token []byte) *http.Request {
	req, err := http.NewRequest("GET", "http://abc.xyz", nil)
	if err != nil {
		b.Fatal(err)
	}
	jwt.SetHeaderAuthorization(req, token)
	return req
}

// BenchmarkWithRunMode_MultiTokenAndScope a bench mark which runs parallel and
// creates token for different store scopes. This means that the underlying map
// in the Service struct much perform many scope switches to return the correct
// scope.
// 200000	      9436 ns/op	    3620 B/op	      47 allocs/op
func BenchmarkWithRunMode_MultiTokenAndScope(b *testing.B) {
	cfg := cfgmock.NewService()
	jwts := jwt.MustNew(
		jwt.WithRootConfig(cfg),
		jwt.WithExpiration(time.Second*15),
		jwt.WithExpiration(time.Second*25, scope.Website.WithID(1)),
		jwt.WithKey(csjwt.WithPasswordRandom(), scope.Website.WithID(1)),
		jwt.WithTemplateToken(func() csjwt.Token {
			return csjwt.Token{
				Header: csjwt.NewHead(),
				Claims: jwtclaim.NewStore(),
			}
		}, scope.Website.WithID(1)),
	)

	// below two lines comment out enables the null black list
	jwts.Blacklist = containable.NewInMemory()

	var generateToken = func(storeCode string) []byte {
		s := jwtclaim.NewStore()
		s.Store = storeCode
		token, err := jwts.NewToken(scope.Website.WithID(1), s)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		return token.Raw
	}

	// generate 9k tokens randomly distributed over those three website scopes.
	const tokenCount = 9000
	var tokens [tokenCount][]byte
	// var storeCodes = [...]string{"au", "de", "at", "uk", "nz"}
	var storeCodes = [...]string{"de", "at", "uk"}
	for i := range tokens {
		tokens[i] = generateToken(storeCodes[rand.Intn(len(storeCodes))])

		// just add garbage to the blacklist
		tbl := generateToken(strconv.Itoa(i))
		if err := jwts.Blacklist.Set(tbl, time.Millisecond*time.Microsecond*time.Duration(i)); err != nil {
			b.Fatalf("%+v", err)
		}
	}

	storeSrv := storemock.NewEurozzyService(cfg)
	jwtHandler := jwts.WithRunMode(scope.Website.WithID(1), // every store with website ID 1
		storeSrv)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, ok := jwt.FromContext(ctx)
		if !ok {
			b.Fatal("Token not found in context") // epic fail
		}
		w.WriteHeader(http.StatusUnavailableForLegalReasons)

		websiteID, storeID, ok := scope.FromContext(r.Context())
		if !ok {
			b.Fatal("Context Scope not found")
		}
		if websiteID != 1 {
			b.Fatalf("websiteID Have %d, Want %d", websiteID, 1)
		}

		switch storeID {
		case 1, 2, 4:
		// do nothing, everything OK
		default:
			b.Fatalf("Unexpected StoreID: %d", storeID)
		}
	}))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			w := httptest.NewRecorder() // 3 allocs
			jwtHandler.ServeHTTP(w, getRequestWithToken(b, tokens[i%tokenCount]))

			if w.Code != http.StatusUnavailableForLegalReasons && w.Code != http.StatusUnauthorized {
				b.Fatalf("Response Code want %d; have %d\n%s", http.StatusUnavailableForLegalReasons, w.Code, w.Body)
			}
			i++
		}
	})
	// b.Log("GC Pause:", cstesting.GCPause())
}
