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
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

const (
	// These strings are stored in the core_config_data table; can be optimized
	ScopeStores   = "stores"
	ScopeWebsites = "websites"
	ScopeDefault  = "default"
	ScopeStore    = "store"
	ScopeGroup    = "group"
	ScopeWebsite  = "website"
)

type (
	// IDX internal index used for iota in code generation and for not mixing up indexes
	IDX uint

	StoreManager struct {
		s *StoreBucket
		g *GroupBucket
		w *WebsiteBucket
	}
)

func NewStoreManager(s *StoreBucket, g *GroupBucket, w *WebsiteBucket) *StoreManager {
	return &StoreManager{
		g: g.SetStores(s).SetWebSite(w),
		s: s,
		w: w.SetGroups(g).SetStores(s),
	}
}

// IsSingleStoreModeEnabled @todo implement
// @see magento2/app/code/Magento/Store/Model/StoreManager.php uses the config from the database
func (sm *StoreManager) IsSingleStoreModeEnabled() bool {
	return false
}

// IsSingleStoreMode @todo implement
func (sm *StoreManager) IsSingleStoreMode() bool {
	return false
}
func (sm *StoreManager) HasSingleStore() bool {
	return false
}
func (sm *StoreManager) Store() *StoreBucket {
	return sm.s
}
func (sm *StoreManager) Website() *WebsiteBucket {
	return sm.w
}
func (sm *StoreManager) Group() *GroupBucket {
	return sm.g
}
func (sm *StoreManager) GetDefaultStoreView() *StoreBucket {
	return sm.s
}
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
