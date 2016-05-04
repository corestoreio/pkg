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

type (
	// Storager implements the requirements to get new websites, groups and store views.
	// This interface is used in the StoreService
	Storager interface {
		// Website creates a new Website pointer from an ID or code including all of its
		// groups and all related stores. It panics when the integrity is incorrect.
		// If ID and code are available then the non-empty code has precedence.
		Website(scope.WebsiteIDer) (*Website, error)
		// Websites creates a slice containing all pointers to Websites with its associated
		// groups and stores. It panics when the integrity is incorrect.
		Websites() (WebsiteSlice, error)
		// Group creates a new Group which contains all related stores and its website.
		// Only the argument ID can be used to get a specific Group.
		Group(scope.GroupIDer) (*Group, error)
		// Groups creates a slice containing all pointers to Groups with its associated
		// stores and websites. It panics when the integrity is incorrect.
		Groups() (GroupSlice, error)
		// Store creates a new Store containing its group and its website.
		// If ID and code are available then the non-empty code has precedence.
		Store(scope.StoreIDer) (*Store, error)
		// Stores creates a new store slice. Can return an error when the website or
		// the group cannot be found.
		Stores() (StoreSlice, error)
		// DefaultStoreView traverses through the websites to find the default website and gets
		// the default group which has the default store id assigned to. Only one website can be the default one.
		DefaultStoreView() (*Store, error)
		// ReInit reloads the websites, groups and stores from the database.
		ReInit(dbr.SessionRunner, ...dbr.SelectCb) error
	}

	// Storage contains a mutex and the raw slices from the database.
	storage struct {
		// cr parent config service. can only be set once.
		cr       config.Getter
		mu       sync.RWMutex
		websites TableWebsiteSlice
		groups   TableGroupSlice
		stores   TableStoreSlice
		// optionError use by functional option arguments to indicate that one
		// option has triggered an error and hence the other can options can
		// skip their process.
		optionError error
	}
)

// check if interface has been implemented
var _ Storager = (*storage)(nil)

// NewStorage creates a new storage object which handles the raw data from the
// three database tables for website, group and store. You can either provide
// the raw data separately for each type or pass an option to load it from
// the database.
//		sto, err = store.NewStorage(
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
// 		sto, err = store.NewStorage( store.WithDatabaseInit(dbrSession) )
func NewStorage(opts ...StorageOption) (Storager, error) {
	s := &storage{}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	if s.optionError != nil {
		return nil, s.optionError
	}
	return s, nil
}

