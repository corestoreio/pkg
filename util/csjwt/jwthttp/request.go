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

package jwthttp

import (
	"net/http"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/csjwt"
)

// HTTPHeaderAuthorization identifies the bearer token in this header key
const HTTPHeaderAuthorization = `Authorization`

// HTTPFormInputName default name for the HTML form field name
const HTTPFormInputName = `access_token`

// NewVerification creates new verification parser with the default signing
// method HS256, if availableSigners slice argument is empty. Nil arguments are
// forbidden.
func NewVerification(availableSigners ...csjwt.Signer) *Verification {
	return &Verification{
		Verification: csjwt.Verification{
			Methods: availableSigners,
		},
	}
}

// Verification allows to parse and verify a token with custom options.
type Verification struct {
	csjwt.Verification
	// FormInputName defines the name of the HTML form input type in which the
	// token has been stored. If empty, the form the gets ignored.
	FormInputName string
	// CookieName defines the name of the cookie where the token has been
	// stored. If empty, cookie parsing gets ignored.
	CookieName string
	// ExtractTokenFn for extracting a token from an HTTP request. The
	// ExtractToken method should return a token string or an error.
	// This function can be nil
	ExtractTokenFn func(*http.Request) (string, error)
}

// ParseFromRequest same as Parse but extracts the token from a request. First
// it searches for the token bearer in the header HTTPHeaderAuthorization. If
// not found the request POST form gets parsed and the FormInputName gets used
// to lookup the token value.
func (vf *Verification) ParseFromRequest(dst *csjwt.Token, keyFunc csjwt.Keyfunc, req *http.Request) error {
	if vf.ExtractTokenFn != nil {
		tkn, err := vf.ExtractTokenFn(req)
		if err != nil {
			return errors.WithStack(err)
		}
		return vf.Parse(dst, []byte(tkn), keyFunc)
	}

	// Look for an Authorization header
	if ah := req.Header.Get(HTTPHeaderAuthorization); ah != "" {
		// Should be a bearer token
		auth := []byte(ah)
		if csjwt.StartsWithBearer(auth) {
			return vf.Parse(dst, auth[7:], keyFunc)
		}
	}

	if vf.CookieName != "" {
		if err := vf.parseCookie(dst, keyFunc, req); err != nil && err != http.ErrNoCookie {
			return errors.WithStack(err)
		}
		if dst.Valid {
			return nil
		}
		// try next, the form
	}

	if vf.FormInputName != "" {
		return vf.parseForm(dst, keyFunc, req)
	}

	return errors.NotFound.Newf(errTokenNotInRequest)
}

func (vf *Verification) parseCookie(dst *csjwt.Token, keyFunc csjwt.Keyfunc, req *http.Request) error {
	keks, err := req.Cookie(vf.CookieName)
	if keks != nil && keks.Value != "" {
		return vf.Parse(dst, []byte(keks.Value), keyFunc)
	}
	return err // error can be http.ErrNoCookie
}

func (vf *Verification) parseForm(dst *csjwt.Token, keyFunc csjwt.Keyfunc, req *http.Request) error {
	_ = req.ParseMultipartForm(10e6) // ignore errors
	if tokStr := req.Form.Get(vf.FormInputName); tokStr != "" {
		return vf.Parse(dst, []byte(tokStr), keyFunc)
	}
	return errors.NotFound.Newf(errTokenNotInRequest)
}

const (
	errTokenNotInRequest = `[csjwt] token not present in request`
)
