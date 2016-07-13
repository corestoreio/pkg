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

package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// Option defines a function argument for the Cors type to apply options.
type Option func(*Service) error

// OptionFactoryFunc a closure around a scoped configuration to figure out which
// options should be returned depending on the scope brought to you during
// a request.
type OptionFactoryFunc func(config.Scoped) []Option

// WithDefaultConfig applies the default CORS configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Allowed Methods: GET, POST
//		- Allowed Headers: Origin, Accept, Content-Type
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	return withDefaultConfig(scp, id)
}

// WithExposedHeaders indicates which headers are safe to expose to the
// API of a CORS API specification.
func WithExposedHeaders(scp scope.Scope, id int64, headers ...string) Option {
	h := scope.NewHash(scp, id)
	exposedHeaders := convert(headers, http.CanonicalHeaderKey)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.exposedHeaders = exposedHeaders
		sc.scopeHash = h
		s.scopeCache[h] = sc
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
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.allowedOriginsAll = allowedOriginsAll
		sc.allowedOrigins = allowedOrigins
		sc.allowedWOrigins = allowedWOrigins
		sc.scopeHash = h
		s.scopeCache[h] = sc
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
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.allowOriginFunc = f
		sc.scopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithAllowedMethods is a list of methods the client is allowed to use with
// cross-domain requests. Default value is simple methods (GET and POST)
func WithAllowedMethods(scp scope.Scope, id int64, methods ...string) Option {
	h := scope.NewHash(scp, id)
	am := convert(methods, strings.ToUpper)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		// Allowed Methods
		// Note: for origins and methods matching, the spec requires a case-sensitive matching.
		// As it may error prone, we chose to ignore the spec here.

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.allowedMethods = am
		sc.scopeHash = h
		s.scopeCache[h] = sc
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
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.allowedHeadersAll = allowedHeadersAll
		sc.allowedHeaders = allowedHeaders
		sc.scopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithAllowCredentials convenient helper function.
// AllowCredentials indicates whether the request can include user credentials like
// cookies, HTTP authentication or client side SSL certificates.
func WithAllowCredentials(scp scope.Scope, id int64, ok bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.allowCredentials = ok
		sc.scopeHash = h
		s.scopeCache[h] = sc
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
			age = strconv.FormatFloat(sec, 'f', 0, 64)
		} else {
			return errors.NewNotValidf(errInvalidDurations, sec)
		}

		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.maxAge = age
		sc.scopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithOptionsPassthrough convenient helper function.
// OptionsPassthrough instructs preflight to let other potential next handlers to
// process the OPTIONS method. Turn this on if your application handles OPTIONS.
func WithOptionsPassthrough(scp scope.Scope, id int64, ok bool) Option {
	h := scope.NewHash(scp, id)
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.optionsPassthrough = ok
		sc.scopeHash = h
		s.scopeCache[h] = sc
		return nil
	}
}

// WithLogger applies a logger to the default scope which gets inherited to
// subsequent scopes. Mainly used for debugging.
func WithLogger(l log.Logger) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		s.Log = l
		for _, sc := range s.scopeCache {
			sc.log = l
		}
		return nil
	}
}

// withLoggerInit only sets the logger during init process and avoids
// overwriting existing settings.
func withLoggerInit(l log.Logger) Option {
	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		if s.Log == nil {
			s.Log = l
		}
		for _, sc := range s.scopeCache {
			if sc.log == nil {
				sc.log = l
			}
		}
		return nil
	}
}
