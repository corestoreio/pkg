// Copyright 2015 CoreGroup Authors
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
	"runtime"
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
)

var testStorage = store.NewStorage(
	store.TableWebsiteSlice{
		&store.TableWebsite{WebsiteID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "euro", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Europe", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
		&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
	},
	store.TableGroupSlice{
		&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		&store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
	},
	store.TableStoreSlice{
		&store.TableStore{StoreID: 0, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		&store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "uk", Valid: true}}, WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
		&store.TableStore{StoreID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "at", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
		&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		&store.TableStore{StoreID: 3, Code: dbr.NullString{NullString: sql.NullString{String: "ch", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
	},
)

func TestStorageWebsite(t *testing.T) {

	tests := []struct {
		have      store.Retriever
		err       error
		wantWCode string
	}{
		{nil, store.ErrWebsiteNotFound, ""},
		{store.ID(2015), store.ErrWebsiteNotFound, ""},
		{store.ID(1), nil, "euro"},
		{store.Code("asia"), store.ErrWebsiteNotFound, ""},
		{store.Code("oz"), nil, "oz"},
		{mockIDCode{1, "oz"}, nil, "oz"},
		{mockIDCode{1, "ozzz"}, store.ErrWebsiteNotFound, ""},
	}
	for _, test := range tests {
		w, err := testStorage.Website(test.have)
		if test.err != nil {
			assert.Nil(t, w)
			assert.EqualError(t, test.err, err.Error())
		} else {
			assert.NotNil(t, w, "Test: %#v", test)
			assert.NoError(t, err, "Test: %#v", test)
			assert.Equal(t, test.wantWCode, w.Data().Code.String, "Test: %#v", test)
		}
	}

	w, err := testStorage.Website(store.Code("euro"))
	assert.NoError(t, err)
	assert.NotNil(t, w)

	dGroup, err := w.DefaultGroup()
	assert.NoError(t, err)
	assert.EqualValues(t, "DACH Group", dGroup.Data().Name)

	groups, err := w.Groups()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.Int64Slice{1, 2}, groups.IDs())

	stores, err := w.Stores()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.StringSlice{"de", "uk", "at", "ch"}, stores.Codes())
}

var benchmarkStorageWebsite *store.Website
var benchmarkStorageWebsiteDefaultGroup *store.Group

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageWebsiteGetDefaultGroup	  200000	      6081 ns/op	    1712 B/op	      45 allocs/op
func BenchmarkStorageWebsiteGetDefaultGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageWebsite, err = testStorage.Website(store.Code("euro"))
		if err != nil {
			b.Error(err)
		}

		benchmarkStorageWebsiteDefaultGroup, err = benchmarkStorageWebsite.DefaultGroup()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestStorageWebsites(t *testing.T) {
	websites, err := testStorage.Websites()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.StringSlice{"admin", "euro", "oz"}, websites.Codes())
	assert.EqualValues(t, utils.Int64Slice{0, 1, 2}, websites.IDs())

	var ids = []struct {
		g utils.Int64Slice
		s utils.Int64Slice
	}{
		{utils.Int64Slice{0}, utils.Int64Slice{0}},             //admin
		{utils.Int64Slice{1, 2}, utils.Int64Slice{1, 4, 2, 3}}, // dach
		{utils.Int64Slice{3}, utils.Int64Slice{5, 6}},          // oz
	}

	for i, w := range websites {
		groups, err := w.Groups()
		assert.NoError(t, err)
		assert.EqualValues(t, ids[i].g, groups.IDs())

		stores, err := w.Stores()
		assert.NoError(t, err)
		assert.EqualValues(t, ids[i].s, stores.IDs())
	}
}

func TestWebsiteSliceFilter(t *testing.T) {
	websites, err := testStorage.Websites()
	assert.NoError(t, err)
	assert.True(t, websites.Len() == 3)

	gs := websites.Filter(func(w *store.Website) bool {
		return w.Data().WebsiteID > 0
	})
	assert.EqualValues(t, utils.Int64Slice{1, 2}, gs.IDs())
}

func TestStorageGroup(t *testing.T) {

	tests := []struct {
		id       store.Retriever
		err      error
		wantName string
	}{
		{nil, store.ErrGroupNotFound, ""},
		{store.ID(2015), store.ErrGroupNotFound, ""},
		{store.ID(1), nil, "DACH Group"},
	}
	for _, test := range tests {
		g, err := testStorage.Group(test.id)
		if test.err != nil {
			assert.Nil(t, g)
			assert.EqualError(t, test.err, err.Error())
		} else {
			assert.NotNil(t, g)
			assert.NoError(t, err)
			assert.Equal(t, test.wantName, g.Data().Name)
		}
	}

	g, err := testStorage.Group(store.ID(3))
	assert.NoError(t, err)
	assert.NotNil(t, g)

	dStore, err := g.DefaultStore()
	assert.NoError(t, err)
	assert.EqualValues(t, "au", dStore.Data().Code.String)

	assert.EqualValues(t, "oz", g.Website().Data().Code.String)

	stores, err := g.Stores()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.StringSlice{"au", "nz"}, stores.Codes())
}

