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
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/errors"
)

// factory contains the raw slices from the database and can read from the
// database. It creates for each call to each of its method receivers new
// Stores, Groups or Websites.
type factory struct {
	// baseConfig parent config service. can only be set once.
	baseConfig config.Getter
	mu         sync.RWMutex
	websites   TableWebsiteSlice
	groups     TableGroupSlice
	stores     TableStoreSlice
}

// newFactory creates a new object which handles the raw data from the three
// database tables for website, group and store. You can either provide the raw
// data separately for each type or pass an option to load it from the database.
// To set the raw data either call the WithTable*() functions or use ReInit()
// and a DB connection.
func newFactory(cfg config.Getter, opts ...Option) (*factory, error) {
	s := &factory{
		baseConfig: cfg,
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(s); err != nil {
				return nil, errors.Wrap(err, "[store] NewStorage Option")
			}
		}
	}
	return s, nil
}

// website returns a TableWebsite by using the id.
func (st factory) website(id int64) (*TableWebsite, bool) {
	return st.websites.FindByWebsiteID(id)
}

// Website creates a new Website  from an ID including all of its groups
// and all related stores. Returns a NotFound error behaviour.
func (st factory) Website(id int64) (Website, error) {
	w, found := st.website(id)
	if !found {
		return Website{}, errors.NewNotFoundf("[store] WebsiteID %d", id)
	}
	return NewWebsite(st.baseConfig, w, st.groups, st.stores)
}

// Websites creates a slice containing all new pointers to Websites with its
// associated new groups and new store pointers. It returns an error if the
// integrity is incorrect or NotFound errors.
func (st factory) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(st.websites), len(st.websites))
	for i, w := range st.websites {
		var err error
		websites[i], err = NewWebsite(st.baseConfig, w, st.groups, st.stores)
		if err != nil {
			return nil, errors.Wrapf(err, "[store] Storage.Websites. WebsiteID: %d", w.WebsiteID)
		}
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument.
func (st factory) group(id int64) (*TableGroup, bool) {
	return st.groups.FindByGroupID(id)
}

// Group creates a new Group  for an ID which contains all related store-
// and its website-pointers.
func (st factory) Group(id int64) (Group, error) {
	g, found := st.group(id)
	if !found {
		return Group{}, errors.NewNotFoundf("[store] Group %d", id)
	}

	w, found := st.website(g.WebsiteID)
	if !found {
		return Group{}, errors.NewNotFoundf("[store] Website. WebsiteID %d GroupID %v", g.WebsiteID, id)
	}
	return NewGroup(st.baseConfig, g, w, st.stores)
}

// Groups creates a slice containing all pointers to Groups with its associated
// new store- and new website-pointers. It returns an error if the integrity is
// incorrect or a NotFound error.
func (st factory) Groups() (GroupSlice, error) {
	groups := make(GroupSlice, len(st.groups), len(st.groups))
	for i, g := range st.groups {
		w, found := st.website(g.WebsiteID)
		if !found {
			return nil, errors.NewNotFoundf("[store] WebsiteID %d", g.WebsiteID)
		}
		var err error
		groups[i], err = NewGroup(st.baseConfig, g, w, st.stores)
		if err != nil {
			return nil, errors.Wrapf(err, "[store] GroupID %d WebsiteID %d", g.GroupID, g.WebsiteID)
		}
	}
	return groups, nil
}

// store returns a TableStore by an id.
func (st factory) store(id int64) (*TableStore, bool) {
	return st.stores.FindByStoreID(id)
}

// Store creates a new Store  containing its group and its website.
// Returns an error if the integrity is incorrect. May return a NotFound error
// behaviour.
func (st factory) Store(id int64) (Store, error) {
	var ns Store
	s, found := st.store(id)
	if !found {
		return ns, errors.NewNotFoundf("[store] Store: %d", id)
	}
	w, found := st.website(s.WebsiteID)
	if !found {
		return ns, errors.NewNotFoundf("[store] WebsiteID: %d", s.WebsiteID)
	}
	g, found := st.group(s.GroupID)
	if !found {
		return ns, errors.NewNotFoundf("[store] GroupID: %d", s.GroupID)
	}
	var err error
	ns, err = NewStore(st.baseConfig, s, w, g)
	if err != nil {
		return ns, errors.Wrapf(err, "[store] StoreID %d WebsiteID %d GroupID %d", s.StoreID, w.WebsiteID, g.GroupID)
	}
	if err := ns.Website.SetGroupsStores(st.groups, st.stores); err != nil {
		return ns, errors.Wrap(err, "")
	}
	if err := ns.Group.SetWebsiteStores(st.baseConfig, w, st.stores); err != nil {
		return ns, errors.Wrap(err, "[store] Storage.Store.Group.SetWebsiteStores")
	}
	return ns, nil
}

// Stores creates a new store slice with all of its new Group and new Website
// pointers. Can return an error when the website or the group cannot be found.
func (st factory) Stores() (StoreSlice, error) {
	stores := make(StoreSlice, len(st.stores), len(st.stores))
	for i, s := range st.stores {
		var err error
		if stores[i], err = st.Store(s.StoreID); err != nil {
			return nil, errors.Wrapf(err, "[store] StoreID %d", s.StoreID)
		}
	}
	return stores, nil
}

// DefaultStoreID traverses through the websites to find the default website
// and gets the default group which has the default store id assigned to. Only
// one website can be the default one.
func (st factory) DefaultStoreID() (int64, error) {
	for _, w := range st.websites {
		if w.IsDefault.Bool && w.IsDefault.Valid {
			g, found := st.group(w.DefaultGroupID)
			if !found {
				return 0, errors.NewNotFoundf("[store] WebsiteID %d DefaultGroupID %d", w.WebsiteID, w.DefaultGroupID)
			}
			return g.DefaultStoreID, nil
		}
	}
	return 0, errors.NewNotFoundf(errStoreDefaultNotFound)
}

// LoadFromDB reloads all websites, groups and stores concurrently from the
// database. On error  all internal slices will be reset to nil.
func (st *factory) LoadFromDB(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	errc := make(chan error)
	defer close(errc)
	// not sure about those three go
	go func() {
		for i := range st.websites {
			st.websites[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.websites = nil
		_, err := st.websites.SQLSelect(dbrSess, cbs...)
		errc <- errors.Wrap(err, "[store] SQLSelect websites")
	}()

	go func() {
		for i := range st.groups {
			st.groups[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.groups = nil
		_, err := st.groups.SQLSelect(dbrSess, cbs...)
		errc <- errors.Wrap(err, "[store] SQLSelect groups")
	}()

	go func() {
		for i := range st.stores {
			st.stores[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.stores = nil
		_, err := st.stores.SQLSelect(dbrSess, cbs...)
		errc <- errors.Wrap(err, "[store] SQLSelect stores")
	}()

	for i := 0; i < 3; i++ {
		if err := <-errc; err != nil {
			// in case of error clear all
			st.websites = nil
			st.groups = nil
			st.stores = nil
			return err
		}
	}
	return nil
}
