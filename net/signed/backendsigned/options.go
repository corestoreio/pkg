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

package backendsigned

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// PrepareOptionFactory creates a closure around the type Backend. The closure
// will be used during a scoped request to figure out the configuration
// depending on the incoming scope. An option array will be returned by the
// closure.
func (be *Configuration) PrepareOptionFactory() signed.OptionFactoryFunc {
	return func(sg config.Scoped) []signed.Option {
		var (
			opts [5]signed.Option
			i    int // used as index in opts
		)
		h := scope.MakeTypeID(sg.ScopeID())

		// i think there is a bug with the scpHash because scpHash returns that
		// hash in which scope the value has been found. for example we have 2
		// websites with each 2 stores. those 4 stores are configured for
		// signing. so we have for each store scope in the scopeCache map an
		// entry. if we now disable the signing for a website then it won't be
		// disabled for the store scopes because it won't fall back to website
		// as we have found an entry in the store scope in the ScopeCache map.
		// hence the solution would be to use the hash value from sg.Scope()

		disabled, _, err := be.Disabled.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] Disabled.Get"))
		}
		opts[i] = signed.WithDisable(h, disabled)
		i++
		if disabled {
			return opts[:i]
		}

		inTrailer, _, err := be.InTrailer.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] InTrailer.Get"))
		}
		opts[i] = signed.WithTrailer(h, inTrailer)
		i++

		methods, _, err := be.AllowedMethods.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] AllowedMethods.Get"))
		}
		opts[i] = signed.WithAllowedMethods(h, methods...)
		i++

		key, _, err := be.Key.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] Key.Obscure.Get"))
		}
		alg, _, err := be.Algorithm.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] Algorithm.Str.Get"))
		}
		opts[i] = signed.WithHash(h, alg, key)
		i++

		keyID, _, err := be.KeyID.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] KeyID.Str.Get"))
		}
		header, _, err := be.HTTPHeaderType.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] HTTPHeaderType.Str.Get"))
		}

		var hpw signed.HeaderParseWriter
		switch header {
		// case "transparent":
		// todo: transparent must be implemented via a new package using also the signed.OptionFactoryFunc; same like ratelimit
		case "hmac":
			hpw = signed.NewContentHMAC(alg)
		case "signature":
			hpw = signed.NewContentSignature(keyID, alg)
		default:
			return signed.OptionsError(errors.NewNotImplementedf("[backendsigned] HTTPHeaderType %q not implemented", header))
		}
		opts[i] = signed.WithHeaderHandler(h, hpw)
		i++

		return opts[:]
	}
}
