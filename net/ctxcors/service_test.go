// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ctxcors_test

import (
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/ctxcors"
	"github.com/corestoreio/csfw/net/ctxcors/internal/corstest"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func reqWithStore(method string) *http.Request {
	req, err := http.NewRequest(method, "http://corestore.io/foo", nil)
	if err != nil {
		panic(err)
	}

	return req.WithContext(
		store.WithContextRequestedStore(req.Context(), storemock.MustNewStoreAU(cfgmock.NewService())),
	)
}

func TestMustNew(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.IsNotValid(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = ctxcors.MustNew(ctxcors.WithMaxAge(scope.Default, 0, -2*time.Second))
}

func TestNoConfig(t *testing.T) {
	s := ctxcors.MustNew()
	req := reqWithStore("GET")
	corstest.TestNoConfig(t, s, req)
}

func TestMatchAllOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "*"),
	)
	req := reqWithStore("GET")
	corstest.TestMatchAllOrigin(t, s, req)
}

func TestAllowedOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestAllowedOrigin(t, s, req)
}

func TestWildcardOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://*.bar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestWildcardOrigin(t, s, req)
}

func TestDisallowedOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestDisallowedOrigin(t, s, req)
}

func TestDisallowedWildcardOrigin(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://*.bar.com"),
	)
	req := reqWithStore("GET")
	corstest.TestDisallowedWildcardOrigin(t, s, req)
}

func TestAllowedOriginFunc(t *testing.T) {
	r, _ := regexp.Compile("^http://foo")
	s := ctxcors.MustNew(
		ctxcors.WithAllowOriginFunc(scope.Default, 0, func(o string) bool {
			return r.MatchString(o)
		}),
	)
	req := reqWithStore("GET")
	corstest.TestAllowedOriginFunc(t, s, req)
}

func TestAllowedMethod(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedMethod(t, s, req)
}

func TestAllowedMethodPassthrough(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
		ctxcors.WithOptionsPassthrough(scope.Default, 0, true),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedMethodPassthrough(t, s, req)
}

func TestDisallowedMethod(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedMethods(scope.Default, 0, "PUT", "DELETE"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedMethod(t, s, req)
}

func TestAllowedHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedHeader(t, s, req)
}

func TestAllowedWildcardHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedHeaders(scope.Default, 0, "*"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestAllowedWildcardHeader(t, s, req)
}

func TestDisallowedHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestDisallowedHeader(t, s, req)
}

func TestOriginHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
	)
	req := reqWithStore("OPTIONS")
	corstest.TestOriginHeader(t, s, req)
}

func TestExposedHeader(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithExposedHeaders(scope.Default, 0, "X-Header-1", "x-header-2"),
	)

	req := reqWithStore("GET")
	corstest.TestExposedHeader(t, s, req)
}

func TestAllowedCredentials(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithAllowCredentials(scope.Default, 0, true),
	)

	req := reqWithStore("OPTIONS")
	corstest.TestAllowedCredentials(t, s, req)
}

func TestMaxAge(t *testing.T) {
	s := ctxcors.MustNew(
		ctxcors.WithAllowedOrigins(scope.Default, 0, "http://foobar.com"),
		ctxcors.WithMaxAge(scope.Default, 0, time.Second*30),
	)

	req := reqWithStore("OPTIONS")
	corstest.TestMaxAge(t, s, req)
}
