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
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

type (
	StoreManager struct {
		s *StoreBucket
		g *GroupBucket
		w *WebsiteBucket
	}
)

// NewStoreManager creaets a new store manager which handles websites, store groups and stores.
func NewStoreManager(tws TableWebsiteSlice, tgs TableGroupSlice, tss TableStoreSlice) *StoreManager {
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
func (sm *StoreManager) Website() *WebsiteBucket {
	return sm.w
}

// Websites returns a slice of website buckets
func (sm *StoreManager) Websites() []*WebsiteBucket {
	return nil
}

// Group returns the group bucket
func (sm *StoreManager) Group() *GroupBucket {
	return sm.g
}

// Store returns the store view bucket
func (sm *StoreManager) Store() *StoreBucket {
	return sm.s
}

// GetDefaultStoreView returns the default store view bucket
func (sm *StoreManager) GetDefaultStoreView() *StoreBucket {
	return sm.s
}

// ReinitStores reloads the website, store group and store view data from the database @todo
func (sm *StoreManager) ReinitStores() error {
	return nil
}

// indexMap for faster access to the website, store group, store structs instead of
// iterating over the slices.
type indexMap struct {
	sync.RWMutex
	id   map[int64]int
	code map[string]int
}

// populateGroup fills the map (itself) with the group ids and the index of the slice. Thread safe.
func (im *indexMap) populateGroup(s TableGroupSlice) *indexMap {
	im.Lock()
	defer im.Unlock()
	im.id = make(map[int64]int)
	for i := 0; i < len(s); i++ {
		im.id[s[i].GroupID] = i
	}
	return im
}

// populateStore fills the map (itself) with the store ids and codes and the index of the slice. Thread safe.
func (im *indexMap) populateStore(s TableStoreSlice) *indexMap {
	im.Lock()
	defer im.Unlock()
	im.id = make(map[int64]int)
	im.code = make(map[string]int)
	for i := 0; i < len(s); i++ {
		im.id[s[i].StoreID] = i
		im.code[s[i].Code.String] = i
	}
	return im
}

// populateWebsite fills the map (itself) with the website ids and codes and the index of the slice. Thread safe.
func (im *indexMap) populateWebsite(s TableWebsiteSlice) *indexMap {
	im.Lock()
	defer im.Unlock()
	im.id = make(map[int64]int)
	im.code = make(map[string]int)
	for i := 0; i < len(s); i++ {
		im.id[s[i].WebsiteID] = i
		im.code[s[i].Code.String] = i
	}
	return im
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
