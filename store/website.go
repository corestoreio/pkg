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
	"encoding/json"

	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/errors"
)

// Website represents the overall parent structure of its children Group and
// Store. A website defines the default group ID. A website can contain custom
// configuration settings which overrides the default scope but get itself
// overridden by the Store scope.
type Website struct {
	// Config contains the scoped configuration which cannot be changed once the
	// object has been created.
	Config config.Scoped
	// Data raw website data from DB table. If nil, website object is invalid
	Data *TableWebsite
	// Groups contains a slice to all groups associated to one website. This slice
	// can be nil.
	Groups GroupSlice
	// Stores contains a slice to all stores associated to one website. This slice
	// can be nil.
	Stores StoreSlice
}

// NewWebsite creates a new Website with its depended groups and stores.
func NewWebsite(cfg config.Getter, tw *TableWebsite, tgs TableGroupSlice, tss TableStoreSlice) (Website, error) {
	// tw cannot be nil, but tgs and tss can be
	w := Website{
		Data: tw,
	}
	if tw != nil {
		w.Config = cfg.NewScoped(tw.WebsiteID, 0) // todo 0 could be a bug because admin store view, should be maybe -1
	}
	if err := w.SetGroupsStores(tgs, tss); err != nil {
		return Website{}, errors.Wrap(err, "[store] NewWebsite.SetWebsiteGroupsStores")
	}
	return w, nil
}

// MustNewWebsite same as NewWebsite but panics on error.
func MustNewWebsite(cfg config.Getter, tw *TableWebsite, tgs TableGroupSlice, tss TableStoreSlice) Website {
	w, err := NewWebsite(cfg, tw, tgs, tss)
	if err != nil {
		panic(err)
	}
	return w
}

// SetGroupsStores uses a group slice and a table store slice to set the groups
// associated to this website and the stores associated to this website. It
// returns an error if the data integrity is incorrect.
func (w *Website) SetGroupsStores(tgs TableGroupSlice, tss TableStoreSlice) error {
	if tgs == nil && tss == nil {
		// avoid recursion and stack overflow
		return nil
	}

	groups := tgs.Filter(func(tg *TableGroup) bool {
		return tg != nil && tg.WebsiteID == w.ID()
	})

	w.Groups = make(GroupSlice, groups.Len(), groups.Len())
	for i, g := range groups {
		var err error
		w.Groups[i], err = NewGroup(w.Config.Root, g, nil, tss) // passing nil to limit the recursion
		if err != nil {
			return errors.Wrapf(err, "[store] NewGroup. Group %#v Website Data: %#v", g, w.Data)
		}
	}
	stores := tss.FilterByWebsiteID(w.ID())
	w.Stores = make(StoreSlice, stores.Len(), stores.Len())
	for i, s := range stores {
		var err error
		w.Stores[i], err = NewStore(w.Config.Root, s, nil, nil)
		if err != nil {
			return errors.Wrapf(err, "[store] NewStore. Store %#v Website Data %#v", s, w.Data)
		}
	}
	return w.Validate()
}

// Validate checks the internal integrity. May panic when the data has not been
// set. Empty Groups or Stores are valid settings.
func (w Website) Validate() error {
	for _, g := range w.Groups {
		if w.Data != nil && g.Data != nil && w.ID() != g.Data.WebsiteID {
			return errors.NewNotValidf("[store] Website %d != Group.WebsiteID %d", w.ID(), g.Data.WebsiteID)
		}
	}
	for _, s := range w.Stores {
		if w.ID() != s.Data.WebsiteID {
			return errors.NewNotValidf("[store] Website ID %d != Store Website ID %d", w.ID(), s.Data.WebsiteID)
		}
	}
	if w.Config.WebsiteID != w.ID() {
		return errors.NewNotValidf("[store] Website ID %d != Config Website ID %d", w.ID(), w.Config.WebsiteID)
	}
	return nil
}

// ID returns the website ID. If Data is nil, returns -1.
func (w Website) ID() int64 {
	if w.Data == nil {
		return -1
	}
	return w.Data.WebsiteID
}

// Code returns the website code. Returns an empty string if Data is nil.
func (w Website) Code() string {
	if w.Data == nil {
		return ""
	}
	return w.Data.Code.String
}

// Name returns the website name. Returns an empty string if Data is nil.
func (w Website) Name() string {
	if w.Data == nil {
		return ""
	}
	return w.Data.Name.String
}

// DefaultGroupID returns the associated default group ID. If Data is nil,
// returns -1.
func (w Website) DefaultGroupID() int64 {
	if w.Data == nil {
		return -1
	}
	return w.Data.DefaultGroupID
}

// DefaultGroup returns the default Group or an error if not found.
func (w Website) DefaultGroup() (Group, error) {
	for _, g := range w.Groups {
		if w.Data.DefaultGroupID == g.Data.GroupID {
			return g, nil
		}
	}
	return Group{}, errors.NewNotFoundf(errWebsiteDefaultGroupNotFound)
}

// DefaultStoreID returns the default store ID associated to the underlying
// group.
func (w Website) DefaultStoreID() (int64, error) {
	g, err := w.DefaultGroup()
	if err != nil {
		return 0, errors.Wrap(err, "[store] Website.DefaultStoreID")
	}
	return g.DefaultStoreID(), nil
}

// DefaultStore returns the default store associated to the underlying group.
func (w Website) DefaultStore() (Store, error) {
	g, err := w.DefaultGroup()
	if err != nil {
		return Store{}, errors.Wrap(err, "[store] Website.DefaultGroup")
	}
	if ds, ok := w.Stores.FindByID(g.DefaultStoreID()); ok {
		return ds, nil
	}
	return Store{}, errors.NewNotFoundf(errWebsiteStoreDefaultNotFound)
}

// MarshalJSON satisfies interface for JSON marshalling. The TableWebsite
// struct will be encoded to JSON.
func (w Website) MarshalJSON() ([]byte, error) {
	// @todo while generating the TableStore structs we can generate the ffjson code ...
	return json.Marshal(w.Data)
}

/*
	@todo implement Magento\Store\Model\Website
*/
