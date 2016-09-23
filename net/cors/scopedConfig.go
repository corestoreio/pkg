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

// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cors

import (
	"net/http"
	"strings"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
)

// ScopedConfig scoped based configuration and should not be embedded into your
// own types. Call ScopedConfig.ScopeID to know to which scope this
// configuration has been bound to.
type ScopedConfig struct {
	scopedConfigGeneric
	log log.Logger

	// Settings general CORS settings
	Settings
}

// IsValid a configuration for a scope is only then valid when
//	- ScopeID set
//	- min 1x allowedMethods set
//	- Logger not nil
func (sc *ScopedConfig) IsValid() error {
	if err := sc.isValid(); err != nil {
		return errors.Wrap(err, "[cors] scopedConfig.isValid as an lastErr")
	}
	// AML = allowed method length: Max 7, also useful for testing ;-)
	if aml := len(sc.AllowedMethods); sc.ScopeID > 0 && aml > 0 && aml <= 7 && sc.log != nil {
		return nil
	}
	return errors.NewNotValidf(errScopedConfigNotValid, sc.ScopeID, sc.AllowedMethods, sc.log == nil)
}

// newScopedConfig creates a new object with the minimum needed configuration.
func newScopedConfig() *ScopedConfig {
	return &ScopedConfig{
		scopedConfigGeneric: newScopedConfigGeneric(),
		log:                 log.BlackHole{}, // disabled info and debug logging
		Settings: Settings{
			// Default is spec's "simple" methods
			AllowedMethods: []string{"GET", "POST"},
			// Use sensible defaults
			AllowedHeaders: []string{"Origin", "Accept", "Content-Type"},
		},
	}
}

// handlePreflight handles pre-flight CORS requests
func (sc *ScopedConfig) handlePreflight(w http.ResponseWriter, r *http.Request) {
	sc.log = log.BlackHole{}

	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != methodOptions {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handlePreflight.aborted", log.String("method", r.Method))
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
			sc.log.Debug("cors.handlePreflight.aborted.empty.origin", log.String("method", r.Method))
		}
		return
	}
	if !sc.isOriginAllowed(origin) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handlePreflight.aborted.notAllowed.origin", log.String("method", r.Method), log.String("origin", origin), log.Strings("allowedOrigins", sc.AllowedOrigins...))
		}
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if !sc.isMethodAllowed(reqMethod) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handlePreflight.aborted.notAllowed.reqMethod", log.String("method", r.Method), log.String("reqMethod", reqMethod))
		}
		return
	}
	reqHeaders := parseHeaderList(r.Header.Get("Access-Control-Request-Headers"))
	if !sc.areHeadersAllowed(reqHeaders) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handlePreflight.aborted.notAllowed.reqHeaders", log.String("method", r.Method), log.Strings("reqHeaders", reqHeaders...))
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
	if sc.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if sc.MaxAge != "" {
		headers.Set("Access-Control-Max-Age", sc.MaxAge)
	}
	if sc.log.IsDebug() {
		sc.log.Debug("cors.handlePreflight.response.headers", log.String("method", r.Method), log.Object("headers", headers))
	}
}

// handleActualRequest handles simple cross-origin requests, actual request or redirects
func (sc *ScopedConfig) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method == methodOptions {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handleActualRequest.aborted.options", log.String("method", r.Method))
		}
		return
	}
	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.Add("Vary", "Origin")
	if origin == "" {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handleActualRequest.aborted.empty.origin", log.String("method", r.Method))
		}
		return
	}
	if !sc.isOriginAllowed(origin) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handleActualRequest.aborted.notAllowed.origin", log.String("method", r.Method), log.String("origin", origin))
		}
		return
	}

	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !sc.isMethodAllowed(r.Method) {
		if sc.log.IsDebug() {
			sc.log.Debug("cors.handleActualRequest.aborted.notAllowed.method", log.String("method", r.Method))
		}
		return
	}
	headers.Set("Access-Control-Allow-Origin", origin)
	if len(sc.ExposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(sc.ExposedHeaders, ", "))
	}
	if sc.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if sc.log.IsDebug() {
		sc.log.Debug("cors.handleActualRequest.response.headers", log.Object("headers", headers))
	}
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain requests
// on the endpoint
func (sc *ScopedConfig) isOriginAllowed(origin string) bool {
	if sc.AllowOriginFunc != nil {
		return sc.AllowOriginFunc(origin)
	}
	if sc.AllowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, o := range sc.AllowedOrigins {
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
func (sc *ScopedConfig) isMethodAllowed(method string) bool {
	if len(sc.AllowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == methodOptions {
		// Always allow preflight requests
		return true
	}
	for _, m := range sc.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if a given list of headers are allowed to used within
// a cross-domain request.
func (sc *ScopedConfig) areHeadersAllowed(requestedHeaders []string) bool {
	if sc.AllowedHeadersAll || len(requestedHeaders) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range sc.AllowedHeaders {
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
