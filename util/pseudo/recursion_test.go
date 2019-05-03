// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package pseudo_test

import (
	"encoding/json"
	"testing"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
)

type Store struct {
	StoreID      uint32        `max_len:"5"`   // store_id smallint(5) unsigned NOT NULL PRI  auto_increment "Store ID"
	Code         null.String   `max_len:"64"`  // code varchar(64) NULL UNI DEFAULT 'NULL'  "Code"
	WebsiteID    uint32        `max_len:"5"`   // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Website ID"
	GroupID      uint32        `max_len:"5"`   // group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Group ID"
	Name         string        `max_len:"255"` // name varchar(255) NOT NULL    "Store Name"
	SortOrder    uint32        `max_len:"5"`   // sort_order smallint(5) unsigned NOT NULL  DEFAULT '0'  "Store Sort Order"
	IsActive     bool          `max_len:"5"`   // is_active smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Store Activity"
	StoreGroup   *StoreGroup   // 1:1 store.group_id => store_group.group_id
	StoreWebsite *StoreWebsite // 1:1 store.website_id => store_website.website_id
}
type StoreCollection struct {
	Data []*Store `json:"data,omitempty"`
}
type StoreGroup struct {
	GroupID        uint32        `max_len:"5"`   // group_id smallint(5) unsigned NOT NULL PRI  auto_increment "Group ID"
	WebsiteID      uint32        `max_len:"5"`   // website_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Website ID"
	Name           string        `max_len:"255"` // name varchar(255) NOT NULL    "Store Group Name"
	RootCategoryID uint32        `max_len:"10"`  // root_category_id int(10) unsigned NOT NULL  DEFAULT '0'  "Root Category ID"
	DefaultStoreID uint32        `max_len:"5"`   // default_store_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Default Store ID"
	Code           null.String   `max_len:"64"`  // code varchar(64) NULL UNI DEFAULT 'NULL'  "Store group unique code"
	StoreWebsite   *StoreWebsite // 1:1 store_group.website_id => store_website.website_id
}

type StoreWebsite struct {
	WebsiteID      uint32           `max_len:"5"`   // website_id smallint(5) unsigned NOT NULL PRI  auto_increment "Website ID"
	Code           null.String      `max_len:"64"`  // code varchar(64) NULL UNI DEFAULT 'NULL'  "Code"
	Name           null.String      `max_len:"128"` // name varchar(128) NULL  DEFAULT 'NULL'  "Website Name"
	SortOrder      uint32           `max_len:"5"`   // sort_order smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Sort Order"
	DefaultGroupID uint32           `max_len:"5"`   // default_group_id smallint(5) unsigned NOT NULL MUL DEFAULT '0'  "Default Group ID"
	IsDefault      bool             `max_len:"5"`   // is_default smallint(5) unsigned NOT NULL  DEFAULT '0'  "Defines Is Website Default"
	Store          *StoreCollection // Reversed 1:M store_website.website_id => store.website_id
}

func init() {
	null.MustSetJSONMarshaler(json.Marshal, nil)
}

func TestRecursion(t *testing.T) {
	ps, err := pseudo.NewService(100, &pseudo.Options{})
	assert.NoError(t, err)

	t.Run("Store", func(t *testing.T) {
		st := new(Store)
		err = ps.FakeData(st)
		assert.NoError(t, err)

		// just check that it works an no stack overflow happens.
		data, err := json.Marshal(st)
		assert.NoError(t, err)
		assert.LenBetween(t, data, 2000, 4000)
	})

	t.Run("StoreWebsite", func(t *testing.T) {
		st := new(StoreWebsite)
		err = ps.FakeData(st)
		assert.NoError(t, err)

		// just check that it works an no stack overflow happens.
		data, err := json.Marshal(st)
		assert.NoError(t, err)
		assert.LenBetween(t, data, 170, 220)
	})
}
