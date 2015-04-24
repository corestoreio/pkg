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
	"sync"

	"net/http"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	sContainer struct {
		// w contains the current website select from the scope, internal cache
		w *Website
		// g contains the current group selected from the scope, internal cache
		g *Group
		// s contains the current store selected from the scope, internal cache
		s *Store
	}

	StoreManager struct {
		storage Storager
		sync.RWMutex
		// map key is a hash value
		cacheW map[uint64]*sContainer
		cacheG map[uint64]*sContainer
		cacheS map[uint64]*sContainer
	}
)

var (
	ErrUnsupportedScopeID = errors.New("Unsupported scope id")
)

// NewStoreManager creates a new store manager which handles websites, store groups and stores.
func NewStoreManager(s Storager) *StoreManager {
	return &StoreManager{
		storage: s,
	}
}

// Init @see \Magento\Store\Model\StorageFactory::_reinitStores
func (sm *StoreManager) Init(scopeCode string, scopeType config.ScopeID) (*Store, error) {
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
func (sm *StoreManager) InitByRequest(r *http.Request, scopeType config.ScopeID) {
	var scopeCode string
	// 1. check cookie store
	// 2. check for ___store variable
	if keks, err := r.Cookie(CookieName); err == nil { // if cookie not present ignore it
		scopeCode = keks.Value
	}
	if gs := r.URL.Query().Get(HttpRequestParamStore); gs != "" {
		scopeCode = gs
	}
	_ = scopeCode
	// @todo
	// now init currentStore and cache
	// also delete and re-set a new cookie
}

// IsSingleStoreModeEnabled @todo implement
// @see magento2/app/code/Magento/Store/Model/StoreManager.php uses the config from the database
func (sm *StoreManager) IsSingleStoreModeEnabled(cfg config.ScopeReader) bool {
	return false
}

// IsSingleStoreMode check if Single-Store mode is enabled in configuration.
// This flag only shows that admin does not want to show certain UI components at backend (like store switchers etc)
// if Magento has only one store view but it does not check the store view collection.
func (sm *StoreManager) IsSingleStoreMode() bool {
	return false
}

//
func (sm *StoreManager) HasSingleStore() bool {
	return false
}

// Website returns a website by IDRetriever. If IDRetriever is nil then default website will be returned
func (sm *StoreManager) Website(id IDRetriever, c CodeRetriever) *Website {
	return nil
}

// Websites returns a slice of website buckets
func (sm *StoreManager) Websites() WebsiteSlice {
	return nil
}

// Group returns the group bucket
func (sm *StoreManager) Group(IDRetriever) *Group {
	return nil
}

// Store returns the store view bucket
func (sm *StoreManager) Store(id IDRetriever, c CodeRetriever) *Store {
	return nil
}

// GetDefaultStoreView returns the default store view bucket
func (sm *StoreManager) GetDefaultStoreView() *Store {
	return nil
}

// ReinitStores reloads the website, store group and store view data from the database @todo
func (sm *StoreManager) ReinitStores() error {
	return nil
}

// @todo wrong place for this func here
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
