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

package store

import (
	"context"

	"github.com/corestoreio/csfw/util/errors"
)

var errContextProviderNotFound = errors.NewNotFoundf("[store] Requested Store not found in context.Context")

type ctxRequestedStoreKey struct{}
type ctxRequestedStoreWrapper struct {
	s   *Store
	err error
}

// FromContextRequestedStore returns the requested store.Store from a context
// valid for the current request scope.
// The *store.Store represents the current requested store (via JWT or cookie or REQUEST
// parameter). The returned store.Store identifies the current
// scope.Scope of a request. If it cannot determine a store.Store then the
// not found error will get returned.
func FromContextRequestedStore(ctx context.Context) (*Store, error) {
	st, ok := ctx.Value(ctxRequestedStoreKey{}).(ctxRequestedStoreWrapper)
	if !ok || (st.s == nil && st.err == nil) {
		return nil, errContextProviderNotFound
	}
	return st.s, st.err
}

// WithContextRequestedStore adds the requestedStore to the context.
// This function must only be used one time to set the requested store for one
// request. Usually a store gets initialized by the store.NewService() init,
// JSON web token middleware or cookie and form based middleware.
func WithContextRequestedStore(ctx context.Context, requestedStore *Store, err error) context.Context {
	return context.WithValue(ctx, ctxRequestedStoreKey{}, ctxRequestedStoreWrapper{requestedStore, err})
}
