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
	"context"
	"net/http"
	"time"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
)

// Middleware is a wrapper for the interface http.Handler to create
// middleware functions.
type Middleware func(http.Handler) http.Handler

// MiddlewareSlice a slice full of middleware functions and with function
// receivers attached
type MiddlewareSlice []Middleware

// Chain will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	// Chain middleware with handler in the end
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h) // performance penalty because of bounds checking of the compiler
	}
	return h
}

// ChainFunc will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func ChainFunc(hf http.HandlerFunc, mws ...Middleware) http.Handler {
	return Chain(hf, mws...)
}

// Chain will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func (mws MiddlewareSlice) Chain(h http.Handler) http.Handler {
	return Chain(h, mws...)
}

// Chain will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func (mws MiddlewareSlice) ChainFunc(hf http.HandlerFunc) http.Handler {
	return Chain(hf, mws...)
}

// Append extends a slice, adding the specified Middleware
// as the last ones in the request flow.
//
// Append returns a new slice, leaving the original one untouched.
func (c MiddlewareSlice) Append(mws ...Middleware) MiddlewareSlice {
	newMWS := make(MiddlewareSlice, len(c)+len(mws))
	copy(newMWS, c)
	copy(newMWS[len(c):], mws)
	return newMWS
}

// WithHeader is a middleware that sets multiple HTTP headers. Will panic if kv
// is imbalanced. len(kv)%2 == 0.
func WithHeader(kv ...string) Middleware {
	lkv := len(kv)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for i := 0; i < lkv; i = i + 2 {
				w.Header().Set(kv[i], kv[i+1])
			}
			h.ServeHTTP(w, r)
		})
	}
}

const (
	// MethodOverrideHeader represents a commonly used http header to override a request method.
	MethodOverrideHeader = "X-HTTP-Method-Override"
	// MethodOverrideFormKey represents a commonly used HTML form key to override a request method.
	MethodOverrideFormKey = "_method"
)

// WithXHTTPMethodOverride adds support for the X-HTTP-Method-Override
// header. Submitted value will be checked against known methods. Adding
// HTTPMethodOverrideFormKey to any form will take precedence before
// HTTP header. If an unknown method will be submitted it gets logged as an
// Info log. This function is chainable.
// Suported options are: SetMethodOverrideFormKey() and SetLogger().
func WithXHTTPMethodOverride(opts ...Option) Middleware {
	ob := newOptionBox(opts...)
	errUnknownMethod := errors.NewNotValidf("[mw] Unknown http method")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mo := r.FormValue(ob.methodOverrideFormKey)
			if mo == "" {
				mo = r.Header.Get(MethodOverrideHeader)
			}
			switch mo {
			case "": // do nothing
			case "HEAD", "GET", "POST", "PUT", "PATCH", "DELETE", "TRACE", "OPTIONS":
				r.Method = mo
			default:
				// not sure if an error is here really needed ...
				if ob.log.IsDebug() {
					ob.log.Debug(
						"ctxhttp.SupportXHTTPMethodOverride.switch",
						log.Err(errUnknownMethod),
						log.String("method", mo),
						log.String("form", r.Form.Encode()),
						log.Object("header", r.Header))
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

// WithCloseNotify returns a ctxhttp.Handler cancelling the context when the client
// connection close unexpectedly.
// Supported options are: SetLogger().
func WithCloseNotify(opts ...Option) Middleware {
	ob := newOptionBox(opts...)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
						ob.log.Debug("ctxhttp.WithCloseNotify.cancel", log.Bool("cancelled", true), log.HTTPRequest("request", r))
					}
				}()
			}
			h.ServeHTTP(w, r)
		})
	}
}

// WithTimeout returns a ctxhttp.Handler which adds a timeout to the context.
//
// Child handlers have the responsibility to obey the context deadline and to return
// an appropriate error (or not) response in case of timeout.
func WithTimeout(timeout time.Duration) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := context.WithTimeout(r.Context(), timeout)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}
