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

import (
	"bytes"
	"net/http"
	"unicode"

	"github.com/corestoreio/csfw/util/errors"
)

// HTTPHeaderAuthorization identifies the bearer token in this header key
const HTTPHeaderAuthorization = `Authorization`

// HTTPFormInputName default name for the HTML form field name
const HTTPFormInputName = `access_token`

// Verification allows to parse and verify a token with custom options.
type Verification struct {
	// FormInputName defines the name of the HTML form input type in which the
	// token has been stored. If empty, the form the gets ignored.
	FormInputName string
	// CookieName defines the name of the cookie where the token has been stored. If
	// empty, cookie parsing gets ignored.
	CookieName string
	// Methods for verifying and signing a token
	Methods SignerSlice

	// Decoder interface to pass in a custom decoder parser. Can be nil, falls
	// back to JSON.
	Deserializer
}

// NewVerification creates new verification parser with the default signing
// method HS256, if availableSigners slice argument is empty. Nil arguments are
// forbidden.
func NewVerification(availableSigners ...Signer) *Verification {
	return &Verification{
		Methods:      availableSigners,
		Deserializer: JSONEncoding{},
	}
}

// Parse parses a rawToken into the destination token and may return an error.
// You must make sure to set the correct expected headers and claims in the
// template Token. The Header and Claims field in the destination token must be
// a pointer as the token itself. Error behaviour: Empty, NotFound, NotValid
func (vf *Verification) Parse(dst *Token, rawToken []byte, keyFunc Keyfunc) error {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return errors.NewNotValidf(errTokenInvalidSegmentCounts)
	}

	if dst.Header == nil || dst.Claims == nil {
		return errors.NewNotValidf(errTokenBaseNil)
	}

	dec := vf.Deserializer
	if dec == nil {
		dec = JSONEncoding{}
	}

	dst.Raw = rawToken

	if startsWithBearer(dst.Raw) {
		return errors.NewNotValidf(errTokenShouldNotContainBearer)
	}

	// parse Header
	if err := dec.Deserialize(dst.Raw[:pos[0]], dst.Header); err != nil {
		return errors.NewNotValid(err, errTokenMalformed)
	}

	// parse Claims
	if err := dec.Deserialize(dst.Raw[pos[0]+1:pos[1]], dst.Claims); err != nil {
		return errors.NewNotValid(err, errTokenMalformed)
	}

	// validate Claims
	if err := dst.Claims.Valid(); err != nil {
		return errors.Wrap(err, errValidationClaimsInvalid)
	}

	// Lookup key
	if keyFunc == nil {
		return errors.NewEmptyf(errMissingKeyFunc)
	}
	key, err := keyFunc(dst)
	if err != nil {
		return errors.NewNotValid(err, errTokenUnverifiable)
	}

	// Lookup signature method
	method, err := vf.getMethod(dst)
	if err != nil {
		return errors.Wrap(err, "[csjwt] Verification.Parse.getMethod")
	}

	// Perform validation
	dst.Signature = dst.Raw[pos[1]+1:]
	if err := method.Verify(dst.Raw[:pos[1]], dst.Signature, key); err != nil {
		return errors.NewNotValid(err, errSignatureInvalid)
	}

	dst.Valid = true
	return nil
}

func (vf *Verification) getMethod(t *Token) (Signer, error) {

	if len(vf.Methods) == 0 {
		return nil, errors.NewEmptyf("[csjwt] No methods supplied to the Verfication Method slice")
	}

	alg := t.Alg()
	if alg == "" {
		return nil, errors.NewEmptyf("[csjwt] Cannot find alg entry in token header: %#v", t.Header)
	}

	for _, m := range vf.Methods {
		if m.Alg() == alg {
			return m, nil
		}
	}
	return nil, errors.NewNotFoundf("[csjwt] Algorithm %q not found in method list %q", alg, vf.Methods)
}

// ParseFromRequest same as Parse but extracts the token from a request. First
// it searches for the token bearer in the header HTTPHeaderAuthorization. If
// not found the request POST form gets parsed and the FormInputName gets used
// to lookup the token value.
func (vf *Verification) ParseFromRequest(dst *Token, keyFunc Keyfunc, req *http.Request) error {
	// Look for an Authorization header
	if ah := req.Header.Get(HTTPHeaderAuthorization); ah != "" {
		// Should be a bearer token
		auth := []byte(ah)
		if startsWithBearer(auth) {
			return vf.Parse(dst, auth[7:], keyFunc)
		}
	}

	if vf.CookieName != "" {
		if err := vf.parseCookie(dst, keyFunc, req); err != nil && err != http.ErrNoCookie {
			return errors.Wrap(err, "[csjwt] Verification.ParseFromRequest.parseCookie")
		}
		if dst.Valid {
			return nil
		}
		// try next, the form
	}

	if vf.FormInputName != "" {
		return vf.parseForm(dst, keyFunc, req)
	}

	return errors.NewNotFoundf(errTokenNotInRequest)
}

func (vf *Verification) parseCookie(dst *Token, keyFunc Keyfunc, req *http.Request) error {
	keks, err := req.Cookie(vf.CookieName)
	if keks != nil && keks.Value != "" {
		return vf.Parse(dst, []byte(keks.Value), keyFunc)
	}
	return err // error can be http.ErrNoCookie
}

func (vf *Verification) parseForm(dst *Token, keyFunc Keyfunc, req *http.Request) error {
	_ = req.ParseMultipartForm(10e6) // ignore errors
	if tokStr := req.Form.Get(vf.FormInputName); tokStr != "" {
		return vf.Parse(dst, []byte(tokStr), keyFunc)
	}
	return errors.NewNotFoundf(errTokenNotInRequest)
}

// SplitForVerify splits the token into two parts: the payload and the
// signature. An error gets returned if the number of dots don't match with the
// JWT standard.
func SplitForVerify(rawToken []byte) (signingString, signature []byte, err error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return nil, nil, errors.NewNotValidf(errTokenInvalidSegmentCounts)
	}
	return rawToken[:pos[1]], rawToken[pos[1]+1:], nil
}

// dotPositions returns the position of the dots within the token slice and if
// the amount of dots are valid for a JWT.
func dotPositions(t []byte) (pos [2]int, valid bool) {
	const aDot byte = '.'
	c := 0
	for i, b := range t {
		if b == aDot {
			if c < 2 {
				pos[c] = i
			}
			c++
		}
	}
	if c == 2 {
		valid = true
	}
	return
}

// length of the string "bearer "
const prefixBearerLen = 7

var prefixBearer = []byte(`bearer `)

// startsWithBearer checks if token starts with bearer
func startsWithBearer(token []byte) bool {
	if len(token) <= prefixBearerLen {
		return false
	}
	var havePrefix [prefixBearerLen]byte
	copy(havePrefix[:], token[0:prefixBearerLen])
	for i, b := range havePrefix {
		havePrefix[i] = byte(unicode.ToLower(rune(b)))
	}
	return bytes.Equal(havePrefix[:], prefixBearer)
}