var benchmarkStorageGroup *store.Group
var benchmarkStorageGroupDefaultStore *store.Store

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageGroupGetDefaultStore	 1000000	      1916 ns/op	     464 B/op	      14 allocs/op
func BenchmarkStorageGroupGetDefaultStore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageGroup, err = testStorage.Group(store.ID(3))
		if err != nil {
			b.Error(err)
		}

		benchmarkStorageGroupDefaultStore, err = benchmarkStorageGroup.DefaultStore()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestStorageGroups(t *testing.T) {
	groups, err := testStorage.Groups()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.Int64Slice{3, 1, 0, 2}, groups.IDs())
	assert.True(t, groups.Len() == 4)

	var ids = []utils.Int64Slice{
		utils.Int64Slice{5, 6},    // oz
		utils.Int64Slice{1, 2, 3}, // dach
		utils.Int64Slice{0},       // default
		utils.Int64Slice{4},       // uk
	}

	for i, g := range groups {
		stores, err := g.Stores()
		assert.NoError(t, err)
		assert.EqualValues(t, ids[i], stores.IDs(), "Group %s ID %d", g.Data().Name, g.Data().GroupID)
	}
}

func TestGroupSliceFilter(t *testing.T) {
	groups, err := testStorage.Groups()
	assert.NoError(t, err)
	gs := groups.Filter(func(g *store.Group) bool {
		return g.Data().GroupID > 0
	})
	assert.EqualValues(t, utils.Int64Slice{3, 1, 2}, gs.IDs())
}

func TestStorageGroupNoWebsite(t *testing.T) {
	var tst = store.NewStorage(
		store.TableWebsiteSlice{
			&store.TableWebsite{WebsiteID: 21, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		},
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 5, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		},
	)
	g, err := tst.Group(store.ID(3))
	assert.Nil(t, g)
	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())

	gs, err := tst.Groups()
	assert.Nil(t, gs)
	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())
}

func TestStorageStore(t *testing.T) {

	tests := []struct {
		have     store.Retriever
		err      error
		wantCode string
	}{
		{nil, store.ErrStoreNotFound, ""},
		{store.ID(2015), store.ErrStoreNotFound, ""},
		{store.ID(1), nil, "de"},
		{store.Code("asia"), store.ErrStoreNotFound, ""},
		{store.Code("nz"), nil, "nz"},
		{mockIDCode{4, "nz"}, nil, "nz"},
		{mockIDCode{4, "auuuuu"}, store.ErrStoreNotFound, ""},
	}
	for _, test := range tests {
		s, err := testStorage.Store(test.have)
		if test.err != nil {
			assert.Nil(t, s, "%#v", test)
			assert.EqualError(t, test.err, err.Error())
		} else {
			assert.NotNil(t, s, "%#v", test)
			assert.NoError(t, err, "%#v", test)
			assert.Equal(t, test.wantCode, s.Data().Code.String)
		}
	}

	s, err := testStorage.Store(store.Code("at"))
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.EqualValues(t, "DACH Group", s.Group().Data().Name)

	website := s.Website()
	assert.EqualValues(t, "euro", website.Data().Code.String)
	wg, err := website.DefaultGroup()
	assert.NotNil(t, wg)
	assert.EqualValues(t, "DACH Group", wg.Data().Name)
	wgs, err := wg.DefaultStore()
	assert.NoError(t, err)
	assert.EqualValues(t, "at", wgs.Data().Code.String)
}

var benchmarkStorageStore *store.Store
var benchmarkStorageStoreWebsite *store.Website

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageStoreGetWebsite	 2000000	       656 ns/op	     176 B/op	       6 allocs/op
func BenchmarkStorageStoreGetWebsite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageStore, err = testStorage.Store(store.Code("de"))
		if err != nil {
			b.Error(err)
		}

		benchmarkStorageStoreWebsite = benchmarkStorageStore.Website()
		if benchmarkStorageStoreWebsite == nil {
			b.Error("benchmarkStorageStoreWebsite is nil")
		}
	}
}

