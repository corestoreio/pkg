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
	websiteID int64
	storeID   int64
}

// WithContext adds the store scope with its parent website scope to the
// context. Different middlewares may call this function to set a new scope
// depending on different conditions. For example the JSON web token middleware
// can set a scope because the JWT contains a new store code. Or a geoip
// middleware can set the scope depending on geo location information. These IDs
// will be later used to e.g. read the scoped configuration.
func WithContext(ctx context.Context, websiteID, storeID int64) context.Context {
	return context.WithValue(ctx, ctxScopeKey{}, ctxScopeWrapper{websiteID: websiteID, storeID: storeID})
}

// FromContext returns the requested current store scope and its parent website
// scope from a context. This scope is only valid for the current context in a
// request. A scope gets set via HTTP form, cookie, JSON Web Token or GeoIP or
// other fancy features.
func FromContext(ctx context.Context) (websiteID, storeID int64, ok bool) {
	w, ok := ctx.Value(ctxScopeKey{}).(ctxScopeWrapper)
	return w.websiteID, w.storeID, ok && w.websiteID >= 0 && w.storeID >= 0
}
