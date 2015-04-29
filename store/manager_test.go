// Copyright 2015 CoreStore Authors
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

package store_test

import (
	"database/sql"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/stretchr/testify/assert"
)

func getTestManager(opts ...func(ms *mockStorage)) *store.Manager {
	ms := &mockStorage{}
	for _, opt := range opts {
		opt(ms)
	}
	return store.NewManager(ms)
}

var managerStoreSimpleTest = getTestManager(func(ms *mockStorage) {
	ms.s = func() (*store.Store, error) {
		return store.NewStore(
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		), nil
	}
})

func TestNewManagerStore(t *testing.T) {
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())
	for j := 0; j < 3; j++ {
		s, err := managerStoreSimpleTest.Store(nil, store.Code("notNil"))
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.EqualValues(t, "de", s.Data().Code.String)
	}
	assert.False(t, managerStoreSimpleTest.IsCacheEmpty())
	managerStoreSimpleTest.ClearCache()
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())

	tests := []struct {
		have    store.Retriever
		wantErr error
	}{
		{store.Code("nilSlices"), store.ErrStoreNotFound},
		{store.ID(2), store.ErrStoreNotFound},
		{nil, store.ErrCurrentStoreNotSet},
	}

	managerEmpty := getTestManager()
	for _, test := range tests {
		s, err := managerEmpty.Store(test.have)
		assert.Nil(t, s)
		assert.EqualError(t, test.wantErr, err.Error())
	}
	assert.True(t, managerStoreSimpleTest.IsCacheEmpty())

	managerDefaultStore := getTestManager(func(ms *mockStorage) {
		ms.dsv = func() (*store.Store, error) {
			return store.NewStore(
				&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
				&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
				&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			), nil
		}
	})

	// call it twice to test internal caching
	s, err := managerDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.Data().Code.String)

	s, err = managerDefaultStore.DefaultStoreView()
	assert.NotNil(t, s)
	assert.NoError(t, err)
	assert.NotEmpty(t, s.Data().Code.String)
	assert.False(t, managerDefaultStore.IsCacheEmpty())
	managerDefaultStore.ClearCache()
	assert.True(t, managerDefaultStore.IsCacheEmpty())
}

var benchmarkManagerStore *store.Store

// BenchmarkManagerGetStore	 5000000	       355 ns/op	      24 B/op	       2 allocs/op
func BenchmarkManagerGetStore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkManagerStore, err = managerStoreSimpleTest.Store(nil, store.Code("de"))
		if err != nil {
			b.Error(err)
		}
		if benchmarkManagerStore == nil {
			b.Error("benchmarkManagerStore is nil")
		}
	}
}

/*
	MOCKS
*/

type mockIDCode struct {
	id   int64
	code string
}

func (ic mockIDCode) ID() int64 {
	return ic.id
}
func (ic mockIDCode) Code() string {
	return ic.code
}

type mockStorage struct {
	w   func() (*store.Website, error)
	ws  func() (store.WebsiteSlice, error)
	g   func() (*store.Group, error)
	gs  func() (store.GroupSlice, error)
	s   func() (*store.Store, error)
	dsv func() (*store.Store, error)
	ss  func() (store.StoreSlice, error)
}

var _ store.Storager = (*mockStorage)(nil)

func (ms *mockStorage) Website(_ store.Retriever) (*store.Website, error) {
	if ms.w == nil {
		return nil, store.ErrWebsiteNotFound
	}
	return ms.w()
}
func (ms *mockStorage) Websites() (store.WebsiteSlice, error) {
	if ms.ws == nil {
		return nil, nil
	}
	return ms.ws()
}
func (ms *mockStorage) Group(_ store.Retriever) (*store.Group, error) {
	if ms.g == nil {
		return nil, store.ErrGroupNotFound
	}
	return ms.g()
}
func (ms *mockStorage) Groups() (store.GroupSlice, error) {
	if ms.gs == nil {
		return nil, nil
	}
	return ms.gs()
}
func (ms *mockStorage) Store(_ store.Retriever) (*store.Store, error) {
	if ms.s == nil {
		return nil, store.ErrStoreNotFound
	}
	return ms.s()
}
func (ms *mockStorage) Stores() (store.StoreSlice, error) {
	if ms.ss == nil {
		return nil, nil
	}
	return ms.ss()
}
func (ms *mockStorage) DefaultStoreView() (*store.Store, error) {
	if ms.dsv == nil {
		return nil, store.ErrStoreNotFound
	}
	return ms.dsv()
}
