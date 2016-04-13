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

package ctxjwt

import (
	"github.com/corestoreio/csfw/util/csjwt"
	"golang.org/x/net/context"
)

type (
	keyCtxToken struct{}
	keyCtxErr   struct{}
)

// withContext creates a new context with csjwt.Token attached.
func withContext(ctx context.Context, t csjwt.Token) context.Context {
	return context.WithValue(ctx, keyCtxToken{}, t)
}

// FromContext returns the csjwt.Token in ctx if it exists or an error.
// If there is no token in the context then the error
// ErrContextJWTNotFound gets returned.
func FromContext(ctx context.Context) (t csjwt.Token, err error) {
	var ok bool
	err, ok = ctx.Value(keyCtxErr{}).(error)
	if ok && err != nil {
		return
	}
	err = ErrContextJWTNotFound
	t, ok = ctx.Value(keyCtxToken{}).(csjwt.Token)
	if ok && len(t.Raw) > 5 && t.Valid {
		err = nil
	}
	return
}

// withContextError creates a new context with an error attached.
func withContextError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, keyCtxErr{}, err)
}
