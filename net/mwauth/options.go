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
type Option func(*Service) error

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
	return func(s *Service) error {
		var err error
		if h == scope.DefaultHash {
			s.defaultScopeCache, err = defaultScopedConfig()
			return errors.Wrap(err, "[mwauth] Default Scope with Default Config")
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		s.scopeCache[h], err = defaultScopedConfig()
		return errors.Wrapf(err, "[mwauth] Scope %s with Default Config", h)
	}
}

func WithIsActive(scp scope.Scope, id int64, active bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.enable = active
			return nil
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
		return nil
	}
}

// WithLogger applies a logger to the default scope which gets inherited to
// subsequent scopes.
// Mainly used for debugging.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.defaultScopeCache.log = l
		return nil
	}
}

// WithOptionFactory applies a function which lazily loads the option depending
// on the incoming scope within a request. For example applies the backend
// configuration to the service.
//
// Once this option function has been set all other manually set option functions,
// which accept a scope and a scope ID as an argument, will be overwritten by the
// new values retrieved from the configuration service.
//
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
	return func(s *Service) error {
		s.scpOptionFnc = f
		return nil
	}
}
