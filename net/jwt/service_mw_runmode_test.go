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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/blacklist"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"hash/fnv"
)

func testAuth_WithRunMode(t *testing.T, opts ...jwt.Option) (http.Handler, []byte) {
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(append(opts, jwt.WithConfigGetter(cfg))...)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	theToken, err := jm.NewToken(scope.Default, 0, jwtclaim.Map{
		"xfoo": "baz",
		"zfoo": 4712,
	})
	assert.NoError(t, err)

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := storemock.NewEurozzyService(cfg)

	authHandler := jm.WithRunMode(scope.RunMode{}, srv)(final)
	return authHandler, theToken.Raw
}

func TestService_WithRunMode_NoToken(t *testing.T) {

	//  request calls default unauthorized handler

	authHandler, _ := testAuth_WithRunMode(t,
		jwt.WithErrorHandler(scope.Default, 0,
			func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("Should not get called")
				})
			}),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	req := httptest.NewRequest("GET", "http://auth.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusUnauthorized)+"\n")
}

func TestService_WithRunMode_Custom_UnauthorizedHandler(t *testing.T) {

	// request calls the unauthorized handler of website scope 1 = euro scope

	var calledUnauthorizedHandler bool
	authHandler, _ := testAuth_WithRunMode(t,
		jwt.WithErrorHandler(scope.Default, 0,
			func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("Should not get called")
				})
			}),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		jwt.WithUnauthorizedHandler(scope.Website, 1, func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tk, ok := jwt.FromContext(r.Context())
				assert.False(t, tk.Valid)
				assert.False(t, ok)

				assert.True(t, errors.IsNotFound(err), "%+v", err)
				w.WriteHeader(http.StatusTeapot)
				calledUnauthorizedHandler = true
			})
		}),
	)

	req := httptest.NewRequest("GET", "http://auth.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTeapot, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledUnauthorizedHandler)
}

func TestService_WithRunMode_Disabled(t *testing.T) {
	authHandler, _ := testAuth_WithRunMode(t,
		jwt.WithDisable(scope.Website, 1, true), // 1 == euro website
		jwt.WithErrorHandler(scope.Website, 1,
			func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("Should not get called")
				})
			}),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestService_WithRunMode_SingleUsage(t *testing.T) {
	authHandler, token := testAuth_WithRunMode(t,
		jwt.WithDisable(scope.Website, 1, false),
		jwt.WithSingleTokenUsage(scope.Website, 1, true),
		jwt.WithErrorHandler(scope.Website, 1,
			func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("Should not get called")
				})
			}),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		// default is a null blacklist so we must set one
		jwt.WithBlacklist(blacklist.NewInMemory(fnv.New64a)),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	req = req.WithContext(scope.WithContext(req.Context(), 1, 0))
	jwt.SetHeaderAuthorization(req, token)

	// 1st request ok
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())

	// 2nd request unauthorized
	w = httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusUnauthorized)+"\n")
}

func TestService_WithRunMode_DefaultStoreID_Error(t *testing.T) {
	var calledServiceErrorHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Default, 0,
			func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("Should not get called")
				})
			}),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
				assert.True(t, errors.IsNotImplemented(err))
				calledServiceErrorHandler = true
			})
		}),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		StoreIDError: errors.NewNotImplementedf("Sorry Dude"),
	})(final)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledServiceErrorHandler)
}

func TestService_WithRunMode_StoreIDbyCode_Error(t *testing.T) {
	var calledServiceErrorHandler bool
	cfg := cfgmock.NewService()
	jm := jwt.MustNew(
		jwt.WithConfigGetter(cfg),
		jwt.WithErrorHandler(scope.Website, 778,
			func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusAlreadyReported)
					assert.True(t, errors.IsNotImplemented(err))
					calledServiceErrorHandler = true
				})
			}),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should NOT get called")
			})
		}),
	)
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Should not get called")
	})

	authHandler := jm.WithRunMode(scope.RunMode{}, storemock.Find{
		WebsiteIDDefault: 778,
		IDByCodeError:    errors.NewNotImplementedf("Sorry Dude"),
	})(final)

	claimStore := jwtclaim.NewStore()
	claimStore.Store = "'80s FTW"
	theToken, err := jm.NewToken(scope.Website, 778, claimStore)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	jwt.SetHeaderAuthorization(req, theToken.Raw)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAlreadyReported, w.Code)
	assert.Empty(t, w.Body.String())
	assert.True(t, calledServiceErrorHandler)
}

