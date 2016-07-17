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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"net/http"
	"sync"
	"sync/atomic"
)

// Requester knows how to retrieve a specific store depending on the
// scope.Option.
type Requester interface {
	// RequestedStore figures out the default active store for a storeID. It
	// takes into account that Getter is bound to a specific scope.Scope. It
	// also prevents running a store from another website or store group, if
	// website or store group was specified explicitly. RequestedStore returns
	// either an error or the store. RequestedStore will be mostly used within
	// an HTTP request.
	RequestedStore(int64) (activeStore *Store, err error)
}

// IDbyCode returns for a website code or store code the id. Group scope is not
// supported because the group table does not contain a code string column. A
// not-supported error behaviour gets returned if an invalid scope has been
// provided. Default scope returns always 0.
type CodeToIDMapper interface {
	IDbyCode(scp scope.Scope, code string) (id int64, err error)
}

// depends on the run mode
type AvailabilityChecker interface {
	AllowedStoreIds(scope.Hash) []int64
	DefaultStoreId(scope.Hash) int64
}

// RunModeFunc initialized the runmode and and scope/ID for the current request.
// app/code/Magento/Store/App/FrontController/Plugin/DefaultStore.php
// The returned Hash will be used in interface DefaultAllower
type RunModeFunc func(config.Getter, *http.Request) scope.Hash

type (
	// Service represents type which handles the underlying storage and takes
	// care of the default stores. A Service is bound a specific scope.Scope.
	// Depending on the scope it is possible or not to switch stores. A Service
	// contains also a config.Getter which gets passed to the scope of a
	// Store(), Group() or Website() so that you always have the possibility to
	// access a scoped based configuration value. This Service uses three
	// internal maps to cache the pointers of Website, Group and Store.
	Service struct {

		// backend communicates with the database in reading mode and creates
		// new store, group and website pointers. If nil, panics.
		backend *Storage
		// defaultStore someone must be always the default guy. Handled via atomic
		// package.
		defaultStoreID int64
		// mu protects the following fields
		mu sync.RWMutex
		// in general these caches can be optimized
		websites WebsiteSlice
		groups   GroupSlice
		stores   StoreSlice

		// int64 key identifies a website, group or store
		cacheWebsite map[int64]*Website
		cacheGroup   map[int64]*Group
		cacheStore   map[int64]*Store
	}
)

// NewService creates a new store Service which handles websites, groups and
// stores. A Service can only act on a certain scope (MAGE_RUN_TYPE) and scope
// ID (MAGE_RUN_CODE). Default scope.Scope is always the scope.WebsiteID
// constant. This function is mainly used when booting the app to set the
// environment configuration Also all other calls to any method receiver with
// nil arguments depends on the internal appStore which reflects the default
// store ID.
func NewService(st *Storage) (*Service, error) {
	srv := &Service{
		backend:        st,
		defaultStoreID: -1,
	}
	if err := srv.ApplyStorage(st); err != nil {
		return errors.Wrap(err, "[store] NewService.ApplyStorage")
	}
	return srv, nil
}

// MustNewService same as NewService, but panics on error.
func MustNewService(st *Storage) *Service {
	m, err := NewService(st)
	if err != nil {
		panic(err)
	}
	return m
}

