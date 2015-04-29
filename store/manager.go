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
	"hash"
	"hash/fnv"
	"net/http"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	// Manager implements the Storager interface and on request the StorageMutator interface.
	// Manager uses three internal maps to cache the pointers of Website, Group and Store.
	Manager struct {
		// storage get set of websites, groups and stores and also type assertion to StorageMutator for
		// ReInit and Persisting
		storage Storager
		mu      sync.RWMutex

		// the next six fields are for internal caching
		// map key is a hash value which is generated bei either an int64 or a string
		websiteMap map[uint64]*Website
		groupMap   map[uint64]*Group
		storeMap   map[uint64]*Store
		websites   WebsiteSlice
		groups     GroupSlice
		stores     StoreSlice

		// fnv64 used to calculate the uint64 value of a string, especially website code and store code
		fnv64 hash.Hash64

		// contains the current selected store from init funcs. Cannot be cleared
		currentStore *Store
		// default store cache
		defaultStore *Store
	}
)

var (
	ErrUnsupportedScopeID         = errors.New("Unsupported scope id")
	ErrCurrentStoreNotSet         = errors.New("Current Store is not initialized")
	ErrManagerMutatorNotAvailable = errors.New("Storage Mutator is not implemented")
	ErrHashRetrieverNil           = errors.New("Hash argument is nil")
)

// NewManager creates a new store manager which handles websites, store groups and stores.
func NewManager(s Storager) *Manager {
	return &Manager{
		storage:    s,
		mu:         sync.RWMutex{},
		websiteMap: make(map[uint64]*Website),
		groupMap:   make(map[uint64]*Group),
		storeMap:   make(map[uint64]*Store),
		fnv64:      fnv.New64(),
	}
}

// Init initializes the current store from a scope code and a scope type.
// This func is mainly used when booting the app to set the environment config.
// Also all other calls to any method receiver with nil arguments depends on the current store
// @see \Magento\Store\Model\StorageFactory::_reinitStores
func (sm *Manager) Init(scopeCode string, scopeType config.ScopeID) (*Store, error) {
	switch scopeType {
	case config.ScopeStore:
		// init storage store by store code
		break
	case config.ScopeGroup:
		// init storage store by group id
		break
	case config.ScopeWebsite:
		// init storage store by website code
		break
	default:
		return nil, ErrUnsupportedScopeID
	}
	return nil, nil
}

// Init @see \Magento\Store\Model\StorageFactory::_reinitStores
func (sm *Manager) InitByRequest(r *http.Request, scopeType config.ScopeID) {
	var scopeCode string
	// 1. check cookie store
	// 2. check for ___store variable
	if keks, err := r.Cookie(CookieName); err == nil { // if cookie not present ignore it
		scopeCode = keks.Value
	}
	if gs := r.URL.Query().Get(HTTPRequestParamStore); gs != "" {
		scopeCode = gs
	}
	_ = scopeCode
	// @todo
	// now init currentStore and cache
	// also delete and re-set a new cookie
}

// IsSingleStoreModeEnabled @todo implement
// @see magento2/app/code/Magento/Store/Model/Manager.php uses the config from the database
func (sm *Manager) IsSingleStoreModeEnabled(cfg config.ScopeReader) bool {
	return false
}

// IsSingleStoreMode check if Single-Store mode is enabled in configuration.
// This flag only shows that admin does not want to show certain UI components at backend (like store switchers etc)
// if Magento has only one store view but it does not check the store view collection.
func (sm *Manager) IsSingleStoreMode() bool {
	return false
}

//
func (sm *Manager) HasSingleStore() bool {
	return false
}

