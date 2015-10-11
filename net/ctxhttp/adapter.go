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

	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
)

// More information about context.Context https://joeshaw.org/net-context-and-http-handler/

// Handler allows http Handlers to include a context.
type Handler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request) error
}

// HandlerFunc defines a function that implements the Handler
// interface including the context.
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// ServeHTTPContext calls the ContextHandlerFunc with the given context,
// ResponseWrite and Request.
func (h HandlerFunc) ServeHTTPContext(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	return h(ctx, rw, req)
}

var _ http.Handler = (*HandlerStdLib)(nil)

// HandlerStdLib type allows to use existing http.Handler middleware, as
// long as they run before it does.
type HandlerStdLib struct {
	Ctx     context.Context
	Handler Handler
}

// ServeHTTP calls ServeHTTPContext(ca.ctx, rw, req).
func (ca *HandlerStdLib) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ca.Handler.ServeHTTPContext(ca.Ctx, rw, req)
}

// Adapter is a wrapper for the ctxhttp.Handler
type Adapter func(Handler) Handler

// Adapt function will iterate over all adapters, calling them one by one
// in a chained manner, returning the result of the final adapter.
func Adapt(h Handler, adapters ...Adapter) Handler {
	for _, a := range adapters {
		h = a(h)
	}
	return h
}

// WithHeader is an Adapter that sets multiple HTTP headers. Will panic if kv
// is imbalanced. len(kv)%2 == 0
func WithHeader(kv ...string) Adapter {
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
// Info log.
func WithXHTTPMethodOverride() Adapter {
	return func(h Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
					log.Info("ctxhttp.SupportXHTTPMethodOverride.switch", "err", "Unknown http method", "method", mo, "form", r.Form.Encode(), "header", r.Header)
				}
			}

			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}