func TestService_WithRunMode_IsAllowedStoreID_Error(t *testing.T) {}

//func TestService_WithInitTokenAndStore_InvalidToken(t *testing.T) {
//
//	ctx := store.WithContextRequestedStore(context.Background(), storemock.MustNewStoreAU(cfgmock.NewService()))
//
//	jwts := jwt.MustNew(
//		jwt.WithExpiration(scope.Website, 12, -time.Second),
//		jwt.WithSkew(scope.Website, 12, 0),
//	)
//
//	if err := jwts.Options(jwt.WithErrorHandler(scope.Default, 0,
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusInternalServerError)
//			token, err := jwt.FromContext(r.Context())
//			assert.Nil(t, token.Raw)
//			assert.False(t, token.Valid)
//			assert.True(t, errors.IsNotValid(err), "Error: %+v", err)
//		}),
//	)); err != nil {
//		t.Fatalf("%+v", err)
//	}
//
//	theToken, err := jwts.NewToken(scope.Website, 12, jwtclaim.Map{
//		"xfoo": "invalid",
//		"zfoo": -time.Second,
//	})
//	assert.NoError(t, err)
//	assert.NotEmpty(t, theToken.Raw)
//
//	req, err := http.NewRequest("GET", "http://corestore.io/customer/wishlist", nil)
//	assert.NoError(t, err)
//	jwt.SetHeaderAuthorization(req, theToken.Raw)
//
//	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		t.Fatal("Should not be executed")
//
//	})
//	authHandler := jwts.WithInitTokenAndStore()(finalHandler)
//
//	wRec := httptest.NewRecorder()
//
//	authHandler.ServeHTTP(wRec, req.WithContext(ctx))
//	assert.Equal(t, http.StatusInternalServerError, wRec.Code)
//	assert.Equal(t, ``, wRec.Body.String())
//}
//
//type testRealBL struct {
//	theToken []byte
//	exp      time.Duration
//}
//
//func (b *testRealBL) Set(t []byte, exp time.Duration) error {
//	b.theToken = t
//	b.exp = exp
//	return nil
//}
//func (b *testRealBL) Has(t []byte) bool { return bytes.Equal(b.theToken, t) }
//
//var _ jwt.Blacklister = (*testRealBL)(nil)
//
//func TestService_WithInitTokenAndStore_InBlackList(t *testing.T) {
//
//	cr := cfgmock.NewService()
//	srv := storemock.NewEurozzyService(
//		scope.MustSetByCode(scope.Website, "euro"),
//		store.WithStorageConfig(cr),
//	)
//	dsv, err := srv.Store()
//	ctx := store.WithContextRequestedStore(context.Background(), dsv, err)
//
//	bl := &testRealBL{}
//	jm, err := jwt.New(
//		jwt.WithBlacklist(bl),
//	)
//	assert.NoError(t, err)
//
//	theToken, err := jm.NewToken(scope.Default, 0, &jwtclaim.Standard{})
//	bl.theToken = theToken.Raw
//	assert.NoError(t, err)
//	assert.NotEmpty(t, theToken.Raw)
//
//	req, err := http.NewRequest("GET", "http://auth.xyz", nil)
//	assert.NoError(t, err)
//	jwt.SetHeaderAuthorization(req, theToken.Raw)
//
//	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		_, err := jwt.FromContext(r.Context())
//		assert.True(t, errors.IsNotValid(err))
//		w.WriteHeader(http.StatusUnauthorized)
//	})
//	authHandler := jm.WithInitTokenAndStore()(finalHandler)
//
//	wRec := httptest.NewRecorder()
//	authHandler.ServeHTTP(wRec, req.WithContext(ctx))
//
//	//assert.NotEqual(t, http.StatusTeapot, wRec.Code)
//	assert.Equal(t, http.StatusUnauthorized, wRec.Code)
//}
//
//// todo add test for form with input field: access_token
//