// MustNewStorage same as NewStorage but panics on error.
func MustNewStorage(opts ...StorageOption) Storager {
	s, err := NewStorage(opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// website returns a TableWebsite by using either id or code to find it. If id and code are
// available then the non-empty code has precedence.
func (st *storage) website(r scope.WebsiteIDer) (*TableWebsite, bool) {
	if r == nil {
		return nil, false
	}
	if c, ok := r.(scope.WebsiteCoder); ok && c.WebsiteCode() != "" {
		return st.websites.FindByCode(c.WebsiteCode())
	}
	return st.websites.FindByWebsiteID(r.WebsiteID())
}

// Website creates a new Website according to the interface definition.
func (st *storage) Website(r scope.WebsiteIDer) (*Website, error) {
	w, found := st.website(r)
	if !found {
		return nil, errors.NewNotFoundf("[store] WebsiteIDer %v", r)
	}
	return NewWebsite(w, SetWebsiteConfig(st.cr), SetWebsiteGroupsStores(st.groups, st.stores))
}

// Websites creates a slice of Website pointers according to the interface definition.
func (st *storage) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(st.websites), len(st.websites))
	for i, w := range st.websites {
		var err error
		websites[i], err = NewWebsite(w, SetWebsiteConfig(st.cr), SetWebsiteGroupsStores(st.groups, st.stores))
		if err != nil {
			return nil, errors.Wrapf(err, "[store] Storage.Websites. WebsiteID: %d", w.WebsiteID)
		}
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument.
func (st *storage) group(r scope.GroupIDer) (*TableGroup, bool) {
	if r == nil {
		return nil, false
	}
	return st.groups.FindByGroupID(r.GroupID())
}

// Group creates a new Group which contains all related stores and its website according to the
// interface definition.
func (st *storage) Group(id scope.GroupIDer) (*Group, error) {
	g, found := st.group(id)
	if !found {
		return nil, errors.NewNotFoundf("[store] Group %v", id)
	}

	w, found := st.website(scope.MockID(g.WebsiteID))
	if !found {
		return nil, errors.NewNotFoundf("[store] Website. WebsiteID %d GroupID %v", g.WebsiteID, id)
	}
	return NewGroup(g, SetGroupConfig(st.cr), SetGroupWebsite(w), SetGroupStores(st.stores, nil))
}

// Groups creates a new group slice containing its website all related stores.
// May panic when a website pointer is nil.
func (st *storage) Groups() (GroupSlice, error) {
	groups := make(GroupSlice, len(st.groups), len(st.groups))
	for i, g := range st.groups {
		w, found := st.website(scope.MockID(g.WebsiteID))
		if !found {
			return nil, errors.NewNotFoundf("[store] WebsiteID %d", g.WebsiteID)
		}
		var err error
		groups[i], err = NewGroup(g, SetGroupConfig(st.cr), SetGroupWebsite(w), SetGroupStores(st.stores, nil))
		if err != nil {
			return nil, errors.Wrapf(err, "[store] GroupID %d WebsiteID %d", g.GroupID, g.WebsiteID)
		}
	}
	return groups, nil
}

// store returns a TableStore by an id or code.
// The non-empty code has precedence if available.
func (st *storage) store(r scope.StoreIDer) (*TableStore, bool) {
	if r == nil {
		return nil, false
	}
	if c, ok := r.(scope.StoreCoder); ok && c.StoreCode() != "" {
		return st.stores.FindByCode(c.StoreCode())
	}
	return st.stores.FindByStoreID(r.StoreID())
}

// Store creates a new Store which contains the store, its group and website
// according to the interface definition.
func (st *storage) Store(r scope.StoreIDer) (*Store, error) {
	s, found := st.store(r)
	if !found {
		return nil, errors.NewNotFoundf("[store] Store: %v", r)
	}
	w, found := st.website(scope.MockID(s.WebsiteID))
	if !found {
		return nil, errors.NewNotFoundf("[store] WebsiteID: %d", s.WebsiteID)
	}
	g, found := st.group(scope.MockID(s.GroupID))
	if !found {
		return nil, errors.NewNotFoundf("[store] GroupID: %d", s.GroupID)
	}
	ns, err := NewStore(s, w, g, WithStoreConfig(st.cr))
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

// Stores creates a new store slice. Can return an error when the website or
// the group cannot be found.
func (st *storage) Stores() (StoreSlice, error) {
	stores := make(StoreSlice, len(st.stores), len(st.stores))
	for i, s := range st.stores {
		var err error
		if stores[i], err = st.Store(scope.MockID(s.StoreID)); err != nil {
			return nil, errors.Wrapf(err, "[store] StoreID %d", s.StoreID)
		}
	}
	return stores, nil
}

// DefaultStoreView traverses through the websites to find the default website and gets
// the default group which has the default store id assigned to. Only one website can be the default one.
func (st *storage) DefaultStoreView() (*Store, error) {
	for _, w := range st.websites {
		if w.IsDefault.Bool && w.IsDefault.Valid {
			g, found := st.group(scope.MockID(w.DefaultGroupID))
			if !found {
				return nil, errors.NewNotFoundf("[store] WebsiteID %d DefaultGroupID %d", w.WebsiteID, w.DefaultGroupID)
			}
			return st.Store(scope.MockID(g.DefaultStoreID))
		}
	}
	return nil, errors.NewNotFoundf(errStoreNotFound)
}

// ReInit reloads all websites, groups and stores concurrently from the database. If GOMAXPROCS
// is set to > 1 then in parallel. Returns an error with location or nil. If an error occurs
// then all internal slices will be reset.
func (st *storage) ReInit(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) error {
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
