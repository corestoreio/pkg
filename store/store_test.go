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
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/corestoreio/cspkg/config/cfgmock"
	"github.com/corestoreio/cspkg/store"
	"github.com/corestoreio/cspkg/util/null"
	"github.com/corestoreio/cspkg/util/slices"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/stretchr/testify/assert"
)

var _ log.Marshaler = (*store.Store)(nil)
var _ fmt.Stringer = (*store.Store)(nil)

const TODO_Better_Test_Data = "@todo implement better test data which is equal for each Magento version"

func TestNewStore(t *testing.T) {

	tests := []struct {
		w *store.TableWebsite
		g *store.TableGroup
		s *store.TableStore
	}{
		{
			w: &store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
			g: &store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			s: &store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		},
		{
			w: &store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
			g: &store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
			s: &store.TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
		},
	}
	for _, test := range tests {
		s, err := store.NewStore(cfgmock.NewService(), test.s, test.w, test.g)
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.EqualValues(t, test.w.WebsiteID, s.Website.Data.WebsiteID)
		assert.EqualValues(t, test.g.GroupID, s.Group.Data.GroupID)
		assert.EqualValues(t, test.s.Code, s.Data.Code)
		assert.NotNil(t, s.Group.Website)
		assert.NotEmpty(t, s.Group.Website.ID())
		assert.Nil(t, s.Group.Stores)
		assert.EqualValues(t, test.s.StoreID, s.ID())
		assert.EqualValues(t, test.s.GroupID, s.GroupID())
		assert.EqualValues(t, test.s.WebsiteID, s.WebsiteID())
	}
}