//func TestService_WithInitTokenAndStore_Request(t *testing.T) {
//
//	var newReq = func(i int, token []byte) *http.Request {
//		req, err := http.NewRequest("GET", fmt.Sprintf("https://corestore.io/store/list/%d", i), nil)
//		if err != nil {
//			t.Fatalf("%+v", err)
//		}
//		jwt.SetHeaderAuthorization(req, token)
//		return req
//	}
//
//	srvDE := storemock.NewEurozzyService(scope.Option{Store: scope.MockCode("de")})
//	dsvDE, errDE := srvDE.Store()
//	dsvDE.Website.Config = cfgmock.NewService().NewScoped(1, 0)
//
//	tests := []struct {
//		scpOpt         scope.Option
//		ctx            context.Context
//		tokenStoreCode string
//		wantStoreCode  string
//		wantErrBhf     errors.BehaviourFunc
//	}{
//		{scope.Option{}, store.WithContextRequestedStore(context.Background(), nil), "de", "de", errors.IsNotFound},
//		{scope.Option{}, store.WithContextRequestedStore(context.Background(), dsvDE, errDE), "de", "de", errors.IsNotFound},
//		{scope.Option{Store: scope.MockCode("de")}, nil, "de", "de", nil},
//		{scope.Option{Store: scope.MockCode("at")}, nil, "ch", "at", errors.IsUnauthorized},
//		{scope.Option{Store: scope.MockCode("de")}, nil, "at", "at", nil},
//		{scope.Option{Store: scope.MockCode("de")}, nil, "a$t", "de", errors.IsNotValid},
//		{scope.Option{Store: scope.MockCode("at")}, nil, "", "at", errors.IsNotValid},
//		//
//		{scope.Option{Group: scope.MockID(1)}, nil, "de", "de", nil},
//		{scope.Option{Group: scope.MockID(1)}, nil, "ch", "at", errors.IsUnauthorized},
//		{scope.Option{Group: scope.MockID(1)}, nil, " ch", "at", errors.IsNotValid},
//		{scope.Option{Group: scope.MockID(1)}, nil, "uk", "at", errors.IsUnauthorized},
//
//		{scope.Option{Website: scope.MockID(2)}, nil, "uk", "au", errors.IsUnauthorized},
//		{scope.Option{Website: scope.MockID(2)}, nil, "nz", "nz", nil},
//		{scope.Option{Website: scope.MockID(2)}, nil, "n z", "au", errors.IsNotValid},
//		{scope.Option{Website: scope.MockID(2)}, nil, "au", "au", nil},
//		{scope.Option{Website: scope.MockID(2)}, nil, "", "au", errors.IsNotValid},
//	}
//	for i, test := range tests {
//		if test.ctx == nil {
//			test.ctx = newStoreServiceWithCtx(test.scpOpt)
//		}
//
//		//logBuf := new(bytes.Buffer)
//
//		jwts := jwt.MustNew(
//			//jwt.WithLogger(logw.NewLog(logw.WithWriter(logBuf), logw.WithLevel(logw.LevelDebug))),
//			jwt.WithKey(scope.Default, 0, csjwt.WithPasswordRandom()),
//			jwt.WithStoreFinder(storemock.NewEurozzyService(test.scpOpt)),
//		)
//
//		token, err := jwts.NewToken(scope.Default, 0, jwtclaim.Map{
//			jwt.StoreParamName: test.tokenStoreCode,
//		})
//		if err != nil {
//			t.Fatalf("%+v", err)
//		}
//
//		if test.wantErrBhf != nil {
//			if err := jwts.Options(jwt.WithErrorHandler(scope.Default, 0,
//				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//					_, err := jwt.FromContext(r.Context())
//					assert.True(t, test.wantErrBhf(err), "Index %d => %s", i, err)
//				}),
//			)); err != nil {
//				t.Fatalf("%+v", err)
//			}
//		}
//		cstesting.NewHTTPParallelUsers(2, 5, 200, time.Millisecond).ServeHTTP(
//			newReq(i, token.Raw).WithContext(test.ctx),
//			jwts.WithInitTokenAndStore()(finalInitStoreHandler(t, i, test.wantStoreCode)),
//		)
//
//		//const searchLogEntry1 = `jwt.Service.ConfigByScopedGetter.IsValid`
//		//if have, want := strings.Count(logBuf.String(), searchLogEntry1), 1; have != want {
//		//	t.Errorf("%s: Have: %v Want: %v", searchLogEntry1, have, want)
//		//}
//	}
//}
//
//// TestService_WithInitTokenAndStore_StoreServiceNil tests to forward a valid
//// request with a vailid token to the next handler in the chain despite the
//// StoreService is nil and the token contains a store code not associated with
//// the current initialized store.
//func TestService_WithInitTokenAndStore_StoreServiceNil(t *testing.T) {
//
//	ctx := store.WithContextRequestedStore(context.Background(), storemock.MustNewStoreAU(cfgmock.NewService()))
//
//	jwts := jwt.MustNew(
//		jwt.WithExpiration(scope.Website, 12, time.Second),
//	)
//
//	if err := jwts.Options(jwt.WithErrorHandler(scope.Default, 0,
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			t.Fatal("Should not be executed this error handler")
//		}),
//	)); err != nil {
//		t.Fatalf("%+v", err)
//	}
//
//	claimStore := jwtclaim.NewStore()
//	claimStore.Store = "de"
//	claimStore.Audience = "eCommerce"
//	theToken, err := jwts.NewToken(scope.Website, 12, claimStore)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, theToken.Raw)
//
//	req, err := http.NewRequest("GET", "http://corestore.io/customer/wishlist", nil)
//	assert.NoError(t, err)
//	jwt.SetHeaderAuthorization(req, theToken.Raw)
//
//	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusAccepted)
//		tk, err := jwt.FromContext(r.Context())
//		if err != nil {
//			t.Fatal(err)
//		}
//		haveSt, err := tk.Claims.Get(jwtclaim.KeyStore)
//		if err != nil {
//			t.Fatal(err)
//		}
//		assert.Exactly(t, "de", haveSt)
//
//		reqStore, err := store.FromContextRequestedStore(r.Context())
//		if err != nil {
//			t.Fatal(err)
//		}
//		assert.Exactly(t, "au", reqStore.StoreCode())
//	})
//	authHandler := jwts.WithInitTokenAndStore()(finalHandler)
//
//	hpu := cstesting.NewHTTPParallelUsers(2, 2, 100, time.Nanosecond)
//	hpu.AssertResponse = func(rec *httptest.ResponseRecorder) {
//		assert.Equal(t, http.StatusAccepted, rec.Code)
//		assert.Equal(t, ``, rec.Body.String())
//	}
//	hpu.ServeHTTP(req.WithContext(ctx), authHandler)
//}
//
//// TestService_WithInitTokenAndStore_Disabled
//// 1. a request with a valid token
//// should do nothing and just call the next handler in the chain because the
//// token service has been disabled therefore no validation and checking takes
//// place.
//// 2. valid request with website oz must be passed through with an
//// invalid token because JWT disabled
//func TestService_WithInitTokenAndStore_Disabled(t *testing.T) {
//
//	jm, err := jwt.New(jwt.WithDisable(scope.Website, 2, true))
//	if err != nil {
//		t.Fatalf("%+v", err)
//	}
//
//	mw := jm.WithInitTokenAndStore()
//
//	// 1. valid request with website euro and token must be validated
//	{
//		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			tk, err := jwt.FromContext(r.Context())
//			if err != nil {
//				t.Fatal("MW Error", err)
//			}
//			assert.True(t, tk.Valid)
//			http.Error(w, http.StatusText(http.StatusMultipleChoices), http.StatusMultipleChoices)
//		}))
//
//		req, err := http.NewRequest("GET", "http://auth.xyz", nil)
//		if err != nil {
//			t.Fatalf("%+v", err)
//		}
//		theToken, err := jm.NewToken(scope.Default, 0, jwtclaim.Map{
//			"noStore": "bar",
//			"zfoo":    4711,
//		})
//		if err != nil {
//			t.Fatalf("%+v", err)
//		}
//		jwt.SetHeaderAuthorization(req, theToken.Raw)
//
//		w := httptest.NewRecorder()
//		handler.ServeHTTP(w, req.WithContext(newStoreServiceWithCtx(scope.MustSetByCode(scope.Website, "euro"))))
//		assert.Equal(t, http.StatusMultipleChoices, w.Code)
//	}
//
//	// 2. valid request with website oz must be passed through with an invalid token because JWT disabled
//	{
//		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			_, err := jwt.FromContext(r.Context())
//			assert.True(t, errors.IsNotFound(err))
//			assert.Exactly(t, `Bearer Invalid Token`, r.Header.Get("Authorization"))
//			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
//		}))
//
//		req, err := http.NewRequest("GET", "http://auth.xyz", nil)
//		if err != nil {
//			t.Fatalf("%+v", err)
//		}
//		jwt.SetHeaderAuthorization(req, []byte(`Invalid Token`))
//
//		w := httptest.NewRecorder()
//		handler.ServeHTTP(w, req.WithContext(newStoreServiceWithCtx(scope.MustSetByCode(scope.Website, "oz"))))
//		assert.Equal(t, http.StatusConflict, w.Code)
//	}
//}
