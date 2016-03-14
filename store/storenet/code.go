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
)

// CodeFromClaim returns a valid store code from a JSON web token or ErrStoreNotFound.
// Token argument is a map like being used by jwt.Token.Claims.
func CodeFromClaim(token map[string]interface{}) (o scope.Option, err error) {
	err = store.ErrStoreNotFound
	if 0 == len(token) {
		return
	}

	tokVal, ok := token[ParamName]
	scopeCode, okcs := tokVal.(string)

	if okcs && ok {
		return setByCode(scopeCode)
	}
	return
}

// CodeFromCookie returns from a Request the value of the store cookie or
// an ErrStoreNotFound.
func CodeFromCookie(req *http.Request) (o scope.Option, err error) {
	err = store.ErrStoreNotFound
	if nil == req {
		return
	}
	var keks *http.Cookie
	keks, err = req.Cookie(ParamName)
	if err != nil {
		return
	}
	return setByCode(keks.Value)
}

// CodeFromRequestGET returns from a Request form the value of the store code or
// an ErrStoreNotFound.
func CodeFromRequestGET(req *http.Request) (o scope.Option, err error) {
	err = store.ErrStoreNotFound
	if req == nil {
		return
	}
	return setByCode(req.URL.Query().Get(HTTPRequestParamStore))
}

func setByCode(scopeCode string) (o scope.Option, err error) {
	err = store.CodeIsValid(scopeCode)
	if err == nil {
		o, err = scope.SetByCode(scopeCode, scope.StoreID)
	}
	return
}

// CodeAddToClaim adds the store code to a JSON web token.
// tokenClaim may be *jwt.Token.Claim
func CodeAddToClaim(s *store.Store, tokenClaim map[string]interface{}) {
	tokenClaim[ParamName] = s.Data.Code.String
}
