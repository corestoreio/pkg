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

package storemock

import (
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
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

// WithContextMustService creates a new StoreService wrapped in a context.Background().
// Panics on error.
func WithContextMustService(so scope.Option, opts ...func(ms *Storage)) context.Context {
	sm, err := NewService(so, opts...)
	if err != nil {
		panic(err)
	}
	return store.WithContextProvider(context.Background(), sm)
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
		return nil, errors.NewNotFoundf("[storemock] Website is nil")
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
		return nil, errors.NewNotFoundf("[storemock] Group is nil")
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
		return nil, errors.NewNotFoundf("[storemock] Store is nil")
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
	return nil, errors.NewNotFoundf("[storemock] Store")
}
func (ms *Storage) ReInit(dbr.SessionRunner, ...dbr.SelectCb) error {
	return nil
}

// NewEurozzyService creates a fully initialized store.Service with 3 websites,
// 4 groups and 7 stores used for testing. Panics on error.
// Website 1 contains Europe and website 2 contains Australia/New Zealand.
func NewEurozzyService(so scope.Option, opts ...store.StorageOption) *store.Service {
	// Yes weird naming, but feel free to provide a better name 8-)

	defaultOpts := []store.StorageOption{
		store.SetStorageWebsites(
			&store.TableWebsite{WebsiteID: 0, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
		),
		store.SetStorageGroups(
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
		),
		store.SetStorageStores(
			&store.TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
			&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 4, Code: dbr.NewNullString("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 2, Code: dbr.NewNullString("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
			&store.TableStore{IsActive: false, StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30},
		),
	}

	return store.MustNewService(so, store.MustNewStorage(append(defaultOpts, opts...)...))
}
