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

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/store/scope"
)

// TriggerFunc defines the condition if the ProviderFunc
// should be executed. An trigger can be for example a certain path or an IP
// address.
type TriggerFunc func(r *http.Request) bool

// ProviderFunc checks if a request is allowed to proceed. It returns nil
// on success. If you compare usernames and passwords make sure to use
// subtle.ConstantTimeCompare(). If callNext returns true the next authenticator
// gets called despite an occurred error, which gets dropped silently. If all
// ProviderFuncs return true to call the next, then the last function call
// gets a force checked error.
type ProviderFunc func(scopeID scope.TypeID, r *http.Request) (callNext bool, err error)

var defaultUnauthorizedHandler = mw.ErrorWithStatusCode(http.StatusUnauthorized)

// ScopedConfig contains the configuration for a specific scope.
type ScopedConfig struct {
	scopedConfigGeneric
	caseSensitivePath bool
	combineTrigger    bool
	// triggers a list of functions which checks if the authenticator
	// should be triggered. the first function to return true triggers the
	// authProvider.
	triggers
	// providers a list of functions which may return an error if authentication
	// fails or even call the next provider in the chain.
	providers
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
	if 0 == len(sc.providers) || 0 == len(sc.triggers) {
		return errors.NotValid.Newf(errScopedConfigNotValid, sc.ScopeID, len(sc.triggers) == 0, len(sc.providers) == 0)
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
	if sc.Disabled {
		return nil
	}
	// an empty or nil authTriggers always returns true to trigger the
	// authentication provider.
	if !sc.triggers.do(sc.combineTrigger, r) {
		return nil
	}
	if err := sc.providers.do(sc.ScopeID, r); err != nil {
		return errors.Wrapf(err, "[auth] Access denied to %q", r.URL.Path)
	}
	return nil
}

type authProvider struct {
	prio int
	ProviderFunc
}
type providers []authProvider

func (ap providers) sort()              { sort.Stable(ap) }
func (ap providers) Len() int           { return len(ap) }
func (ap providers) Less(i, j int) bool { return ap[i].prio < ap[j].prio }
func (ap providers) Swap(i, j int)      { ap[i], ap[j] = ap[j], ap[i] }

// do iterates over the Authenticators and checks if it should call the next
// Authenticator and drop silently the error. Even if the last Authenticator
// forces to call the next non-existent Authenticator this function detects that
// and checks the very last error.
func (ap providers) do(scopeID scope.TypeID, r *http.Request) error {
	if len(ap) == 0 {
		return errors.NotImplemented.Newf("[auth] No authentication provider available")
	}

	nc := 1 // nc == next counter
	for i, apf := range ap {
		next, err := apf.ProviderFunc(scopeID, r)
		if next && nc < len(ap) {
			nc++
			continue
		}
		if err != nil {
			return errors.Wrapf(err, "[auth] Authentication failed at index %d", i)
		}
	}
	return nil
}

// authTrigger checks if the authentication should be triggered. Contains a
// priority in which these trigger functions gets executed.
type authTrigger struct {
	prio int
	TriggerFunc
}

type triggers []authTrigger

func (ap triggers) sort()              { sort.Stable(ap) }
func (ap triggers) Len() int           { return len(ap) }
func (ap triggers) Less(i, j int) bool { return ap[i].prio < ap[j].prio }
func (ap triggers) Swap(i, j int)      { ap[i], ap[j] = ap[j], ap[i] }

func (ap triggers) do(combined bool, r *http.Request) bool {
	if len(ap) == 0 {
		return true
	}
	var i int
	for _, check := range ap {
		ok := check.TriggerFunc(r)
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
