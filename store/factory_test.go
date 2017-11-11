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
	"testing"

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/util/null"
	"github.com/corestoreio/cspkg/util/slices"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

// todo inspect the high allocs

var testFactory = mustNewFactory(
	cfgmock.NewService(),
	WithTableWebsites(
		&TableWebsite{WebsiteID: 0, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
		&TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		&TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
	),
	WithTableGroups(
		&TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		&TableGroup{GroupID: 1, WebsiteID: 1, Name: "DACH Group", RootCategoryID: 2, DefaultStoreID: 2},
		&TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		&TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
	),
	WithTableStores(
		&TableStore{StoreID: 0, Code: null.StringFrom("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
		&TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 4, Code: null.StringFrom("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
		&TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
		&TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		&TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
	),
)

func TestFactoryWebsite(t *testing.T) {

	tests := []struct {
		have       int64
		wantErrBhf errors.BehaviourFunc
		wantWCode  string
	}{
		{-1, errors.IsNotFound, ""},
		{2015, errors.IsNotFound, ""},
		{1, nil, "euro"},
	}
	for i, test := range tests {
		w, err := testFactory.Website(test.have)
		if test.wantErrBhf != nil {
			assert.Error(t, w.Validate(), "Index %d", i)
			assert.True(t, test.wantErrBhf(err), "Index %d Error: %s", i, err)
		} else {
			assert.NotNil(t, w, "Index %d", i)
			assert.NoError(t, err, "Index %d", i)
			assert.Equal(t, test.wantWCode, w.Data.Code.String, "Index %d", i)
		}
	}

	w, err := testFactory.Website(1)
	assert.NoError(t, err)
	assert.NotNil(t, w)

	dGroup, err := w.DefaultGroup()
	assert.NoError(t, err)
	assert.EqualValues(t, "DACH Group", dGroup.Data.Name)

	assert.NotNil(t, w.Groups)
	assert.EqualValues(t, slices.Int64{1, 2}, w.Groups.IDs())

	assert.NotNil(t, w.Stores)
	assert.EqualValues(t, slices.String{"de", "uk", "at", "ch"}, w.Stores.Codes())
}

func TestFactoryWebsites(t *testing.T) {

	websites, err := testFactory.Websites()
	assert.NoError(t, err)
	assert.EqualValues(t, slices.String{"admin", "euro", "oz"}, websites.Codes())
	assert.EqualValues(t, slices.Int64{0, 1, 2}, websites.IDs())

	var ids = []struct {
		g slices.Int64
		s slices.Int64
	}{
		{slices.Int64{0}, slices.Int64{0}},             //admin
		{slices.Int64{1, 2}, slices.Int64{1, 4, 2, 3}}, // dach
		{slices.Int64{3}, slices.Int64{5, 6}},          // oz
	}

	for i, w := range websites {
		assert.NotNil(t, w.Groups)
		assert.EqualValues(t, ids[i].g, w.Groups.IDs())

		assert.NotNil(t, w.Stores)
		assert.EqualValues(t, ids[i].s, w.Stores.IDs())
	}
}

func TestWebsiteSliceFilter(t *testing.T) {

	websites, err := testFactory.Websites()
	assert.NoError(t, err)
	assert.True(t, websites.Len() == 3)

	gs := websites.Filter(func(w Website) bool {
		return w.Data.WebsiteID > 0
	})
	assert.EqualValues(t, slices.Int64{1, 2}, gs.IDs())
}

func TestFactoryGroup(t *testing.T) {

	tests := []struct {
		id         int64
		wantErrBhf errors.BehaviourFunc
		wantName   string
	}{
		{-1, errors.IsNotFound, ""},
		{2015, errors.IsNotFound, ""},
		{1, nil, "DACH Group"},
	}
	for i, test := range tests {
		g, err := testFactory.Group(test.id)
		if test.wantErrBhf != nil {
			assert.NoError(t, g.Validate())
			assert.True(t, test.wantErrBhf(err), "Index %d Error: %s", i, err)
		} else {
			assert.NotNil(t, g, "Index %d", i)
			assert.NoError(t, err, "Index %d", i)
			assert.Equal(t, test.wantName, g.Data.Name, "Index %d", i)
		}
	}

	g, err := testFactory.Group(3)
	assert.NoError(t, err)
	assert.NotNil(t, g)

	dStore, err := g.DefaultStore()
	assert.NoError(t, err)
	assert.EqualValues(t, "au", dStore.Data.Code.String)

	assert.EqualValues(t, "oz", g.Website.Data.Code.String)

	assert.NotNil(t, g.Stores)
	assert.EqualValues(t, slices.String{"au", "nz"}, g.Stores.Codes())
}

func TestFactoryGroups(t *testing.T) {

	groups, err := testFactory.Groups()
	assert.NoError(t, err)
	assert.EqualValues(t, slices.Int64{3, 1, 0, 2}, groups.IDs())
	assert.True(t, groups.Len() == 4)

	var ids = []slices.Int64{
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

	groups, err := testFactory.Groups()
	assert.NoError(t, err)
	gs := groups.Filter(func(g Group) bool {
		return g.Data.GroupID > 0
	})
	assert.EqualValues(t, slices.Int64{3, 1, 2}, gs.IDs())
}

func TestFactoryGroupNoWebsite(t *testing.T) {

	var tst = mustNewFactory(
		cfgmock.NewService(),
		WithTableWebsites(
			&TableWebsite{WebsiteID: 21, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
		),
		WithTableGroups(
			&TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
		WithTableStores(
			&TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)
	g, err := tst.Group(3)
	assert.NoError(t, g.Validate())
	assert.True(t, errors.IsNotFound(err), err.Error())

	gs, err := tst.Groups()
	assert.Nil(t, gs)
	assert.True(t, errors.IsNotFound(err), err.Error())
}

func TestFactoryStore(t *testing.T) {

	tests := []struct {
		have       int64
		wantErrBhf errors.BehaviourFunc
		wantCode   string
	}{
		{-1, errors.IsNotFound, ""},
		{2015, errors.IsNotFound, ""},
		{1, nil, "de"},
	}
	for i, test := range tests {
		s, err := testFactory.Store(test.have)
		if test.wantErrBhf != nil {
			errV := s.Validate()
			assert.Error(t, errV, "Index %d: %+v", i, errV)
			assert.True(t, test.wantErrBhf(err), "Index: %d Error: %s", i, err)
		} else {
			assert.NotNil(t, s, "Index %d", i)
			assert.NoError(t, err, "Index %d", i)
			assert.Equal(t, test.wantCode, s.Data.Code.String, "Index %d", i)
		}
	}

	s, err := testFactory.Store(2)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.Exactly(t, "DACH Group", s.Group.Data.Name)

	assert.Exactly(t, "euro", s.Website.Data.Code.String)
	wg, err := s.Website.DefaultGroup()
	assert.NotNil(t, wg)
	assert.Exactly(t, "DACH Group", wg.Data.Name)

	wgs, err := testFactory.Store(wg.DefaultStoreID())
	assert.NoError(t, err, "%+v", err)
	assert.Exactly(t, "at", wgs.Code(), " WebsiteGroup Stores: %#v", wg.Stores)
}

func TestFactoryStores(t *testing.T) {

	stores, err := testFactory.Stores()
	assert.NoError(t, err)
	assert.EqualValues(t, slices.String{"admin", "au", "de", "uk", "at", "nz", "ch"}, stores.Codes())
	assert.EqualValues(t, slices.Int64{0, 5, 1, 4, 2, 6, 3}, stores.IDs())

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

	st, err := testFactory.DefaultStoreID()
	assert.NoError(t, err)
	assert.Exactly(t, int64(2), st)

	tst := mustNewFactory(
		cfgmock.NewService(),
		WithTableWebsites(
			&TableWebsite{WebsiteID: 21, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
		),
		WithTableGroups(
			&TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
		WithTableStores(
			&TableStore{StoreID: 4, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)
	dSt, err := tst.DefaultStoreID()
	assert.Empty(t, dSt)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)

	var tst2 = mustNewFactory(
		cfgmock.NewService(),
		WithTableWebsites(
			&TableWebsite{WebsiteID: 21, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(true)},
		),
		WithTableGroups(
			&TableGroup{GroupID: 33, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
		WithTableStores(),
	)
	dSt2, err := tst2.DefaultStoreID()
	assert.Empty(t, dSt2)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

func TestFactoryStoreErrors(t *testing.T) {

	var nsw = mustNewFactory(
		cfgmock.NewService(),
		WithTableWebsites(),
		WithTableGroups(),
		WithTableStores(
			&TableStore{StoreID: 4, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)
	stw, err := nsw.Store(6)
	assert.Error(t, stw.Validate())
	assert.True(t, errors.IsNotFound(err), err.Error())

	stws, err := nsw.Stores()
	assert.Nil(t, stws)
	assert.True(t, errors.IsNotFound(err), err.Error())

	var nsg = mustNewFactory(
		cfgmock.NewService(),
		WithTableWebsites(
			&TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
		),
		WithTableGroups(
			&TableGroup{GroupID: 13, WebsiteID: 12, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 4},
		),
		WithTableStores(
			&TableStore{StoreID: 4, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
		),
	)

	stg, err := nsg.Store(6)
	assert.Error(t, stg.Validate())
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)

	stgs, err := nsg.Stores()
	assert.Nil(t, stgs)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}
