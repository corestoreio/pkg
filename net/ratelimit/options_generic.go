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
	"io"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/sync/singleflight"
	"github.com/corestoreio/errors"
)

// Auto generated: Do not edit. See net/internal/scopedService package for more details.

// Option can be used as an argument in NewService to configure it with
// different settings.
type Option func(*Service) error

// OptionsError helper function to be used within the backend package or other
// sub-packages whose functions may return an OptionFactoryFunc.
func OptionsError(err error) []Option {
	return []Option{func(s *Service) error {
		return err // no need to mask here, not interesting.
	}}
}

// withDefaultConfig triggers the default settings for a specific ScopeID.
func withDefaultConfig(scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		target, parents := scope.TypeIDs(scopeIDs).TargetAndParents()
		sc = newScopedConfig(target, parents[0])
		return s.updateScopedConfig(sc)
	}
}

// WithErrorHandler adds a custom error handler. Gets called in the http.Handler
// after the scope can be extracted from the context.Context and the
// configuration has been found and is valid. The default error handler prints
// the error to the user and returns a http.StatusServiceUnavailable.
//
// The variadic "scopeIDs" argument define to which scope the value gets applied
// and from which parent scope should be inherited. Setting no "scopeIDs" sets
// the value to the default scope. Setting one scope.TypeID defines the primary
// scope to which the value will be applied. Subsequent scope.TypeID are
// defining the fall back parent scopes to inherit the default or previously
// applied configuration from.
func WithErrorHandler(eh mw.ErrorHandler, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.ErrorHandler = eh
		return s.updateScopedConfig(sc)
	}
}

// WithDisable disables the current service and calls the next HTTP handler.
//
// The variadic "scopeIDs" argument define to which scope the value gets applied
// and from which parent scope should be inherited. Setting no "scopeIDs" sets
// the value to the default scope. Setting one scope.TypeID defines the primary
// scope to which the value will be applied. Subsequent scope.TypeID are
// defining the fall back parent scopes to inherit the default or previously
// applied configuration from.
func WithDisable(isDisabled bool, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.Disabled = isDisabled
		return s.updateScopedConfig(sc)
	}
}

// WithMarkPartiallyApplied if set to true marks a configuration for a scope
// as partially applied with functional options set via source code. The
// internal service knows that it must trigger additionally the
// OptionFactoryFunc to load configuration from a backend. Useful in the case
// where parts of the configurations are coming from backend storages and other
// parts like http handler have been set via code. This function should only be
// applied in case you work with WithOptionFactory().
//
// The variadic "scopeIDs" argument define to which scope the value gets applied
// and from which parent scope should be inherited. Setting no "scopeIDs" sets
// the value to the default scope. Setting one scope.TypeID defines the primary
// scope to which the value will be applied. Subsequent scope.TypeID are
// defining the fall back parent scopes to inherit the default or previously
// applied configuration from.
func WithMarkPartiallyApplied(partially bool, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.lastErr = nil
		if partially {
			sc.lastErr = errors.NewTemporaryf(errConfigMarkedAsPartiallyLoaded, sc.ScopeID)
		}
		return s.updateScopedConfig(sc)
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

// WithRootConfig sets the root configuration service to retrieve the scoped
// base configuration. If you set the option WithOptionFactory() then the option
// WithRootConfig() does not need to be set as it won't get used.
func WithRootConfig(cg config.Getter) Option {
	_ = cg.NewScoped(0, 0) // let it panic as early as possible if cg is nil
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.RootConfig = cg
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

// WithLogger convenient helper function to apply a logger to the Service type.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		s.Log = l
		return nil
	}
}

// OptionFactoryFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during a
// request.
type OptionFactoryFunc func(config.Scoped) []Option

// WithOptionFactory applies a function which lazily loads the options from a
// slow backend (config.Getter) depending on the incoming scope within a
// request. For example applies the backend configuration to the service.
//
// Once this option function has been set all other manually set option
// functions, which accept a scope and a scope ID as an argument, will NOT be
// overwritten by the new values retrieved from the configuration service.
//
//	cfgStruct, err := backendratelimit.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	be := backendratelimit.New(cfgStruct)
//
//	srv := ratelimit.MustNewService(
//		ratelimit.WithOptionFactory(be.PrepareOptions()),
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
// in the backendratelimit.Configuration type.
type OptionFactories struct {
	rwmu sync.RWMutex
	// register where the key defines the name as specified in the
	// configuration path what/ever/path. The key equals the
	// 3rd party package name.
	register map[string]OptionFactoryFunc
}

// Register adds another functional option factory to the internal register.
// Overwrites existing entries.
func (of *OptionFactories) Register(name string, factory OptionFactoryFunc) {
	of.rwmu.Lock()
	defer of.rwmu.Unlock()
	of.register[name] = factory
}

// Names returns an unordered list of names of all registered functional option
// factories.
func (of *OptionFactories) Names() []string {
	of.rwmu.RLock()
	defer of.rwmu.RUnlock()
	var names = make([]string, len(of.register))
	i := 0
	for n := range of.register {
		names[i] = n
		i++
	}
	return names
}

// Deregister removes a functional option factory from the internal register.
func (of *OptionFactories) Deregister(name string) {
	of.rwmu.Lock()
	defer of.rwmu.Unlock()
	delete(of.register, name)
}

// Lookup returns a functional option factory identified by name or an error if
// the entry doesn't exists. May return a NotFound error behaviour.
func (of *OptionFactories) Lookup(name string) (OptionFactoryFunc, error) {
	of.rwmu.RLock()
	defer of.rwmu.RUnlock()
	if off, ok := of.register[name]; ok { // off = OptionFactoryFunc ;-)
		return off, nil
	}
	return nil, errors.NewNotFoundf("[ratelimit] Requested OptionFactoryFunc %q not registered.", name)
}
