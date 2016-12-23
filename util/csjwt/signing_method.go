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

	"github.com/corestoreio/errors"
)

// Signer interface to add new methods for signing or verifying tokens.
type Signer interface {
	// Verify returns nil if signature is valid
	Verify(signingString, signature []byte, key Key) error
	// Sign returns encoded signature or error
	Sign(signingString []byte, key Key) ([]byte, error)
	// Alg returns the alg identifier for this method (example: 'HS256')
	Alg() string
}

// All available algorithms which are supported by this package.
const (
	ES256      = `ES256`
	ES384      = `ES384`
	ES512      = `ES512`
	HS256      = `HS256`
	HS384      = `HS384`
	HS512      = `HS512`
	PS256      = `PS256`
	PS384      = `PS384`
	PS512      = `PS512`
	RS256      = `RS256`
	RS384      = `RS384`
	RS512      = `RS512`
	ES         = `ES`
	HS         = `HS`
	PS         = `PS`
	RS         = `RS`
	Blake2b256 = `blk2b256`
	Blake2b512 = `blk2b512`
)

// SigningMethodFactory creates a new signing method by an algorithm. Supported
// algorithms are: ES, HS, PS and RS, all within 256-512 and they do not need a
// symmetric key. Returns an error for an unknown signing method.
func SigningMethodFactory(alg string) (s Signer, err error) {
	switch alg {

	case ES256:
		s = NewSigningMethodES256()
	case ES384:
		s = NewSigningMethodES384()
	case ES512:
		s = NewSigningMethodES512()

	case HS256:
		s = NewSigningMethodHS256()
	case HS384:
		s = NewSigningMethodHS384()
	case HS512:
		s = NewSigningMethodHS512()

	case PS256:
		s = NewSigningMethodPS256()
	case PS384:
		s = NewSigningMethodPS384()
	case PS512:
		s = NewSigningMethodPS512()

	case RS256:
		s = NewSigningMethodRS256()
	case RS384:
		s = NewSigningMethodRS384()
	case RS512:
		s = NewSigningMethodRS512()

	}
	if s == nil {
		err = errors.NewNotSupportedf("[csjwt] Unknown signing algorithm %q", alg)
	}
	return s, err
}

// MustSigningMethodFactory same as SigningMethodFactory but panics on error.
// You should only use the Must* functions during init process or testing.
func MustSigningMethodFactory(alg string) Signer {
	s, err := SigningMethodFactory(alg)
	if err != nil {
		panic(err)
	}
	return s
}

// SignerSlice helper type
type SignerSlice []Signer

// String returns a list of algorithms, comma separated
func (ms SignerSlice) String() string {
	var buf bytes.Buffer
	for i, m := range ms {
		_, _ = buf.WriteString(m.Alg())
		if i < len(ms)-1 {
			_, _ = buf.WriteString(`, `)
		}
	}
	return buf.String()
}

// Contains checks if the algorithm has already been added
func (ms SignerSlice) Contains(alg string) bool {
	for _, m := range ms {
		if m.Alg() == alg {
			return true
		}
	}
	return false
}
