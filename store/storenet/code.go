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
)

// ParamName use in Cookie and JWT important when the user selects a different
// store within the current website/group context. This name will be used in a
// cookie or as key value in a token to permanently save the new selected store
// code.
const ParamName = `store`

// HTTPRequestParamStore name of the GET parameter to set a new store in a
// current website/group context.
const HTTPRequestParamStore = `___store`

// CodeFromRequest returns from a request form, with the field name of the value
// of constant HTTPRequestParamStore, the value of the store code. If empty it
// falls back to the cookie name defined in constant ParamName. An error of
// NotValid gets returned if the code cannot be validated.
func CodeFromRequest(req *http.Request) (code string, valid bool) {
	code = req.URL.Query().Get(HTTPRequestParamStore)
	if code == "" {
		code, _ = CodeFromCookie(req)
	}

	if err := store.CodeIsValid(code); err == nil {
		valid = true
	}
	return
}

func CodeFromCookie(req *http.Request) (code string, valid bool) {
	if keks, err := req.Cookie(ParamName); err == nil {
		code = keks.Value
	}
	if err := store.CodeIsValid(code); err == nil {
		valid = true
	}
	return
}
