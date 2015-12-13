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

package ctxmw

import (
	"net/http"
	"time"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/httputil"
	"golang.org/x/net/context"
)

// WithHeader is a middleware that sets multiple HTTP headers. Will panic if kv
// is imbalanced. len(kv)%2 == 0.
func WithHeader(kv ...string) ctxhttp.Middleware {
	lkv := len(kv)
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			for i := 0; i < lkv; i = i + 2 {
				w.Header().Set(kv[i], kv[i+1])
			}
			return hf(ctx, w, r)
		}
	}
}

// WithXHTTPMethodOverride adds support for the X-HTTP-Method-Override
// header. Submitted value will be checked against known methods. Adding
// HTTPMethodOverrideFormKey to any form will take precedence before
// HTTP header. If an unknown method will be submitted it gets logged as an
// Info log. This function is chainable.
func WithXHTTPMethodOverride() ctxhttp.Middleware {
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
				if PkgLog.IsDebug() {
					PkgLog.Debug("ctxhttp.SupportXHTTPMethodOverride.switch", "err", "Unknown http method", "method", mo, "form", r.Form.Encode(), "header", r.Header)
				}
			}
			return hf(ctx, w, r)
		}
	}
}

// WithCloseNotify returns a ctxhttp.Handler cancelling the context when the client
// connection close unexpectedly.
func WithCloseNotify() ctxhttp.Middleware {
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
			return hf(ctx, w, r)
		}
	}
}

// WithTimeout returns a ctxhttp.Handler which adds a timeout to the context.
//
// Child handlers have the responsibility to obey the context deadline and to return
// an appropriate error (or not) response in case of timeout.
func WithTimeout(timeout time.Duration) ctxhttp.Middleware {
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, _ = context.WithTimeout(ctx, timeout)
			return hf(ctx, w, r)
		}
	}
}
