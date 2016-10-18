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
	"strings"

	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// WithDefaultConfig applies the default configuration settings for
// a specific scope.
//
// Default values are:
//		- authentication returns always access denied
//		- all resources protected
func WithDefaultConfig(scopeIDs ...scope.TypeID) Option {
	return withDefaultConfig(scopeIDs...)
}

// WithUnauthorizedHandler sets the handler which calls the interface to request
// data from a user after the authentication failed.
func WithUnauthorizedHandler(uah mw.ErrorHandler, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.UnauthorizedHandler = uah
		return s.updateScopedConfig(sc)
	}
}

// WithUnauthorizedRedirect redirects if the authorization fails.
func WithUnauthorizedRedirect(url string, code int, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.UnauthorizedHandler = func(_ error) http.Handler {
			return http.RedirectHandler(url, code)
		}
		return s.updateScopedConfig(sc)
	}
}

func matchPath(caseSensitivePath bool, r *http.Request, other string) bool {
	if caseSensitivePath {
		return strings.HasPrefix(r.URL.Path, other)
	}
	return strings.HasPrefix(strings.ToLower(r.URL.Path), strings.ToLower(other))
}

// WithResourceACLs enables to define specific URL paths to be black- and/or
// white listed. Matching for black- and white lists checks if the URL path has
// the provided string of a list as a prefix.
//		auth.WithResources(nil,nil) // blocks everything
//		auth.WithResources([]string{"/"}, []string{}) // blocks everything
//		auth.WithResources([]string{"/"}, []string{"/catalog"}) // blocks everything except the routes starting with /catalog.
// Providing no scopeIDs applies the resource ACL to the default scope ID. The
// string based ACL checks will always be executed before REGEX based ACL
// checks, if both functional options have been provided.
func WithResourceACLs(blacklist, whitelist []string, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		isCaseSensitive := sc.caseSensitivePath // copy the value to avoid races
		sc.triggers = append(sc.triggers, authTrigger{
			prio: -10,
			TriggerFunc: func(r *http.Request) bool {
				blocked := len(blacklist) == 0
				for _, b := range blacklist {
					if matchPath(isCaseSensitive, r, b) {
						blocked = true
					}
				}
				if blocked {
					for _, w := range whitelist {
						if matchPath(isCaseSensitive, r, w) {
							return false
						}
					}
				}
				return blocked
			},
		})
		sc.triggers.sort()
		return s.updateScopedConfig(sc)
	}
}

func strSliceToRegexSlice(sl []string) ([]*regexp.Regexp, error) {
	rs := make([]*regexp.Regexp, 0, len(sl))
	for i, b := range sl {
		if b == "" {
			continue
		}
		r, err := regexp.Compile(b)
		if err != nil {
			return nil, errors.NewFatalf("[auth] Failed to compile regex %q at index %d", b, i)
		}
		rs = append(rs, r)
	}
	return rs, nil
}

// WithResourceRegexpACLs same as WithResourceACLs but uses the slow
// pre-compiled and more powerful regexes.
func WithResourceRegexpACLs(block, whitelist []string, scopeIDs ...scope.TypeID) Option {
	br, err := strSliceToRegexSlice(block)
	if err != nil {
		return func(s *Service) error {
			return errors.Wrap(err, "[auth] WithResourcesRegexp black list")
		}
	}
	wlr, err := strSliceToRegexSlice(whitelist)
	if err != nil {
		return func(s *Service) error {
			return errors.Wrap(err, "[auth] WithResourcesRegexp white list")
		}
	}

	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.triggers = append(sc.triggers, authTrigger{
			prio: -5,
			TriggerFunc: func(r *http.Request) bool {
				block := len(br) == 0
				for _, blockr := range br {
					if blockr.MatchString(r.URL.Path) {
						block = true
					}
				}
				if block {
					for _, whiter := range wlr {
						if whiter.MatchString(r.URL.Path) {
							return false
						}
					}
				}
				return block
			},
		})
		sc.triggers.sort()
		return s.updateScopedConfig(sc)
	}
}

