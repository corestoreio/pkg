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

package auth_test

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/net/auth"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/hashpool"
	"github.com/stretchr/testify/require"
)

func init() {
	if err := hashpool.Register("sha256", sha256.New); err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

func TestScopedConfig_Authenticate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc             string
		opts             []auth.Option
		req              *http.Request
		wantConfigErrBhf errors.Kind
		wantAuthErrBhf   errors.Kind
	}{
		{
			"Invalid default scope config",
			nil,
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NotValid,
			errors.NoKind,
		},
		{
			"REGEX blockList error",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"][0-9"}, nil)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.Fatal,
			errors.NoKind,
		},
		{
			"REGEX allowList error",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{""}, []string{"][0-9"})},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.Fatal,
			errors.NoKind,
		},
		{
			"Blocks all resources and never authenticates (nils)",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"REGEXP Blocks all resources and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs(nil, nil), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Blocks all resources and never authenticates (slash [Not the G'n'R guitar player])",
			[]auth.Option{auth.WithResourceACLs([]string{"/"}, []string{}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Blocks all resources and never authenticates calls multiple authenticators; all return true",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithInvalidAuth(true), auth.WithInvalidAuth(true), auth.WithInvalidAuth(true)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Blocks all resources and never authenticates calls multiple authenticators but last returns false",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithInvalidAuth(true), auth.WithInvalidAuth(true), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"REGEX Blocks all resources and never authenticates (slash [Not the G'n'R guitar player])",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"/"}, []string{}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Blocks all resources and never authenticates; but disabled",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithInvalidAuth(false), auth.WithDisable(true)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"Blocks all resources, except /catalog prefix route and never authenticates",
			[]auth.Option{auth.WithResourceACLs([]string{"/"}, []string{"/catalog"}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/catalog/product", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"REGEXP Blocks all resources, except /catalog prefix route and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/.+"}, []string{"^/cata[a-z]+/pro[a-z]+"}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/catalog/product", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and never authenticates",
			[]auth.Option{auth.WithResourceACLs(nil, []string{"/catalog"}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/customer/catalog", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"REGEX Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs(nil, []string{"^/cata[a-z]+"}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/customer/catalog", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and always authenticates",
			[]auth.Option{auth.WithResourceACLs(nil, []string{"/catalog"}), auth.WithValidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/catalog", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"REGEX Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and always authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs(nil, []string{"/catalog"}), auth.WithValidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/calalog", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"Blocks one resource /customer, we call /customer/wishlist and always authenticates",
			[]auth.Option{auth.WithResourceACLs([]string{"/customer"}, nil), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/customer/wishlist", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"REGEX Blocks one resource /customer, we call /customer/wishlist and always authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/custom[a-z]+"}, nil), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/customel/wishlist", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Blocks one resource /customer, we call /customer/forgetpassword which is allowListed",
			[]auth.Option{auth.WithResourceACLs([]string{"/customer"}, []string{"/customer/resetpassword", "/customer/forgetpassword"}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/customer/forgetpassword?param=1", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"REGEX Blocks one resource /customer, we call /customer/forgetpassword which is allowListed",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/customer"}, []string{"^/[a-z]+/resetpassword", "^/[a-z]+/forgetpassword"}), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/customer/forgetpassword?param=1", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"Blocks one resource /customer, we call /catalog/category and never authenticates",
			[]auth.Option{auth.WithResourceACLs([]string{"/customer"}, nil), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"REGEX Blocks one resource /customer, we call /catalog/category and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/customer"}, nil), auth.WithInvalidAuth(false)},
			httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"Basic Auth: Blocks all resources. Unauthorized",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithSimpleBasicAuth("user1", "pass2", "R3alm")},
			httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Basic Auth: Blocks all resources. Authorized",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithSimpleBasicAuth("user1", "pass2", "R3alm")},
			func() *http.Request {
				r := httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil)
				r.SetBasicAuth("user1", "pass2")
				return r
			}(),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"Basic Auth: Blocks all resources. Authorization failed",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithSimpleBasicAuth("user1", "pass2", "R3alm")},
			func() *http.Request {
				r := httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil)
				r.SetBasicAuth("user2", "pass3")
				return r
			}(),
			errors.NoKind,
			errors.Unauthorized,
		},
		{
			"Basic Auth: Blocks all resources. Basic Failed but ValidAuth always succeeds",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithSimpleBasicAuth("user1", "pass2", "R3alm"), auth.WithValidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil),
			errors.NoKind,
			errors.NoKind,
		},
		{
			"Basic Auth: Blocks all resources. Unauthorized and always fails",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithSimpleBasicAuth("user1", "pass2", "R3alm"), auth.WithInvalidAuth(true)},
			func() *http.Request {
				r := httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil)
				r.SetBasicAuth("user1", "uuups")
				return r
			}(),
			errors.NoKind,
			errors.Unauthorized,
		},
	}
	for i, test := range tests {
		srv, haveErr := auth.New(nil, test.opts...)
		if test.wantConfigErrBhf > 0 && haveErr != nil {
			assert.True(t, test.wantConfigErrBhf.Match(haveErr), "Index %d %q\n%+v", i, test.desc, haveErr)
			continue
		}
		require.NoError(t, haveErr, "Index %d %q\n%+v", i, test.desc, haveErr)

		scpCfg, haveErr := srv.ConfigByScopeID(scope.DefaultTypeID, scope.DefaultTypeID)
		if test.wantConfigErrBhf > 0 {
			assert.True(t, test.wantConfigErrBhf.Match(haveErr), "Index %d %q\n%+v", i, test.desc, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d %q\n%+v", i, test.desc, haveErr)

		haveErr = scpCfg.Authenticate(test.req)
		if test.wantAuthErrBhf > 0 {
			assert.True(t, test.wantAuthErrBhf.Match(haveErr), "Index %d %q\n%+v", i, test.desc, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d %q\n%+v", i, test.desc, haveErr)
		}
	}
}
