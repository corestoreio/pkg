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
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
)

// DefaultGroupID defines the default group id which is always 0.
const DefaultGroupID int64 = 0

// Group defines the root category id and default store id for a set of stores.
// A group is assigned to one website and a group can have multiple stores.
// A group does not have any kind of configuration setting.
type Group struct {
	cr config.Getter // internal root config.Reader which can be overridden

	// Data contains the raw group data.
	Data *TableGroup
	// Stores contains a slice to all stores associated to this group. Can be nil.
	Stores StoreSlice
	// Website contains the Website which belongs to this group. Can be nil.
	Website *Website
	*cserr.MultiErr
}

// GroupOption can be used as an argument in NewGroup to configure a group.
type GroupOption func(*Group)

// ErrGroup* are general errors when handling with the Group type.
// They are self explanatory.
var (
	ErrGroupNotFound               = errors.New("Group not found")
	ErrGroupDefaultStoreNotFound   = errors.New("Group default store not found")
	ErrGroupWebsiteNotFound        = errors.New("Group Website not found or nil or ID do not match")
	ErrGroupWebsiteIntegrityFailed = errors.New("Groups WebsiteID does not match the Websites ID")
)

// SetGroupConfig sets the config.Reader to the Group. Default reader is
// config.DefaultManager. You should call this function before calling other
// option functions otherwise your preferred config.Reader won't be inherited
// to a Website or a Store.
func SetGroupConfig(cr config.Getter) GroupOption { return func(g *Group) { g.cr = cr } }

// SetGroupWebsite assigns a website to a group. If website ID does not match
// the group website ID then add error will be generated.
func SetGroupWebsite(tw *TableWebsite) GroupOption {
	return func(g *Group) {
		if g.Data == nil {
			g.AppendErrors(ErrGroupNotFound)
			return
		}
		if tw != nil && g.Data.WebsiteID != tw.WebsiteID {
			g.AppendErrors(ErrGroupWebsiteNotFound)
			return
		}
		if tw != nil {
			var err error
			g.Website, err = NewWebsite(tw, SetWebsiteConfig(g.cr))
			g.AppendErrors(err)
		}
	}
}

// SetGroupStores uses the full store collection to extract the stores which are
// assigned to a group. Either Website must be set before calling SetGroupStores() or
// the second argument may not be nil. Does nothing if tss variable is nil.
func SetGroupStores(tss TableStoreSlice, w *TableWebsite) GroupOption {
	return func(g *Group) {
		if tss == nil {
			g.Stores = nil
			return
		}
		if g.Website == nil && w == nil {
			g.AppendErrors(ErrGroupWebsiteNotFound)
			return
		}
		if w == nil {
			w = g.Website.Data
		}
		if w.WebsiteID != g.Data.WebsiteID {
			g.AppendErrors(ErrGroupWebsiteIntegrityFailed)
			return
		}
		for _, s := range tss.FilterByGroupID(g.Data.GroupID) {
			ns, err := NewStore(s, w, g.Data, WithStoreConfig(g.cr))
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.SetGroupStores.NewStore", "err", err, "s", s, "w", w, "g.Data", g.Data)
				}
				g.AppendErrors(errors.Mask(err))
				return
			}
			g.Stores = append(g.Stores, ns)
		}
	}
}

// NewGroup creates a new Group. Returns an error if 1st argument is nil.
// Config will only be set if there has been a Website provided via
// an option argument,
func NewGroup(tg *TableGroup, opts ...GroupOption) (*Group, error) {
	if tg == nil {
		return nil, ErrArgumentCannotBeNil
	}

	g := &Group{
		cr:   config.DefaultService,
		Data: tg,
	}
	return g.ApplyOptions(opts...)
}

// MustNewGroup creates a NewGroup but panics on error.
func MustNewGroup(tg *TableGroup, opts ...GroupOption) *Group {
	g, err := NewGroup(tg, opts...)
	if err != nil {
		panic(err)
	}
	return g
}

// ApplyOptions sets the options to a Group.
func (g *Group) ApplyOptions(opts ...GroupOption) (*Group, error) {
	for _, opt := range opts {
		if opt != nil {
			opt(g)
		}
	}
	if g.HasErrors() {
		return nil, g
	}
	if g.Website != nil {
		_, err := g.Website.ApplyOptions(SetWebsiteConfig(g.cr))
		if err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("store.Group.ApplyOptions.Website.ApplyOptions", "err", err, "g", g)
			}
			return nil, errors.Mask(err)
		}
	}
	return g, nil
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
func (g *Group) DefaultStore() (*Store, error) {
	for _, sb := range g.Stores {
		if sb.Data.StoreID == g.Data.DefaultStoreID {
			return sb, nil
		}
	}
	return nil, ErrGroupDefaultStoreNotFound
}

/*
	@todo implement Magento\Store\Model\Group
*/
