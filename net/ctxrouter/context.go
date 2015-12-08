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

package ctxrouter

import (
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type ctxKeyParams struct{}

// FromContextParams returns the Params slice from a context. It is guaranteed
// that the return value is non-nil.
func FromContextParams(ctx context.Context) Params {
	if p, ok := ctx.Value(ctxKeyParams{}).(Params); ok {
		return p
	}
	return Params{}
}

// WithContextParams puts Params into a context.
func WithContextParams(ctx context.Context, p Params) context.Context {
	return context.WithValue(ctx, ctxKeyParams{}, p)
}

type ctxKeyPanic struct{}

// FromContextPanic returns the value of a panic. You are responsible
// to extract the correct type from the interface{}.
func FromContextPanic(ctx context.Context) interface{} {
	return ctx.Value(ctxKeyPanic{})
}

// WithContextPanic puts a panic into a context.
func WithContextPanic(ctx context.Context, p interface{}) context.Context {
	return context.WithValue(ctx, ctxKeyPanic{}, p)
}

type ctxKeyError struct{}

// FromContextError returns the error interface. The usual error
// handling patterns apply. This function should be used within
// Router.ErrorHandler
func FromContextError(ctx context.Context) error {
	if err, ok := ctx.Value(ctxKeyError{}).(error); ok {
		return err
	}
	return nil
}

// withContextError puts an error into a context.
func withContextError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, ctxKeyError{}, err)
}

type webSocketKey struct{}

// FromContextWebsocket extracts a websocket connection from the context.
// Returns false even when the socket is nil. A true return value is guaranteed
// that the socket is not nil.
func FromContextWebsocket(ctx context.Context) (ws *websocket.Conn, ok bool) {
	ws, ok = ctx.Value(webSocketKey{}).(*websocket.Conn)
	if ok && ws == nil {
		ok = false
	}
	return
}

func withContextWebsocket(ctx context.Context, ws *websocket.Conn) context.Context {
	return context.WithValue(ctx, webSocketKey{}, ws)
}