// Website returns the cached Website pointer from an ID or code including all of its
// groups and all related stores. It panics when the integrity is incorrect.
// If ID and code are available then the non-empty code has precedence.
// If no argument has been supplied then the Website of the internal current store
// will be returned. If more than one argument has been provided it returns an error.
func (sm *Manager) Website(r ...Retriever) (*Website, error) {
	notR := notRetriever(r...)
	switch {
	case notR && sm.currentStore == nil:
		return nil, ErrCurrentStoreNotSet
	case notR && sm.currentStore != nil:
		return sm.currentStore.Website(), nil
	}

	key, err := sm.hash(r[0])
	if err != nil {
		return nil, err
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	if w, ok := sm.websiteMap[key]; ok && w != nil {
		return w, nil
	}

	w, err := sm.storage.Website(r[0])
	sm.websiteMap[key] = w
	return sm.websiteMap[key], errgo.Mask(err)
}

// Websites returns a cached slice containing all pointers to Websites with its associated
// groups and stores. It panics when the integrity is incorrect.
func (sm *Manager) Websites() (WebsiteSlice, error) {
	if sm.websites != nil {
		return sm.websites, nil
	}
	var err error
	sm.websites, err = sm.storage.Websites()
	return sm.websites, err
}

// Group returns a cached Group which contains all related stores and its website.
// Only the argument ID is supported.
// If no argument has been supplied then the Group of the internal current store
// will be returned. If more than one argument has been provided it returns an error.
func (sm *Manager) Group(r ...Retriever) (*Group, error) {
	notR := notRetriever(r...)
	switch {
	case notR && sm.currentStore == nil:
		return nil, ErrCurrentStoreNotSet
	case notR && sm.currentStore != nil:
		return sm.currentStore.Group(), nil
	}

	key, err := sm.hash(r[0])
	if err != nil {
		return nil, err
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	if g, ok := sm.groupMap[key]; ok && g != nil {
		return g, nil
	}

	g, err := sm.storage.Group(r[0])
	sm.groupMap[key] = g
	return sm.groupMap[key], errgo.Mask(err)
}

// Groups returns a cached slice containing all pointers to Groups with its associated
// stores and websites. It panics when the integrity is incorrect.
func (sm *Manager) Groups() (GroupSlice, error) {
	if sm.groups != nil {
		return sm.groups, nil
	}
	var err error
	sm.groups, err = sm.storage.Groups()
	return sm.groups, err
}

// Store returns the cached Store view containing its group and its website.
// If ID and code are available then the non-empty code has precedence.
// If no argument has been supplied then the current store
// will be returned. If more than one argument has been provided it returns an error.
func (sm *Manager) Store(r ...Retriever) (*Store, error) {
	notR := notRetriever(r...)
	switch {
	case notR && sm.currentStore == nil:
		return nil, ErrCurrentStoreNotSet
	case notR && sm.currentStore != nil:
		return sm.currentStore, nil
	}

	key, err := sm.hash(r[0])
	if err != nil {
		return nil, err
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	if s, ok := sm.storeMap[key]; ok && s != nil {
		return s, nil
	}

	s, err := sm.storage.Store(r[0])
	sm.storeMap[key] = s
	return sm.storeMap[key], errgo.Mask(err)
}

// Stores returns a cached Store slice. Can return an error when the website or
// the group cannot be found.
func (sm *Manager) Stores() (StoreSlice, error) {
	if sm.stores != nil {
		return sm.stores, nil
	}
	var err error
	sm.stores, err = sm.storage.Stores()
	return sm.stores, err
}

// DefaultStoreView returns the default store view.
func (sm *Manager) DefaultStoreView() (*Store, error) {
	if sm.defaultStore != nil {
		return sm.defaultStore, nil
	}
	var err error
	sm.defaultStore, err = sm.storage.DefaultStoreView()
	return sm.defaultStore, err
}

// ReInit reloads the website, store group and store view data from the database @todo
func (sm *Manager) ReInit(dbrSess dbr.SessionRunner) error {
	if mut, ok := sm.storage.(StorageMutator); ok {
		defer sm.ClearCache() // hmmm .... defer ...
		return mut.ReInit(dbrSess)
	}
	return ErrManagerMutatorNotAvailable
}

// hash generates the key for the map from either an id int64 or a code string.
// If both arguments are nil it returns 0 which is default for website, group or store.
func (sm *Manager) hash(r Retriever) (uint64, error) {
	uz := uint64(0)
	if r == nil {
		return uz, ErrHashRetrieverNil
	}
	if c, ok := r.(CodeRetriever); ok && c.Code() != "" {
		sm.fnv64.Reset()
		_, err := sm.fnv64.Write([]byte(c.Code()))
		if err != nil {
			return uz, errgo.Mask(err)
		}
		return sm.fnv64.Sum64(), nil
	}
	return uint64(r.ID()), nil
}

// ClearCache resets the internal caches which stores the pointers to a Website, Group or Store and
// all related slices. Please use with caution. ReInit() also uses this method.
func (sm *Manager) ClearCache() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.websiteMap) > 0 {
		for k := range sm.websiteMap {
			delete(sm.websiteMap, k)
		}
	}
	if len(sm.groupMap) > 0 {
		for k := range sm.groupMap {
			delete(sm.groupMap, k)
		}
	}
	if len(sm.storeMap) > 0 {
		for k := range sm.storeMap {
			delete(sm.storeMap, k)
		}
	}
	sm.websites = nil
	sm.groups = nil
	sm.stores = nil
	sm.defaultStore = nil
	// do not clear currentStore as this one depends on the init funcs
}

// IsCacheEmpty returns true if the internal cache is empty.
func (sm *Manager) IsCacheEmpty() bool {
	return len(sm.websiteMap) == 0 && len(sm.groupMap) == 0 && len(sm.storeMap) == 0 &&
		sm.websites == nil && sm.groups == nil && sm.stores == nil && sm.defaultStore == nil
}

// notRetriever checks if variadic Retriever is nil or has more than two entries.
func notRetriever(r ...Retriever) bool {
	lr := len(r)
	return r == nil || (lr == 1 && r[0] == nil) || lr > 1
}

// loadSlice internal global helper func to execute a SQL select. @todo refactor and remove dependency of GetTableS...
func loadSlice(dbrSess dbr.SessionRunner, table csdb.Index, dest interface{}, cbs ...csdb.DbrSelectCb) (int, error) {
	ts, err := GetTableStructure(table)
	if err != nil {
		return 0, errgo.Mask(err)
	}

	sb, err := ts.Select(dbrSess)
	if err != nil {
		return 0, errgo.Mask(err)
	}

	for _, cb := range cbs {
		sb = cb(sb)
	}
	return sb.LoadStructs(dest)
}
