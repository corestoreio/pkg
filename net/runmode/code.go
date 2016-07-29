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

package runmode

import (
	"net/http"

	"github.com/corestoreio/csfw/store"
)

// CodeExtracter knows how to extract a website/store code from
// a request.
type CodeExtracter interface {
	// FromRequest returns the code and valid which has three values: 0 not
	// valid, 10 valid and code found in GET query string, 20 valid and code
	// found in cookie.
	FromRequest(req *http.Request) (code string, valid int8)
}

// FieldName use in Cookies and JSON Web Tokens (JWT) to identify an active
// store besides from the default loaded store. Default value.
const FieldName = `store`

// URLFieldName name of the GET parameter to set a new store in a current
// website/group context. Default value.
const URLFieldName = `___store`

// ExtractCode can extract the website or store code from an HTTP Request. This
// code is then responsible for initiating the runMode.
type ExtractCode struct {
	// FieldName optional custom name, defaults to constant FieldName
	FieldName string
	// URLFieldName optional custom name, defaults to constant URLFieldName
	URLFieldName string
}

// FromRequest returns from a GET request with a query string the value of the
// website/store code. If no code can be found in the query string, this
// function falls back to the cookie name defined in field FieldName. Valid has
// three values: 0 not valid, 10 valid and code found in GET query string, 20
// valid and code found in cookie.
func (c ExtractCode) FromRequest(req *http.Request) (code string, valid int8) {
	// todo find a better solution for the valid type
	hps := URLFieldName
	if c.URLFieldName != "" {
		hps = c.URLFieldName
	}
	code = req.URL.Query().Get(hps)
	if code == "" {
		return c.fromCookie(req)
	}
	if err := store.CodeIsValid(code); err == nil {
		valid = 10
	}
	return
}

// fromCookie extracts a store from a cookie using the field name FieldName as
// an identifier.
func (c ExtractCode) fromCookie(req *http.Request) (code string, valid int8) {
	p := FieldName
	if c.FieldName != "" {
		p = c.FieldName
	}
	if keks, err := req.Cookie(p); err == nil {
		code = keks.Value
	}
	if err := store.CodeIsValid(code); err == nil {
		valid = 20
	}
	return
}
