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

package store_test

import (
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

// todo inspect the high allocs

var testStorage = store.MustNewStorage(
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
		&store.TableStore{StoreID: 3, Code: dbr.NewNullString("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
	),
)

func TestStorageWebsite(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have       scope.WebsiteIDer
		wantErrBhf errors.BehaviourFunc
		wantWCode  string
	}{
		{nil, errors.IsNotFound, ""},
		{scope.MockID(2015), errors.IsNotFound, ""},
		{scope.MockID(1), nil, "euro"},
		{scope.MockCode("asia"), errors.IsNotFound, ""},
		{scope.MockCode("oz"), nil, "oz"},
		{mockIDCode{1, "oz"}, nil, "oz"},
		{mockIDCode{1, "ozzz"}, errors.IsNotFound, ""},
	}
	for i, test := range tests {
		w, err := testStorage.Website(test.have)
		if test.wantErrBhf != nil {
			assert.Nil(t, w)
			assert.True(t, test.wantErrBhf(err), "Index %d Error: %s", i, err)
		} else {
			assert.NotNil(t, w, "Index %d", i)
			assert.NoError(t, err, "Index %d", i)
			assert.Equal(t, test.wantWCode, w.Data.Code.String, "Index %d", i)
		}
	}

	w, err := testStorage.Website(scope.MockCode("euro"))
	assert.NoError(t, err)
	assert.NotNil(t, w)

	dGroup, err := w.DefaultGroup()
	assert.NoError(t, err)
	assert.EqualValues(t, "DACH Group", dGroup.Data.Name)

	assert.NotNil(t, w.Groups)
	assert.EqualValues(t, util.Int64Slice{1, 2}, w.Groups.IDs())

	assert.NotNil(t, w.Stores)
	assert.EqualValues(t, util.StringSlice{"de", "uk", "at", "ch"}, w.Stores.Codes())
}

var benchmarkStorageWebsite *store.Website
var benchmarkStorageWebsiteDefaultGroup *store.Group

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageWebsiteGetDefaultGroup	  200000	      6081 ns/op	    1712 B/op	      45 allocs/op
// BenchmarkStorageWebsiteGetDefaultGroup-4	   50000	     26210 ns/op	   10608 B/op	     229 allocs/op
func BenchmarkStorageWebsiteGetDefaultGroup(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageWebsite, err = testStorage.Website(scope.MockCode("euro"))
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
	t.Parallel()
	websites, err := testStorage.Websites()
	assert.NoError(t, err)
	assert.EqualValues(t, util.StringSlice{"admin", "euro", "oz"}, websites.Codes())
	assert.EqualValues(t, util.Int64Slice{0, 1, 2}, websites.IDs())

	var ids = []struct {
		g util.Int64Slice
		s util.Int64Slice
	}{
		{util.Int64Slice{0}, util.Int64Slice{0}},             //admin
		{util.Int64Slice{1, 2}, util.Int64Slice{1, 4, 2, 3}}, // dach
		{util.Int64Slice{3}, util.Int64Slice{5, 6}},          // oz
	}

	for i, w := range websites {
		assert.NotNil(t, w.Groups)
		assert.EqualValues(t, ids[i].g, w.Groups.IDs())

		assert.NotNil(t, w.Stores)
		assert.EqualValues(t, ids[i].s, w.Stores.IDs())
	}
}

func TestWebsiteSliceFilter(t *testing.T) {
	t.Parallel()
	websites, err := testStorage.Websites()
	assert.NoError(t, err)
	assert.True(t, websites.Len() == 3)

	gs := websites.Filter(func(w *store.Website) bool {
		return w.Data.WebsiteID > 0
	})
	assert.EqualValues(t, util.Int64Slice{1, 2}, gs.IDs())
}

func TestStorageGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id         scope.GroupIDer
		wantErrBhf errors.BehaviourFunc
		wantName   string
	}{
		{nil, errors.IsNotFound, ""},
		{scope.MockID(2015), errors.IsNotFound, ""},
		{scope.MockID(1), nil, "DACH Group"},
	}
	for i, test := range tests {
		g, err := testStorage.Group(test.id)
		if test.wantErrBhf != nil {
			assert.Nil(t, g)
			assert.True(t, test.wantErrBhf(err), "Index %d Error: %s", i, err)
		} else {
			assert.NotNil(t, g, "Index %d", i)
			assert.NoError(t, err, "Index %d", i)
			assert.Equal(t, test.wantName, g.Data.Name, "Index %d", i)
		}
	}

	g, err := testStorage.Group(scope.MockID(3))
	assert.NoError(t, err)
	assert.NotNil(t, g)

	dStore, err := g.DefaultStore()
	assert.NoError(t, err)
	assert.EqualValues(t, "au", dStore.Data.Code.String)

	assert.EqualValues(t, "oz", g.Website.Data.Code.String)

	assert.NotNil(t, g.Stores)
	assert.EqualValues(t, util.StringSlice{"au", "nz"}, g.Stores.Codes())
}