// WithCombineTriggers setting to true forces all authentication triggers to
// return true. Otherwise the first trigger which returns true, triggers the
// authentication providers. Default value: false.
func WithCombineTriggers(combine bool, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.combineTrigger = combine
		return s.updateScopedConfig(sc)
	}
}

// WithTrigger sets the authentication trigger function which implements a
// condition to check if the list of authentication providers should be called.
// Subsequent calls of this functional option will add more TriggerFuncs to the
// internal list. If not trigger has been applied the authentication providers
// will always be called.
func WithTrigger(tf TriggerFunc, priority int, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.triggers = append(sc.triggers, authTrigger{
			prio:        priority,
			TriggerFunc: tf,
		})
		return s.updateScopedConfig(sc)
	}
}

// WithProvider sets the authentication provider function which checks if a
// request should be considered valid to call the next HTTP handler on err ==
// nil or even call the next provider. Subsequent calls of this functions will
// add more ProviderFuncs to the internal list. This internal list cannot yet be
// cleared or reset.
func WithProvider(pf ProviderFunc, priority int, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.providers = append(sc.providers, authProvider{
			prio:         priority,
			ProviderFunc: pf,
		})
		return s.updateScopedConfig(sc)
	}
}

// WithSimpleBasicAuth sets a single username/password for a scope. Username and
// password must be provided as "plain text" arguments. This basic auth handler
// calls the next authentication provider if the authentication fails. Username
// and password will be compared in constant time.
func WithSimpleBasicAuth(username, password, realm string, scopeIDs ...scope.TypeID) Option {
	ba256, err := basicAuthValidator("sha256", username, password)
	if err != nil {
		return func(s *Service) error {
			return errors.Wrap(err, "[auth] WithSimpleBasicAuth basicAuthHashed")
		}
	}
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.providers = append(sc.providers, authProvider{
			prio:         10,
			ProviderFunc: basicAuth(ba256),
		})
		sc.UnauthorizedHandler = basicAuthHandler(realm)
		return s.updateScopedConfig(sc)
	}
}

// WithBasicAuth provides the basic authentication header but allows to set a
// custom function to compare the input data of username and password.
func WithBasicAuth(baf BasicAuthFunc, realm string, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		sc.providers = append(sc.providers, authProvider{
			prio:         1,
			ProviderFunc: basicAuth(baf),
		})
		sc.UnauthorizedHandler = basicAuthHandler(realm)
		return s.updateScopedConfig(sc)
	}
}

// prioIncrement only used for testing to trigger the sorting. This variable
// should not trigger any race conditions.
var prioIncrement = 1000

// WithInvalidAuth authentication will always fail. Mainly used for testing ;-)
func WithInvalidAuth(callNext bool, scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		prioIncrement++
		sc.providers = append(sc.providers, authProvider{
			prio: prioIncrement,
			ProviderFunc: func(scopeID scope.TypeID, r *http.Request) (bool, error) {
				return callNext, errors.NewUnauthorizedf("[auth] Access denied in Scope %s for path %q", scopeID, r.URL.Path)
			},
		})
		sc.providers.sort()
		return s.updateScopedConfig(sc)
	}
}

// WithValidAuth authentication will always succeed. Mainly used for testing ;-)
func WithValidAuth(scopeIDs ...scope.TypeID) Option {
	return func(s *Service) error {
		sc := s.findScopedConfig(scopeIDs...)
		prioIncrement++
		sc.providers = append(sc.providers, authProvider{
			prio: prioIncrement,
			ProviderFunc: func(_ scope.TypeID, _ *http.Request) (bool, error) {
				return false, nil
			},
		})
		sc.providers.sort()
		return s.updateScopedConfig(sc)
	}
}
