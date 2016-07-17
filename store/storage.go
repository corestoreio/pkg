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
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// Storage contains the raw slices from the database and can read from the
// database. It creates for each call to each of its method receivers new
// pointers to Stores, Groups or Websites.
type Storage struct {
	// baseConfig parent config service. can only be set once.
	baseConfig config.Getter
	mu         sync.RWMutex
	websites   TableWebsiteSlice
	groups     TableGroupSlice
	stores     TableStoreSlice
}

// NewStorage creates a new storage object which handles the raw data from the
// three database tables for website, group and store. You can either provide
// the raw data separately for each type or pass an option to load it from the
// database. Passing no function option causes panics on nil.
//		sto, err = store.NewStorage(
//			cfg,
//			store.SetStorageWebsites(
//				&store.TableWebsite{WebsiteID: 0, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
//				...
//			),
//			store.SetStorageGroups(
//				&store.TableGroup{GroupID: 3, WebsiteID: 2, Name: "Australia", RootCategoryID: 2, DefaultStoreID: 5},
//				...
//			),
//			store.SetStorageStores(
//				&store.TableStore{StoreID: 0, Code: dbr.NewNullString("admin"), WebsiteID: 0, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
//				...
//			),
//		)
//		// or alternatively:
// 		sto, err = store.NewStorage(cfg, store.WithDatabaseInit(dbrSession) )
func NewStorage(cfg config.Getter, opts ...StorageOption) (*Storage, error) {
	s := &Storage{
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

// MustNewStorage same as NewStorage but panics on error.
func MustNewStorage(cfg config.Getter, opts ...StorageOption) *Storage {
	s, err := NewStorage(cfg, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// website returns a TableWebsite by using the id.
func (st *Storage) website(id int64) (*TableWebsite, bool) {
	return st.websites.FindByWebsiteID(id)
}

// Website creates a new Website pointer from an ID including all of its groups
// and all related stores. Returns a NotFound error behaviour.
func (st *Storage) Website(id int64) (*Website, error) {
	w, found := st.website(id)
	if !found {
		return nil, errors.NewNotFoundf("[store] WebsiteID %d", id)
	}
	return NewWebsite(st.baseConfig, w, SetWebsiteGroupsStores(st.groups, st.stores))
}

// Websites creates a slice containing all new pointers to Websites with its
// associated new groups and new store pointers. It returns an error if the
// integrity is incorrect or NotFound errors.
func (st *Storage) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(st.websites), len(st.websites))
	for i, w := range st.websites {
		var err error
		websites[i], err = NewWebsite(st.baseConfig, w, SetWebsiteGroupsStores(st.groups, st.stores))
		if err != nil {
			return nil, errors.Wrapf(err, "[store] Storage.Websites. WebsiteID: %d", w.WebsiteID)
		}
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument.
func (st *Storage) group(id int64) (*TableGroup, bool) {
	return st.groups.FindByGroupID(id)
}

// Group creates a new Group pointer for an ID which contains all related store-
// and its website-pointers.
func (st *Storage) Group(id int64) (*Group, error) {
	g, found := st.group(id)
	if !found {
		return nil, errors.NewNotFoundf("[store] Group %d", id)
	}

	w, found := st.website(g.WebsiteID)
	if !found {
		return nil, errors.NewNotFoundf("[store] Website. WebsiteID %d GroupID %v", g.WebsiteID, id)
	}
	return NewGroup(st.baseConfig, g, SetGroupWebsite(w), SetGroupStores(st.stores, nil))
}

// Groups creates a slice containing all pointers to Groups with its associated
// new store- and new website-pointers. It returns an error if the integrity is
// incorrect or a NotFound error.
func (st *Storage) Groups() (GroupSlice, error) {
	groups := make(GroupSlice, len(st.groups), len(st.groups))
	for i, g := range st.groups {
		w, found := st.website(g.WebsiteID)
		if !found {
			return nil, errors.NewNotFoundf("[store] WebsiteID %d", g.WebsiteID)
		}
		var err error
		groups[i], err = NewGroup(st.baseConfig, g, SetGroupWebsite(w), SetGroupStores(st.stores, nil))
		if err != nil {
			return nil, errors.Wrapf(err, "[store] GroupID %d WebsiteID %d", g.GroupID, g.WebsiteID)
		}
	}
	return groups, nil
}

// store returns a TableStore by an id.
func (st *Storage) store(id int64) (*TableStore, bool) {
	return st.stores.FindByStoreID(id)
}

// Store creates a new Store pointer containing its group and its website.
// Returns an error if the integrity is incorrect. May return a NotFound error
// behaviour.
func (st *Storage) Store(id int64) (*Store, error) {
	s, found := st.store(id)
	if !found {
		return nil, errors.NewNotFoundf("[store] Store: %d", id)
	}
	w, found := st.website(s.WebsiteID)
	if !found {
		return nil, errors.NewNotFoundf("[store] WebsiteID: %d", s.WebsiteID)
	}
	g, found := st.group(s.GroupID)
	if !found {
		return nil, errors.NewNotFoundf("[store] GroupID: %d", s.GroupID)
	}
	ns, err := NewStore(st.baseConfig, s, w, g)
	if err != nil {
		return nil, errors.Wrapf(err, "[store] StoreID %d WebsiteID %d GroupID %d", s.StoreID, w.WebsiteID, g.GroupID)
	}
	if err := ns.Website.Options(SetWebsiteGroupsStores(st.groups, st.stores)); err != nil {
		return nil, errors.Wrap(err, "")
	}
	if err := ns.Group.Options(SetGroupStores(st.stores, w)); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return ns, nil
}

// Stores creates a new store slice with all of its new Group and new Website
// pointers. Can return an error when the website or the group cannot be found.
func (st *Storage) Stores() (StoreSlice, error) {
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
func (st *Storage) DefaultStoreID() (int64, error) {
	for _, w := range st.websites {
		if w.IsDefault.Bool && w.IsDefault.Valid {
			g, found := st.group(w.DefaultGroupID)
			if !found {
				return nil, errors.NewNotFoundf("[store] WebsiteID %d DefaultGroupID %d", w.WebsiteID, w.DefaultGroupID)
			}
			return g.DefaultStoreID, nil
		}
	}
	return nil, errors.NewNotFoundf(errStoreDefaultNotFound)
}

// ReInit reloads all websites, groups and stores concurrently from the
// database. On error  all internal slices will be reset to nil.
func (st *Storage) ReInit(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) error {
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
		errc <- errors.Wrap(err, "[store] websites")
	}()

	go func() {
		for i := range st.groups {
			st.groups[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.groups = nil
		_, err := st.groups.SQLSelect(dbrSess, cbs...)
		errc <- errors.Wrap(err, "[store] groups")
	}()

	go func() {
		for i := range st.stores {
			st.stores[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.stores = nil
		_, err := st.stores.SQLSelect(dbrSess, cbs...)
		errc <- errors.Wrap(err, "[store] stores")
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
