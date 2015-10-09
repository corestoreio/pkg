// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"errors"

	"encoding/json"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/utils"
)

const (
	// DefaultGroupID defines the default group id which is always 0.
	DefaultGroupID int64 = 0
)

// Group defines the root category id and default store id for a set of stores.
// A group is assigned to one website and a group can have multiple stores.
// A group does not have any kind of configuration setting.
type Group struct {
	cr config.Reader
	// Data contains the raw group data.
	Data *TableGroup
	// Stores contains a slice to all stores associated to this group. Can be nil.
	Stores StoreSlice
	// Website contains the Website which belongs to this group. Can be nil.
	Website *Website
}

// GroupSlice collection of Group. GroupSlice has some nice method receivers.
type GroupSlice []*Group

// GroupOption can be used as an argument in NewGroup to configure a group.
type GroupOption func(*Group)

var (
	ErrGroupNotFound             = errors.New("Group not found")
	ErrGroupDefaultStoreNotFound = errors.New("Group default store not found")
	// ErrGroupWebsiteNotFound the Website struct is nil so we cannot assign the stores to a group.
	ErrGroupWebsiteNotFound = errors.New("Group Website not found or nil or ID do not match")
)

// SetGroupConfig sets the configuration Reader to the Group.
// Default reader is config.DefaultManager
func SetGroupConfig(cr config.Reader) GroupOption {
	return func(g *Group) { g.cr = cr }
}

// SetGroupWebsite assigns a website to a group. If website ID does not match
// the group website ID then this function panics.
func SetGroupWebsite(tw *TableWebsite) GroupOption {
	return func(g *Group) {
		if g.Data == nil {
			panic(ErrGroupNotFound)
		}
		if tw != nil && g.Data.WebsiteID != tw.WebsiteID {
			panic(ErrGroupWebsiteNotFound)
		}
		if tw != nil {
			g.Website = NewWebsite(tw)
		}
	}
}

// NewGroup initializes a new Group with the config.DefaultManager
func NewGroup(tg *TableGroup, opts ...GroupOption) *Group {
	if tg == nil {
		panic(ErrStoreNewArgNil)
	}

	g := &Group{
		cr:   config.DefaultManager,
		Data: tg,
	}
	g.ApplyOptions(opts...)
	if g.Website != nil {
		g.Website.ApplyOptions(SetWebsiteConfig(g.cr))
	}
	return g
}

var _ scope.GroupIDer = (*Group)(nil)

// GroupID satisfies interface scope.GroupIDer and returns the group ID.
func (g *Group) GroupID() int64 {
	return g.Data.GroupID
}

// ApplyOptions sets the options to a Group.
func (g *Group) ApplyOptions(opts ...GroupOption) *Group {
	for _, opt := range opts {
		if opt != nil {
			opt(g)
		}
	}
	return g
}

// MarshalJSON satisfies interface for JSON marshalling. The TableWebsite
// struct will be encoded to JSON.
func (g *Group) MarshalJSON() ([]byte, error) {
	// @todo while generating the TableStore structs we can generate the ffjson code ...
	return json.Marshal(g.Data)
}

// DefaultStore returns the default Store or an error.
func (g *Group) DefaultStore() (*Store, error) {
	for _, sb := range g.Stores {
		if sb.Data.StoreID == g.Data.DefaultStoreID {
			return sb, nil
		}
	}
	return nil, ErrGroupDefaultStoreNotFound
}

// SetStores uses the full store collection to extract the stores which are
// assigned to a group. Either Website must be set before calling SetStores() or
// the second argument must be set i.e. 2nd argument can be nil. Panics if both
// values are nil. If both are set, the 2nd argument will be considered.
func (g *Group) SetStores(tss TableStoreSlice, w *TableWebsite) *Group {
	if tss == nil {
		g.Stores = nil
		return g
	}
	if g.Website == nil && w == nil {
		panic(ErrGroupWebsiteNotFound)
	}
	if w == nil {
		w = g.Website.Data
	}
	if w.WebsiteID != g.Data.WebsiteID {
		panic(ErrGroupWebsiteNotFound)
	}
	for _, s := range tss.FilterByGroupID(g.Data.GroupID) {
		g.Stores = append(g.Stores, NewStore(s, w, g.Data, SetStoreConfig(g.cr)))
	}
	return g
}

/*
	@todo implement Magento\Store\Model\Group
*/

/*
	GroupSlice method receivers
*/

// Len returns the length
func (s GroupSlice) Len() int { return len(s) }

// Filter returns a new slice filtered by predicate f
func (s GroupSlice) Filter(f func(*Group) bool) GroupSlice {
	var gs GroupSlice
	for _, v := range s {
		if v != nil && f(v) {
			gs = append(gs, v)
		}
	}
	return gs
}

// IDs returns an Int64Slice with all store ids
func (s GroupSlice) IDs() utils.Int64Slice {
	if len(s) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, g := range s {
		if g != nil {
			ids.Append(g.Data.GroupID)
		}
	}
	return ids
}
