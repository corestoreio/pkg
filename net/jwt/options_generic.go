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

package jwt

import (
	"io"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/csfw/util/errors"
)

// Auto generated: Do not edit. See net/internal/scopedService package for more details.

// Option can be used as an argument in NewService to configure it with
// different settings.
type Option func(*Service) error

// OptionFactoryFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during a
// request.
type OptionFactoryFunc func(config.Scoped) []Option

// OptionsError helper function to be used within the backend package or other
// sub-packages whose functions may return an OptionFactoryFunc.
func OptionsError(err error) []Option {
	return []Option{func(s *Service) error {
		return err // no need to mask here, not interesting.
	}}
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

// WithErrorHandler adds a custom error handler. Gets called after the scope can
// be extracted from the context.Context and the configuration has been found
// and is valid. The default error handler prints the error to the user and
// returns a http.StatusServiceUnavailable.
func WithErrorHandler(scp scope.Scope, id int64, eh mw.ErrorHandler) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.ErrorHandler = eh
		sc.ScopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithServiceErrorHandler sets the error handler on the Service object.
// Convenient helper function.
func WithServiceErrorHandler(eh mw.ErrorHandler) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.ErrorHandler = eh
		return nil
	}
}

// WithRootConfig sets the root configuration service. While using any HTTP
// related functions or middlewares you must set the config.Getter.
func WithRootConfig(cg config.Getter) Option {
	_ = cg.NewScoped(0, 0) // let it panic as early as possible if cg is nil
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.rootConfig = cg
		return nil
	}
}

// WithDebugLog creates a new standard library based logger with debug mode
// enabled. The passed writer must be thread safe.
func WithDebugLog(w io.Writer) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.Log = logw.NewLog(logw.WithWriter(w), logw.WithLevel(logw.LevelDebug))
		return nil
	}
}

// WithOptionFactory applies a function which lazily loads the options from a
// slow backend (config.Getter) depending on the incoming scope within a
// request. For example applies the backend configuration to the service.
//
// Once this option function has been set all other manually set option
// functions, which accept a scope and a scope ID as an argument, will NOT be
// overwritten by the new values retrieved from the configuration service.
//
//	cfgStruct, err := backendjwt.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendjwt.New(cfgStruct)
//
//	srv := jwt.MustNewService(
//		jwt.WithOptionFactory(backendjwt.PrepareOptions(pb)),
//	)
func WithOptionFactory(f OptionFactoryFunc) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.optionInflight = new(singleflight.Group)
		s.optionFactory = f
		return nil
	}
}

// NewOptionFactories creates a new struct and initializes the internal map for
// the registration of different option factories.
func NewOptionFactories() *OptionFactories {
	return &OptionFactories{
		register: make(map[string]OptionFactoryFunc),
	}
}

// OptionFactories allows to register multiple OptionFactoryFunc identified by
// their names. Those OptionFactoryFuncs will be loaded in the backend package
// depending on the configured name under a certain path. This type is embedded
// in the backendjwt.Backend package.
type OptionFactories struct {
	rwmu sync.RWMutex
	// register where the key defines the name as specified in the
	// configuration path what/ever/path. The key equals the
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
	return nil, errors.NewNotFoundf("[jwt] Requested OptionFactoryFunc %q not registered.", name)
}
