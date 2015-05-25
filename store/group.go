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

// Package store implements the handling of websites, groups and stores
package store

import (
	"errors"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
)

const (
	// DefaultGroupID defines the default group id which is always 0.
	DefaultGroupID int64 = 0
)

type (

	// Group contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface GroupGetter.
	Group struct {
		cr config.Reader
		// g group data
		g *TableGroup
		// stores contains a slice to all stores associated to this group.
		// This slice can be nil
		stores StoreSlice
		// w contains the Website which belongs to this group. Can be nil.
		w *Website
	}
	// GroupSlice collection of Group. GroupSlice has some nice method receivers.
	GroupSlice []*Group

	GroupOption func(*Group)
)

var (
	// ErrGroupNotFound when the group has not been found.
	ErrGroupNotFound = errors.New("Group not found")
	// ErrGroupStoresNotAvailable not really an error but more an info when the stores has not been set
	// this usually occurs when the group has been set on a website or a store.
	ErrGroupStoresNotAvailable = errors.New("Group stores not available")
	// ErrGroupDefaultStoreNotFound default store cannot be found.
	ErrGroupDefaultStoreNotFound = errors.New("Group default store not found")
	// ErrGroupWebsiteNotFound the Website struct is nil so we cannot assign the stores to a group.
	ErrGroupWebsiteNotFound = errors.New("Group Website not found or nil or ID do not match")
)

// SetGroupConfig adds a configuration Reader to the Group. Optional.
// Default reader is config.DefaultManager
func SetGroupConfig(cr config.Reader) GroupOption {
	return func(g *Group) { g.cr = cr }
}

func SetGroupWebsite(tw *TableWebsite) GroupOption {
	return func(g *Group) {
		if g.Data() == nil {
			panic(ErrGroupNotFound)
		}
		if g.Data().WebsiteID != tw.WebsiteID {
			panic(ErrGroupWebsiteNotFound)
		}
		g.w = NewWebsite(tw)
	}
}

// NewGroup returns a new pointer to a Group. Second argument can be nil.
func NewGroup(tg *TableGroup, opts ...GroupOption) *Group {
	if tg == nil {
		panic(ErrStoreNewArgNil)
	}

	g := &Group{
		cr: config.DefaultManager,
		g:  tg,
	}
	g.ApplyOptions(opts...)
	return g
}

// ApplyOptions sets the options
func (g *Group) ApplyOptions(opts ...GroupOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(g)
		}
	}
}

// Data returns the TableGroup data which is raw database data.
func (g *Group) Data() *TableGroup {
	return g.g
}

// Website returns the website associated to this group or nil.
func (g *Group) Website() *Website {
	return g.w
}

// DefaultStore returns the default Store or an error.
func (g *Group) DefaultStore() (*Store, error) {
	for _, sb := range g.stores {
		if sb.Data().StoreID == g.g.DefaultStoreID {
			return sb, nil
		}
	}
	return nil, ErrGroupDefaultStoreNotFound
}

// Stores returns all stores associated to a group or an error if stores are not available.
func (g *Group) Stores() (StoreSlice, error) {
	if len(g.stores) > 0 {
		return g.stores, nil
	}
	return nil, ErrGroupStoresNotAvailable
}

// SetStores uses the full store collection to extract the stores which are
// assigned to a group. Either Website must be set before calling SetStores() or
// the second argument must be set i.e. 2nd argument can be nil. Panics if both
// values are nil. If both are set, the 2nd argument will be considered.
func (g *Group) SetStores(tss TableStoreSlice, w *TableWebsite) *Group {
	if tss == nil {
		g.stores = nil
		return g
	}
	if g.Website() == nil && w == nil {
		panic(ErrGroupWebsiteNotFound)
	}
	if w == nil {
		w = g.Website().Data()
	}
	if w.WebsiteID != g.Data().WebsiteID {
		panic(ErrGroupWebsiteNotFound)
	}
	for _, s := range tss.FilterByGroupID(g.g.GroupID) {
		g.stores = append(g.stores, NewStore(s, SetStoreGroup(g.g), SetStoreWebsite(w), SetStoreConfig(g.cr)))
	}
	return g
}

/*
	@todo implement Magento\Store\Model\Group
*/

/*
	GroupSlice method receivers
*/

// Len returns the length
func (s GroupSlice) Len() int { return len(s) }

// Filter returns a new slice filtered by predicate f
func (s GroupSlice) Filter(f func(*Group) bool) GroupSlice {
	var gs GroupSlice
	for _, v := range s {
		if v != nil && f(v) {
			gs = append(gs, v)
		}
	}
	return gs
}

// IDs returns an Int64Slice with all store ids
func (s GroupSlice) IDs() utils.Int64Slice {
	if len(s) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, g := range s {
		if g != nil {
			ids.Append(g.Data().GroupID)
		}
	}
	return ids
}

/*
	TableGroup and TableGroupSlice method receivers
*/

// Load uses a dbr session to load all data from the core_store_group table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Group/Collection.php::_beforeLoad()
func (s *TableGroupSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndexGroup, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.name ASC")
	})...)
}

// Len returns the length
func (s TableGroupSlice) Len() int { return len(s) }

// FindByID returns a TableGroup if found by id or an error
func (s TableGroupSlice) FindByID(id int64) (*TableGroup, error) {
	for _, g := range s {
		if g.GroupID == id {
			return g, nil
		}
	}
	return nil, ErrGroupNotFound
}

// FilterByWebsiteID returns a new slice with all groups belonging to a website id
func (s TableGroupSlice) FilterByWebsiteID(id int64) TableGroupSlice {
	return s.Filter(func(w *TableGroup) bool {
		return w.WebsiteID == id
	})
}

// Filter returns a new slice filtered by predicate f
func (s TableGroupSlice) Filter(f func(*TableGroup) bool) TableGroupSlice {
	var tgs TableGroupSlice
	for _, v := range s {
		if v != nil && f(v) {
			tgs = append(tgs, v)
		}
	}
	return tgs
}

// IDs returns an Int64Slice with all group ids
func (s TableGroupSlice) IDs() utils.Int64Slice {
	id := make(utils.Int64Slice, len(s))
	for i, store := range s {
		id[i] = store.GroupID
	}
	return id
}
