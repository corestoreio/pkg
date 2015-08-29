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

import "net/http"

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
