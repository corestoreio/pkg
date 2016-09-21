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
	"net/http"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
)

type Authenticator interface {
	// Authenticate authenticates a request and returns nil on success.
	// You must use subtle.ConstantTimeCompare()
	Authenticate(h scope.TypeID, r *http.Request) error
}

// ScopedConfig private internal scoped based configuration
type ScopedConfig struct {
	scopedConfigGeneric

	log      log.Logger
	disabled bool
	// if nil fall back to default scope
	Authenticator
	loginHandler  http.Handler // e.g. basic auth browser popup
	deniedHandler http.Handler
}

// IsValid a configuration for a scope is only then valid when the Key has been
// supplied, a non-nil signing method and a non-nil Verifier.
func (sc ScopedConfig) IsValid() bool {
	return sc.ScopeHash > 0 && sc.Authenticator != nil && sc.disabled
}

func defaultScopedConfig() (ScopedConfig, error) {
	return ScopedConfig{
		ScopeHash: scope.DefaultTypeID,
		log:       log.BlackHole{}, // disabled info and debug logging
	}, nil
}
