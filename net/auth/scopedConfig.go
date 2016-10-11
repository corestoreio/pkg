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
	"sort"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// AuthenticationTriggerFunc defines the condition if the AuthenticationFunc
// should be executed. An trigger can be for example a certain path or an IP
// address.
type AuthenticationTriggerFunc func(r *http.Request) bool

// AuthenticationFunc checks if a request is allowed to proceed. It returns nil
// on success. If you compare usernames and passwords make sure to use
// subtle.ConstantTimeCompare(). If callNext returns true the next authenticator
// gets called despite an occurred error, which gets dropped silently.
type AuthenticationFunc func(scopeID scope.TypeID, r *http.Request) (callNext bool, err error)

var defaultUnauthorizedHandler = mw.ErrorWithStatusCode(http.StatusUnauthorized)

// ScopedConfig contains the configuration for a specific scope.
type ScopedConfig struct {
	scopedConfigGeneric
	caseSensitivePath bool
	combineTrigger    bool
	// shouldDoAuthChecks a list of functions which checks if the authenticator
	// should be triggered. the first function to return true triggers the
	// authProvider.
	shouldDoAuthChecks
	// authProviders a list of functions which may return an error if authentication
	// fails.
	authProviders
	UnauthorizedHandler mw.ErrorHandler
}

// isValid check if the scoped configuration is valid when:
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
	if 0 == len(sc.authProviders) || 0 == len(sc.shouldDoAuthChecks) {
		return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeID, len(sc.shouldDoAuthChecks) == 0, len(sc.authProviders) == 0)
	}
	return nil
}

func newScopedConfig(target, parent scope.TypeID) *ScopedConfig {
	return &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(target, parent),
		caseSensitivePath:   true,
		combineTrigger:      false,
		UnauthorizedHandler: defaultUnauthorizedHandler,
	}
}

// Authenticate validates if a request is allowed to pass.
func (sc ScopedConfig) Authenticate(r *http.Request) error {
	// sc must not be a pointer, the data we're copying is small and it avoids
	// race conditions. we can change back to a pointer if tests and benchmarks
	// says different things.
	if len(sc.authProviders) == 0 || len(sc.shouldDoAuthChecks) == 0 {
		return errors.NewNotImplementedf("[auth] Authentication checker or provider not available")
	}
	if sc.Disabled {
		return nil
	}
	if !sc.shouldDoAuthChecks.triggerAuth(sc.combineTrigger, r) {
		return nil
	}
	if err := sc.authProviders.do(sc.ScopeID, r); err != nil {
		return errors.Wrapf(err, "[auth] Access denied to %q", r.URL.Path)
	}
	return nil
}

type authProvider struct {
	prio int
	AuthenticationFunc
}
type authProviders []authProvider

func (ap authProviders) sort()              { sort.Stable(ap) }
func (ap authProviders) Len() int           { return len(ap) }
func (ap authProviders) Less(i, j int) bool { return ap[i].prio < ap[j].prio }
func (ap authProviders) Swap(i, j int)      { ap[i], ap[j] = ap[j], ap[i] }

func (ap authProviders) do(scopeID scope.TypeID, r *http.Request) error {
	for i, apf := range ap {
		if err := apf.AuthenticationFunc(scopeID, r); err != nil {
			return errors.Wrapf(err, "[auth] Authentication failed at index %d", i)
		}
	}
	return nil
}

type shouldAuth struct {
	prio int
	AuthenticationTriggerFunc
}
type shouldDoAuthChecks []shouldAuth

func (ap shouldDoAuthChecks) sort()              { sort.Stable(ap) }
func (ap shouldDoAuthChecks) Len() int           { return len(ap) }
func (ap shouldDoAuthChecks) Less(i, j int) bool { return ap[i].prio < ap[j].prio }
func (ap shouldDoAuthChecks) Swap(i, j int)      { ap[i], ap[j] = ap[j], ap[i] }

func (ap shouldDoAuthChecks) triggerAuth(combined bool, r *http.Request) bool {
	var i int
	for _, check := range ap {
		ok := check.AuthenticationTriggerFunc(r)
		if ok && !combined {
			return true
		}
		if ok && combined {
			i++
		}
	}
	if combined {
		return i > 0 && i == len(ap)
	}
	return false
}
