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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/util/errors"
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
	w := Website{
		Config: cfg.NewScoped(tw.WebsiteID, 0),
		Data:   tw,
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

// Validate checks the internal integrity. May panic when the data has not been
// set. Empty Groups or Stores are valid settings.
func (w Website) Validate() error {
	for _, g := range w.Groups {
		if w.GroupID() != g.Data.GroupID {
			return errors.NewNotValidf("[store] Website.Validate: Website Group ID %d does not match Store Group ID %d", w.GroupID(), g.Data.GroupID)
		}
	}
	for _, s := range w.Stores {
		if w.ID() != s.Data.WebsiteID {
			return errors.NewNotValidf("[store] Website.Validate: Website ID %d does not match Store Website ID %d", w.ID(), s.Data.WebsiteID)
		}
	}
	if w.Config.WebsiteID != w.ID() {
		return errors.NewNotValidf("[store] Website.Validate: Config Website ID %d does not match Website ID %d", w.Config.WebsiteID, w.ID())
	}
	return nil
}

// SetGroupsStores uses a group slice and a table store slice to set the groups
// associated to this website and the stores associated to this website. It
// returns an error if the data integrity is incorrect.
func (w *Website) SetGroupsStores(tgs TableGroupSlice, tss TableStoreSlice) error {

	groups := tgs.Filter(func(tg *TableGroup) bool {
		return tg.WebsiteID == w.Data.WebsiteID
	})

	w.Groups = make(GroupSlice, groups.Len(), groups.Len())
	for i, g := range groups {
		var err error
		w.Groups[i], err = NewGroup(w.Config.Root, g, w.Data, tss)
		if err != nil {
			return errors.Wrapf(err, "[store] NewGroup. Group %#v Website Data: %#v", g, w.Data)
		}
	}
	stores := tss.FilterByWebsiteID(w.Data.WebsiteID)
	w.Stores = make(StoreSlice, stores.Len(), stores.Len())
	for i, s := range stores {
		group, found := tgs.FindByGroupID(s.GroupID)
		if !found {
			return errors.NewNotFoundf("[store] Website Integrity error. A store %#v must be assigned to a group.\nGroupSlice: %#v\n\n", s, tgs)
		}
		var err error
		w.Stores[i], err = NewStore(w.Config.Root, s, w.Data, group)
		if err != nil {
			return errors.Wrapf(err, "[store] NewStore. Store %#v Website Data %#v Group %#v", s, w.Data, group)
		}
	}
	return w.Validate()
}

// ID returns the website ID.
func (w Website) ID() int64 { return w.Data.WebsiteID }

// Code returns the website code.
func (w Website) Code() string { return w.Data.Code.String }

// GroupID returns the associated group ID.
func (w Website) GroupID() int64 { return w.Data.DefaultGroupID }

// DefaultStoreID returns the default store ID associated to the underlying
// group.
func (w Website) DefaultStoreID() (int64, error) {
	g, err := w.DefaultGroup()
	if err != nil {
		return 0, errors.Wrap(err, "[store] Website.DefaultStoreID")
	}
	return g.Data.DefaultStoreID, nil
}

// MarshalJSON satisfies interface for JSON marshalling. The TableWebsite
// struct will be encoded to JSON.
func (w Website) MarshalJSON() ([]byte, error) {
	// @todo while generating the TableStore structs we can generate the ffjson code ...
	return json.Marshal(w.Data)
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

// DefaultStore returns the default store associated to the underlying group.
func (w Website) DefaultStore() (Store, error) {
	g, err := w.DefaultGroup()
	if err != nil {
		return Store{}, errors.Wrap(err, "[store] Website.DefaultGroup")
	}
	return g.DefaultStore()
}

/*
	@todo implement Magento\Store\Model\Website
*/
