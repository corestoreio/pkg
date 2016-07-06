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

package cors

import (
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
)

// auto generated: do not edit. See net/gen eric package

// scopedConfigGeneric private internal scoped based configuration used for
// embedding into scopedConfig type.
type scopedConfigGeneric struct {
	// lastErr used during selecting the config from the scopeCache map and gets
	// filled if an entry cannot be found.
	lastErr error
	// scopeHash defines the scope to which this configuration is bound to.
	scopeHash scope.Hash
}

func newScopedConfigError(err error) scopedConfig {
	return scopedConfig{
		scopedConfigGeneric: scopedConfigGeneric{
			lastErr: err,
		},
	}
}

// optionInheritDefault looks up if the default configuration exists and if not
// creates a newScopedConfig(). This function can only be used within a
// functional option because it expects that it runs within a acquired lock.
func optionInheritDefault(s *Service) *scopedConfig {
	if sc, ok := s.scopeCache[scope.DefaultHash]; ok && sc != nil {
		shallowCopy := new(scopedConfig)
		*shallowCopy = *sc
		return shallowCopy
	}
	return newScopedConfig()
}

// WithOptionFactory applies a function which lazily loads the option depending
// on the incoming scope within a request. For example applies the backend
// configuration to the service.
//
// Once this option function has been set all other manually set option functions,
// which accept a scope and a scope ID as an argument, will be overwritten by the
// new values retrieved from the configuration service.
//
//	cfgStruct, err := backendpkg.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendpkg.New(cfgStruct)
//
//	srv := pkg.MustNewService(
//		pkg.WithOptionFactory(backendpkg.PrepareOptions(pb)),
//	)
func WithOptionFactory(f OptionFactoryFunc) Option {
	return func(s *Service) error {
		s.oFactory.Group = new(singleflight.Group)
		s.oFactory.Group.DisableForget = true
		s.oFactory.OptionFactoryFunc = f
		return nil
	}
}
