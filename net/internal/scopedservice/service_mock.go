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

import "github.com/corestoreio/csfw/store/scope"

// Service DO NOT USE
type Service struct {
	service
}

// New DO NOT USE
func New(opts ...Option) (*Service, error) {
	return newService(opts...)
}

// ScopedConfig DO NOT USE
type ScopedConfig struct {
	scopedConfigGeneric
	value string
}

// IsValid returns nil if the scoped configuration os valid otherwise a detailed
// error.
func (sc *ScopedConfig) IsValid() error {
	if sc.lastErr != nil {
		return sc.lastErr
	}
	return nil
}

func newScopedConfig() *ScopedConfig {
	return &ScopedConfig{
		value: "Hello Default Gophers",
	}
}

// WithDefaultConfig DO NOT USE
func WithDefaultConfig(h scope.TypeID) Option {
	return func(s *Service) error {
		return withDefaultConfig(h)(s)
	}
}

func withValue(h scope.TypeID, val string) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.value = val
		sc.ScopeID = h
		s.scopeCache[h] = sc
		return nil
	}
}
