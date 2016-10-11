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

	"regexp"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// Authenticator ...
type Authenticator interface {
	// Authenticate authenticates a request and returns nil on success.
	// You must use subtle.ConstantTimeCompare()
	Authenticate(scopeID scope.TypeID, r *http.Request) error
}

var defaultUnauthorizedHandler = mw.ErrorWithStatusCode(http.StatusUnauthorized)

// ScopedConfig contains the configuration for a specific scope.
type ScopedConfig struct {
	scopedConfigGeneric
	Authenticator
	// Resources protects all mentioned routes. If empty protects everything.
	Resources []string
	// ResourcesWhiteList disables authentication for all mentioned routes.
	ResourcesWhiteList []string
	// ResourcesRegExp protects all mentioned routes matched by a regular
	// expression. If empty protects everything.
	ResourcesRegExp []*regexp.Regexp
	// ResourcesRegExpWhiteList disables authentication all mentioned routes matched
	// by a regular expression.
	ResourcesRegExpWhiteList []*regexp.Regexp
	UnauthorizedHandler      mw.ErrorHandler
}

// IsValid check if the scoped configuration is valid when:
//		- Authenticator
//		- UnauthorizedHandler
// has been set and no other previous error has occurred.
func (sc *ScopedConfig) isValid() error {
	if err := sc.isValidPreCheck(); err != nil {
		return errors.Wrap(err, "[auth] ScopedConfig.isValid as an lastErr")
	}
	if sc.Disabled {
		return nil
	}
	if sc.Authenticator == nil {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeID, sc.Authenticator == nil)
	}
	return nil
}

func newScopedConfig(target, parent scope.TypeID) *ScopedConfig {
	return &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(target, parent),
		UnauthorizedHandler: defaultUnauthorizedHandler,
	}
}

func (sc *ScopedConfig) authenticate(r *http.Request) error {
	return nil
}
