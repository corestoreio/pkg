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

package mw

import (
	"net/http"
	"time"

	"context"

	"github.com/corestoreio/csfw/net/httputil"
)

// Middleware is a wrapper for the function http.HandlerFunc to create
// middleware functions.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// MiddlewareSlice a slice full of middleware functions and with function
// receivers attached
type MiddlewareSlice []Middleware

// Chain will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func Chain(h http.HandlerFunc, mws ...Middleware) http.HandlerFunc {
	// Chain middleware with handler in the end
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Chain will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func (mws MiddlewareSlice) Chain(h http.HandlerFunc) http.HandlerFunc {
	return Chain(h, mws...)
}

// WithHeader is a middleware that sets multiple HTTP headers. Will panic if kv
// is imbalanced. len(kv)%2 == 0.
func WithHeader(kv ...string) Middleware {
	lkv := len(kv)
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for i := 0; i < lkv; i = i + 2 {
				w.Header().Set(kv[i], kv[i+1])
			}
			hf(w, r)
		}
	}
}

// WithXHTTPMethodOverride adds support for the X-HTTP-Method-Override
// header. Submitted value will be checked against known methods. Adding
// HTTPMethodOverrideFormKey to any form will take precedence before
// HTTP header. If an unknown method will be submitted it gets logged as an
// Info log. This function is chainable.
func WithXHTTPMethodOverride(opts ...Option) Middleware {
	ob := newOptionBox(opts...)
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			mo := r.FormValue(httputil.MethodOverrideFormKey)
			if mo == "" {
				mo = r.Header.Get(httputil.MethodOverrideHeader)
			}
			switch mo {
			case "": // do nothing
			case httputil.MethodHead, httputil.MethodGet, httputil.MethodPost, httputil.MethodPut, httputil.MethodPatch, httputil.MethodDelete, httputil.MethodTrace, httputil.MethodOptions:
				r.Method = mo
			default:
				// not sure if an error is here really needed ...
				if ob.log.IsDebug() {
					ob.log.Debug("ctxhttp.SupportXHTTPMethodOverride.switch", "err", "Unknown http method", "method", mo, "form", r.Form.Encode(), "header", r.Header)
				}
			}
			hf(w, r)
		}
	}
}

// WithCloseNotify returns a ctxhttp.Handler cancelling the context when the client
// connection close unexpectedly.
func WithCloseNotify(opts ...Option) Middleware {
	ob := newOptionBox(opts...)
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Cancel the context if the client closes the connection
			if wcn, ok := w.(http.CloseNotifier); ok {
				ctx, cancel := context.WithCancel(r.Context())
				defer cancel()
				r = r.WithContext(ctx)

				notify := wcn.CloseNotify()
				go func() {
					<-notify
					cancel()
					if ob.log.IsDebug() {
						ob.log.Debug("ctxhttp.WithCloseNotify.cancel", "cancelled", true, "request", r)
					}
				}()
			}
			hf(w, r)
		}
	}
}

// WithTimeout returns a ctxhttp.Handler which adds a timeout to the context.
//
// Child handlers have the responsibility to obey the context deadline and to return
// an appropriate error (or not) response in case of timeout.
func WithTimeout(timeout time.Duration) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := context.WithTimeout(r.Context(), timeout)
			r = r.WithContext(ctx)
			hf(w, r)
		}
	}
}
