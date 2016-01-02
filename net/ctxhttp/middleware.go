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

// Middleware is a wrapper for the function ctxhttp.HandlerFunc to create
// middleware functions.
type Middleware func(HandlerFunc) HandlerFunc

// MiddlewareSlice a slice full of middleware functions and with function
// receivers attached
type MiddlewareSlice []Middleware

// Chain will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func Chain(h HandlerFunc, mws ...Middleware) HandlerFunc {
	// Chain middleware with handler in the end
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Chain will iterate over all middleware functions, calling them one by one
// in a chained manner, returning the result of the final middleware.
func (mws MiddlewareSlice) Chain(h HandlerFunc) HandlerFunc {
	return Chain(h, mws...)
}