var benchmarkStorageGroup *store.Group
var benchmarkStorageGroupDefaultStore *store.Store

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageGroupGetDefaultStore	 1000000	      1916 ns/op	     464 B/op	      14 allocs/op
// BenchmarkStorageGroupGetDefaultStore-4  	  300000	      5387 ns/op	    2880 B/op	      64 allocs/op
func BenchmarkStorageGroupGetDefaultStore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageGroup, err = testStorage.Group(scope.MockID(3))
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
	t.Parallel()
	groups, err := testStorage.Groups()
	assert.NoError(t, err)
	assert.EqualValues(t, util.Int64Slice{3, 1, 0, 2}, groups.IDs())
	assert.True(t, groups.Len() == 4)

	var ids = []util.Int64Slice{
		{5, 6},    // oz
		{1, 2, 3}, // dach
		{0},       // default
		{4},       // uk
	}

	for i, g := range groups {
		assert.NotNil(t, g.Stores)
		assert.EqualValues(t, ids[i], g.Stores.IDs(), "Group %s ID %d", g.Data.Name, g.Data.GroupID)
	}
}

func TestGroupSliceFilter(t *testing.T) {
	t.Parallel()
	groups, err := testStorage.Groups()
	assert.NoError(t, err)
	gs := groups.Filter(func(g *store.Group) bool {
		return g.Data.GroupID > 0
	})
	assert.EqualValues(t, util.Int64Slice{3, 1, 2}, gs.IDs())
}

func TestStorageGroupNoWebsite(t *testing.T) {
	t.Parallel()
	var tst = store.MustNewStorage(
		store.SetStorageWebsites(
			&store.TableWebsite{WebsiteID: 21, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
		),
		store.SetStorageGroups(
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
		store.SetStorageStores(
			&store.TableStore{StoreID: 5, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)
	g, err := tst.Group(scope.MockID(3))
	assert.Nil(t, g)
	assert.True(t, errors.IsNotFound(err), err.Error())

	gs, err := tst.Groups()
	assert.Nil(t, gs)
	assert.True(t, errors.IsNotFound(err), err.Error())
}

func TestStorageStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		have       scope.StoreIDer
		wantErrBhf errors.BehaviourFunc
		wantCode   string
	}{
		{nil, errors.IsNotFound, ""},
		{scope.MockID(2015), errors.IsNotFound, ""},
		{scope.MockID(1), nil, "de"},
		{scope.MockCode("asia"), errors.IsNotFound, ""},
		{scope.MockCode("nz"), nil, "nz"},
		{mockIDCode{4, "nz"}, nil, "nz"},
		{mockIDCode{4, "auuuuu"}, errors.IsNotFound, ""},
	}
	for i, test := range tests {
		s, err := testStorage.Store(test.have)
		if test.wantErrBhf != nil {
			assert.Nil(t, s, "%#v", test)
			assert.True(t, test.wantErrBhf(err), "Index: %d Error: %s", i, err)
		} else {
			assert.NotNil(t, s, "Index %d", i)
			assert.NoError(t, err, "Index %d", i)
			assert.Equal(t, test.wantCode, s.Data.Code.String, "Index %d", i)
		}
	}

	s, err := testStorage.Store(scope.MockCode("at"))
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.EqualValues(t, "DACH Group", s.Group.Data.Name)

	assert.EqualValues(t, "euro", s.Website.Data.Code.String)
	wg, err := s.Website.DefaultGroup()
	assert.NotNil(t, wg)
	assert.EqualValues(t, "DACH Group", wg.Data.Name)
	wgs, err := wg.DefaultStore()
	assert.NoError(t, err)
	assert.EqualValues(t, "at", wgs.Data.Code.String)
}

var benchmarkStorageStore *store.Store
var benchmarkStorageStoreWebsite *store.Website

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageStoreGetWebsite	 2000000	       656 ns/op	     176 B/op	       6 allocs/op
// BenchmarkStorageStoreGetWebsite-4       	   50000	     32968 ns/op	   15280 B/op	     334 allocs/op
func BenchmarkStorageStoreGetWebsite(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageStore, err = testStorage.Store(scope.MockCode("de"))
		if err != nil {
			b.Error(err)
		}

		benchmarkStorageStoreWebsite = benchmarkStorageStore.Website
		if benchmarkStorageStoreWebsite == nil {
			b.Error("benchmarkStorageStoreWebsite is nil")
		}
	}
}

func TestStorageStores(t *testing.T) {
	t.Parallel()
	stores, err := testStorage.Stores()
	assert.NoError(t, err)
	assert.EqualValues(t, util.StringSlice{"admin", "au", "de", "uk", "at", "nz", "ch"}, stores.Codes())
	assert.EqualValues(t, util.Int64Slice{0, 5, 1, 4, 2, 6, 3}, stores.IDs())

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
		assert.EqualValues(t, ids[i].g, s.Group.Data.Name)
		assert.EqualValues(t, ids[i].w, s.Website.Data.Code.String)
	}
}

