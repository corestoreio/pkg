// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package store

import (
	"net/http"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
)

// CodeFromClaim returns a valid store code from a JSON web token or ErrStoreNotFound.
// Token argument is a map like being used by jwt.Token.Claims.
func CodeFromClaim(token map[string]interface{}) (o scope.Option, err error) {
	err = ErrStoreNotFound
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
	err = ErrStoreNotFound
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
	err = ErrStoreNotFound
	if req == nil {
		return
	}
	return setByCode(req.URL.Query().Get(HTTPRequestParamStore))
}

func setByCode(scopeCode string) (o scope.Option, err error) {
	err = CodeIsValid(scopeCode)
	if err == nil {
		o, err = scope.SetByCode(scopeCode, scope.StoreID)
	}
	return
}

// CodeIsValid checks if a store code is valid. Returns an ErrStoreCodeEmpty
// or an ErrStoreCodeInvalid if the first letter is not a-zA-Z and followed by
// a-zA-Z0-9_ or store code length is greater than 32 characters.
func CodeIsValid(c string) error {
	if c == "" || len(c) > 32 {
		return ErrStoreCodeInvalid
	}
	c1 := c[0]
	if false == ((c1 >= 'a' && c1 <= 'z') || (c1 >= 'A' && c1 <= 'Z')) {
		return ErrStoreCodeInvalid
	}
	if false == util.StrIsAlNum(c) {
		return ErrStoreCodeInvalid
	}
	return nil
}
