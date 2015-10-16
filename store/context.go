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

package store

import "golang.org/x/net/context"

// ctxKey type is unexported to prevent collisions with context keys defined in
// other packages.
type ctxKey uint

// Key* defines the keys to access a value in a context.Context
const (
	ctxKeyManagerReader ctxKey = iota
)

// ContextMustManagerReader returns a store.ManagerReader from a context.
func FromContextManagerReader(ctx context.Context) (r ManagerReader, ok bool) {
	r, ok = ctx.Value(ctxKeyManagerReader).(ManagerReader)
	return
}

// NewContextManagerReader adds a ManagerReader to the context.
func NewContextManagerReader(ctx context.Context, r ManagerReader) context.Context {
	return context.WithValue(ctx, ctxKeyManagerReader, r)
}
