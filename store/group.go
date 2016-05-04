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
// A group is assigned to one website and a group can have multiple stores.
// A group does not have any kind of configuration setting.
type Group struct {
	// cr internal root config.Getter which will be applied to stores and websites
	cr config.Getter
	// Data contains the raw group data.
	Data *TableGroup
	// Stores contains a slice to all stores associated to this group. Can be nil.
	Stores StoreSlice
	// Website contains the Website which belongs to this group. Can be nil.
	Website *Website
	// optionError use by functional option arguments to indicate that one
	// option has triggered an error and hence the other can options can
	// skip their process.
	optionError error
}

// NewGroup creates a new Group. Returns an error if 1st argument is nil.
// Config will only be set if there has been a Website provided via
// an option argument.
// Error behaviour: Empty
func NewGroup(tg *TableGroup, opts ...GroupOption) (*Group, error) {
	g := &Group{
		Data: tg,
	}
	if err := g.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[store] NewGroup Options")
	}
	return g, nil
}

// MustNewGroup creates a NewGroup but panics on error.
func MustNewGroup(tg *TableGroup, opts ...GroupOption) *Group {
	g, err := NewGroup(tg, opts...)
	if err != nil {
		panic(err)
	}
	return g
}

// Options sets the options to a Group.
func (g *Group) Options(opts ...GroupOption) error {
	for _, opt := range opts {
		opt(g)
	}
	if g.optionError != nil {
		// clear error or next call to Options() will fail.
		defer func() { g.optionError = nil }()
		return g.optionError
	}
	if g.Website != nil {
		if err := g.Website.Options(SetWebsiteConfig(g.cr)); err != nil {
			return errors.Wrapf(err, "[store] Group %#v", g)
		}
	}
	return nil
}

// GroupID satisfies interface scope.GroupIDer and returns the group ID.
func (g *Group) GroupID() int64 {
	return g.Data.GroupID
}

// StoreID satisfies interface scope.StoreIDer and returns the default store ID.
func (g *Group) StoreID() int64 {
	return g.Data.DefaultStoreID
}

// MarshalJSON satisfies interface for JSON marshalling. The TableWebsite
// struct will be encoded to JSON.
func (g *Group) MarshalJSON() ([]byte, error) {
	// @todo while generating the TableStore structs we can generate the ffjson code ...
	return json.Marshal(g.Data)
}

// DefaultStore returns the default *Store or an error. If an error will be returned of
// type ErrGroupDefaultStoreNotFound you can then access Data field to get the
// DefaultStoreID. The returned *Store does not contain that much data to other
// Website or Groups.
// Error behaviour: NotFound
func (g *Group) DefaultStore() (*Store, error) {
	for _, sb := range g.Stores {
		if sb.Data.StoreID == g.Data.DefaultStoreID {
			return sb, nil
		}
	}
	return nil, errors.NewNotFoundf(errGroupDefaultStoreNotFound)
}

/*
	@todo implement Magento\Store\Model\Group
*/
