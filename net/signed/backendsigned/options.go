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

		disabled, err := be.Disabled.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] Disabled.Get"))
		}
		opts[i] = signed.WithDisable(disabled, sg.ScopeIDs()...)
		i++
		if disabled {
			return opts[:i]
		}

		inTrailer, err := be.InTrailer.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] InTrailer.Get"))
		}
		opts[i] = signed.WithTrailer(inTrailer, sg.ScopeIDs()...)
		i++

		methods, err := be.AllowedMethods.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] AllowedMethods.Get"))
		}
		opts[i] = signed.WithAllowedMethods(methods, sg.ScopeIDs()...)
		i++

		key, err := be.Key.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] Key.Obscure.Get"))
		}
		alg, err := be.Algorithm.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] Algorithm.Str.Get"))
		}
		opts[i] = signed.WithHash(alg, key, sg.ScopeIDs()...)
		i++

		keyID, err := be.KeyID.Get(sg)
		if err != nil {
			return signed.OptionsError(errors.Wrap(err, "[backendsigned] KeyID.Str.Get"))
		}
		header, err := be.HTTPHeaderType.Get(sg)
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
		opts[i] = signed.WithHeaderHandler(hpw, sg.ScopeIDs()...)
		i++

		return opts[:]
	}
}
