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
			opts  [6]jwt.Option
			i     int // used as index in opts
			scp   scope.Scope
			scpID int64
		)

		off, h, err := be.NetJwtDisabled.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtDisabled.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = jwt.WithDisable(scp, scpID, off)
		i++

		exp, h, err := be.NetJwtExpiration.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtExpiration.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = jwt.WithExpiration(scp, scpID, exp)
		i++

		skew, h, err := be.NetJwtSkew.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtSkew.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = jwt.WithSkew(scp, scpID, skew)
		i++

		isSU, h, err := be.NetJwtSingleTokenUsage.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtSingleUsage.Get"))
		}
		scp, scpID = h.Unpack()
		opts[i] = jwt.WithSingleTokenUsage(scp, scpID, isSU)
		i++

		// todo: avoid the next code and use OptionFactories to apply a signing method. Example in ratelimit package.

		signingMethod, h, err := be.NetJwtSigningMethod.Get(sg)
		if err != nil {
			return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtSigningMethod.Get"))
		}
		scp, scpID = h.Unpack()

		var key csjwt.Key

		switch signingMethod.Alg() {
		case csjwt.RS256, csjwt.RS384, csjwt.RS512:

			rsaKey, _, err := be.NetJwtRSAKey.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtRSAKey.Get"))
			}
			rsaPW, _, err := be.NetJwtRSAKeyPassword.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtRSAKeyPassword.Get"))
			}
			key = csjwt.WithRSAPrivateKeyFromPEM(rsaKey, rsaPW)

		case csjwt.ES256, csjwt.ES384, csjwt.ES512:

			ecdsaKey, _, err := be.NetJwtECDSAKey.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtECDSAKey.Get"))
			}
			ecdsaPW, _, err := be.NetJwtECDSAKeyPassword.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtECDSAKeyPassword.Get"))
			}
			key = csjwt.WithECPrivateKeyFromPEM(ecdsaKey, ecdsaPW)

		case csjwt.HS256, csjwt.HS384, csjwt.HS512:

			password, _, err := be.NetJwtHmacPassword.Get(sg)
			if err != nil {
				return jwt.OptionsError(errors.Wrap(err, "[backendjwt] NetJwtHmacPassword.Get"))
			}
			key = csjwt.WithPassword(password)

		default:
			return jwt.OptionsError(errors.Errorf("[jwt] Unknown signing method: %q", signingMethod.Alg()))
		}

		// WithSigningMethod must be added at the end of the slice to overwrite default signing methods
		opts[i] = jwt.WithKey(scp, scpID, key)
		i++
		opts[i] = jwt.WithSigningMethod(scp, scpID, signingMethod)
		return opts[:]
	}
}
