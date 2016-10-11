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

package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/auth"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScopedConfig_Authenticate(t *testing.T) {
	tests := []struct {
		desc             string
		opts             []auth.Option
		req              *http.Request
		wantConfigErrBhf errors.BehaviourFunc
		wantAuthErrBhf   errors.BehaviourFunc
	}{
		{
			"Invalid default scope config",
			nil,
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.IsNotValid,
			nil,
		},
		{
			"REGEX blacklist error",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"][0-9"}, nil)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.IsFatal,
			nil,
		},
		{
			"REGEX whitelist error",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{""}, []string{"][0-9"})},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			errors.IsFatal,
			nil,
		},
		{
			"Blocks all resources and never authenticates (nils)",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"REGEXP Blocks all resources and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs(nil, nil), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"Blocks all resources and never authenticates (slash [Not the G'n'R guitar player])",
			[]auth.Option{auth.WithResourceACLs([]string{"/"}, []string{}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"REGEX Blocks all resources and never authenticates (slash [Not the G'n'R guitar player])",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"/"}, []string{}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"Blocks all resources and never authenticates; but disabled",
			[]auth.Option{auth.WithResourceACLs(nil, nil), auth.WithInvalidAuth(), auth.WithDisable(true)},
			httptest.NewRequest("GET", "http://corestore.io/anyroute", nil),
			nil,
			nil,
		},
		{
			"Blocks all resources, except /catalog prefix route and never authenticates",
			[]auth.Option{auth.WithResourceACLs([]string{"/"}, []string{"/catalog"}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/catalog/product", nil),
			nil,
			nil,
		},
		{
			"REGEXP Blocks all resources, except /catalog prefix route and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/.+"}, []string{"^/cata[a-z]+/pro[a-z]+"}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/catalog/product", nil),
			nil,
			nil,
		},
		{
			"Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and never authenticates",
			[]auth.Option{auth.WithResourceACLs(nil, []string{"/catalog"}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/catalog", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"REGEX Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs(nil, []string{"^/cata[a-z]+"}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/catalog", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and always authenticates",
			[]auth.Option{auth.WithResourceACLs(nil, []string{"/catalog"}), auth.WithValidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/catalog", nil),
			nil,
			nil,
		},
		{
			"REGEX Blocks all resources, except /catalog prefix route but we call different route /customer/catalog and always authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs(nil, []string{"/catalog"}), auth.WithValidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/calalog", nil),
			nil,
			nil,
		},
		{
			"Blocks one resource /customer, we call /customer/wishlist and always authenticates",
			[]auth.Option{auth.WithResourceACLs([]string{"/customer"}, nil), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/wishlist", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"REGEX Blocks one resource /customer, we call /customer/wishlist and always authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/custom[a-z]+"}, nil), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customel/wishlist", nil),
			nil,
			errors.IsUnauthorized,
		},
		{
			"Blocks one resource /customer, we call /customer/forgetpassword which is whitelisted",
			[]auth.Option{auth.WithResourceACLs([]string{"/customer"}, []string{"/customer/resetpassword", "/customer/forgetpassword"}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/forgetpassword?param=1", nil),
			nil,
			nil,
		},
		{
			"REGEX Blocks one resource /customer, we call /customer/forgetpassword which is whitelisted",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/customer"}, []string{"^/[a-z]+/resetpassword", "^/[a-z]+/forgetpassword"}), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/customer/forgetpassword?param=1", nil),
			nil,
			nil,
		},
		{
			"Blocks one resource /customer, we call /catalog/category and never authenticates",
			[]auth.Option{auth.WithResourceACLs([]string{"/customer"}, nil), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil),
			nil,
			nil,
		},
		{
			"REGEX Blocks one resource /customer, we call /catalog/category and never authenticates",
			[]auth.Option{auth.WithResourceRegexpACLs([]string{"^/customer"}, nil), auth.WithInvalidAuth()},
			httptest.NewRequest("GET", "http://corestore.io/catalog/category", nil),
			nil,
			nil,
		},
	}
	for i, test := range tests {
		srv, haveErr := auth.New(test.opts...)
		if test.wantConfigErrBhf != nil && haveErr != nil {
			assert.True(t, test.wantConfigErrBhf(haveErr), "Index %d %q\n%+v", i, test.desc, haveErr)
			continue
		}
		require.NoError(t, haveErr, "Index %d %q\n%+v", i, test.desc, haveErr)

		scpCfg, haveErr := srv.ConfigByScopeID(scope.DefaultTypeID, scope.DefaultTypeID)
		if test.wantConfigErrBhf != nil {
			assert.True(t, test.wantConfigErrBhf(haveErr), "Index %d %q\n%+v", i, test.desc, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d %q\n%+v", i, test.desc, haveErr)

		haveErr = scpCfg.Authenticate(test.req)
		if test.wantAuthErrBhf != nil {
			assert.True(t, test.wantAuthErrBhf(haveErr), "Index %d %q\n%+v", i, test.desc, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d %q\n%+v", i, test.desc, haveErr)
		}
	}
}
