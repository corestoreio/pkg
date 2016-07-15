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

package scope

import (
	"context"
)

type ctxScopeKey struct{}

type ctxScopeWrapper struct {
	c Hash // current
	p Hash // parent
}

// FromContext returns the requested current scope and its parent from a
// context. This scope is only valid for the current request. A scope gets set
// via HTTP form, cookie, JSON Web Token or GeoIP or other fancy features. The
// returned bool checks also if the current and parent Hash are valid in their
// hierarchical relation.
func FromContext(ctx context.Context) (current Hash, parent Hash, ok bool) {
	w, ok := ctx.Value(ctxScopeKey{}).(ctxScopeWrapper)
	if !ok {
		return 0, 0, false
	}
	return w.c, w.p, w.c.ValidParent(w.p)
}

// WithContext adds the requestedStore to the context.
// This function must only be used one time to set the requested store for one
// request. Usually a store gets initialized by the store.NewService() init,
// JSON web token middleware or cookie and form based middleware.
// Only one error can be passed, multiple errors get ignored from the 2nd position.
func WithContext(ctx context.Context, current Hash, parent Hash) context.Context {
	return context.WithValue(ctx, ctxScopeKey{}, ctxScopeWrapper{c: current, p: parent})
}
