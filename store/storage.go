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

package store

import (
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/juju/errgo"
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

	// Storage contains a mutex and the raw slices from the database. @todo maybe make private?
	Storage struct {
		lastErrors []error
		cr         config.Getter
		mu         sync.RWMutex
		websites   TableWebsiteSlice
		groups     TableGroupSlice
		stores     TableStoreSlice
	}

	// StorageOption option func for NewStorage()
	StorageOption func(*Storage)
)

// check if interface has been implemented
var _ Storager = (*Storage)(nil)

// SetStorageWebsites adds the TableWebsiteSlice to the Storage. By default, the slice is nil.
func SetStorageWebsites(tws ...*TableWebsite) StorageOption {
	return func(s *Storage) { s.websites = TableWebsiteSlice(tws) }
}

// SetStorageGroups adds the TableGroupSlice to the Storage. By default, the slice is nil.
func SetStorageGroups(tgs ...*TableGroup) StorageOption {
	return func(s *Storage) { s.groups = TableGroupSlice(tgs) }
}

// SetStorageStores adds the TableStoreSlice to the Storage. By default, the slice is nil.
func SetStorageStores(tss ...*TableStore) StorageOption {
	return func(s *Storage) { s.stores = TableStoreSlice(tss) }
}

// SetStorageConfig sets the configuration Reader. Optional.
// Default reader is config.DefaultManager
func SetStorageConfig(cr config.Getter) StorageOption {
	return func(s *Storage) { s.cr = cr }
}

// WithDatabaseInit triggers the ReInit function to load the data from the
// database.
func WithDatabaseInit(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) StorageOption {
	return func(s *Storage) {
		if err := s.ReInit(dbrSess, cbs...); err != nil {
			s.lastErrors = append(s.lastErrors, err)
		}
	}
}

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
func NewStorage(opts ...StorageOption) (*Storage, error) {
	s := &Storage{
		cr: config.DefaultService,
		mu: sync.RWMutex{},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	if len(s.lastErrors) > 0 {
		return nil, s
	}
	return s, nil
}

// MustNewStorage same as NewStorage but panics on error.
func MustNewStorage(opts ...StorageOption) *Storage {
	s, err := NewStorage(opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// Error returns an error string
func (st *Storage) Error() string {
	return util.Errors(st.lastErrors...)
}

// website returns a TableWebsite by using either id or code to find it. If id and code are
// available then the non-empty code has precedence.
func (st *Storage) website(r scope.WebsiteIDer) (*TableWebsite, error) {
	if r == nil {
		return nil, ErrWebsiteNotFound
	}
	if c, ok := r.(scope.WebsiteCoder); ok && c.WebsiteCode() != "" {
		return st.websites.FindByCode(c.WebsiteCode())
	}
	return st.websites.FindByWebsiteID(r.WebsiteID())
}

// Website creates a new Website according to the interface definition.
func (st *Storage) Website(r scope.WebsiteIDer) (*Website, error) {
	w, err := st.website(r)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return NewWebsite(w, SetWebsiteConfig(st.cr), SetWebsiteGroupsStores(st.groups, st.stores))
}

// Websites creates a slice of Website pointers according to the interface definition.
func (st *Storage) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(st.websites), len(st.websites))
	for i, w := range st.websites {
		var err error
		websites[i], err = NewWebsite(w, SetWebsiteConfig(st.cr), SetWebsiteGroupsStores(st.groups, st.stores))
		if err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("store.Storage.Websites.NewWebsite", "err", err, "w", w, "websites", st.websites)
			}
			return nil, errgo.Mask(err)
		}
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument. If no argument or more than
// one has been supplied it returns an error.
func (st *Storage) group(r scope.GroupIDer) (*TableGroup, error) {
	if r == nil {
		return nil, ErrGroupNotFound
	}
	return st.groups.FindByGroupID(r.GroupID())
}

// Group creates a new Group which contains all related stores and its website according to the
// interface definition.
func (st *Storage) Group(id scope.GroupIDer) (*Group, error) {
	g, err := st.group(id)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	w, err := st.website(scope.MockID(g.WebsiteID))
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("store.Storage.Group.website", "err", err, "websiteID", g.WebsiteID, "groupID", id.GroupID())
		}
		return nil, errgo.Mask(err)
	}
	return NewGroup(g, SetGroupConfig(st.cr), SetGroupWebsite(w), SetGroupStores(st.stores, nil))
}

