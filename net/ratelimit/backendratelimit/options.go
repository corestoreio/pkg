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

package backendratelimit

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/util/errors"
)

// PrepareOptionFactory creates a closure around the type Backend. The closure will be
// used during a scoped request to figure out the configuration depending on the
// incoming scope. An option array will be returned by the closure.
func (be *Configuration) PrepareOptionFactory() ratelimit.OptionFactoryFunc {
	return func(sg config.Scoped) []ratelimit.Option {
		var opts []ratelimit.Option

		disabled, err := be.Disabled.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[backendratelimit] RateLimitDisabled.Get"))
		}
		opts = append(opts, ratelimit.WithDisable(disabled, sg.ScopeIDs()...))
		if disabled {
			return opts
		}

		name, err := be.GCRAName.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[backendratelimit] RateLimitGCRAName.Get"))
		}

		// name contains the configured ratelimit calculation/storage engine. in this case either
		// memstore or redigostore. Of course you can plugin your own engine.
		off, err := be.Lookup(name) // off = OptionFactoryFunc
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[backendratelimit] Backend.Lookup"))
		}
		return append(opts, off(sg)...)
	}
}
