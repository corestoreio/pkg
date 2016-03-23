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

package ctxthrottled

import (
	"github.com/corestoreio/csfw/config"
	"github.com/juju/errors"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

// DefaultRequests number of requests allowed per time period.
// Used when *PkgBackend has not been provided.
var DefaultRequests = 100

// DefaultBurst defines the number of requests that
// will be allowed to exceed the rate in a single burst and must be
// greater than or equal to zero.
// Used when *PkgBackend has not been provided.
var DefaultBurst = 20

// DefaultDuration per second (s), minute (i), hour (h), day (d)
// Used when *PkgBackend has not been provided.
var DefaultDuration = "h"

const MemStoreMaxKeys = 65536

// NewGCRAMemStore creates the default memory based GCRA rate limiter.
// It uses the PkgBackend models to create a ratelimiter for each scope.
func NewGCRAMemStore(maxKeys int) RateLimiterFactory {
	return func(be *PkgBackend, sg config.ScopedGetter) (throttled.RateLimiter, error) {

		rlStore, err := memstore.New(maxKeys)
		if err != nil {
			return nil, err
		}

		rq, err := rateQuota(be, sg)
		if err != nil {
			return nil, err
		}

		rl, err := throttled.NewGCRARateLimiter(rlStore, rq)
		if err != nil {
			return nil, err
		}

		return rl, nil
	}
}

// rateQuota creates a new quota for the GCRARateLimiter
func rateQuota(be *PkgBackend, sg config.ScopedGetter) (rq throttled.RateQuota, err error) {

	if be == nil {
		return throttled.RateQuota{
			MaxRate:  calculateRate(DefaultDuration, DefaultRequests),
			MaxBurst: DefaultBurst,
		}, nil
	}

	burst, err := be.RateLimitBurst.Get(sg)
	if err != nil {
		err = errors.Mask(err)
		return
	}
	request, err := be.RateLimitRequests.Get(sg)
	if err != nil {
		err = errors.Mask(err)
		return
	}
	if request == 0 {
		request = DefaultRequests
	}

	rate, err := be.RateLimitDuration.Get(sg, request)
	err = errors.Mask(err)

	rq.MaxRate = rate
	rq.MaxBurst = burst
	return
}
