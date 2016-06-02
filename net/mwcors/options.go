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

package mwcors

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

// WithDefaultConfig applies the default CORS configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Allowed Methods: GET, POST
//		- Allowed Headers: Origin, Accept, Content-Type
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) (err error) {
		if h == scope.DefaultHash {
			s.defaultScopeCache, err = defaultScopedConfig()
			return errors.Wrap(err, "[mwcors] Default Scope with Default Config")
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		s.scopeCache[h], err = defaultScopedConfig()
		return errors.Wrapf(err, "[mwcors] Scope %s with Default Config", h)
	}
}

// WithExposedHeaders indicates which headers are safe to expose to the
// API of a CORS API specification.
func WithExposedHeaders(scp scope.Scope, id int64, headers ...string) Option {
	h := scope.NewHash(scp, id)
	exposedHeaders := convert(headers, http.CanonicalHeaderKey)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.exposedHeaders = exposedHeaders
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.exposedHeaders = exposedHeaders

		if sc, ok := s.scopeCache[h]; ok {
			sc.exposedHeaders = scNew.exposedHeaders
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

func convertAllowedOrigins(domains ...string) (allowedOriginsAll bool, allowedOrigins []string, allowedWOrigins []wildcard) {
	if len(domains) == 0 {
		// Default is all origins
		allowedOriginsAll = true
		return
	}

	for _, origin := range domains {
		// Normalize
		origin = strings.ToLower(origin)
		if origin == "*" {
			// If "*" is present in the list, turn the whole list into a match all
			allowedOriginsAll = true
			allowedOrigins = nil
			allowedWOrigins = nil
			return
		} else if i := strings.IndexByte(origin, '*'); i >= 0 {
			// Split the origin in two: start and end string without the *
			w := wildcard{origin[0:i], origin[i+1:]}
			allowedWOrigins = append(allowedWOrigins, w)
		} else {
			allowedOrigins = append(allowedOrigins, origin)
		}
	}
	return
}

// WithAllowedOrigins is a list of origins a cross-domain request can be executed from.
// If the special "*" value is present in the list, all origins will be allowed.
// An origin may contain a wildcard (*) to replace 0 or more characters
// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penality.
// Only one wildcard can be used per origin.
// Default value is ["*"]
func WithAllowedOrigins(scp scope.Scope, id int64, domains ...string) Option {
	h := scope.NewHash(scp, id)
	allowedOriginsAll, allowedOrigins, allowedWOrigins := convertAllowedOrigins(domains...)

	// Note: for origins and methods matching, the spec requires a case-sensitive matching.
	// As it may error prone, we chose to ignore the spec here.
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.allowedOriginsAll = allowedOriginsAll
			s.defaultScopeCache.allowedOrigins = allowedOrigins
			s.defaultScopeCache.allowedWOrigins = allowedWOrigins
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.allowedOriginsAll = allowedOriginsAll
		scNew.allowedOrigins = allowedOrigins
		scNew.allowedWOrigins = allowedWOrigins

		if sc, ok := s.scopeCache[h]; ok {
			sc.allowedOriginsAll = scNew.allowedOriginsAll
			sc.allowedOrigins = scNew.allowedOrigins
			sc.allowedWOrigins = scNew.allowedWOrigins
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithAllowOriginFunc convenient helper function.
// AllowOriginFunc is a custom function to validate the origin. It take the origin
// as argument and returns true if allowed or false otherwise. If this option is
// set, the content of AllowedOrigins is ignored.
func WithAllowOriginFunc(scp scope.Scope, id int64, f func(origin string) bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.allowOriginFunc = f
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.allowOriginFunc = f

		if sc, ok := s.scopeCache[h]; ok {
			sc.allowOriginFunc = scNew.allowOriginFunc
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithAllowedMethods is a list of methods the client is allowed to use with
// cross-domain requests. Default value is simple methods (GET and POST)
func WithAllowedMethods(scp scope.Scope, id int64, methods ...string) Option {
	h := scope.NewHash(scp, id)
	am := convert(methods, strings.ToUpper)
	return func(s *Service) error {
		// Allowed Methods
		// Note: for origins and methods matching, the spec requires a case-sensitive matching.
		// As it may error prone, we chose to ignore the spec here.

		if h == scope.DefaultHash {
			s.defaultScopeCache.allowedMethods = am
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.allowedMethods = am

		if sc, ok := s.scopeCache[h]; ok {
			sc.allowedMethods = scNew.allowedMethods
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

func convertAllowedHeaders(headers ...string) (allowedHeadersAll bool, allowedHeaders []string) {
	allowedHeaders = convert(append(headers, "Origin"), http.CanonicalHeaderKey)
	// Origin is always appended as some browsers will always request for this header at preflight
	//c.allowedHeaders = convert(append(headers, "Origin"), http.CanonicalHeaderKey)
	for _, h := range headers {
		if h == "*" {
			allowedHeadersAll = true
			allowedHeaders = nil
			return
		}
	}
	return
}

// WithAllowedHeaders is list of non simple headers the client is allowed to use with
// cross-domain requests.
// If the special "*" value is present in the list, all headers will be allowed.
// Default value is [] but "Origin" is always appended to the list.
func WithAllowedHeaders(scp scope.Scope, id int64, headers ...string) Option {
	h := scope.NewHash(scp, id)
	allowedHeadersAll, allowedHeaders := convertAllowedHeaders(headers...)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.allowedHeadersAll = allowedHeadersAll
			s.defaultScopeCache.allowedHeaders = allowedHeaders
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.allowedHeadersAll = allowedHeadersAll
		scNew.allowedHeaders = allowedHeaders

		if sc, ok := s.scopeCache[h]; ok {
			sc.allowedHeadersAll = scNew.allowedHeadersAll
			sc.allowedHeaders = scNew.allowedHeaders
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithAllowCredentials convenient helper function.
// AllowCredentials indicates whether the request can include user credentials like
// cookies, HTTP authentication or client side SSL certificates.
func WithAllowCredentials(scp scope.Scope, id int64, ok bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.allowCredentials = ok
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.allowCredentials = ok

		if sc, ok := s.scopeCache[h]; ok {
			sc.allowCredentials = scNew.allowCredentials
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithMaxAge indicates how long (in seconds) the results of a preflight request
// can be cached
func WithMaxAge(scp scope.Scope, id int64, seconds time.Duration) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if seconds < 1 {
			return nil
		}

		var age string
		if sec := seconds.Seconds(); sec > 0 {
			age = fmt.Sprintf("%.0f", sec)
		} else {
			return errors.NewNotValidf(errInvalidDurations, sec)
		}

		if h == scope.DefaultHash {
			s.defaultScopeCache.maxAge = age
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.maxAge = age

		if sc, ok := s.scopeCache[h]; ok {
			sc.maxAge = scNew.maxAge
			scNew = sc
		}
		scNew.scopeHash = h
		s.scopeCache[h] = scNew
		return nil
	}
}

// WithOptionsPassthrough convenient helper function.
// OptionsPassthrough instructs preflight to let other potential next handlers to
// process the OPTIONS method. Turn this on if your application handles OPTIONS.
func WithOptionsPassthrough(scp scope.Scope, id int64, ok bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		if h == scope.DefaultHash {
			s.defaultScopeCache.optionsPassthrough = ok
			return nil
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		// inherit default config
		scNew := s.defaultScopeCache
		scNew.optionsPassthrough = ok

		if sc, ok := s.scopeCache[h]; ok {
			sc.optionsPassthrough = scNew.optionsPassthrough
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
//	cfgStruct, err := backendcors.NewConfigStructure()
//	if err != nil {
//		panic(err)
//	}
//	pb := backendcors.New(cfgStruct)
//
//	cors := mwcors.MustNewService(
//		mwcors.WithOptionFactory(backendcors.PrepareOptions(pb)),
//	)
func WithOptionFactory(f ScopedOptionFunc) Option {
	return func(s *Service) error {
		s.scpOptionFnc = f
		return nil
	}
}
