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

package store_test

import (
	"database/sql"
	"testing"

	"bytes"

	"encoding/json"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store"
	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	s := store.NewStore(
		&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		&store.TableGroup{GroupID: 1, WebsiteID: 0, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	)

	var buf bytes.Buffer
	assert.NoError(t, s.ToJSON(&buf))

	assert.Equal(t, `{"StoreID":1,"Code":"de","WebsiteID":1,"GroupID":1,"Name":"Germany","SortOrder":10,"IsActive":true}`, buf.String())

	var ds store.TableStore
	dec := json.NewDecoder(&buf)
	dec.Decode(&ds)

	assert.Equal(t, "de", ds.Code.String)
	assert.Equal(t, "Germany", ds.Name)
	assert.Equal(t, int64(1), ds.WebsiteID)

}
