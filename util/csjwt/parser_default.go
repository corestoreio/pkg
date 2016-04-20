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

package csjwt

import "net/http"

var (
	defaultHSVerification *Verification
)

func init() {
	defaultHSVerification = NewVerification(NewSigningMethodHS256(), NewSigningMethodHS384(), NewSigningMethodHS512())
	defaultHSVerification.CookieName = HTTPFormInputName
	defaultHSVerification.FormInputName = HTTPFormInputName
}

// Parse parses a rawToken into the template token and returns the fully parsed and
// verified Token, or an error. You must make sure to set the correct expected
// headers and claims in the template Token. The Header and Claims field in the
// template token must be a pointer.
//
// Default configuration with defined CookieName of constant HTTPFormInputName,
// defined FormInputNam with value of constant HTTPFormInputName and supported
// Signers of HS256, HS384 and HS512.
func Parse(template Token, rawToken []byte, keyFunc Keyfunc) (Token, error) {
	return defaultHSVerification.Parse(template, rawToken, keyFunc)
}

// ParseFromRequest same as Parse but extracts the token from a request.
// First it searches for the token bearer in the header HTTPHeaderAuthorization.
// If not found, the cookie gets parsed and if not found then the request POST
// form gets parsed.
//
// Default configuration with defined CookieName of constant HTTPFormInputName,
// defined FormInputNam with value of constant HTTPFormInputName and supported
// Signers of HS256, HS384 and HS512.
func ParseFromRequest(template Token, keyFunc Keyfunc, req *http.Request) (Token, error) {
	return defaultHSVerification.ParseFromRequest(template, keyFunc, req)
}
