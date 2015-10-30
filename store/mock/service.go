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

package mock

import (
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"golang.org/x/net/context"
)

// NewService creates a new StoreService
func NewService(so scope.Option, opts ...func(ms *Storage)) (*store.Service, error) {
	ms := &Storage{}
	for _, opt := range opts {
		opt(ms)
	}
	return store.NewService(so, ms)
}

// MustNewService creates a new StoreService but panics on error
func MustNewService(so scope.Option, opts ...func(ms *Storage)) *store.Service {
	ms := &Storage{}
	for _, opt := range opts {
		opt(ms)
	}
	return store.MustNewService(so, ms)
}

// NewContextService creates a new StoreService wrapped in a context.Context
func NewContextService(so scope.Option, opts ...func(ms *Storage)) context.Context {
	var sm *store.Service
	{
		var err error
		if sm, err = NewService(so, opts...); err != nil {
			panic(err)
		}
	}
	return store.NewContextReader(context.Background(), sm, nil)
}

// Storage main underlying data container
type Storage struct {
	MockWebsite      func() (*store.Website, error)
	MockWebsiteSlice func() (store.WebsiteSlice, error)
	MockGroup        func() (*store.Group, error)
	MockGroupSlice   func() (store.GroupSlice, error)
	MockStore        func() (*store.Store, error)
	MockDefaultStore func() (*store.Store, error)
	MockStoreSlice   func() (store.StoreSlice, error)
}

var _ store.Storager = (*Storage)(nil)

func (ms *Storage) Website(_ scope.WebsiteIDer) (*store.Website, error) {
	if ms.MockWebsite == nil {
		return nil, store.ErrWebsiteNotFound
	}
	return ms.MockWebsite()
}
func (ms *Storage) Websites() (store.WebsiteSlice, error) {
	if ms.MockWebsiteSlice == nil {
		return nil, nil
	}
	return ms.MockWebsiteSlice()
}
func (ms *Storage) Group(_ scope.GroupIDer) (*store.Group, error) {
	if ms.MockGroup == nil {
		return nil, store.ErrGroupNotFound
	}
	return ms.MockGroup()
}
func (ms *Storage) Groups() (store.GroupSlice, error) {
	if ms.MockGroupSlice == nil {
		return nil, nil
	}
	return ms.MockGroupSlice()
}
func (ms *Storage) Store(_ scope.StoreIDer) (*store.Store, error) {
	if ms.MockStore == nil {
		return nil, store.ErrStoreNotFound
	}
	return ms.MockStore()
}

func (ms *Storage) Stores() (store.StoreSlice, error) {
	if ms.MockStoreSlice == nil {
		return nil, nil
	}
	return ms.MockStoreSlice()
}
func (ms *Storage) DefaultStoreView() (*store.Store, error) {
	if ms.MockDefaultStore != nil {
		return ms.MockDefaultStore()
	}
	if ms.MockStore != nil {
		return ms.MockStore()
	}
	return nil, store.ErrStoreNotFound
}
func (ms *Storage) ReInit(dbr.SessionRunner, ...csdb.DbrSelectCb) error {
	return nil
}
