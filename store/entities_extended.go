// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/errors"
)

type sortNaturallyWebsites struct {
	*StoreWebsiteCollection
}

func (sl sortNaturallyWebsites) Less(i, j int) bool {
	switch {
	// 1st
	case sl.Data[i].SortOrder < sl.Data[j].SortOrder:
		return true
	case sl.Data[i].SortOrder > sl.Data[j].SortOrder:
		return false

		// 2nd
	case sl.Data[i].WebsiteID < sl.Data[j].WebsiteID:
		return true
	case sl.Data[i].WebsiteID > sl.Data[j].WebsiteID:
		return false

		// 3rd
	case sl.Data[i].DefaultGroupID < sl.Data[j].DefaultGroupID:
		return true
	case sl.Data[i].DefaultGroupID > sl.Data[j].DefaultGroupID:
		return false

	default:
		return false
	}
}

type sortNaturallyGroups struct {
	*StoreGroupCollection
}

func (sl sortNaturallyGroups) Less(i, j int) bool {
	switch {
	// 1st
	case sl.Data[i].WebsiteID < sl.Data[j].WebsiteID:
		return true
	case sl.Data[i].WebsiteID > sl.Data[j].WebsiteID:
		return false

		// 2nd
	case sl.Data[i].DefaultStoreID < sl.Data[j].DefaultStoreID:
		return true
	case sl.Data[i].DefaultStoreID > sl.Data[j].DefaultStoreID:
		return false

		// 3rd
	case sl.Data[i].GroupID < sl.Data[j].GroupID:
		return true
	case sl.Data[i].GroupID > sl.Data[j].GroupID:
		return false

	default:
		return false
	}
}

type sortNaturallyStores struct {
	*StoreCollection
}

func (sl sortNaturallyStores) Less(i, j int) bool {
	switch {
	// 1st
	case sl.Data[i].WebsiteID < sl.Data[j].WebsiteID:
		return true
	case sl.Data[i].WebsiteID > sl.Data[j].WebsiteID:
		return false

		// 2nd
	case sl.Data[i].GroupID < sl.Data[j].GroupID:
		return true
	case sl.Data[i].GroupID > sl.Data[j].GroupID:
		return false

		// 3rd
	case sl.Data[i].SortOrder < sl.Data[j].SortOrder:
		return true
	case sl.Data[i].SortOrder > sl.Data[j].SortOrder:
		return false

		// 4th
	case sl.Data[i].StoreID < sl.Data[j].StoreID:
		return true
	case sl.Data[i].StoreID > sl.Data[j].StoreID:
		return false

	default:
		return false
	}
}

// Default returns the default website. The returned pointer is owned by
// StoreWebsiteCollection.
func (cc *StoreWebsiteCollection) Default() (*StoreWebsite, error) {
	for _, e := range cc.Data { // assuming already correctly sorted
		if e.IsDefault {
			return e, nil
		}
	}
	return nil, errors.NotFound.Newf("[store] Default Website in slice not found")
}

// DefaultStore returns the first default active store.
func (e *StoreWebsite) DefaultStore() (*Store, error) {

	if e.StoreGroup == nil {
		return nil, errors.NotSupported.Newf("[store] StoreGroup is nil")
	}

	for _, g := range e.StoreGroup.Data {
		if e.DefaultGroupID == g.GroupID && g.DefaultStoreID > 0 {
			for _, s := range e.Store.Data {
				if g.DefaultStoreID == s.StoreID && s.WebsiteID == e.WebsiteID && s.GroupID == g.GroupID && s.IsActive {
					return s, nil
				}
			}
		}
	}
	return nil, errors.NotFound.Newf("[store] DefaultStore for Website %d not found", e.WebsiteID)
}
