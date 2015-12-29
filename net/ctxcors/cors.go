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

package ctxcors

import (
	"net/http"
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/log"
	"golang.org/x/net/context"
)

// Cors describes the CrossOriginResourceSharing which is used to create a
// Container Filter that implements CORS. Cross-origin resource sharing (CORS)
// is a mechanism that allows JavaScript on a web page to make XMLHttpRequests
// to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request
type Cors struct {
	config config.Getter

	// Log is a logger mainly for debugging. Default Logger writes to a black hole.
	Log log.Logger
	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool
	// Normalized list of plain allowed origins
	allowedOrigins []string
	// List of allowed origins containing wildcards
	allowedWOrigins []wildcard

	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool

	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool
	// Normalized list of allowed headers
	allowedHeaders []string
	// Normalized list of allowed methods
	allowedMethods []string
	// Normalized list of exposed headers
	exposedHeaders []string

	// maxAge in seconds will be added to the header, if set.
	maxAge string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
}

// New creates a new Cors handler with the provided options.
func New(opts ...Option) *Cors {
	c := &Cors{
		// Default is spec's "simple" methods
		allowedMethods: []string{"GET", "POST"},
		// Use sensible defaults
		allowedHeaders: []string{"Origin", "Accept", "Content-Type"},
		Log:            log.BlackHole{}, // debug and info logging disabled
	}
	return c.applyOpts(opts...)
}

func (c *Cors) applyOpts(opts ...Option) *Cors {
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}

// WithCORS to be used as a middleware for ctxhttp.Handler. Arguments can be used
// to apply the last time any options. This middleware does not take into account
// different configurations for different store scopes. The applied configuration
// is used for the all store scopes.
func (c *Cors) WithCORS(opts ...Option) ctxhttp.Middleware {
	c.applyOpts(opts...)
	csc := c.initCache()

	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var cc = c.current(csc, ctx)

			if r.Method == "OPTIONS" {
				if cc.Log.IsDebug() {
					cc.Log.Debug("ctxcors.Cors.WithCORS.handlePreflight", "method", r.Method, "OptionsPassthrough", cc.OptionsPassthrough)
				}
				cc.handlePreflight(w, r)
				// Preflight requests are standalone and should stop the chain as some other
				// middleware may not handle OPTIONS requests correctly. One typical example
				// is authentication middleware ; OPTIONS requests won't carry authentication
				// headers (see #1)
				if cc.OptionsPassthrough {
					return hf(ctx, w, r)
				}
				return nil
			}
			if cc.Log.IsDebug() {
				cc.Log.Debug("ctxcors.Cors.WithCORS.handleActualRequest", "method", r.Method, "cors", cc)
			}
			cc.handleActualRequest(w, r)
			return hf(ctx, w, r)
		}
	}
}

// current returns a non-nil pointer to a Cors. current is used within a request.
func (c *Cors) current(csc *corsScopeCache, ctx context.Context) *Cors {
	if c.config == nil || csc == nil {
		return c
	}

	_, st, err := store.FromContextReader(ctx)
	if err != nil {
		if c.Log.IsInfo() {
			c.Log.Info("ctxcors.Cors.current.store.FromContextReader", "err", err)
		}
		return c
	}

	var cc *Cors // cc == current CORS config the current request
	if cc = csc.get(st.WebsiteID()); cc == nil {
		cc = csc.insert(st.WebsiteID())
	}
	// todo: run a defer or goroutine to check if config changes
	// and if so delete the entry from the map
	return cc
}

// initCache if config.Getter has been set returns an initialized internal
// cache for different Cors configurations. Returns nil if config.Getter
// is not in use.
func (c *Cors) initCache() (cs *corsScopeCache) {
	if c.config != nil {
		cs = newCorsScopeCache(c.config, scope.WebsiteID, c)
	}
	return
}

// handlePreflight handles pre-flight CORS requests
func (c *Cors) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != "OPTIONS" {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handlePreflight.aborted", "method", r.Method)
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
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handlePreflight.aborted.empty.origin", "method", r.Method)
		}
		return
	}
	if false == c.isOriginAllowed(origin) {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handlePreflight.aborted.notAllowed.origin", "method", r.Method, "origin", origin)
		}
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if false == c.isMethodAllowed(reqMethod) {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handlePreflight.aborted.notAllowed.reqMethod", "method", r.Method, "reqMethod", reqMethod)
		}
		return
	}
	reqHeaders := parseHeaderList(r.Header.Get("Access-Control-Request-Headers"))
	if false == c.areHeadersAllowed(reqHeaders) {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handlePreflight.aborted.notAllowed.reqHeaders", "method", r.Method, "reqHeaders", reqHeaders)
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
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.maxAge != "" {
		headers.Set("Access-Control-Max-Age", c.maxAge)
	}
	if c.Log.IsDebug() {
		c.Log.Debug("ctxcors.Cors.handlePreflight.response.headers", "method", r.Method, "headers", headers)
	}
}

// handleActualRequest handles simple cross-origin requests, actual request or redirects
func (c *Cors) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method == "OPTIONS" {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handleActualRequest.aborted.options", "method", r.Method)
		}
		return
	}
	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.Add("Vary", "Origin")
	if origin == "" {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handleActualRequest.aborted.empty.origin", "method", r.Method)
		}
		return
	}
	if !c.isOriginAllowed(origin) {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handleActualRequest.aborted.notAllowed.origin", "method", r.Method, "origin", origin)
		}
		return
	}

	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !c.isMethodAllowed(r.Method) {
		if c.Log.IsDebug() {
			c.Log.Debug("ctxcors.Cors.handleActualRequest.aborted.notAllowed.method", "method", r.Method)
		}
		return
	}
	headers.Set("Access-Control-Allow-Origin", origin)
	if len(c.exposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(c.exposedHeaders, ", "))
	}
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.Log.IsDebug() {
		c.Log.Debug("ctxcors.Cors.handleActualRequest.response.headers", "headers", headers)
	}
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain requests
// on the endpoint
func (c *Cors) isOriginAllowed(origin string) bool {
	if c.AllowOriginFunc != nil {
		return c.AllowOriginFunc(origin)
	}
	if c.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, o := range c.allowedOrigins {
		if o == origin {
			return true
		}
	}
	for _, w := range c.allowedWOrigins {
		if w.match(origin) {
			return true
		}
	}
	return false
}

// isMethodAllowed checks if a given method can be used as part of a cross-domain request
// on the endpoing
func (c *Cors) isMethodAllowed(method string) bool {
	if len(c.allowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == "OPTIONS" {
		// Always allow preflight requests
		return true
	}
	for _, m := range c.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if a given list of headers are allowed to used within
// a cross-domain request.
func (c *Cors) areHeadersAllowed(requestedHeaders []string) bool {
	if c.allowedHeadersAll || len(requestedHeaders) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range c.allowedHeaders {
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
