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
	"crypto"
	"crypto/hmac"
	"hash"

	"github.com/corestoreio/csfw/util/hashpool"
	"github.com/corestoreio/errors"
)

// SigningMethodHSFast implements the HMAC-SHA family of pre-warmed signing
// methods. Less allocations, bytes and a little bit faster but maybe the
// underlying mutex can become the bottleneck.
type SigningMethodHSFast struct {
	Name string
	ht   hashpool.Tank
}

func newHSFast(a string, h crypto.Hash, key Key) (Signer, error) {
	if key.Error != nil {
		return nil, errors.Wrap(key.Error, "[csjwt] newHMACFast.key")
	}
	if len(key.hmacPassword) == 0 {
		return nil, errors.NewEmptyf(errHmacPasswordEmpty)
	}
	// Can we use the specified hashing method?
	if !h.Available() {
		return nil, errors.NewNotImplementedf(errHmacHashUnavailable)
	}
	return &SigningMethodHSFast{
		Name: a,
		ht: hashpool.New(func() hash.Hash {
			return hmac.New(h.New, key.hmacPassword)
		}),
	}, nil
}

// NewSigningMethodHS256Fast creates a new HMAC-SHA hash with a preset password
// and does not register it globally. It uses internally a sync.Pool hashes.
func NewSigningMethodHS256Fast(key Key) (Signer, error) {
	return newHSFast(HS256, crypto.SHA256, key)
}

// NewSigningMethodHS384Fast creates a new HMAC-SHA hash with a preset password
// and does not register it globally. It uses internally a sync.Pool hashes.
func NewSigningMethodHS384Fast(key Key) (Signer, error) {
	return newHSFast(HS384, crypto.SHA384, key)
}

// NewSigningMethodHS512Fast creates a new HMAC-SHA hash with a preset password
// and does not register it globally. It uses internally a sync.Pool hashes.
func NewSigningMethodHS512Fast(key Key) (Signer, error) {
	return newHSFast(HS512, crypto.SHA512, key)
}

func (m *SigningMethodHSFast) Alg() string {
	return m.Name
}

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
// Error behaviour: NotImplemented, WriteFailed, NotValid
func (m *SigningMethodHSFast) Verify(signingString, signature []byte, _ Key) error {

	// Decode signature, for comparison
	sig, err := DecodeSegment(signature)
	if err != nil {
		return errors.Wrap(err, "[csjwt] SigningMethodHMACFast.Verify.DecodeSegment")
	}

	// This signing method is symmetric, so we validate the signature by
	// reproducing the signature from the signing string and key, then comparing
	// that against the provided signature.
	hasher := m.ht.Get()
	defer m.ht.Put(hasher)

	if _, err := hasher.Write(signingString); err != nil {
		return errors.NewWriteFailed(err, "[csjwt] SigningMethodHMACFast.Verify.hasher.Write")
	}

	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return errors.NewNotValidf(errHmacSignatureInvalid)
	}

	// No validation errors.  Signature is good.
	return nil
}

// Sign implements the Sign method from SigningMethod interface.
// Error behaviour: WriteFailed
func (m *SigningMethodHSFast) Sign(signingString []byte, _ Key) ([]byte, error) {

	hasher := m.ht.Get()
	defer m.ht.Put(hasher)

	if _, err := hasher.Write(signingString); err != nil {
		return nil, errors.NewWriteFailed(err, "[csjwt] SigningMethodHMACFast.Sign.hasher.Write")
	}

	return EncodeSegment(hasher.Sum(nil)), nil
}
