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

// generated via https://github.com/ChimeraCoder/gojson json2struct
type TestToJSONStore struct {
	Code      string `json:"Code"`
	GroupID   int    `json:"GroupID"`
	IsActive  bool   `json:"IsActive"`
	Name      string `json:"Name"`
	SortOrder int    `json:"SortOrder"`
	StoreID   int    `json:"StoreID"`
	WebsiteID int    `json:"WebsiteID"`
}

func TestToJSON(t *testing.T) {
	s := store.NewStore(
		&store.TableStore{StoreID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "de", Valid: true}}, WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NullString{NullString: sql.NullString{String: "admin", Valid: true}}, Name: dbr.NullString{NullString: sql.NullString{String: "Admin", Valid: true}}, SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NullBool{NullBool: sql.NullBool{Bool: false, Valid: true}}},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 0},
	)

	var buf bytes.Buffer
	assert.NoError(t, s.ToJSON(&buf))

	tsd := TestToJSONStore{}
	json.Unmarshal(buf.Bytes(), &tsd)

	want := TestToJSONStore{Code: "de", GroupID: 1, IsActive: true, Name: "Germany", SortOrder: 10, StoreID: 1, WebsiteID: 1}

	assert.Equal(t, want, tsd)

	var ds store.TableStore
	dec := json.NewDecoder(&buf)
	dec.Decode(&ds)

	assert.Equal(t, "de", ds.Code.String)
	assert.Equal(t, "Germany", ds.Name)
	assert.Equal(t, int64(1), ds.WebsiteID)

}
