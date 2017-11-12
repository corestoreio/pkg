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
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/net/ratelimit"
)

// Configuration just exported for the sake of documentation. See fields for more
// information. Please call the New() function for creating a new Backend
// object. Only the New() function will set the paths to the fields.
type Configuration struct {
	*ratelimit.OptionFactories

	// Disabled set to true to disable the rate limiting.
	//
	// Path: net/ratelimit/disabled
	Disabled cfgmodel.Bool

	// Burst defines the number of requests that will be allowed to
	// exceed the rate in a single burst and must be greater than or equal to
	// zero.
	//
	// Path: net/ratelimit/burst
	Burst cfgmodel.Int

	// Requests number of requests allowed per time period
	//
	// Path: net/ratelimit/requests
	Requests cfgmodel.Int

	// Duration per second (s), minute (i), hour (h), day (d)
	//
	// Path: net/ratelimit/duration
	Duration cfgmodel.Str

	// GCRAName sets the name which GCRA can be used. The GCRA must be
	// registered prior to calling the middleware handler. The name is usually
	// the package name. For example net/ratelimit/memstore or
	// net/ratelimit/redigostore. Leaving this configuration value empty or
	// setting a not registered name causes the middleware handler to panic.
	//
	// Path: net/ratelimit_storage/gcra_name
	GCRAName cfgmodel.Str

	// StorageGCRAMaxMemoryKeys If maxKeys > 0 (enabled), the number of
	// different keys is restricted to the specified amount. In this case, it
	// uses an LRU algorithm to evict older keys to make room for newer ones.
	//
	// Path: net/ratelimit_storage/enable_gcra_memory
	StorageGCRAMaxMemoryKeys cfgmodel.Int

	// StorageGCRARedis a valid Redis URL enables Redis as GCRA key
	// storage. URLs should follow the draft IANA specification for the scheme
	// (https://www.iana.org/assignments/uri-schemes/prov/redis).
	//
	//
	// For example:
	// 		redis://localhost:6379/3
	// 		redis://:6380/0 => connects to localhost:6380
	// 		redis:// => connects to localhost:6379 with DB 0
	// 		redis://empty:myPassword@clusterName.xxxxxx.0001.usw2.cache.amazonaws.com:6379/0
	//
	// Path: net/ratelimit_storage/enable_gcra_redis
	StorageGCRARedis cfgmodel.Str
}

// New initializes the backend configuration models containing the cfgpath.Route
// variable to the appropriate entries in the storage. The argument SectionSlice
// and opts will be applied to all models.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Configuration {
	be := &Configuration{
		OptionFactories: ratelimit.NewOptionFactories(),
	}

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	be.Disabled = cfgmodel.NewBool(`net/ratelimit/disabled`, opts...)
	be.Burst = cfgmodel.NewInt(`net/ratelimit/burst`, opts...)
	be.Requests = cfgmodel.NewInt(`net/ratelimit/requests`, opts...)
	be.Duration = cfgmodel.NewStr(`net/ratelimit/duration`, append(opts, cfgmodel.WithSourceByString(
		"s", "Second",
		"i", "Minute",
		"h", "Hour",
		"d", "Day",
	))...)
	be.GCRAName = cfgmodel.NewStr(`net/ratelimit_storage/gcra_name`, opts...)
	be.StorageGCRAMaxMemoryKeys = cfgmodel.NewInt(`net/ratelimit_storage/enable_gcra_memory`, opts...)
	be.StorageGCRARedis = cfgmodel.NewStr(`net/ratelimit_storage/enable_gcra_redis`, opts...)

	return be
}
