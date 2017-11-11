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

	"github.com/corestoreio/cspkg/util/hashpool"
	"github.com/corestoreio/errors"
	_ "golang.org/x/crypto/blake2b"
)

// NewSigningMethodBlake2b256 creates a new HMAC-Blake2b hash with a preset
// password and 32-byte checksum. Blake2b uses SIMD optimizations via ASM code.
// https://blake2.net/
func NewSigningMethodBlake2b256(key Key) (Signer, error) {
	if key.Error != nil {
		return nil, errors.Wrap(key.Error, "[csjwt] NewBlake2b256.key")
	}
	if len(key.hmacPassword) == 0 {
		return nil, errors.NewEmptyf(errHmacPasswordEmpty)
	}

	return &SigningMethodHSFast{
		Name: Blake2b256,
		ht: hashpool.New(func() hash.Hash {
			return hmac.New(crypto.BLAKE2b_256.New, key.hmacPassword)
		}),
	}, nil
}

// NewSigningMethodBlake2b512 creates a new HMAC-Blake2b hash with a preset
// password and 64-byte checksum. Blake2b uses SIMD optimizations via ASM code.
// https://blake2.net/
func NewSigningMethodBlake2b512(key Key) (Signer, error) {
	if key.Error != nil {
		return nil, errors.Wrap(key.Error, "[csjwt] NewBlake2b512.key")
	}
	if len(key.hmacPassword) == 0 {
		return nil, errors.NewEmptyf(errHmacPasswordEmpty)
	}

	return &SigningMethodHSFast{
		Name: Blake2b512,
		ht:   hashpool.New(crypto.BLAKE2b_512.New),
	}, nil
}
