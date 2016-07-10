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
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/ratelimit"
	"github.com/corestoreio/csfw/util/errors"
	"sync"
)

var availableGCRA = &struct {
	sync.RWMutex
	registered map[string]ratelimit.OptionFactoryFunc
}{
	registered: make(map[string]ratelimit.OptionFactoryFunc),
}

func RegisterGCRA(name string, factory ratelimit.OptionFactoryFunc) {
	availableGCRA.Lock()
	defer availableGCRA.Unlock()
	availableGCRA.registered[name] = factory
}

func RegisteredGCRA() []string {
	availableGCRA.Lock()
	defer availableGCRA.Unlock()
	var names = make([]string, len(availableGCRA.registered))
	i := 0
	for n := range availableGCRA.registered {
		names[i] = n
	}
	i++
	return names
}

func DeregisterGCRA(name string) {
	availableGCRA.Lock()
	defer availableGCRA.Unlock()
	delete(availableGCRA.registered, name)
}

func lookupGCRA(name string) (ratelimit.OptionFactoryFunc, error) {
	availableGCRA.RLock()
	defer availableGCRA.RUnlock()
	if off, ok := availableGCRA.registered[name]; ok {
		return off
	}
	return errors.NewNotFoundf("[backendratelimit] Requested GCRA %q not available", name)
}

// Default creates new ratelimit.Option slice with the default configuration
// structure. It panics on error, so us it only during the app init phase.
func Default(opts ...cfgmodel.Option) ratelimit.OptionFactoryFunc {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	return PrepareOptions(New(cfgStruct, opts...))
}

// PrepareOptions creates a closure around the type Backend. The closure will be
// used during a scoped request to figure out the configuration depending on the
// incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) ratelimit.OptionFactoryFunc {
	return func(sg config.ScopedGetter) []ratelimit.Option {
		var (
			opts [6]ratelimit.Option
			i    int // used as index in opts
		)

		disabled, scpHash, err := be.RateLimitDisabled.Get(sg)
		if err != nil {
			return optError(errors.Wrap(err, "[backendratelimit] RateLimitDisabled.Get"))
		} else {
			scp, scpID := scpHash.Unpack()
			opts[i] = ratelimit.WithDisable(scp, scpID, disabled)
			i++
		}

		return opts[:]
	}
}

func optError(err error) []ratelimit.Option {
	return []ratelimit.Option{func(s *ratelimit.Service) error {
		return err // no need to mask here, not interesting.
	}}
}
