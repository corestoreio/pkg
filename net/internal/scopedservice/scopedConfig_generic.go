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
	"net/http"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
)

// Auto generated: Do not edit. See net/internal/scopedService package for more details.

var defaultErrorHandler = mw.ErrorWithStatusCode(http.StatusServiceUnavailable)

// scopedConfigGeneric private internal scoped based configuration used for
// embedding into scopedConfig type. This type and its parent type ScopedConfig
// should be embedded.
type scopedConfigGeneric struct {
	// lastErr used during selecting the config from the scopeCache map and
	// inflight package.
	lastErr error
	// ScopeHash defines the scope to which this configuration is bound to.
	ScopeHash scope.Hash

	// todo think about adding config.Scoped

	// ErrorHandler gets called whenever a programmer makes an error. The
	// default handler prints the error to the client and returns
	// http.StatusServiceUnavailable
	mw.ErrorHandler
}

// newScopedConfigError easy helper to create an error
func newScopedConfigError(err error) ScopedConfig {
	return ScopedConfig{
		scopedConfigGeneric: scopedConfigGeneric{
			lastErr: err,
		},
	}
}

// newScopedConfigGeneric creates a new non-pointer generic config with a
// default scope and an error handler which returns status service unavailable.
// This function must be embedded in the targeted package newScopedConfig().
func newScopedConfigGeneric() scopedConfigGeneric {
	return scopedConfigGeneric{
		ScopeHash:    scope.DefaultHash,
		ErrorHandler: defaultErrorHandler,
	}
}

// optionInheritDefault looks up if the default configuration exists and if not
// creates a newScopedConfig(). This function can only be used within a
// functional option because it expects that it runs within an acquired lock
// because of the map.
func optionInheritDefault(s *Service) *ScopedConfig {
	if sc, ok := s.scopeCache[scope.DefaultHash]; ok && sc != nil {
		shallowCopy := new(ScopedConfig)
		*shallowCopy = *sc
		return shallowCopy
	}
	return newScopedConfig()
}
