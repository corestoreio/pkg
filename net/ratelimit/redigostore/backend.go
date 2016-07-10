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

package redigostore

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/net/ratelimit/backendratelimit"
	"github.com/corestoreio/csfw/util/errors"
)

func RegisterOptionFacory(be *backendratelimit.Backend) (string, ratelimit.OptionFactoryFunc) {
	return "redigostore", func(sg config.ScopedGetter) []ratelimit.Option {

		burst, _, err := be.RateLimitBurst.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[redigostore] RateLimitBurst.Get"))
		}
		req, _, err := be.RateLimitRequests.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[redigostore] RateLimitRequests.Get"))
		}
		durRaw, _, err := be.RateLimitDuration.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[redigostore] RateLimitDuration.Get"))
		}

		if len(durRaw) != 1 {
			return optError(errors.NewFatalf("[redigostore] RateLimitDuration invalid character count: %q. Should be one character long.", durRaw))
		}

		dur := rune(durRaw[0])

		redisURL, scpHash, err := be.RateLimitStorageGCRARedis.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[redigostore] RateLimitStorageGcraRedis.Get"))
		}
		if redisURL != "" {
			scp, scpID := scpHash.Unpack()
			return WithGCRA(scp, scpID, redisURL, dur, req, burst)
		}
		return errors.NewEmptyf("[redigostore] Redis not active because RateLimitStorageGCRARedis is not set.")
	}
}

func optError(err error) []ratelimit.Option {
	return []ratelimit.Option{func(s *ratelimit.Service) error {
		return err // no need to mask here, not interesting.
	}}
}
