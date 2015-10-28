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

package ctxhttp

import (
	"net/http"

	"github.com/corestoreio/csfw/net/httputils"
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
	"time"
)

// Middleware is a wrapper for the ctxhttp.Handler to create middleware functions.
type Middleware func(Handler) Handler

// Chain function will iterate over all middleware, calling them one by one
// in a chained manner, returning the result of the final middleware.
func Chain(h Handler, mws ...Middleware) Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}

// WithHeader is a middleware that sets multiple HTTP headers. Will panic if kv
// is imbalanced. len(kv)%2 == 0.
func WithHeader(kv ...string) Middleware {
	lkv := len(kv)
	return func(h Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			for i := 0; i < lkv; i = i + 2 {
				w.Header().Set(kv[i], kv[i+1])
			}
			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// WithXHTTPMethodOverride adds support for the X-HTTP-Method-Override
// header. Submitted value will be checked against known methods. Adding
// HTTPMethodOverrideFormKey to any form will take precedence before
// HTTP header. If an unknown method will be submitted it gets logged as an
// Info log. This function is chainable.
func WithXHTTPMethodOverride() Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			mo := r.FormValue(httputils.MethodOverrideFormKey)
			if mo == "" {
				mo = r.Header.Get(httputils.MethodOverrideHeader)
			}
			switch mo {
			case "": // do nothing
			case httputils.MethodHead, httputils.MethodGet, httputils.MethodPost, httputils.MethodPut, httputils.MethodPatch, httputils.MethodDelete, httputils.MethodTrace, httputils.MethodOptions:
				r.Method = mo
			default:
				// not sure if an error is here really needed ...
				if log.IsInfo() {
					log.Info("ctxhttp.SupportXHTTPMethodOverride.switch", "err", "Unknown http method", "method", mo, "form", r.Form.Encode(), "header", r.Header)
				}
			}
			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// WithCloseNotify returns a Handler cancelling the context when the client
// connection close unexpectedly.
func WithCloseNotify() Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Cancel the context if the client closes the connection
			if wcn, ok := w.(http.CloseNotifier); ok {

				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()

				notify := wcn.CloseNotify()
				go func() {
					<-notify
					cancel()
				}()
			}
			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// WithTimeout returns a Handler which adds a timeout to the context.
//
// Child handlers have the responsibility to obey the context deadline and to return
// an appropriate error (or not) response in case of timeout.
func WithTimeout(timeout time.Duration) Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, _ = context.WithTimeout(ctx, timeout)
			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}
