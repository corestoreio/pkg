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

package store

import (
	"errors"
	"fmt"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
)

const (
	// DefaultWebsiteID is always 0
	DefaultWebsiteID int64 = 0
)

type (
	// Website contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface WebsiteGetter.
	Website struct {
		w *TableWebsite

		// groups contains a slice to all groups associated to one website. This slice can be nil.
		groups GroupSlice
		// stores contains a slice to all stores associated to one website. This slice can be nil.
		stores StoreSlice
	}
	// WebsiteSlice contains pointer to Website struct and some nifty method receivers.
	WebsiteSlice []*Website
)

var (
	// ErrWebsiteNotFound when the website has not been found within a slice
	ErrWebsiteNotFound = errors.New("Website not found")
	// ErrWebsiteDefaultGroupNotFound the default group cannot be found
	ErrWebsiteDefaultGroupNotFound = errors.New("Website Default Group not found")
	// ErrWebsiteGroupsNotAvailable Groups are in the current context not available and nil
	ErrWebsiteGroupsNotAvailable = errors.New("Website Groups not available")
	// ErrWebsiteStoresNotAvailable Stores are in the current context not available and nil
	ErrWebsiteStoresNotAvailable = errors.New("Website Stores not available")
)

// NewWebsite returns a new pointer to a Website.
func NewWebsite(w *TableWebsite) *Website {
	if w == nil {
		panic(ErrStoreNewArgNil)
	}

	return &Website{
		w: w,
	}
}

// Data returns the data from the database
func (wb *Website) Data() *TableWebsite { return wb.w }

// DefaultGroup returns the default Group or an error if not found
func (wb *Website) DefaultGroup() (*Group, error) {
	for _, g := range wb.groups {
		if wb.w.DefaultGroupID == g.Data().GroupID {
			return g, nil
		}
	}
	return nil, ErrWebsiteDefaultGroupNotFound
}

// DefaultStore returns the default store which via the default group.
func (wb *Website) DefaultStore() (*Store, error) {
	g, err := wb.DefaultGroup()
	if err != nil {
		return nil, err
	}
	return g.DefaultStore()
}

// Stores returns all stores associated to this website or an error when the stores
// are not available aka not needed.
func (wb *Website) Stores() (StoreSlice, error) {
	if len(wb.stores) > 0 {
		return wb.stores, nil
	}
	return nil, ErrWebsiteStoresNotAvailable
}

// Groups returns all groups associated to this website or an error when the groups
// are not available aka not needed.
func (wb *Website) Groups() (GroupSlice, error) {
	if len(wb.groups) > 0 {
		return wb.groups, nil
	}
	return nil, ErrWebsiteGroupsNotAvailable
}

// SetGroupsStores uses a group slice and a table slice to set the groups associated to this website
// and the stores associated to this website. It panics if the integrity is incorrect.
func (wb *Website) SetGroupsStores(tgs TableGroupSlice, tss TableStoreSlice) *Website {
	groups := tgs.FilterByWebsiteID(wb.w.WebsiteID)
	wb.groups = make(GroupSlice, groups.Len(), groups.Len())
	for i, g := range groups {
		wb.groups[i] = NewGroup(g, wb.w).SetStores(tss, nil)
	}
	stores := tss.FilterByWebsiteID(wb.w.WebsiteID)
	wb.stores = make(StoreSlice, stores.Len(), stores.Len())
	for i, s := range stores {
		group, err := tgs.FindByID(s.GroupID)
		if err != nil {
			panic(fmt.Sprintf("Integrity error. A store %#v must be assigned to a group.\nGroupSlice: %#v\n\n", s, tgs))
		}
		wb.stores[i] = NewStore(wb.w, group, s)
	}
	return wb
}

/*
	WebsiteSlice method receivers
*/

// Len returns the length
func (ws WebsiteSlice) Len() int { return len(ws) }

// Filter returns a new slice filtered by predicate f
func (ws WebsiteSlice) Filter(f func(*Website) bool) WebsiteSlice {
	var nws WebsiteSlice
	for _, v := range ws {
		if v != nil && f(v) {
			nws = append(nws, v)
		}
	}
	return nws
}

// Codes returns a StringSlice with all website codes
func (ws WebsiteSlice) Codes() utils.StringSlice {
	if len(ws) == 0 {
		return nil
	}
	var c utils.StringSlice
	for _, w := range ws {
		if w != nil {
			c.Append(w.Data().Code.String)
		}
	}
	return c
}

// IDs returns an Int64Slice with all website ids
func (ws WebsiteSlice) IDs() utils.Int64Slice {
	if len(ws) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, w := range ws {
		if w != nil {
			ids.Append(w.Data().WebsiteID)
		}
	}
	return ids
}

/*
	TableWebsite and TableWebsiteSlice method receivers
*/

// Load uses a dbr session to load all data from the core_website table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Website/Collection.php::Load()
func (s *TableWebsiteSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndexWebsite, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.sort_order ASC").OrderBy("main_table.name ASC")
	})...)
}

// Len returns the length
func (s TableWebsiteSlice) Len() int { return len(s) }

// FindByID returns a TableWebsite if found by id or an error
func (s TableWebsiteSlice) FindByID(id int64) (*TableWebsite, error) {
	for _, w := range s {
		if w != nil && w.WebsiteID == id {
			return w, nil
		}
	}
	return nil, ErrWebsiteNotFound
}

// FindByCode returns a TableWebsite if found by code or an error
func (s TableWebsiteSlice) FindByCode(code string) (*TableWebsite, error) {
	for _, w := range s {
		if w != nil && w.Code.Valid && w.Code.String == code {
			return w, nil
		}
	}
	return nil, ErrWebsiteNotFound
}

// Filter returns a new slice filtered by predicate f
func (s TableWebsiteSlice) Filter(f func(*TableWebsite) bool) TableWebsiteSlice {
	var tws TableWebsiteSlice
	for _, w := range s {
		if w != nil && f(w) {
			tws = append(tws, w)
		}
	}
	return tws
}

// @todo review Magento code because of column is_default
//func (s TableWebsite) IsDefault() bool {
//	return s.WebsiteID == DefaultWebsiteId
//}

/*
	@todo implement Magento\Store\Model\Website
*/
