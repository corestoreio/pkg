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

package scopedservice

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/util/errors"
)

// Auto generated: Do not edit. See net/internal/scopedService package for more details.

// Option can be used as an argument in NewService to configure it with
// different settings.
type Option func(*Service) error

// OptionFactoryFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during a
// request.
type OptionFactoryFunc func(config.ScopedGetter) []Option

// OptionsError helper function to be used within the backend package or other
// sub-packages whose functions may return an OptionFactoryFunc.
func OptionsError(err error) []Option {
	return []Option{func(s *Service) error {
		return err // no need to mask here, not interesting.
	}}
}

// NewOptionFactories creates a new struct and inits the internal map.
func NewOptionFactories() *OptionFactories {
	return &OptionFactories{
		register: make(map[string]OptionFactoryFunc),
	}
}

// OptionFactories allows to register multiple OptionFactoryFunc identified by
// their names. Those OptionFactoryFuncs will be loaded in the backend package
// depending on the configured name under a certain path. This type is embedded
// in the backendscopedservice.Backend package.
type OptionFactories struct {
	rwmu sync.RWMutex
	// register where the key defines the name as specified in the
	// configuration path net/ratelimit_storage/gcra_name. The key equals the
	// 3rd party package name.
	register map[string]OptionFactoryFunc
}

// Register adds another functional option factory to the internal register.
// Overwrites existing entries.
func (be *OptionFactories) Register(name string, factory OptionFactoryFunc) {
	be.rwmu.Lock()
	defer be.rwmu.Unlock()
	be.register[name] = factory
}

// Names returns an unordered list of names of all registered functional option
// factories.
func (be *OptionFactories) Names() []string {
	be.rwmu.RLock()
	defer be.rwmu.RUnlock()
	var names = make([]string, len(be.register))
	i := 0
	for n := range be.register {
		names[i] = n
	}
	i++
	return names
}

// Deregister removes a functional option factory from the internal register.
func (be *OptionFactories) Deregister(name string) {
	be.rwmu.Lock()
	defer be.rwmu.Unlock()
	delete(be.register, name)
}

// Lookup returns a functional option factory identified by name or an error if
// the entry doesn't exists. May return a NotFound error behaviour.
func (be *OptionFactories) Lookup(name string) (OptionFactoryFunc, error) {
	be.rwmu.RLock()
	defer be.rwmu.RUnlock()
	if off, ok := be.register[name]; ok { // off = OptionFactoryFunc ;-)
		return off, nil
	}
	return nil, errors.NewNotFoundf("[backendratelimit] Requested OptionFactoryFunc %q not registered.", name)
}
