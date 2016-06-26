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
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type Backend struct {
	cfgmodel.PkgBackend

	// RateLimitDisabled set to true to disable the rate limiting.
	//
	// Path: net/ratelimit/disabled
	RateLimitDisabled cfgmodel.Bool

	// RateLimitBurst defines the number of requests that
	// will be allowed to exceed the rate in a single burst and must be
	// greater than or equal to zero.
	//
	// Path: net/ratelimit/burst
	RateLimitBurst cfgmodel.Int

	// RateLimitRequests number of requests allowed per time period
	//
	// Path: net/ratelimit/requests
	RateLimitRequests cfgmodel.Int

	// RateLimitDuration per second (s), minute (i), hour (h), day (d)
	//
	// Path: net/ratelimit/duration
	RateLimitDuration ConfigRate

	// RateLimitStorageEnableGcraMemory set to true to enable in-memory rate
	// limit key storage.
	//
	// Path: net/ratelimit_storage/gcra_memory_enable
	RateLimitStorageGcraMemoryEnable cfgmodel.Bool

	// RateLimitStorageGcraMaxMemoryKeys If maxKeys > 0, the number of different
	// keys is restricted to the specified amount. In this case, it uses an LRU
	// algorithm to evict older keys to make room for newer ones. If maxKeys <=
	// 0, there is no limit on the number of keys, which may use an unbounded
	// amount of memory.
	//
	// Path: net/ratelimit_storage/enable_gcra_memory
	RateLimitStorageGcraMaxMemoryKeys cfgmodel.Int
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries. The function Load() will be executed to
// apply the SectionSlice to all models. See Load() for more details.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	return (&Backend{}).Load(cfgStruct, opts...)
}

// Load creates the configuration models for each PkgBackend field. Internal
// mutex will protect the fields during loading. The argument SectionSlice will
// be applied to all models.
func (pp *Backend) Load(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	pp.Lock()
	defer pp.Unlock()

	opt := cfgmodel.WithFieldFromSectionSlice(cfgStruct)

	pp.RateLimitBurst = cfgmodel.NewInt(`net/ratelimit/burst`, opt)
	pp.RateLimitRequests = cfgmodel.NewInt(`net/ratelimit/requests`, opt)
	pp.RateLimitDuration = NewConfigDuration(`net/ratelimit/duration`, opt)

	return pp
}
