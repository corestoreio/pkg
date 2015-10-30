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

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/utils"
)

// StoreCodeFromClaim returns a valid store code from a JSON web token or ErrStoreNotFound.
// Token argument is a map like being used by jwt.Token.Claims.
func StoreCodeFromClaim(token map[string]interface{}) (o scope.Option, err error) {
	err = ErrStoreNotFound
	if 0 == len(token) {
		return
	}

	tokVal, ok := token[CookieName]
	scopeCode, okcs := tokVal.(string)

	if okcs && ok {
		return setByCode(scopeCode)
	}
	return
}

// StoreCodeFromCookie returns from a Request the value of the store cookie or
// an ErrStoreNotFound.
func StoreCodeFromCookie(req *http.Request) (o scope.Option, err error) {
	err = ErrStoreNotFound
	if nil == req {
		return
	}
	var keks *http.Cookie
	keks, err = req.Cookie(CookieName)
	if err != nil {
		return
	}
	return setByCode(keks.Value)
}

// StoreCodeFromRequestGET returns from a Request form the value of the store code or
// an ErrStoreNotFound.
func StoreCodeFromRequestGET(req *http.Request) (o scope.Option, err error) {
	err = ErrStoreNotFound
	if req == nil {
		return
	}
	return setByCode(req.URL.Query().Get(HTTPRequestParamStore))
}

func setByCode(scopeCode string) (o scope.Option, err error) {
	err = ValidateStoreCode(scopeCode)
	if err == nil {
		o, err = scope.SetByCode(scopeCode, scope.StoreID)
	}
	return
}

// ValidateStoreCode checks if a store code is valid. Returns an ErrStoreCodeEmpty
// or an ErrStoreCodeInvalid if the first letter is not a-zA-Z and followed by
// a-zA-Z0-9_ or store code length is greater than 32 characters.
func ValidateStoreCode(c string) error {
	if c == "" {
		return ErrStoreCodeEmpty
	}
	if len(c) > 32 {
		return ErrStoreCodeInvalid
	}
	c1 := c[0]
	if false == ((c1 >= 'a' && c1 <= 'z') || (c1 >= 'A' && c1 <= 'Z')) {
		return ErrStoreCodeInvalid
	}
	if false == utils.StrIsAlNum(c) {
		return ErrStoreCodeInvalid
	}
	return nil
}
