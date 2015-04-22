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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	StoreManager struct {
		// scopeCode is either a store_code or a website_code
		scopeCode string
		// scopeType is either store, group or website. If group casts scopeCode to int.
		// Default scope must not be used.
		scopeType config.ScopeID
		// s contains the current store selected from the scope, internal cache
		s *Store
		// g contains the current group selected from the scope, internal cache
		g *Group
		// w contains the current website select from the scope, internal cache
		w *Website
	}
)

// NewStoreManager creates a new store manager which handles websites, store groups and stores.
func NewStoreManager() *StoreManager {
	return nil
	//	return &StoreManager{
	//		g: g.SetStores(s).SetWebSite(w),
	//		s: s,
	//		w: w.SetGroups(g).SetStores(s),
	//	}
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

// Website returns the website bucket
func (sm *StoreManager) Website() *Website {
	return sm.w
}

// Websites returns a slice of website buckets
func (sm *StoreManager) Websites() []*Website {
	return nil
}

// Group returns the group bucket
func (sm *StoreManager) Group() *Group {
	return sm.g
}

// Store returns the store view bucket
func (sm *StoreManager) Store() *Store {
	return sm.s
}

// GetDefaultStoreView returns the default store view bucket
func (sm *StoreManager) GetDefaultStoreView() *Store {
	return sm.s
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
