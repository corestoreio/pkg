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

package signedsha

import (
	"crypto/sha256"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// OptionName identifies this package within the register of the
// backendsigned.Backend type.
const OptionName = `sha`

// WithContentHMACSHA256 applies the SHA256 hash with your symmetric key.
func WithContentHMAC256(scp scope.Scope, id int64, key []byte) signed.Option {
	return func(s *signed.Service) error {
		if err := signed.WithHash(scp, id, sha256.New, key)(s); err != nil {
			return errors.Wrap(err, "[signed] WithContentHMAC_SHA256.WithHash")
		}
		sig := signed.NewHMAC("sha256")
		return signed.WithHeaderHandler(scp, id, sig)(s)
	}
}

func WithContentSignature256(scp scope.Scope, id int64, keyID string, key []byte) signed.Option {
	return func(s *signed.Service) error {
		if err := signed.WithHash(scp, id, sha256.New, key)(s); err != nil {
			return errors.Wrap(err, "[signed] WithContentHMAC_SHA256.WithHash")
		}
		sig := signed.NewSignature(keyID, "sha256")
		return signed.WithHeaderHandler(scp, id, sig)(s)
	}
}

// WithContentHMACSHA256 applies the SHA256 hash with your symmetric key.
func WithTransparent256(scp scope.Scope, id int64, key []byte) signed.Option {
	return func(s *signed.Service) error {
		// incorrect code
		if err := signed.WithHash(scp, id, sha256.New, key)(s); err != nil {
			return errors.Wrap(err, "[signed] WithContentHMAC_SHA256.WithHash")
		}
		sig := signed.NewHMAC("sha256")
		return signed.WithHeaderHandler(scp, id, sig)(s)
	}
}

// NewOptionFactory creates a new option factory function for the signedsha in the
// backend package to be used for automatic scope based configuration
// initialization. Configuration values are read from argument `be`.
func NewOptionFactory(key cfgmodel.Obscure, alg cfgmodel.Str, header cfgmodel.Str) (string, signed.OptionFactoryFunc) {
	return OptionName, func(sg config.Scoped) []signed.Option {

		passwd, _, err := key.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[signedsha] Key.Obscure.Get"))
		}
		httpHeader, _, err := header.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[signedsha] Header.Str.Get"))
		}

		_ = httpHeader

		algVal, scpHash, err := alg.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[signedsha] Alg.Str.Get"))
		}

		scp, scpID := scpHash.Unpack()
		switch algVal {
		case "sha256":
			return []signed.Option{
				WithContentHMAC256(scp, scpID, passwd),
			}
		case "sha512":
		}
		return signed.OptionsError(errors.NewNotSupportedf("[signedsha] Algorithm %q not suppored", algVal))
	}
}