func TestNewStoreErrorIncorrectGroup(t *testing.T) {

	s, err := store.NewStore(
		cfgmock.NewService(),
		&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		&store.TableGroup{GroupID: 2, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
	)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	err = s.Validate()
	assert.True(t, errors.IsNotValid(err), "%+v", err)
}

func TestNewStoreErrorIncorrectWebsite(t *testing.T) {

	s, err := store.NewStore(
		cfgmock.NewService(),
		&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("euro"), Name: null.StringFrom("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: null.BoolFrom(true)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "UK Group", RootCategoryID: 2, DefaultStoreID: 4},
	)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
	err = s.Validate()
	assert.True(t, errors.IsNotValid(err), "%+v", err)
}

func TestStoreSlice(t *testing.T) {

	storeSlice := store.StoreSlice{
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
			&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
		),
		store.MustNewStore(
			cfgmock.NewService(),
			&store.TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
			&store.TableWebsite{WebsiteID: 2, Code: null.StringFrom("oz"), Name: null.StringFrom("OZ"), SortOrder: 20, DefaultGroupID: 3, IsDefault: null.BoolFrom(false)},
			&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
		),
	}
	assert.True(t, storeSlice.Len() == 2)
	assert.EqualValues(t, slices.Int64{1, 5}, storeSlice.IDs())
	assert.EqualValues(t, slices.String{"de", "au"}, storeSlice.Codes())

	storeSlice2 := storeSlice.Filter(func(s store.Store) bool {
		return s.Website.Data.WebsiteID == 2
	})
	assert.True(t, storeSlice2.Len() == 1)
	assert.Equal(t, "au", storeSlice2[0].Data.Code.String)
	assert.EqualValues(t, slices.Int64{5}, storeSlice2.IDs())
	assert.EqualValues(t, slices.String{"au"}, storeSlice2.Codes())

	assert.Nil(t, (store.StoreSlice{}).IDs())
	assert.Nil(t, (store.StoreSlice{}).Codes())
}

var testStores = store.TableStoreSlice{
	&store.TableStore{StoreID: 0, Code: null.StringFrom("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
	&store.TableStore{StoreID: 5, Code: null.StringFrom("au"), WebsiteID: 2, GroupID: 3, Name: "Australia", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 4, Code: null.StringFrom("uk"), WebsiteID: 1, GroupID: 2, Name: "UK", SortOrder: 10, IsActive: true},
	&store.TableStore{StoreID: 2, Code: null.StringFrom("at"), WebsiteID: 1, GroupID: 1, Name: "Österreich", SortOrder: 20, IsActive: true},
	&store.TableStore{StoreID: 6, Code: null.StringFrom("nz"), WebsiteID: 2, GroupID: 3, Name: "Kiwi", SortOrder: 30, IsActive: true},
	&store.TableStore{StoreID: 3, Code: null.StringFrom("ch"), WebsiteID: 1, GroupID: 1, Name: "Schweiz", SortOrder: 30, IsActive: true},
}

func TestTableStoreSliceFindByID(t *testing.T) {

	const eLen = 7
	assert.Exactly(t, eLen, testStores.Len())

	s1, found := testStores.FindByStoreID(999)
	assert.Nil(t, s1)
	assert.False(t, found)

	s2, found := testStores.FindByStoreID(6)
	assert.NotNil(t, s2)
	assert.True(t, found)
	assert.Equal(t, int64(6), s2.StoreID)
}

func TestTableStoreSliceFindByCode(t *testing.T) {

	s1, found := testStores.FindByCode("corestore")
	assert.Nil(t, s1)
	assert.False(t, found)

	s2, found := testStores.FindByCode("ch")
	assert.NotNil(t, s2)
	assert.True(t, found)
	assert.Equal(t, "ch", s2.Code.String)
}

func TestTableStoreSliceFilterByGroupID(t *testing.T) {

	gStores := testStores.FilterByGroupID(3)
	assert.NotNil(t, gStores)
	assert.Len(t, gStores, 2)
	gStores2 := testStores.FilterByGroupID(32)
	assert.NotNil(t, gStores2)
	assert.Len(t, gStores2, 0)
}

func TestTableStoreSliceFilterByWebsiteID(t *testing.T) {

	gStores := testStores.FilterByWebsiteID(0)
	assert.NotNil(t, gStores)
	assert.Len(t, gStores, 1)
	gStores2 := testStores.FilterByWebsiteID(32)
	assert.NotNil(t, gStores2)
	assert.Len(t, gStores2, 0)

	var ts = store.TableStoreSlice{}
	tsRes := ts.FilterByGroupID(2)
	assert.NotNil(t, tsRes)
	assert.Len(t, tsRes, 0)
}

func TestTableStoreSliceCodes(t *testing.T) {

	t.Log(TODO_Better_Test_Data)

	codes := testStores.Extract().Code()
	assert.NotNil(t, codes)
	assert.Equal(t, []string{"admin", "au", "de", "uk", "at", "nz", "ch"}, codes)

	var ts = store.TableStoreSlice{}
	assert.Empty(t, ts.Extract().Code())
}

func TestTableStoreSliceIDs(t *testing.T) {

	t.Log(TODO_Better_Test_Data)

	ids := testStores.Extract().StoreID()
	assert.NotNil(t, ids)
	assert.Equal(t, []int64{0, 5, 1, 4, 2, 6, 3}, ids)

	var ts = store.TableStoreSlice{}
	assert.Empty(t, ts.Extract().StoreID())
}

func TestStore_MarshalJSON(t *testing.T) {
	s := store.MustNewStore(
		cfgmock.NewService(),
		&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	)

	jdata, err := json.Marshal(s)
	assert.NoError(t, err)
	have := []byte(`{"StoreID":1,"Code":"de","WebsiteID":1,"GroupID":1,"Name":"Germany","SortOrder":10,"IsActive":true}`)
	assert.Equal(t, have, jdata, "Have: %s\nWant: %s", have, jdata)
}

func TestStore_MarshalLog(t *testing.T) {
	s := store.MustNewStore(
		cfgmock.NewService(),
		&store.TableStore{StoreID: 1, Code: null.StringFrom("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: null.StringFrom("admin"), Name: null.StringFrom("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: null.BoolFrom(false)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	)
	buf := bytes.Buffer{}
	lg := logw.NewLog(logw.WithWriter(&buf), logw.WithLevel(logw.LevelDebug))

	lg.Debug("storeTest", log.Marshal("aStore1", s))

	have := `store_id: 1 store_code: "de"`
	assert.Contains(t, buf.String(), have)
}
