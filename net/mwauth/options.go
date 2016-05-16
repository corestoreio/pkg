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

package mwauth

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
)

// Option defines a function argument for the Cors type to apply options.
type Option func(*Service)

// ScopedOptionFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during
// a request.
type ScopedOptionFunc func(config.ScopedGetter) []Option

// WithDefaultConfig applies the default configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- ?
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {
		if s.optionError != nil {
			return
		}

		if h == scope.DefaultHash {
			s.defaultScopeCache, s.optionError = defaultScopedConfig()
			s.optionError = errors.Wrap(s.optionError, "[mwauth] Default Scope with Default Config")
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		s.scopeCache[h], s.optionError = defaultScopedConfig()
		s.optionError = errors.Wrapf(s.optionError, "[mwauth] Scope %s with Default Config", h)
	}
}

func WithIsActive(scp scope.Scope, id int64, active bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) {
		if h == scope.DefaultHash {
			s.defaultScopeCache.enable = active
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.enable = active

		if sc, ok := s.scopeCache[h]; ok {
			sc.enable = scNew.enable
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
	}
}

// WithLogger applies a logger to the default scope which gets inherited to
// subsequent scopes.
// Mainly used for debugging.
func WithLogger(l log.Logger) Option {
	return func(s *Service) {
		s.defaultScopeCache.log = l
	}
}

// WithOptionFactory applies a function which lazily loads the option depending
// on the incoming scope within a request. For example applies the backend
// configuration to the service.
// Once this option function has been set all other option functions are not really
// needed.
//	cfgStruct, err := backendauth.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendauth.New(cfgStruct)
//
//	cors := mwauth.MustNewService(
//		mwauth.WithOptionFactory(backendauth.PrepareOptions(pb)),
//	)
func WithOptionFactory(f ScopedOptionFunc) Option {
	return func(s *Service) {
		s.scpOptionFnc = f
	}
}
