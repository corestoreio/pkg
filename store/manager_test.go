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

//Storager interface {
//Website(id IDRetriever, c CodeRetriever) (*Website, error)
//Websites() (WebsiteSlice, error)
//Group(id IDRetriever) (*Group, error)
//Groups() (GroupSlice, error)
//Store(id IDRetriever, c CodeRetriever) (*Store, error)
//Stores() (StoreSlice, error)
//DefaultStoreView() (*Store, error)
//}

type mockStorage struct {
	w *store.Website
	g *store.Group
	s *store.Store
}

var _ store.Storager = (*mockStorage)(nil)

func (ts *mockStorage) Website(id store.IDRetriever, c store.CodeRetriever) (*store.Website, error) {
	return ts.w, nil
}
func (ts *mockStorage) Websites() (store.WebsiteSlice, error) { return nil, nil }
func (ts *mockStorage) Group(id store.IDRetriever) (*store.Group, error) {
	return ts.g, nil
}
func (ts *mockStorage) Groups() (store.GroupSlice, error) { return nil, nil }
func (ts *mockStorage) Store(id store.IDRetriever, c store.CodeRetriever) (*store.Store, error) {
	return ts.s, nil
}
func (ts *mockStorage) Stores() (store.StoreSlice, error)       { return nil, nil }
func (ts *mockStorage) DefaultStoreView() (*store.Store, error) { return nil, nil }

func TestNewManager(t *testing.T) {
	ms := &mockStorage{}
	mgnr := store.NewManager(ms)
	ms.s = store.NewStore(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
	)

	s, err := mgnr.Store(nil, store.Code("de"))
	assert.NoError(t, err)
	assert.NotNil(t, s)
	t.Logf("\nLOG: %#v\n", s)
}

var benchmarkManagerStore *store.Store

// BenchmarkManagerGetStore	 5000000	       355 ns/op	      24 B/op	       2 allocs/op
func BenchmarkManagerGetStore(b *testing.B) {
	mngr := store.NewManager(&mockStorage{
		s: store.NewStore(
			&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
			&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		),
	})

	for i := 0; i < b.N; i++ {
		var err error
		benchmarkManagerStore, err = mngr.Store(nil, store.Code("de"))
		if err != nil {
			b.Error(err)
		}
		if benchmarkManagerStore == nil {
			b.Error("benchmarkManagerStore is nil")
		}
	}
}
