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

import "github.com/corestoreio/cspkg/store/scope"

const configDefaultString = "Hello Default Gophers"
const configDefaultInt int = 42

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
	string
	int
}

// isValid returns nil if the scoped configuration os valid otherwise a detailed
// error.
func (sc *ScopedConfig) isValid() error {
	return sc.isValidPreCheck()
}

func newScopedConfig(target, parent scope.TypeID) *ScopedConfig {
	return &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(target, parent),
		string:              configDefaultString,
		int:                 configDefaultInt,
	}
}

// WithDefaultConfig DO NOT USE
func WithDefaultConfig(h scope.TypeID) Option {
	return func(s *Service) error {
		return withDefaultConfig(h)(s)
	}
}

func withString(val string, ids ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(ids...)
		sc.string = val
		return s.updateScopedConfig(sc)
	}
}

func withInt(val int, ids ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(ids...)
		sc.int = val
		return s.updateScopedConfig(sc)
	}
}
