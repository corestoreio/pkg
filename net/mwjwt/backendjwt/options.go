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

package backendjwt

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/mwjwt"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

// Default creates new mwjwt.Option slice with the default configuration
// structure and a noop encryptor/decryptor IF no option arguments have been
// provided. It panics on error, so us it only during the app init phase.
func Default(opts ...cfgmodel.Option) mwjwt.ScopedOptionFunc {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	if len(opts) == 0 {
		opts = append(opts, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))
	}

	return PrepareOptions(New(cfgStruct, opts...))
}

// PrepareOptions creates a closure around the type Backend. The closure will
// be used during a scoped request to figure out the configuration depending on
// the incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) mwjwt.ScopedOptionFunc {

	return func(sg config.ScopedGetter) (opts []mwjwt.Option) {

		scp, id := sg.Scope()

		exp, err := be.NetJwtExpiration.Get(sg)
		if err != nil {
			return append(opts, func(s *mwjwt.Service) {
				s.AddError(errors.Wrap(err, "[backendjwt] NetJwtExpiration.Get"))
			})
		}
		opts = append(opts, mwjwt.WithExpiration(scp, id, exp))

		skew, err := be.NetJwtSkew.Get(sg)
		if err != nil {
			return append(opts, func(s *mwjwt.Service) {
				s.AddError(errors.Wrap(err, "[backendjwt] NetJwtSkew.Get"))
			})
		}
		opts = append(opts, mwjwt.WithSkew(scp, id, skew))

		isJTI, err := be.NetJwtEnableJTI.Get(sg)
		if err != nil {
			return append(opts, func(s *mwjwt.Service) {
				s.AddError(errors.Wrap(err, "[backendjwt] NetJwtEnableJTI.Get"))
			})
		}
		opts = append(opts, mwjwt.WithTokenID(scp, id, isJTI))

		signingMethod, err := be.NetJwtSigningMethod.Get(sg)
		if err != nil {
			return append(opts, func(s *mwjwt.Service) {
				s.AddError(errors.Wrap(err, "[backendjwt] NetJwtSigningMethod.Get"))
			})
		}

		var key csjwt.Key

		switch signingMethod.Alg() {
		case csjwt.RS256, csjwt.RS384, csjwt.RS512:

			rsaKey, err := be.NetJwtRSAKey.Get(sg)
			if err != nil {
				return append(opts, func(s *mwjwt.Service) {
					s.AddError(errors.Wrap(err, "[backendjwt] NetJwtRSAKey.Get"))
				})
			}
			rsaPW, err := be.NetJwtRSAKeyPassword.Get(sg)
			if err != nil {
				return append(opts, func(s *mwjwt.Service) {
					s.AddError(errors.Wrap(err, "[backendjwt] NetJwtRSAKeyPassword.Get"))
				})
			}
			key = csjwt.WithRSAPrivateKeyFromPEM(rsaKey, rsaPW)

		case csjwt.ES256, csjwt.ES384, csjwt.ES512:

			ecdsaKey, err := be.NetJwtECDSAKey.Get(sg)
			if err != nil {
				return append(opts, func(s *mwjwt.Service) {
					s.AddError(errors.Wrap(err, "[backendjwt] NetJwtECDSAKey.Get"))
				})
			}
			ecdsaPW, err := be.NetJwtECDSAKeyPassword.Get(sg)
			if err != nil {
				return append(opts, func(s *mwjwt.Service) {
					s.AddError(errors.Wrap(err, "[backendjwt] NetJwtECDSAKeyPassword.Get"))
				})
			}
			key = csjwt.WithECPrivateKeyFromPEM(ecdsaKey, ecdsaPW)

		case csjwt.HS256, csjwt.HS384, csjwt.HS512:

			password, err := be.NetJwtHmacPassword.Get(sg)
			if err != nil {
				return append(opts, func(s *mwjwt.Service) {
					s.AddError(errors.Wrap(err, "[backendjwt] NetJwtHmacPassword.Get"))
				})
			}
			key = csjwt.WithPassword(password)

		default:
			opts = append(opts, func(s *mwjwt.Service) {
				s.AddError(errors.Errorf("[mwjwt] Unknown signing method: %q", signingMethod.Alg()))
			})
		}

		// WithSigningMethod must be added at the end of the slice to overwrite default signing methods
		return append(opts, mwjwt.WithKey(scp, id, key), mwjwt.WithSigningMethod(scp, id, signingMethod))
	}
}
