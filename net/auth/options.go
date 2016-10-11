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

package auth

import (
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
)

// WithDefaultConfig applies the default configuration settings based for
// a specific scope.
//
// Default values are:
//		- basic auth
func WithDefaultConfig(scopeIDs ...scope.TypeID) Option {
	return withDefaultConfig(scopeIDs...)
}

func WithUnauthorizedHandler(uah mw.ErrorHandler, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.UnauthorizedHandler = uah
		return s.updateScopedConfig(sc)
	}
}

func WithSimpleBasicAuth(username, password, realm string, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		// sc.UnauthorizedHandler = uah
		return s.updateScopedConfig(sc)
	}
}

func WithBasicAuth(authFunc func(username, password string) bool, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		// sc.UnauthorizedHandler = uah
		return s.updateScopedConfig(sc)
	}
}
