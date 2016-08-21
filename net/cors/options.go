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
	"strings"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/store/scope"
)

// Settings general settings for the cors service. Those settings will be
// applied via functional options on a per scope basis.
type Settings struct {
	// AllowedOrigins is a list of origins a cross-domain request can be
	// executed from. If the special "*" value is present in the list, all
	// origins will be allowed. An origin may contain a wildcard (*) to replace
	// 0 or more characters (i.e.: http://*.domain.com). Usage of wildcards
	// implies a small performance penality. Only one wildcard can be used per
	// origin. Default value is ["*"]. Normalized list of plain allowed origins.
	AllowedOrigins []string
	// allowedWOrigins a list of allowed origins containing wildcards
	allowedWOrigins []wildcard
	// AllowedHeaders normalized list of allowed headers the client is allowed
	// to use with cross-domain requests. If the special "*" value is present in
	// the list, all headers will be allowed. Default value is [] but "Origin"
	// is always appended to the list.
	AllowedHeaders []string
	// AllowedMethods normalized list of methods, the client is allowed to use
	// with cross-domain requests. Default value is simple methods (GET and
	// POST)
	AllowedMethods []string
	// ExposedHeaders indicates which headers are safe to expose to the
	// API of a CORS API specification. Normalized list of exposed headers.
	ExposedHeaders []string

	// MaxAge in seconds will be added to the header, if set. Indicates how long
	// (in seconds) the results of a preflight request can be cached.
	MaxAge string

	// AllowOriginFunc is a custom function to validate the origin. It take the
	// origin as argument and returns true if allowed or false otherwise. If
	// this option is set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool

	// AllowedOriginsAll set to true when allowed origins contains a "*"
	AllowedOriginsAll bool

	// AllowedHeadersAll set to true when allowed headers contains a "*"
	AllowedHeadersAll bool

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
}

// WithDefaultConfig applies the default CORS configuration settings based for
// a specific scope. This function overwrites any previous set options.
//
// Default values are:
//		- Allowed Methods: GET, POST
//		- Allowed Headers: Origin, Accept, Content-Type
func WithDefaultConfig(scp scope.Scope, id int64) Option {
	return withDefaultConfig(scp, id)
}

func WithSettings(scp scope.Scope, id int64, stng Settings) Option {
	h := scope.NewHash(scp, id)
	exposedHeaders := convert(stng.ExposedHeaders, http.CanonicalHeaderKey)
	allowedOriginsAll, allowedOrigins, allowedWOrigins := convertAllowedOrigins(stng.AllowedOrigins...)
	am := convert(stng.AllowedMethods, strings.ToUpper)

	allowedHeadersAll, allowedHeaders := convertAllowedHeaders(stng.AllowedHeaders...)

	return func(s *Service) error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		sc := s.scopeCache[h]
		if sc == nil {
			sc = optionInheritDefault(s)
		}
		sc.ExposedHeaders = exposedHeaders

		// Note: for origins and methods matching, the spec requires a
		// case-sensitive matching. As it may error prone, we chose to ignore
		// the spec here.
		sc.AllowedOriginsAll = allowedOriginsAll
		sc.AllowedOrigins = allowedOrigins
		sc.allowedWOrigins = allowedWOrigins

		sc.AllowOriginFunc = stng.AllowOriginFunc
		sc.AllowedMethods = am

		sc.AllowedHeadersAll = allowedHeadersAll
		sc.AllowedHeaders = allowedHeaders

		sc.AllowCredentials = stng.AllowCredentials
		sc.MaxAge = stng.MaxAge
		sc.OptionsPassthrough = stng.OptionsPassthrough

		sc.ScopeHash = h
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
