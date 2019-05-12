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
	*StoreWebsites
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
	*StoreGroups
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
	*Stores
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
// StoreWebsites.
func (cc *StoreWebsites) Default() (*StoreWebsite, error) {
	for _, e := range cc.Data { // assuming already correctly sorted
		if e.IsDefault {
			return e, nil
		}
	}
	return nil, errors.NotFound.Newf("[store] Default Website in slice not found")
}

func (e *StoreWebsite) DefaultGroup() (*StoreGroup, error) {
	if e.StoreGroups == nil {
		return nil, errors.NotSupported.Newf("[store] StoreWebsite %d does not have StoreGroups, which is nil", e.WebsiteID)
	}
	for _, g := range e.StoreGroups.Data {
		if e.DefaultGroupID > 0 && e.DefaultGroupID == g.GroupID {
			return g, nil
		}
	}
	return nil, errors.NotFound.Newf("[store] DefaultGroup for Website %d not found", e.WebsiteID)
}

// DefaultStore returns the first default active store.
func (e *StoreWebsite) DefaultStore() (*Store, error) {
	if e.StoreGroups == nil {
		return nil, errors.NotSupported.Newf("[store] StoreWebsite %d does not have StoreGroups, which is nil", e.WebsiteID)
	}
	for _, g := range e.StoreGroups.Data {
		if e.DefaultGroupID == g.GroupID && g.DefaultStoreID > 0 {
			for _, s := range e.Stores.Data {
				if g.DefaultStoreID == s.StoreID && s.WebsiteID == e.WebsiteID && s.GroupID == g.GroupID && s.IsActive {
					return s, nil
				}
			}
		}
	}
	return nil, errors.NotFound.Newf("[store] DefaultStore for Website %d not found", e.WebsiteID)
}

func init() {
	validateStore = func(e *Store) (err error) {
		switch {
		case e.StoreWebsite != nil && e.StoreWebsite.WebsiteID != e.WebsiteID:
			err = errors.NotValid.Newf("[store] Store %d: WebsiteID %d != Website.ID %d", e.StoreID, e.WebsiteID, e.StoreWebsite.WebsiteID)

		case e.StoreGroup != nil && e.StoreGroup.WebsiteID != e.WebsiteID:
			err = errors.NotValid.Newf("[store] Store %d: Group.WebsiteID %d != Website.ID %d", e.StoreID, e.StoreGroup.WebsiteID, e.WebsiteID)

		case e.StoreGroup != nil && e.StoreGroup.GroupID != e.GroupID:
			err = errors.NotValid.Newf("[store] Store %d: Store.GroupID %d != Group.ID %d", e.StoreID, e.GroupID, e.StoreGroup.GroupID)

		case e.Code == "":
			err = errors.NotValid.Newf("[store] Store %d: Empty code", e.StoreID)

		}
		return err
	}

	validateStoreGroup = func(g *StoreGroup) (err error) {
		switch {
		case g.StoreWebsite != nil && g.StoreWebsite.WebsiteID != g.WebsiteID:
			err = errors.NotValid.Newf("[store] Group %d: WebsiteID %d != Website.ID %d", g.GroupID, g.WebsiteID, g.StoreWebsite.WebsiteID)

		case g.Code == "":
			err = errors.NotValid.Newf("[store] Group %d: Empty code", g.GroupID)

		}
		return err
	}

	validateStoreWebsite = func(w *StoreWebsite) (err error) {
		if w.Stores != nil {
			for _, st := range w.Stores.Data {
				if st.WebsiteID != w.WebsiteID {
					return errors.NotValid.Newf("[store] Website %d: Stores.WebsiteID %d != Website.ID %d", w.WebsiteID, st.WebsiteID, w.WebsiteID)
				}
			}
		}
		if w.StoreGroups != nil {
			for _, g := range w.StoreGroups.Data {
				if g.WebsiteID != w.WebsiteID {
					return errors.NotValid.Newf("[store] Website %d: StoreGroups.WebsiteID %d != Website.ID %d", w.WebsiteID, g.WebsiteID, w.WebsiteID)
				}
			}
		}

		if w.Code == "" {
			return errors.NotValid.Newf("[store] Website %d: Empty code", w.WebsiteID)
		}
		return nil
	}
}
