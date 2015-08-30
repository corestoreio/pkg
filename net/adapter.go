// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package net

import (
	"github.com/corestoreio/csfw/utils/log"
	"net/http"
)

const (
	// HTTPMethodOverrideHeader represents a commonly used http header to override a request method.
	HTTPMethodOverrideHeader = "X-HTTP-Method-Override"
	// HTTPMethodOverrideFormKey represents a commonly used HTML form key to override a request method.
	HTTPMethodOverrideFormKey = "_method"
)

// HTTPMethodxxx defines the available methods which this framework supports
const (
	HTTPMethodHead    = `HEAD`
	HTTPMethodGet     = "GET"
	HTTPMethodPost    = "POST"
	HTTPMethodPut     = "PUT"
	HTTPMethodPatch   = "PATCH"
	HTTPMethodDelete  = "DELETE"
	HTTPMethodTrace   = "TRACE"
	HTTPMethodOptions = "OPTIONS"
)

// @see https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81

// Adapter is a wrapper for the http.Handler
type Adapter func(http.Handler) http.Handler

// Adapt function will iterate over all adapters, calling them one by one
// in a chained manner, returning the result of the final adapter.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, a := range adapters {
		h = a(h)
	}
	return h
}

// WithHeader is an Adapter that sets an HTTP handler.
func WithHeader(key, value string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			h.ServeHTTP(w, r)
		})
	}
}

// SupportXHTTPMethodOverride adds support for the X-HTTP-Method-Override
// header. Submitted value will be checked against known methods. Adding
// HTTPMethodOverrideFormKey to any form will take precedence before
// HTTP header. If an unknown method will be submitted it gets logged as an
// Info log.
func SupportXHTTPMethodOverride() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mo := r.FormValue(HTTPMethodOverrideFormKey)
			if mo == "" {
				mo = r.Header.Get(HTTPMethodOverrideHeader)
			}
			switch mo {
			case "": // do nothing
			case HTTPMethodHead, HTTPMethodGet, HTTPMethodPost, HTTPMethodPut, HTTPMethodPatch, HTTPMethodDelete, HTTPMethodTrace, HTTPMethodOptions:
				r.Method = mo
			default:
				// not sure if an error is here really needed ...
				if log.IsInfo() {
					log.Info("net.SupportXHTTPMethodOverride.switch", "err", "Unknown http method", "method", mo, "form", r.Form.Encode(), "header", r.Header)
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