func TestDefaultStoreView(t *testing.T) {
	t.Parallel()
	st, err := testStorage.DefaultStoreView()
	assert.NoError(t, err)
	assert.EqualValues(t, "at", st.Data.Code.String)

	var tst = store.MustNewStorage(
		store.SetStorageWebsites(
			&store.TableWebsite{WebsiteID: 21, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
		),
		store.SetStorageGroups(
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
		store.SetStorageStores(
			&store.TableStore{StoreID: 4, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)
	dSt, err := tst.DefaultStoreView()
	assert.Nil(t, dSt)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)

	var tst2 = store.MustNewStorage(
		store.SetStorageWebsites(
			&store.TableWebsite{WebsiteID: 21, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(true)},
		),
		store.SetStorageGroups(
			&store.TableGroup{GroupID: 33, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
		store.SetStorageStores(),
	)
	dSt2, err := tst2.DefaultStoreView()
	assert.Nil(t, dSt2)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

func TestStorageStoreErrors(t *testing.T) {
	t.Parallel()

	var nsw = store.MustNewStorage(
		store.SetStorageWebsites(),
		store.SetStorageGroups(),
		store.SetStorageStores(
			&store.TableStore{StoreID: 4, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)
	stw, err := nsw.Store(scope.MockCode("nz"))
	assert.Nil(t, stw)
	assert.True(t, errors.IsNotFound(err), err.Error())

	stws, err := nsw.Stores()
	assert.Nil(t, stws)
	assert.True(t, errors.IsNotFound(err), err.Error())

	var nsg = store.MustNewStorage(
		store.SetStorageWebsites(
			&store.TableWebsite{WebsiteID: 2, Code: dbr.NewNullString("oz"), Name: dbr.NewNullString("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: dbr.NewNullBool(false)},
		),
		store.SetStorageGroups(
			&store.TableGroup{GroupID: 13, WebsiteID: 12, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 4},
		),
		store.SetStorageStores(
			&store.TableStore{StoreID: 4, Code: dbr.NewNullString("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableStore{StoreID: 6, Code: dbr.NewNullString("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)

	stg, err := nsg.Store(scope.MockCode("nz"))
	assert.Nil(t, stg)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)

	stgs, err := nsg.Stores()
	assert.Nil(t, stgs)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

// MBA mid 2012 CPU: Intel Core i5-3427U CPU @ 1.80GHz
// BenchmarkStorageDefaultStoreView	 2000000	       724 ns/op	     176 B/op	       7 allocs/op
// BenchmarkStorageDefaultStoreView-4      	   50000	     40856 ns/op	   15296 B/op	     335 allocs/op
func BenchmarkStorageDefaultStoreView(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkStorageStore, err = testStorage.DefaultStoreView()
		if err != nil {
			b.Error(err)
		}
	}
}

func TestStorageReInit(t *testing.T) {
	t.Parallel()

	// quick implement, use mock of dbr.SessionRunner and remove connection
	dbc := csdb.MustConnectTest()
	defer func() { assert.NoError(t, dbc.Close()) }()

	nsg := store.MustNewStorage(nil, nil, nil)
	assert.NoError(t, nsg.ReInit(dbc.NewSession()))

	stores, err := nsg.Stores()
	assert.NoError(t, err)
	assert.True(t, stores.Len() > 0, "Expecting at least one store loaded from DB")
	for _, s := range stores {
		assert.NotEmpty(t, s.Data.Code.String, "Store: %#v", s.Data)
	}

	groups, err := nsg.Groups()
	assert.True(t, groups.Len() > 0, "Expecting at least one group loaded from DB")
	assert.NoError(t, err)
	for _, g := range groups {
		assert.NotEmpty(t, g.Data.Name, "Group: %#v", g.Data)
	}

	websites, err := nsg.Websites()
	assert.True(t, websites.Len() > 0, "Expecting at least one website loaded from DB")
	assert.NoError(t, err)
	for _, w := range websites {
		assert.NotEmpty(t, w.Data.Code.String, "Website: %#v", w.Data)
	}
}
