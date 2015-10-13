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

	"fmt"

	"github.com/corestoreio/csfw/net/httputils"
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

var _ http.Handler = (*Adapter)(nil)

// AdapterErrorFunc specifies the error handler for the Adapter
type AdapterErrFunc func(http.ResponseWriter, *http.Request, error)

// Adapter type allows to use existing http.Handler middleware, as
// long as they run before it does.
type Adapter struct {
	Ctx       context.Context // Root Context
	Handler   Handler         // will be called
	ErrorFunc AdapterErrFunc  // gets called when Handler returns an error
}

// DefaultAdapterErrFunc logs the error (if Debug is enabled) and sends a
// 400 StatusBadRequest. You can replace this variable with your own default
// version. This function gets called in ServeHTTP of the Adapter type.
var DefaultAdapterErrFunc AdapterErrFunc = defaultAdapterErrFunc

// ServeHTTP calls ServeHTTPContext(ca.ctx, rw, req) and on error calls the
// ErrorFunc.
func (ca *Adapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if err := ca.Handler.ServeHTTPContext(ca.Ctx, rw, req); err != nil {
		ca.ErrorFunc(rw, req, err)
	}
}

// NewAdapter creates a new Adapter with the default AdapterErrorFunc.
func NewAdapter(ctx context.Context, h Handler) *Adapter {
	return &Adapter{
		Ctx:       ctx,
		Handler:   h,
		ErrorFunc: DefaultAdapterErrFunc,
	}
}

func defaultAdapterErrFunc(rw http.ResponseWriter, req *http.Request, err error) {
	if log.IsDebug() {
		log.Error("ctxhttp.AdapterErrorFunc", "err", err, "req", req, "url", req.URL)
	}
	code := http.StatusBadRequest
	http.Error(rw, fmt.Sprintf(
		"%s\nApp Error: %s",
		http.StatusText(code),
		err,
	), code)
}

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
