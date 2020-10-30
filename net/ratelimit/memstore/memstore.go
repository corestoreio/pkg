// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package memstore

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/net/ratelimit"
	"github.com/corestoreio/pkg/store/scope"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

// OptionName identifies this package within the register of the
// backendratelimit.Backend type.
const OptionName = `memstore`

// NewOptionFactory creates a new option factory function for the memstore in the
// backend package to be used for automatic scope based configuration
// initialization. Configuration values are read from package backendratelimit.Configuration.
func NewOptionFactory(burst, requests cfgmodel.Int, duration cfgmodel.Str, gcraMaxMemoryKeys cfgmodel.Int) (optionName string, _ ratelimit.OptionFactoryFunc) {
	return OptionName, func(sg config.Scoped) []ratelimit.Option {
		burst, err := burst.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[memstore] RateLimitBurst.Get"))
		}
		req, err := requests.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[memstore] RateLimitRequests.Get"))
		}
		durRaw, err := duration.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[memstore] RateLimitDuration.Get"))
		}

		if len(durRaw) != 1 {
			return ratelimit.OptionsError(errors.NewFatalf("[memstore] RateLimitDuration invalid character count: %q. Should be one character long.", durRaw))
		}

		dur := rune(durRaw[0])

		useInMemMaxKeys, err := gcraMaxMemoryKeys.Get(sg)
		if err != nil {
			return ratelimit.OptionsError(errors.Wrap(err, "[memstore] RateLimitStorageGcraMaxMemoryKeys.Get"))
		} else if useInMemMaxKeys > 0 {
			return []ratelimit.Option{
				WithGCRA(useInMemMaxKeys, dur, req, burst, sg.ScopeIDs()...),
			}
		}
		return ratelimit.OptionsError(errors.NewEmptyf("[memstore] Memstore not active because RateLimitStorageGcraMaxMemoryKeys is %d.", useInMemMaxKeys))
	}
}

// WithGCRA creates a memory based GCRA rate limiter.
// Duration: (s second,i minute,h hour,d day).
// GCRA => https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm
// This function implements a debug log.
func WithGCRA(maxKeys int, duration rune, requests, burst int, scopeIDs ...scope.TypeID) ratelimit.Option {
	return func(s *ratelimit.Service) error {
		rlStore, err := memstore.New(maxKeys)
		if err != nil {
			return errors.NewFatalf("[memstore] memstore.New MaxKeys(%d): %s", maxKeys, err)
		}
		if s.Log.IsDebug() {
			s.Log.Debug("ratelimit.memstore.WithGCRA",
				log.Stringer("scope", scope.TypeIDs(scopeIDs)),
				log.Int("max_keys", maxKeys),
				log.String("duration", string(duration)),
				log.Int("requests", requests),
				log.Int("burst", burst),
			)
		}
		return ratelimit.WithGCRAStore(rlStore, duration, requests, burst, scopeIDs...)(s)
	}
}
