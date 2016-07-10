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

package ratelimit

import (
	"fmt"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
)

// Auto generated: Do not edit. See net/internal/scopedService package for more details.

type preventEmbedder interface {
	prevent()
}

// scopedConfigGeneric private internal scoped based configuration used for
// embedding into scopedConfig type.
type scopedConfigGeneric struct {
	preventEmbedder
	// lastErr used during selecting the config from the scopeCache map and infligh
	// package.
	lastErr error
	// ScopeHash defines the scope to which this configuration is bound to.
	ScopeHash scope.Hash
}

func (scg scopedConfigGeneric) GoString() string {
	return fmt.Sprintf("scopedConfigGeneric{lastErr: %q, ScopeHash: %s}", scg.lastErr, scg.ScopeHash.GoString())
}

// newScopedConfigError easy helper to create an error
func newScopedConfigError(err error) ScopedConfig {
	return ScopedConfig{
		scopedConfigGeneric: scopedConfigGeneric{
			lastErr: err,
		},
	}
}

// optionInheritDefault looks up if the default configuration exists and if not
// creates a newScopedConfig(). This function can only be used within a
// functional option because it expects that it runs within an acquired lock
// because of the map.
func optionInheritDefault(s *Service) *ScopedConfig {
	if sc, ok := s.scopeCache[scope.DefaultHash]; ok && sc != nil {
		shallowCopy := new(ScopedConfig)
		*shallowCopy = *sc
		return shallowCopy
	}
	return newScopedConfig()
}

// withDefaultConfig triggers the default settings
func withDefaultConfig(scp scope.Scope, id int64) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		sc := optionInheritDefault(s)
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithOptionFactory applies a function which lazily loads the options from a
// slow backend depending on the incoming scope within a request. For example
// applies the backend configuration to the service.
//
// Once this option function has been set all other manually set option
// functions, which accept a scope and a scope ID as an argument, will NOT be
// overwritten by the new values retrieved from the configuration service.
//
//	cfgStruct, err := backendratelimit.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendratelimit.New(cfgStruct)
//
//	srv := ratelimit.MustNewService(
//		ratelimit.WithOptionFactory(backendratelimit.PrepareOptions(pb)),
//	)
func WithOptionFactory(f OptionFactoryFunc) Option {
	return func(s *Service) error {
		s.optionInflight = new(singleflight.Group)
		s.optionFactory = f
		return nil
	}
}