func TestStorageStores(t *testing.T) {
	stores, err := testStorage.Stores()
	assert.NoError(t, err)
	assert.EqualValues(t, utils.StringSlice{"admin", "au", "de", "uk", "at", "nz", "ch"}, stores.Codes())
	assert.EqualValues(t, utils.Int64Slice{0, 5, 1, 4, 2, 6, 3}, stores.IDs())

	var ids = []struct {
		g string
		w string
	}{
		{"Default", "admin"},
		{"Australia", "oz"},
		{"DACH Group", "euro"},
		{"UK Group", "euro"},
		{"DACH Group", "euro"},
		{"Australia", "oz"},
		{"DACH Group", "euro"},
	}

	for i, s := range stores {
		assert.EqualValues(t, ids[i].g, s.Group().Data().Name)
		assert.EqualValues(t, ids[i].w, s.Website().Data().Code.String)
	}
}

func TestDefaultStoreView(t *testing.T) {
	st, err := testStorage.DefaultStoreView()
	assert.NoError(t, err)
	assert.EqualValues(t, "at", st.Data().Code.String)

	var tst = store.NewStorage(
		store.TableWebsiteSlice{
			&store.TableWebsite{WebsiteID: 21, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		},
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		},
	)
	dSt, err := tst.DefaultStoreView()
	assert.Nil(t, dSt)
	assert.EqualError(t, store.ErrStoreNotFound, err.Error())

	var tst2 = store.NewStorage(
		store.TableWebsiteSlice{
			&store.TableWebsite{WebsiteID: 21, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}},
		},
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 33, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		},
		store.TableStoreSlice{},
	)
	dSt2, err := tst2.DefaultStoreView()
	assert.Nil(t, dSt2)
	assert.EqualError(t, store.ErrGroupNotFound, err.Error())
}

func TestStorageStoreErrors(t *testing.T) {

	var nsw = store.NewStorage(
		store.TableWebsiteSlice{},
		store.TableGroupSlice{},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		},
	)
	stw, err := nsw.Store(store.Code("nz"))
	assert.Nil(t, stw)
	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())

	stws, err := nsw.Stores()
	assert.Nil(t, stws)
	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())

	var nsg = store.NewStorage(
		store.TableWebsiteSlice{
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		},
		store.TableGroupSlice{
			&store.TableGroup{GroupID: 13, WebsiteID: 12, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 4},
		},
		store.TableStoreSlice{
			&store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "au", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NullString{NullString: sql.NullString{String: "nz", Valid: true}}, WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		},
	)

	stg, err := nsg.Store(store.Code("nz"))
	assert.Nil(t, stg)
	assert.EqualError(t, store.ErrGroupNotFound, err.Error())

	stgs, err := nsg.Stores()
	assert.Nil(t, stgs)
	assert.EqualError(t, store.ErrGroupNotFound, err.Error())

}

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageDefaultStoreView	 2000000	       724 ns/op	     176 B/op	       7 allocs/op
func BenchmarkStorageDefaultStoreView(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageStore, err = testStorage.DefaultStoreView()
		if err != nil {
			b.Error(err)
		}

	}
}

func TestStorageReInit(t *testing.T) {
	numCPU := runtime.NumCPU()
	prevCPU := runtime.GOMAXPROCS(numCPU)
	t.Logf("GOMAXPROCS was: %d now: %d", prevCPU, numCPU)
	defer runtime.GOMAXPROCS(prevCPU)

	// quick implement, use mock of dbr.SessionRunner and remove connection
	db := csdb.MustConnectTest()
	defer db.Close()
	dbrConn := dbr.NewConnection(db, nil)

	nsg := store.NewStorage(nil, nil, nil)
	if err := nsg.ReInit(dbrConn.NewSession(nil)); err != nil {
		t.Error(err)
	}

	stores, err := nsg.Stores()
	assert.NoError(t, err)
	assert.True(t, stores.Len() > 0, "Expecting at least one store loaded from DB")
	for _, s := range stores {
		assert.NotEmpty(t, s.Data().Code.String, "Store: %#v", s.Data())
	}

	groups, err := nsg.Groups()
	assert.True(t, groups.Len() > 0, "Expecting at least one group loaded from DB")
	assert.NoError(t, err)
	for _, g := range groups {
		assert.NotEmpty(t, g.Data().Name, "Group: %#v", g.Data())
	}

	websites, err := nsg.Websites()
	assert.True(t, websites.Len() > 0, "Expecting at least one website loaded from DB")
	assert.NoError(t, err)
	for _, w := range websites {
		assert.NotEmpty(t, w.Data().Code.String, "Website: %#v", w.Data())
	}
}
