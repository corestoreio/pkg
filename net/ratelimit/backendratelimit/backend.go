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
type PkgBackend struct {
	cfgmodel.PkgBackend

	// RateLimitBurst defines the number of requests that
	// will be allowed to exceed the rate in a single burst and must be
	// greater than or equal to zero.
	//
	// Path: net/shy/burst
	RateLimitBurst cfgmodel.Int

	// RateLimitRequests number of requests allowed per time period
	//
	// Path: net/shy/requests
	RateLimitRequests cfgmodel.Int

	// RateLimitDuration per second (s), minute (i), hour (h), day (d)
	//
	// Path: net/shy/duration
	RateLimitDuration ConfigDuration
}

// NewBackend initializes the global configuration models containing the
// cfgpath.Route variable to the appropriate entry.
// The function Load() will be executed to apply the SectionSlice
// to all models. See Load() for more details.
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).Load(cfgStruct)
}

// Load creates the configuration models for each PkgBackend field.
// Internal mutex will protect the fields during loading.
// The argument SectionSlice will be applied to all models.
func (pp *PkgBackend) Load(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()

	opt := cfgmodel.WithFieldFromSectionSlice(cfgStruct)

	pp.RateLimitBurst = cfgmodel.NewInt(`net/shy/burst`, opt)
	pp.RateLimitRequests = cfgmodel.NewInt(`net/shy/requests`, opt)
	pp.RateLimitDuration = NewConfigDuration(`net/shy/duration`, opt)

	return pp
}
