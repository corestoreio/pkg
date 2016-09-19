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
	"github.com/corestoreio/csfw/net/jwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

// PrepareOptions creates a closure around the type Backend. The closure will be
// used during a scoped request to figure out the configuration depending on the
// incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Configuration) jwt.OptionFactoryFunc {
	return func(sg config.Scoped) []jwt.Option {
		var (
			opts [6]jwt.Option
			i    int // used as index in opts
		)

		off, h, err := be.Disabled.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtDisabled.Get"))
		}
		opts[i] = jwt.WithDisable(h, off)
		i++

		exp, h, err := be.Expiration.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtExpiration.Get"))
		}
		opts[i] = jwt.WithExpiration(h, exp)
		i++

		skew, h, err := be.Skew.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtSkew.Get"))
		}
		opts[i] = jwt.WithSkew(h, skew)
		i++

		isSU, h, err := be.SingleTokenUsage.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtSingleUsage.Get"))
		}
		opts[i] = jwt.WithSingleTokenUsage(h, isSU)
		i++

		// todo: avoid the next code and use OptionFactories to apply a signing method. Example in ratelimit package.

		signingMethod, h, err := be.SigningMethod.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtSigningMethod.Get"))
		}

		// in case we later support store scope the hashes variable protects use
		// from applying the incorrect scope to a functional option.
		var hashes = make(scope.Hashes, 0, 5)
		var key csjwt.Key
		hashes = append(hashes, h)

		switch signingMethod.Alg() {
		case csjwt.RS256, csjwt.RS384, csjwt.RS512:
			rsaKey, h1, err := be.RSAKey.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtRSAKey.Get"))
			}
			rsaPW, h2, err := be.RSAKeyPassword.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtRSAKeyPassword.Get"))
			}
			key = csjwt.WithRSAPrivateKeyFromPEM(rsaKey, rsaPW)
			hashes = append(hashes, h1, h2)
		case csjwt.ES256, csjwt.ES384, csjwt.ES512:

			ecdsaKey, h1, err := be.ECDSAKey.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtECDSAKey.Get"))
			}
			ecdsaPW, h2, err := be.ECDSAKeyPassword.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtECDSAKeyPassword.Get"))
			}
			key = csjwt.WithECPrivateKeyFromPEM(ecdsaKey, ecdsaPW)
			hashes = append(hashes, h1, h2)
		case csjwt.HS256, csjwt.HS384, csjwt.HS512:

			password, h1, err := be.HmacPassword.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtHmacPassword.Get"))
			}
			key = csjwt.WithPassword(password)
			hashes = append(hashes, h1)
		default:
			return jwt.OptionsError(errors.Errorf("[jwt] Unknown signing method: %q", signingMethod.Alg()))
		}

		// figure out the common lowest scope to apply.
		newH, err := hashes.Lowest()
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] Hashes.Lowest"))
		}

		// WithSigningMethod must be added at the end of the slice to overwrite
		// default signing methods
		opts[i] = jwt.WithKey(newH, key)
		i++
		opts[i] = jwt.WithSigningMethod(newH, signingMethod)
		return opts[:]
	}
}
