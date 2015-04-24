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

func TestNewStore(t *testing.T) {

	tests := []struct {
		w *store.TableWebsite
		g *store.TableGroup
		s *store.TableStore
	}{
		{
			w: &store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "base", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Main Website", Valid: true}}, SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}}, IsStaging: false, MasterLogin: dbr.NullString{NullString: sql.NullString{String: "", Valid: false}}, MasterPassword: dbr.NullString{NullString: sql.NullString{String: "", Valid: false}}, Visibility: dbr.NullString{NullString: sql.NullString{String: "", Valid: false}}},
			g: &store.TableGroup{GroupID: 0, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
			s: &store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "default", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "English", SortOrder: 0, IsActive: true},
		},
		{
			w: &store.TableWebsite{WebsiteID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "oz", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "OZ", Valid: true}}, SortOrder: 20, DefaultGroupID: 4, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}, IsStaging: false, MasterLogin: dbr.NullString{NullString: sql.NullString{String: "", Valid: false}}, MasterPassword: dbr.NullString{NullString: sql.NullString{String: "", Valid: false}}, Visibility: dbr.NullString{NullString: sql.NullString{String: "", Valid: false}}},
			g: &store.TableGroup{GroupID: 4, WebsiteID: 4, Name: "AU+NZ", RootCategoryID: 2, DefaultStoreID: 4},
			s: &store.TableStore{StoreID: 4, Code: dbr.NullString{NullString: sql.NullString{String: "ozau", Valid: true}}, WebsiteID: 4, GroupID: 4, Name: "Australia", SortOrder: 0, IsActive: true},
		},
	}
	for _, test := range tests {
		s := store.NewStore(test.w, test.g, test.s)
		assert.NotNil(t, s)
		assert.EqualValues(t, test.w.WebsiteID, s.Website.Data().WebsiteID)
		assert.EqualValues(t, test.g.GroupID, s.Group.Data().GroupID)
		assert.EqualValues(t, test.s.Code, s.Data().Code)
		assert.Nil(t, s.Group.Website)
		gStores, gErr := s.Group.Stores()
		assert.Nil(t, gStores)
		assert.EqualError(t, store.ErrGroupStoresNotAvailable, gErr.Error())
	}
}

func TestNewStorePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.EqualError(t, store.ErrStoreNewArgNil, err.Error())
		}
	}()
	_ = store.NewStore(nil, nil, nil)
}
