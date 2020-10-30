// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/net/jwt"
	"github.com/corestoreio/pkg/storage/containable"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
)

func testAuth_WithToken(t *testing.T, opts ...jwt.Option) (http.Handler, []byte) {
	cfg := cfgmock.NewService()
	jm, err := jwt.New(append(opts, jwt.WithRootConfig(cfg))...)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	jm.Log = log.BlackHole{EnableDebug: true, EnableInfo: true}

	theToken, err := jm.NewToken(scope.DefaultTypeID, jwtclaim.Map{
		"xfoo": "bar",
		"zfoo": 4711,
	})
	assert.NoError(t, err)

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	authHandler := jm.WithToken(final)
	return authHandler, theToken.Raw
}

func TestService_WithToken_EmptyScope(t *testing.T) {
	authHandler, _ := testAuth_WithToken(t,
		jwt.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotImplemented)
				tk, ok := jwt.FromContext(r.Context())
				assert.False(t, tk.Valid)
				assert.False(t, ok)
				assert.True(t, errors.IsNotFound(err), "%+v", err)
			})
		}),
	)

	req := httptest.NewRequest("GET", "http://auth1.xyz", nil)
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestService_WithToken_MissingToken(t *testing.T) {
	authHandler, _ := testAuth_WithToken(t,
		jwt.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}, scope.Website.WithID(1)),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	req = req.WithContext(scope.WithContext(req.Context(), 1, 2))
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), http.StatusText(http.StatusUnauthorized)+"\n")
}

func TestService_WithToken_Disabled(t *testing.T) {
	authHandler, _ := testAuth_WithToken(t,
		jwt.WithDisable(true, scope.Website.WithID(44)),
		jwt.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}, scope.Website.WithID(1)),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	req = req.WithContext(scope.WithContext(req.Context(), 44, 0))
	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestService_WithToken_Success(t *testing.T) {
	authHandler, token := testAuth_WithToken(t,
		jwt.WithDisable(false, scope.Website.WithID(55)),
		jwt.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}, scope.Website.WithID(55)),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	req = req.WithContext(scope.WithContext(req.Context(), 55, 0))
	jwt.SetHeaderAuthorization(req, token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestService_WithToken_SingleUsage(t *testing.T) {
	authHandler, token := testAuth_WithToken(t,
		jwt.WithDisable(false, scope.Website.WithID(66)),
		jwt.WithSingleTokenUsage(true, scope.Website.WithID(66)),
		jwt.WithErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}, scope.Website.WithID(66)),
		jwt.WithServiceErrorHandler(func(err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("Should not get called")
			})
		}),
		// default is a null blockList so we must set one
		jwt.WithBlocklist(set.NewInMemory()),
	)

	req := httptest.NewRequest("GET", "http://auth2.xyz", nil)
	req = req.WithContext(scope.WithContext(req.Context(), 66, 0))
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
