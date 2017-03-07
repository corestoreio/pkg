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
	"github.com/corestoreio/errors"
)

// Group defines the root category id and default store id for a set of stores.
// A group is assigned to one website and a group can have multiple stores. A
// group does not have any kind of configuration setting but hands down the
// BaseConfig to the stores and the Website.
type Group struct {
	// Data contains the raw group data. Cannot be nil
	Data *TableGroup
	// Stores contains a slice to all stores associated to this group. Can be nil.
	Stores StoreSlice
	// Website contains the Website which belongs to this group.
	Website Website
}

// NewGroup creates a new Group with its depended Website and Stores.
func NewGroup(cfg config.Getter, tg *TableGroup, tw *TableWebsite, tss TableStoreSlice) (Group, error) {
	// tg cannot be nil but tw and tss can be nil
	g := Group{
		Data: tg,
	}
	if err := g.SetWebsiteStores(cfg, tw, tss); err != nil {
		return Group{}, errors.Wrap(err, "[store] NewGroup.SetWebsiteStores")
	}
	return g, nil
}

// MustNewGroup creates a NewGroup but panics on error.
func MustNewGroup(cfg config.Getter, tg *TableGroup, tw *TableWebsite, tss TableStoreSlice) Group {
	g, err := NewGroup(cfg, tg, tw, tss)
	if err != nil {
		panic(err)
	}
	return g
}

// SetWebsiteStores applies a raw website and multiple stores belonging to the
// group. Validates the internal integrity afterwards.
func (g *Group) SetWebsiteStores(cfg config.Getter, w *TableWebsite, tss TableStoreSlice) error {
	if w == nil && tss == nil {
		// avoid recursion and stack overflow
		return nil
	}

	var err error
	g.Website, err = NewWebsite(cfg, w, nil, nil)
	if err != nil {
		return errors.Wrap(err, "[store] SetWebsiteStores.NewWebsite")
	}

	for _, s := range tss.FilterByGroupID(g.ID()) {
		ns, err := NewStore(cfg, s, nil, nil)
		if err != nil {
			var wID int64
			if w != nil {
				wID = w.WebsiteID
			}
			return errors.Wrapf(err, "[store] SetWebsiteStores.FilterByGroupID.NewStore. StoreID %d WebsiteID %d Group %v", s.StoreID, wID, g.ID())
		}
		g.Stores = append(g.Stores, ns)
	}
	return g.Validate()
}

// Validate checks the internal integrity. May panic when the data has not been
// set. Empty Website or Stores are valid settings.
func (g Group) Validate() error {

	gw, ww := g.WebsiteID(), g.Website.ID()
	if g.Website.Data != nil && gw != ww {
		return errors.NewNotValidf("[store] Group %d Group.WebsiteID %d != Website.WebsiteID %d", g.ID(), gw, ww)
	}
	for _, s := range g.Stores {
		if g.ID() != s.GroupID() {
			return errors.NewNotValidf("[store] Group %d != Store %d GroupID %d", g.ID(), s.ID(), s.Data.GroupID)
		}
	}
	return nil
}

// ID returns the group ID. If Data is nil returns -1.
func (g Group) ID() int64 {
	if g.Data == nil {
		return -1
	}
	return g.Data.GroupID
}

// WebsiteID returns the website ID. If Data is nil returns -1.
func (g Group) WebsiteID() int64 {
	if g.Data == nil {
		return -1
	}
	return g.Data.WebsiteID
}

// DefaultStoreID returns the default store ID. If Data is nil returns -1.
func (g Group) DefaultStoreID() int64 {
	if g.Data == nil {
		return -1
	}
	return g.Data.DefaultStoreID
}

// Name returns the Group name or empty if Data is nil.
func (g Group) Name() string {
	if g.Data == nil {
		return ""
	}
	return g.Data.Name
}

// MarshalJSON satisfies interface for JSON marshalling. The TableWebsite
// struct will be encoded to JSON.
func (g Group) MarshalJSON() ([]byte, error) {
	// @todo while generating the TableStore structs we can generate the ffjson code ...
	return json.Marshal(g.Data)
}

// DefaultStore returns the default Store or an error of behaviour NotFound.
func (g Group) DefaultStore() (Store, error) {
	for _, s := range g.Stores {
		if s.Data.StoreID == g.Data.DefaultStoreID {
			return s, nil
		}
	}
	return Store{}, errors.NewNotFoundf(errGroupDefaultStoreNotFound, g.Data.DefaultStoreID)
}
