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

package signedblake2

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/minio/blake2b-simd"
)

// OptionName identifies this package within the register of the
// backendsigned.Backend type.
const OptionName = `blake2`

// WithContentHMAC256 applies the very fast Blake2 hashing algorithm.
// The current package has been optimized with ASM with for x64 systems, hence
// Blake2 is faster than SHA.
func WithContentHMAC256(scp scope.Scope, id int64, key []byte) signed.Option {
	return func(s *signed.Service) error {
		if err := signed.WithHash(scp, id, blake2b.New256, key)(s); err != nil {
			return errors.Wrap(err, "[signedblake2] WithContentHMAC_Blake2b256.WithHash")
		}
		sig := signed.NewHMAC("blk2b256")
		return signed.WithHeaderHandler(scp, id, sig)(s)
	}
}

// WithTransparent256 applies the very fast Blake2 hashing algorithm.
// The current package has been optimized with ASM with for x64 systems, hence
// Blake2 is faster than SHA.
func WithTransparent256(scp scope.Scope, id int64, key []byte) signed.Option {
	return func(s *signed.Service) error {
		// incorrect code
		if err := signed.WithHash(scp, id, blake2b.New256, key)(s); err != nil {
			return errors.Wrap(err, "[signedblake2] WithContentHMAC_Blake2b256.WithHash")
		}
		sig := signed.NewHMAC("blk2b256")
		return signed.WithHeaderHandler(scp, id, sig)(s)
	}
}

// NewOptionFactory creates a new option factory function for the signedblake2 in the
// backend package to be used for automatic scope based configuration
// initialization. Configuration values are read from argument `be`.
func NewOptionFactory(key cfgmodel.Obscure, alg cfgmodel.Str, header cfgmodel.Str) (string, signed.OptionFactoryFunc) {
	return OptionName, func(sg config.Scoped) []signed.Option {

		passwd, _, err := key.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[signedblake2] Key.Obscure.Get"))
		}
		httpHeader, _, err := header.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[signedblake2] Header.Str.Get"))
		}

		algVal, scpHash, err := alg.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[signedblake2] Alg.Str.Get"))
		}

		_ = passwd
		_ = httpHeader
		_ = algVal
		_ = scpHash

		return signed.OptionsError(errors.NewEmptyf("[memstore] TODO"))
	}
}
