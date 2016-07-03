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
	"github.com/corestoreio/csfw/util/errors"
)

// scopedConfig private internal scoped based configuration
type scopedConfig struct {
	scopedConfigGeneric

	// allowedOrigins normalized list of plain allowed origins
	allowedOrigins []string
	// List of allowed origins containing wildcards
	allowedWOrigins []wildcard

	// Normalized list of allowed headers
	allowedHeaders []string
	// Normalized list of allowed methods
	allowedMethods []string
	// Normalized list of exposed headers
	exposedHeaders []string

	// maxAge in seconds will be added to the header, if set.
	maxAge string

	// allowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	allowOriginFunc func(origin string) bool

	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool

	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool

	// allowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	allowCredentials bool

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	optionsPassthrough bool

	log log.Logger
}

// scopedConfig cannot use pointer based function receivers otherwise we will
// run into a very rare race condition.

// isValid a configuration for a scope is only then valid when
// - scopeHash set
// - min 1x allowedMethods set
// this is the only pointer receiver
func (sc scopedConfig) isValid() error {
	if sc.lastErr != nil {
		return errors.Wrap(sc.lastErr, "[cors] scopedConfig.isValid as an lastErr")
	}
	if sc.scopeHash > 0 && len(sc.allowedMethods) > 0 && sc.log != nil {
		return nil
	}
	return errors.NewNotValidf(errScopedConfigNotValid, sc.scopeHash)
}

func defaultScopedConfig() *scopedConfig {
	return &scopedConfig{
		scopedConfigGeneric: scopedConfigGeneric{
			scopeHash: scope.DefaultHash,
		},
		// Default is spec's "simple" methods
		allowedMethods: []string{"GET", "POST"},
		// Use sensible defaults
		allowedHeaders: []string{"Origin", "Accept", "Content-Type"},
		log:            log.BlackHole{}, // disabled info and debug logging
	}
}

// handlePreflight handles pre-flight CORS requests
func (sc scopedConfig) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != methodOptions {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handlePreflight.aborted", log.String("method", r.Method))
		}
		return
	}
	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	if origin == "" {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handlePreflight.aborted.empty.origin", log.String("method", r.Method))
		}
		return
	}
	if false == sc.isOriginAllowed(origin) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handlePreflight.aborted.notAllowed.origin", log.String("method", r.Method), log.String("origin", origin), log.Strings("allowedOrigins", sc.allowedOrigins...))
		}
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if false == sc.isMethodAllowed(reqMethod) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handlePreflight.aborted.notAllowed.reqMethod", log.String("method", r.Method), log.String("reqMethod", reqMethod))
		}
		return
	}
	reqHeaders := parseHeaderList(r.Header.Get("Access-Control-Request-Headers"))
	if false == sc.areHeadersAllowed(reqHeaders) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handlePreflight.aborted.notAllowed.reqHeaders", log.String("method", r.Method), log.Strings("reqHeaders", reqHeaders...))
		}
		return
	}
	headers.Set("Access-Control-Allow-Origin", origin)
	// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
	// by Access-Control-Request-Method (if supported) can be enough
	headers.Set("Access-Control-Allow-Methods", strings.ToUpper(reqMethod))
	if len(reqHeaders) > 0 {

		// Spec says: Since the list of headers can be unbounded, simply returning supported headers
		// from Access-Control-Request-Headers can be enough
		headers.Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
	}
	if sc.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if sc.maxAge != "" {
		headers.Set("Access-Control-Max-Age", sc.maxAge)
	}
	if sc.log.IsDebug() {
		sc.log.Debug("cors.Cors.handlePreflight.response.headers", log.String("method", r.Method), log.Object("headers", headers))
	}
}

// handleActualRequest handles simple cross-origin requests, actual request or redirects
func (sc scopedConfig) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method == methodOptions {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handleActualRequest.aborted.options", log.String("method", r.Method))
		}
		return
	}
	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.Add("Vary", "Origin")
	if origin == "" {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handleActualRequest.aborted.empty.origin", log.String("method", r.Method))
		}
		return
	}
	if !sc.isOriginAllowed(origin) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handleActualRequest.aborted.notAllowed.origin", log.String("method", r.Method), log.String("origin", origin))
		}
		return
	}

	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !sc.isMethodAllowed(r.Method) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.Cors.handleActualRequest.aborted.notAllowed.method", log.String("method", r.Method))
		}
		return
	}
	headers.Set("Access-Control-Allow-Origin", origin)
	if len(sc.exposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(sc.exposedHeaders, ", "))
	}
	if sc.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if sc.log.IsDebug() {
		sc.log.Debug("cors.Cors.handleActualRequest.response.headers", log.Object("headers", headers))
	}
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain requests
// on the endpoint
func (sc scopedConfig) isOriginAllowed(origin string) bool {
	if sc.allowOriginFunc != nil {
		return sc.allowOriginFunc(origin)
	}
	if sc.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, o := range sc.allowedOrigins {
		if o == origin {
			return true
		}
	}
	for _, w := range sc.allowedWOrigins {
		if w.match(origin) {
			return true
		}
	}
	return false
}

// isMethodAllowed checks if a given method can be used as part of a cross-domain request
// on the endpoing
func (sc scopedConfig) isMethodAllowed(method string) bool {
	if len(sc.allowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == methodOptions {
		// Always allow preflight requests
		return true
	}
	for _, m := range sc.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if a given list of headers are allowed to used within
// a cross-domain request.
func (sc scopedConfig) areHeadersAllowed(requestedHeaders []string) bool {
	if sc.allowedHeadersAll || len(requestedHeaders) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range sc.allowedHeaders {
			if h == header {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