// Groups creates a new group slice containing its website all related stores.
// May panic when a website pointer is nil.
func (st *Storage) Groups() (GroupSlice, error) {
	groups := make(GroupSlice, len(st.groups), len(st.groups))
	for i, g := range st.groups {
		w, err := st.website(scope.MockID(g.WebsiteID))
		if err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("store.Storage.Groups.website", "err", err, "g", g, "websiteID", g.WebsiteID)
			}
			return nil, errgo.Mask(err)
		}

		groups[i], err = NewGroup(g, SetGroupConfig(st.cr), SetGroupWebsite(w), SetGroupStores(st.stores, nil))
		if err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("store.Storage.Groups.NewGroup", "err", err, "g", g, "websiteID", g.WebsiteID)
			}
			return nil, errgo.Mask(err)
		}
	}
	return groups, nil
}

// store returns a TableStore by an id or code.
// The non-empty code has precedence if available.
func (st *Storage) store(r scope.StoreIDer) (*TableStore, error) {
	if r == nil {
		return nil, ErrStoreNotFound
	}
	if c, ok := r.(scope.StoreCoder); ok && c.StoreCode() != "" {
		return st.stores.FindByCode(c.StoreCode())
	}
	return st.stores.FindByStoreID(r.StoreID())
}

// Store creates a new Store which contains the the store, its group and website
// according to the interface definition.
func (st *Storage) Store(r scope.StoreIDer) (*Store, error) {
	s, err := st.store(r)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	w, err := st.website(scope.MockID(s.WebsiteID))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	g, err := st.group(scope.MockID(s.GroupID))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	ns, err := NewStore(s, w, g, SetStoreConfig(st.cr))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	if _, err := ns.Website.ApplyOptions(SetWebsiteGroupsStores(st.groups, st.stores)); err != nil {
		return nil, errgo.Mask(err)
	}
	if _, err := ns.Group.ApplyOptions(SetGroupStores(st.stores, w)); err != nil {
		return nil, errgo.Mask(err)
	}
	return ns, nil
}

// Stores creates a new store slice. Can return an error when the website or
// the group cannot be found.
func (st *Storage) Stores() (StoreSlice, error) {
	stores := make(StoreSlice, len(st.stores), len(st.stores))
	for i, s := range st.stores {
		var err error
		if stores[i], err = st.Store(scope.MockID(s.StoreID)); err != nil {
			return nil, errgo.Mask(err)
		}
	}
	return stores, nil
}

// DefaultStoreView traverses through the websites to find the default website and gets
// the default group which has the default store id assigned to. Only one website can be the default one.
func (st *Storage) DefaultStoreView() (*Store, error) {
	for _, website := range st.websites {
		if website.IsDefault.Bool && website.IsDefault.Valid {
			g, err := st.group(scope.MockID(website.DefaultGroupID))
			if err != nil {
				return nil, errgo.Mask(err)
			}
			return st.Store(scope.MockID(g.DefaultStoreID))
		}
	}
	return nil, ErrStoreNotFound
}

// ReInit reloads all websites, groups and stores concurrently from the database. If GOMAXPROCS
// is set to > 1 then in parallel. Returns an error with location or nil. If an error occurs
// then all internal slices will be reset.
func (st *Storage) ReInit(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	if dbrSess == nil {
		return errgo.New("dbr.SessionRunner is nil")
	}

	errc := make(chan error)
	defer close(errc)
	// not sure about those three go
	go func() {
		for i := range st.websites {
			st.websites[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.websites = nil
		_, err := st.websites.SQLSelect(dbrSess, cbs...)
		errc <- errgo.Mask(err)
	}()

	go func() {
		for i := range st.groups {
			st.groups[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.groups = nil
		_, err := st.groups.SQLSelect(dbrSess, cbs...)
		errc <- errgo.Mask(err)
	}()

	go func() {
		for i := range st.stores {
			st.stores[i] = nil // I'm not quite sure if that is needed to clear the pointers
		}
		st.stores = nil
		_, err := st.stores.SQLSelect(dbrSess, cbs...)
		errc <- errgo.Mask(err)
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
