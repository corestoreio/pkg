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

package csjwt

import (
	"bytes"
	"encoding"
	"unicode"

	"github.com/corestoreio/errors"
)

// Verification allows to parse and verify a token with custom options.
type Verification struct {
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
		Methods: availableSigners,
	}
}

func (vf *Verification) unmarshal(src []byte, dst interface{}) error {
	src, err := DecodeSegment(src)
	if err != nil {
		return errors.WithStack(err)
	}

	if vf.Deserializer != nil {
		if err := vf.Deserializer.Deserialize(src, dst); err != nil {
			return errors.NotValid.New(err, errTokenMalformed)
		}
		return nil
	}
	// the order of the cases is important as JSON can be embedded but main type
	// has e.g. TextMarshaler.
	switch dt := dst.(type) {
	case encoding.BinaryUnmarshaler:
		if err := dt.UnmarshalBinary(src); err != nil {
			return errors.NotValid.New(err, errTokenMalformed)
		}
	case encoding.TextUnmarshaler:
		if err := dt.UnmarshalText(src); err != nil {
			return errors.NotValid.New(err, errTokenMalformed)
		}
	case interface{ Unmarshal([]byte) error }:
		if err := dt.Unmarshal(src); err != nil {
			return errors.NotValid.New(err, errTokenMalformed)
		}
	case interface{ UnmarshalJSON([]byte) error }:
		if err := dt.UnmarshalJSON(src); err != nil {
			return errors.NotValid.New(err, errTokenMalformed)
		}
	default:
		dec := jsonEncoding{}
		if err := dec.Deserialize(src, dst); err != nil {
			return errors.NotValid.New(err, errTokenMalformed)
		}
	}
	return nil
}

// Parse parses a rawToken into the destination token and may return an error.
// You must make sure to set the correct expected headers and claims in the
// template Token. The Header and Claims field in the destination token must be
// a pointer as the token itself. Error behaviour: Empty, NotFound, NotValid.
// Parse supports custom binary, text, json, protobuf decoding.
func (vf *Verification) Parse(dst *Token, rawToken []byte, keyFunc Keyfunc) error {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return errors.NotValid.Newf(errTokenInvalidSegmentCounts)
	}

	if dst.Header == nil || dst.Claims == nil {
		return errors.NotValid.Newf(errTokenBaseNil)
	}

	dst.Raw = rawToken

	if StartsWithBearer(dst.Raw) {
		return errors.NotValid.Newf(errTokenShouldNotContainBearer)
	}

	if err := vf.unmarshal(dst.Raw[:pos[0]], dst.Header); err != nil {
		return errors.WithStack(err)
	}
	if err := vf.unmarshal(dst.Raw[pos[0]+1:pos[1]], dst.Claims); err != nil {
		return errors.WithStack(err)
	}

	if err := dst.Claims.Valid(); err != nil {
		return errors.Wrap(err, errValidationClaimsInvalid)
	}

	if keyFunc == nil {
		return errors.Empty.Newf(errMissingKeyFunc)
	}
	key, err := keyFunc(dst)
	if err != nil {
		return errors.NotValid.Newf(errTokenUnverifiable, err)
	}

	method, err := vf.getMethod(dst)
	if err != nil {
		return errors.Wrap(err, "[csjwt] Verification.Parse.getMethod")
	}

	dst.Signature = dst.Raw[pos[1]+1:]
	if err := method.Verify(dst.Raw[:pos[1]], dst.Signature, key); err != nil {
		return errors.NotValid.Newf(errSignatureInvalid, err, dst)
	}

	dst.Valid = true
	return nil
}

// ParseUnverified parses a rawToken into the unverified destination token and
// may return an error. Uses this function only during testing.
func (vf *Verification) ParseUnverified(dst *Token, rawToken []byte) error {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return errors.NotValid.Newf(errTokenInvalidSegmentCounts)
	}

	if dst.Header == nil || dst.Claims == nil {
		return errors.NotValid.Newf(errTokenBaseNil)
	}

	dst.Raw = rawToken

	if StartsWithBearer(dst.Raw) {
		return errors.NotValid.Newf(errTokenShouldNotContainBearer)
	}

	if err := vf.unmarshal(dst.Raw[:pos[0]], dst.Header); err != nil {
		return errors.WithStack(err)
	}
	if err := vf.unmarshal(dst.Raw[pos[0]+1:pos[1]], dst.Claims); err != nil {
		return errors.WithStack(err)
	}

	if err := dst.Claims.Valid(); err != nil {
		return errors.Wrap(err, errValidationClaimsInvalid)
	}

	// Lookup signature method
	_, err := vf.getMethod(dst)
	if err != nil {
		return errors.Wrap(err, "[csjwt] Verification.Parse.getMethod")
	}

	dst.Signature = dst.Raw[pos[1]+1:]

	dst.Valid = false
	return nil
}

func (vf *Verification) getMethod(t *Token) (Signer, error) {
	if len(vf.Methods) == 0 {
		return nil, errors.Empty.Newf(errVerificationMethodsEmpty)
	}

	alg := t.Alg()
	if alg == "" {
		return nil, errors.Empty.Newf(errAlgorithmEmpty, t.Header)
	}

	for _, m := range vf.Methods {
		if m.Alg() == alg {
			return m, nil
		}
	}
	return nil, errors.NotFound.Newf(errAlgorithmNotFound, alg, vf.Methods)
}

// SplitForVerify splits the token into two parts: the payload and the
// signature. An error gets returned if the number of dots don't match with the
// JWT standard.
func SplitForVerify(rawToken []byte) (signingString, signature []byte, err error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return nil, nil, errors.NotValid.Newf(errTokenInvalidSegmentCounts)
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

// StartsWithBearer checks if token starts with bearer
func StartsWithBearer(token []byte) bool {
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
