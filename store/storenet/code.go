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

package storenet

import (
	"net/http"

	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// ParamName use in Cookie and JWT important when the user selects a different
// store within the current website/group context. This name will be used in a
// cookie or as key value in a token to permanently save the new selected store
// code.
const ParamName = `store`

// HTTPRequestParamStore name of the GET parameter to set a new store in a
// current website/group context.
const HTTPRequestParamStore = `___store`

// CodeFromCookie returns from a Request the value of the store cookie. Error
// behaviour: NotValid, NotSupported or an http.ErrNoCookie.
func CodeFromCookie(req *http.Request) (scope.Option, error) {
	keks, err := req.Cookie(ParamName)
	if err != nil {
		return scope.Option{}, errors.NewNotFound(err, "[storenet] Cookie")
	}
	return setByCode(keks.Value)
}

// CodeFromRequestGET returns from a Request form the value of the store code or
// an ErrStoreNotFound.
func CodeFromRequestGET(req *http.Request) (o scope.Option, err error) {
	return setByCode(req.URL.Query().Get(HTTPRequestParamStore))
}

func setByCode(scopeCode string) (scope.Option, error) {
	if err := store.CodeIsValid(scopeCode); err != nil {
		return scope.Option{}, errors.Wrap(err, "[storenet] Code not valid")
	}
	o, err := scope.SetByCode(scope.Store, scopeCode)
	if err != nil {
		return scope.Option{}, errors.Wrap(err, "[storenet] Unknown scope")
	}
	return o, nil
}
