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

package ctxhttp

import (
	"net/http"

	"fmt"

	"golang.org/x/net/context"
)

// More information about context.Context https://joeshaw.org/net-context-and-http-handler/

// TODO(CS) with Go 1.7 the net/context package gets merged into stdlib within end of Q1/2016
// https://groups.google.com/forum/#!topic/golang-dev/cQs1z9LrJDU

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

// Adapter type allows to use existing http.Handler middleware, as
// long as they run before it does.
type Adapter struct {
	Ctx       context.Context // Root Context
	Handler   Handler         // will be called
	ErrorFunc AdapterErrFunc  // gets called when Handler returns an error
}

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

// AdapterErrFunc specifies the error handler for the Adapter
type AdapterErrFunc func(http.ResponseWriter, *http.Request, error)

// DefaultAdapterErrFunc sends a 400 StatusBadRequest. You can replace this
// variable with your own default version. This function gets called in
// ServeHTTP of the Adapter type.
var DefaultAdapterErrFunc AdapterErrFunc = defaultAdapterErrFunc

func defaultAdapterErrFunc(rw http.ResponseWriter, req *http.Request, err error) {
	code := http.StatusBadRequest
	http.Error(rw, fmt.Sprintf(
		"%s\nApp Error: %s",
		http.StatusText(code),
		err,
	), code)
}
