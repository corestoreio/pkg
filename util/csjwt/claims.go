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
	"fmt"
	"time"

	"github.com/juju/errors"
)

// Claimer for a type to be a Claims object
type Claimer interface {
	// Valid method that determines if the token is invalid for any supported reason.
	// Returns nil on success
	Valid() error
	// Expires declares when a token expires. A duration smaller or equal
	// to zero means that the token has already expired.
	// Useful when adding a token to a blacklist.
	Expires() time.Duration
	// Set sets a value to the claim and may overwrite existing values
	Set(key string, value interface{}) error
	// Get retrieves a value from the claim.
	Get(key string) (value interface{}, err error)
	// Keys returns a list of all available keys
	Keys() []string
}

// (CS) I personally don't like the Set() and Get() functions but there is no
// other way around it.

const (
	headerAlg = "alg"
	headerTyp = "typ"
)

// Header defines the contract for a type to act like a header. It must be able
// to marshal and unmarshal itself.
// The members of the JSON object represented by the Decoded JWT Header Segment
// describe the signature applied to the JWT Header Segment and the JWT Payload
// Segment and optionally additional properties of the JWT. Implementations MUST
// understand the entire contents of the header; otherwise, the JWT MUST be
// rejected for processing.
type Header interface {
	// Alg returns the name of the underlying algorithm
	Alg() string
	// Typ identifies the type of the JSON web token
	Typ() string
	Set(key, value string) error
	Get(key string) (value string, err error)
}

// Head minimum default header.
// To extend this header please use the struct jwtclaim.HeadSegments
type Head struct {
	// Alg (algorithm) header parameter identifies the cryptographic algorithm
	// used to secure the JWT. A list of reserved alg values is in Table 4. The
	// processing of the "alg" (algorithm) header parameter, if present, requires
	// that the value of the "alg" header parameter MUST be one that is both
	// supported and for which there exists a key for use with that algorithm
	// associated with the issuer of the JWT. This header parameter is REQUIRED.
	Algorithm string `json:"alg,omitempty"`
	// Typ (type) header parameter is used to declare that this data structure
	// is a JWT. If a "typ" parameter is present, it is RECOMMENDED that its
	// value be "JWT". This header parameter is OPTIONAL.
	Type string `json:"typ,omitempty"`
}

// NewHead creates a new minimum default header.
// Arguments alg can be optionally applied one time to define an algorithm
// but in all cases the algorithm gets set by the signing method.
// For test cases you can pass an algorithm argument.
// To extend this header please use the struct jwtclaim.HeadSegments
func NewHead(alg ...string) *Head {
	var a string
	if len(alg) > 0 {
		a = alg[0]
	}
	return &Head{
		Algorithm: a,
		Type:      ContentTypeJWT,
	}
}

// Alg returns the underlying algorithm.
func (s *Head) Alg() string {
	return s.Algorithm
}

// Typ returns the token type.
func (s *Head) Typ() string {
	return s.Type
}

func (s *Head) String() string {
	if s.Algorithm == "" {
		return `csjwt.NewHead()`
	}
	return fmt.Sprintf("csjwt.NewHead(%q)", s.Alg())
}

func (s *Head) GoString() string {
	return s.String()
}

const errHeaderKeyNotSupported = "[csjwt] Header %q not yet supported. Please switch to type jwtclaim.HeadSegments."

// Set sets a value. Key must be one of the constants Header*.
func (s *Head) Set(key, value string) (err error) {
	switch key {
	case headerAlg:
		s.Algorithm = value
	case headerTyp:
		s.Type = value
	default:
		return errors.Errorf(errHeaderKeyNotSupported, key)
	}
	return err
}

// Get returns a value or nil or an error. Key must be one of the constants Header*.
func (s *Head) Get(key string) (value string, err error) {
	switch key {
	case headerAlg:
		return s.Algorithm, nil
	case headerTyp:
		return s.Type, nil
	}
	return "", errors.Errorf(errHeaderKeyNotSupported, key)
}

// MergeClaims merges the sources Claimers into the destination claimer existing token claims and overwrites
// existing entries. Destination Claimer must be a pointer.
func MergeClaims(dst Claimer, srcs ...Claimer) error {
	if dst == nil || len(srcs) == 0 {
		return nil
	}

	for idx, c := range srcs {
		for _, k := range c.Keys() {
			v, err := c.Get(k)
			if err != nil {
				return errors.Errorf("[csjwt] Cannot get Key %q from Claim index %d. Error: %s", k, idx, err)
			}
			if err := dst.Set(k, v); err != nil {
				return errors.Errorf("[csjwt] Cannot set Key %q with value `%v'. Claim index %d. Error: %s", k, v, idx, err)
			}
		}
	}
	return nil
}
