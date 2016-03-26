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
	"errors"

	"github.com/corestoreio/csfw/store/scope"
	"golang.org/x/net/context"
)

// ErrContextProviderNotFound gets returned when store.Provider cannot be found
// in context.Context
var ErrContextProviderNotFound = errors.New("store.Provider not found in context.Context")

type ctxServiceKey struct{}
type ctxServiceWrapper struct {
	service        Provider
	requestedStore *Store
}

// FromContextProvider returns a store.Provider and a store.Store from a context.
// The *store.Store is either the current requested store (via JWT or cookie or REQUEST
// parameter) or if those are not set then the default initialized store when
// instantiating a new Getter. The returned store.Store identifies the current
// scope.Scope of a request. If it cannot determine a store.Store then the
// error ErrStoreNotFound will get returned.
func FromContextProvider(ctx context.Context) (Provider, *Store, error) {
	sw, ok := ctx.Value(ctxServiceKey{}).(ctxServiceWrapper)
	if !ok || sw.service == nil {
		return nil, nil, ErrContextProviderNotFound
	}

	if sw.requestedStore == nil {
		var err error
		sw.requestedStore, err = sw.service.Store()
		if err != nil {
			return nil, nil, err
		}
	}
	return sw.service, sw.requestedStore, nil
}

// WithContextProvider adds a store.Provider and an optional requestedStore to the context.
// requestedStore can be provided 0 or 1 time. If you provide the RequestedStore
// argument then it will override the default RequestedStore from FromContextReader()
func WithContextProvider(ctx context.Context, r Provider, requestedStore ...*Store) context.Context {
	var rs *Store
	if len(requestedStore) == 1 {
		rs = requestedStore[0]
	}
	return context.WithValue(ctx, ctxServiceKey{}, ctxServiceWrapper{
		service:        r,
		requestedStore: rs,
	})
}

// WithContextMustService creates a new StoreService wrapped in a context.Background().
// Convenience function. Panics on error.
func WithContextMustService(so scope.Option, s Storager) context.Context {
	sm, err := NewService(so, s)
	if err != nil {
		panic(err)
	}
	return WithContextProvider(context.Background(), sm)
}
