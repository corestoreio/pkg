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

package backendauth_test

import (
	"net/http"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/net/auth"
	"github.com/corestoreio/csfw/net/auth/backendauth"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
)

type fataler interface {
	Fatal(args ...interface{})
}

func mustToPath(fa fataler, f func(s scope.Type, scopeID int64) (cfgpath.Path, error), s scope.Type, scopeID int64) string {
	p, err := f(s, scopeID)
	if err != nil {
		fa.Fatal(errors.PrintLoc(err))
	}
	return p.String()
}

func reqWithStore(method string, cfgOpt ...cfgmock.OptionFunc) *http.Request {
	req, err := http.NewRequest(method, "http://corestore.io/foo", nil)
	if err != nil {
		panic(err)
	}

	return req.WithContext(
		store.WithContextRequestedStore(req.Context(), storemock.MustNewStoreAU(cfgmock.NewService(cfgOpt...))),
	)
}

func newService() *auth.Service {
	return auth.MustNew(
		auth.WithOptionFactory(backendauth.PrepareOptions(backend)),
	)
}

//
//func TestBackend_Path_Errors(t *testing.T) {
//
//	tests := []struct {
//		toPath func(s scope.Scope, scopeID int64) (cfgpath.Path, error)
//		val    interface{}
//		errBhf errors.BehaviourFunc
//	}{
//		{backend.NetCorsExposedHeaders.ToPath, struct{}{}, errors.IsNotValid},
//		{backend.NetCorsAllowedOrigins.ToPath, struct{}{}, errors.IsNotValid},
//		{backend.NetCorsAllowOriginRegex.ToPath, struct{}{}, errors.IsNotValid},
//		{backend.NetCorsAllowOriginRegex.ToPath, "[a-z+", errors.IsFatal},
//		{backend.NetCorsAllowedMethods.ToPath, struct{}{}, errors.IsNotValid},
//		{backend.NetCorsAllowedHeaders.ToPath, struct{}{}, errors.IsNotValid},
//		{backend.NetCorsAllowCredentials.ToPath, struct{}{}, errors.IsNotValid},
//		{backend.NetCorsOptionsPassthrough.ToPath, struct{}{}, errors.IsNotValid},
//		{backend.NetCorsMaxAge.ToPath, struct{}{}, errors.IsNotValid},
//	}
//	for i, test := range tests {
//
//		scpFnc := backendauth.PrepareOptions(backend)
//		cfgSrv := cfgmock.NewService(cfgmock.WithPV(cfgmock.PathValue{
//			mustToPath(t, test.toPath, scope.Website, 2): test.val,
//		}))
//		cfgScp := cfgSrv.NewScoped(2, 0)
//
//		_, err := auth.New(scpFnc(cfgScp)...)
//		assert.True(t, test.errBhf(err), "Index %d Error: %s", i, err)
//	}
//}
