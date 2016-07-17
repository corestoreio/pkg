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

// DefaultGroupID defines the default group id which is always 0.
const DefaultGroupID int64 = 0

// Group defines the root category id and default store id for a set of stores.
// A group is assigned to one website and a group can have multiple stores. A
// group does not have any kind of configuration setting but hands down the
// BaseConfig to the stores and the Website.
type Group struct {
	// baseConfig base config.Getter which will be applied to stores and websites.
	baseConfig config.Getter
	// Data contains the raw group data. Cannot be nil
	Data *TableGroup
	// Stores contains a slice to all stores associated to this group. Can be nil.
	Stores StoreSlice
	// Website contains the Website which belongs to this group.
	Website Website
}

// NewGroup creates a new Group. Returns an error if 1st argument is nil. Config
// will only be set if there has been a Website provided via an option argument.
// Error behaviour: Empty
func NewGroup(cfg config.Getter, tg *TableGroup, opts ...GroupOption) (Group, error) {
	g := Group{
		baseConfig: cfg,
		Data:       tg,
	}
	if err := g.Options(opts...); err != nil {
		return Group{}, errors.Wrap(err, "[store] NewGroup Options")
	}
	return g, nil
}

// MustNewGroup creates a NewGroup but panics on error.
func MustNewGroup(cfg config.Getter, tg *TableGroup, opts ...GroupOption) Group {
	g, err := NewGroup(cfg, tg, opts...)
	if err != nil {
		panic(err)
	}
	return g
}

// Options applies different options to a Group.
func (g *Group) Options(opts ...GroupOption) error {
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return errors.Wrap(err, "[store] Group.Options")
		}
	}
	return nil
}

// GroupID returns the under
func (g Group) GroupID() int64 {
	return g.Data.GroupID
}

// MarshalJSON satisfies interface for JSON marshalling. The TableWebsite
// struct will be encoded to JSON.
func (g Group) MarshalJSON() ([]byte, error) {
	// @todo while generating the TableStore structs we can generate the ffjson code ...
	return json.Marshal(g.Data)
}

// DefaultStore returns the default *Store or an error. If an error will be returned of
// type ErrGroupDefaultStoreNotFound you can then access Data field to get the
// DefaultStoreID. The returned *Store does not contain that much data to other
// Website or Groups.
// Error behaviour: NotFound
func (g Group) DefaultStore() (Store, error) {
	for _, sb := range g.Stores {
		if sb.Data.StoreID == g.Data.DefaultStoreID {
			return sb, nil
		}
	}
	return Store{}, errors.NewNotFoundf(errGroupDefaultStoreNotFound)
}
