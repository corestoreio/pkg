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
	"github.com/corestoreio/csfw/config"
	"github.com/juju/errors"
)

// GroupOption can be used as an argument in NewGroup to configure a group.
type GroupOption func(*Group)

// SetGroupConfig sets the config.Getter to the Group. You should call this
// function before calling other option functions otherwise your preferred
// config.Getter won't be inherited to a Website or a Store.
func SetGroupConfig(cr config.Getter) GroupOption { return func(g *Group) { g.cr = cr } }

// SetGroupWebsite assigns a website to a group. If website ID does not match
// the group website ID then add error will be generated.
func SetGroupWebsite(tw *TableWebsite) GroupOption {
	return func(g *Group) {
		if g.Data == nil {
			g.MultiErr = g.AppendErrors(ErrGroupNotFound)
			return
		}
		if tw != nil && g.Data.WebsiteID != tw.WebsiteID {
			g.MultiErr = g.AppendErrors(ErrGroupWebsiteNotFound)
			return
		}
		if tw != nil {
			var err error
			g.Website, err = NewWebsite(tw, SetWebsiteConfig(g.cr))
			g.MultiErr = g.AppendErrors(err)
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
			g.MultiErr = g.AppendErrors(ErrGroupWebsiteNotFound)
			return
		}
		if w == nil {
			w = g.Website.Data
		}
		if w.WebsiteID != g.Data.WebsiteID {
			g.MultiErr = g.AppendErrors(ErrGroupWebsiteIntegrityFailed)
			return
		}
		for _, s := range tss.FilterByGroupID(g.Data.GroupID) {
			ns, err := NewStore(s, w, g.Data, WithStoreConfig(g.cr))
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.SetGroupStores.NewStore", "err", err, "s", s, "w", w, "g.Data", g.Data)
				}
				g.MultiErr = g.AppendErrors(errors.Mask(err))
				return
			}
			g.Stores = append(g.Stores, ns)
		}
	}
}