func (s *Service) ApplyStorage(st *Storage) error {
	if s == nil {
		s = new(Service)
		s.defaultStoreID = -1
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.backend = st
	s.cacheWebsite = make(map[int64]*Website)
	s.cacheGroup = make(map[int64]*Group)
	s.cacheStore = make(map[int64]*Store)

	ws, err := st.Websites()
	if err != nil {
		return errors.Wrap(err, "[store] NewService.Websites")
	}
	s.websites = ws
	ws.Each(func(w *Website) {
		s.cacheWebsite[w.Data.WebsiteID] = w
	})

	gs, err := st.Groups()
	if err != nil {
		return errors.Wrap(err, "[store] NewService.Groups")
	}
	s.groups = gs
	gs.Each(func(g *Group) {
		s.cacheGroup[g.Data.GroupID] = g
	})

	ss, err := st.Stores()
	if err != nil {
		return errors.Wrap(err, "[store] NewService.Stores")
	}
	s.stores = ss
	ss.Each(func(str *Store) {
		s.cacheStore[str.Data.StoreID] = str
	})
	return nil
}

func (s *Service) AllowedStoreIds(h scope.Hash) ([]int64, error) {
	scp, id := h.Unpack()
	switch scp {
	case scope.Store:
		var ids = make([]int64, 0, len(s.cacheStore))
		for _, st := range s.cacheStore {
			if st.Data.IsActive {
				ids = append(ids, st.Data.StoreID)
			}
		}
		return ids
	case scope.Group:
	case scope.Website:
	}
	return nil, errors.NewNotSupportedf("[store] Unknown Scope: %q", scp)
}

// findDefaultStoreByScope tries to detect the default store by a given scope option.
// Precedence of detection by passed scope.Option: 1. Store 2. Group 3. Website
func (s *Service) DefaultStoreID(h scope.Hash) (int64, error) {
	scp, id := h.Unpack()
	switch scp {
	case scope.Store:
		st, err := s.Store(id)
		if err != nil {
			return 0, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		if !st.Data.IsActive {
			return 0, errors.NewNotValidf("[store] DefaultStoreID %s the store ID %d is not active", h, st.StoreID())
		}
		return st.StoreID(), nil

	case scope.Group:
		g, err := s.Group(id)
		if err != nil {
			return 0, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		st, err := s.Store(g.StoreID())
		if err != nil {
			return 0, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		if !st.Data.IsActive {
			return 0, errors.NewNotValidf("[store] DefaultStoreID %s the store ID %d is not active", h, st.StoreID())
		}
		return st.StoreID(), nil

	case scope.Website:
		w, err := s.Website(id)
		if err != nil {
			return nil, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		st, err := w.DefaultStore()
		if err != nil {
			return nil, errors.Wrapf(err, "[store] DefaultStoreID Scope %s ID %d", scp, id)
		}
		if !st.Data.IsActive {
			return 0, errors.NewNotValidf("[store] DefaultStoreID %s the store ID %d is not active", h, st.StoreID())
		}
		return st, nil
	}
	return nil, errors.NewNotSupportedf("[store] Unknown Scope: %q", scp)
}

// IDbyCode returns for a website code or store code the id. It iterates over
// the internal cache maps. Group scope is not supported because the group table
// does not contain a code string column. A not-supported error behaviour gets
// returned if an invalid scope has been provided. Default scope returns always
// 0. Implements interface CodeToIDMapper.
func (s *Service) IDbyCode(scp scope.Scope, code string) (int64, error) {
	if code == "" {
		return 0, errors.NewEmptyf("[store] Service IDByCode: Code canot be empty.")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	// todo maybe add map cache
	switch scp {
	case scope.Store:
		for _, cs := range s.cacheStore {
			if cs.StoreCode() == code {
				return cs.StoreID(), nil
			}
		}
		return 0, errors.NewNotFoundf("[store] Code %q not found in %s", code, scp)
	case scope.Website:
		for _, cs := range s.cacheWebsite {
			if cs.WebsiteCode() == code {
				return cs.WebsiteID(), nil
			}
		}
	case scope.Default:
		return 0, nil
	}
	return 0, errors.NewNotSupportedf("[store] Scope %q not supported", scp)
}

// RequestedStore see interface description Getter.RequestedStore.
// Error behaviour: Unauthorized, NotFound, NotSupported
func (s *Service) RequestedStore(id int64) (activeStore *Store, err error) {

	//activeStore, err = sm.findDefaultStoreByScope(scope.Store, id)
	//if err != nil {
	//	return nil, errors.Wrap(err, "[store] findDefaultStoreByScope")
	//}
	//
	////	activeStore, err = sm.newActiveStore(activeStore) // this is the active store from a request.
	//// todo rethink here if we really need a newActiveStore
	//// newActiveStore creates a new Store, Website and Group pointers !!!
	////	if activeStore == nil || err != nil {
	////		// store is not active so ignore
	////		return nil, err
	////	}
	//
	//if false == activeStore.Data.IsActive {
	//	return nil, errors.NewUnauthorizedf(errStoreNotActive)
	//}
	//
	//allowStoreChange := false
	//switch sm.boundToScope {
	//case scope.Store:
	//	allowStoreChange = true
	//	break
	//case scope.Group:
	//	allowStoreChange = activeStore.Data.GroupID == sm.appStore.Data.GroupID
	//	break
	//case scope.Website:
	//	allowStoreChange = activeStore.Data.WebsiteID == sm.appStore.Data.WebsiteID
	//	break
	//}
	//
	//if allowStoreChange {
	//	return activeStore, nil
	//}
	return nil, errors.NewUnauthorizedf(errStoreChangeNotAllowed)
}

// IsSingleStoreMode check if Single-Store mode is enabled in configuration and from Store count < 3.
// This flag only shows that admin does not want to show certain UI components at backend (like store switchers etc)
// if Magento has only one store view but it does not check the store view collection.
//func (sm *Service) IsSingleStoreMode() bool {
//	// refactor and remove dependency to backend.Backend
//	return sm.HasSingleStore() // && backend.Backend.GeneralSingleStoreModeEnabled.Get(sm.cr.NewScoped(0, 0)) // default scope
//}
//
//// HasSingleStore checks if we only have one store view besides the admin store view.
//// Mostly used in models to the set store id and in blocks to not display the store switch.
//func (sm *Service) HasSingleStore() bool {
//	ss, err := sm.Stores()
//	if err != nil {
//		return false
//	}
//	// that means: index 0 is admin store and always present plus one more store view.
//	return ss.Len() < 3
//}

// Website returns the cached Website pointer from an ID or code including all of its
// groups and all related stores. It panics when the integrity is incorrect.
// If ID and code are available then the non-empty code has precedence.
// If no argument has been supplied then the Website of the internal appStore
// will be returned. If more than one argument has been provided it returns an error.
func (s *Service) Website(id int64) (*Website, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if cs, ok := s.cacheWebsite[id]; ok {
		return cs, nil
	}
	return nil, errors.NewNotFoundf("[store] Cannot find Website ID %d", id)
}

// Websites returns a cached slice containing all pointers to Websites with its associated
// groups and stores. It panics when the integrity is incorrect.
func (s *Service) Websites() WebsiteSlice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ws := make(WebsiteSlice, len(s.cacheWebsite))
	i := 0
	for _, cw := range s.cacheWebsite {
		ws[i] = cw
		i++
	}
	return ws
}

// Group returns a cached Group which contains all related stores and its website.
// Only the argument ID is supported.
// If no argument has been supplied then the Group of the internal appStore
// will be returned. If more than one argument has been provided it returns an error.
func (s *Service) Group(id int64) (*Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if cg, ok := s.cacheGroup[id]; ok {
		return cg, nil
	}
	return nil, errors.NewNotFoundf("[store] Cannot find Group ID %d", id)
}

// Groups returns a cached slice containing all pointers to Groups with its associated
// stores and websites. It panics when the integrity is incorrect.
func (s *Service) Groups(h scope.Hash) GroupSlice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	gs := make(GroupSlice, len(s.cacheGroup))
	i := 0
	for _, cg := range s.cacheGroup {
		gs[i] = cg
		i++
	}
	return gs
}

// Store returns the cached Store view containing its group and its website.
// If ID and code are available then the non-empty code has precedence.
// If no argument has been supplied then the appStore
// will be returned. If more than one argument has been provided it returns an error.
func (s *Service) Store(id int64) (*Store, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if cs, ok := s.cacheStore[id]; ok {
		return cs, nil
	}
	return nil, errors.NewNotFoundf("[store] Cannot find Store ID %d", id)
}

// Stores returns a cached Store slice. Can return an error when the website or
// the group cannot be found.
func (s *Service) Stores(h scope.Hash) StoreSlice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ss := make(StoreSlice, len(s.cacheStore))
	i := 0
	for _, cs := range s.cacheStore {
		ss[i] = cs
		i++
	}
	return ss
}

// DefaultStoreView returns the default store view, independent of the
// applied scope.Option while creating the service.
func (s *Service) DefaultStoreView() (*Store, error) {
	if s.defaultStoreID >= 0 {
		s.mu.RLock()
		defer s.mu.RUnlock() // bug
		if cs, ok := s.cacheStore[atomic.LoadInt64(&s.defaultStoreID)]; ok {
			return cs, nil
		}
	}

	id, err := s.backend.DefaultStoreID()
	if err != nil {
		return nil, errors.Wrap(err, "[store] Service.storage.DefaultStoreView")
	}
	atomic.StoreInt64(&s.defaultStoreID, id)
	return s.Store(id)
}

// ReInit reloads the website, store group and store view data from the database.
// After reloading internal cache will be cleared if there are no errors.
func (sm *Service) ReInit(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) error {

	if err := sm.backend.ReInit(dbrSess, cbs...); err != nil {
		return errors.Wrap(err, "[store] ReInit.Backend")
	}
	sm.ClearCache()
	if err := sm.ApplyStorage(sm.backend); err != nil {
		return errors.Wrap(err, "[store] ReInit.ApplyStorage")
	}
	return nil
}

// ClearCache resets the internal caches which stores the pointers to Websites,
// Groups or Stores. The ReInit() also uses this method to clear caches before
// the Storage gets reloaded.
func (sm *Service) ClearCache() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.cacheWebsite) > 0 {
		for k := range sm.cacheWebsite {
			delete(sm.cacheWebsite, k)
		}
	}
	if len(sm.cacheGroup) > 0 {
		for k := range sm.cacheGroup {
			delete(sm.cacheGroup, k)
		}
	}
	if len(sm.cacheStore) > 0 {
		for k := range sm.cacheStore {
			delete(sm.cacheStore, k)
		}
	}
	sm.defaultStoreID = -1
}

// IsCacheEmpty returns true if the internal cache is empty.
func (sm *Service) IsCacheEmpty() bool {
	return len(sm.cacheWebsite) == 0 && len(sm.cacheGroup) == 0 && len(sm.cacheStore) == 0 &&
		sm.defaultStoreID == -1
}
