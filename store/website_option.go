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
	"fmt"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/util/errors"
)

// WebsiteOption can be used as an argument in NewWebsite to configure a website.
type WebsiteOption func(*Website)

// SetWebsiteConfig sets the config.Getter to the Website. You should call this
// function before calling other option functions otherwise your preferred
// config.Getter won't be inherited to a Group or Store.
func SetWebsiteConfig(cr config.Getter) WebsiteOption { return func(w *Website) { w.cr = cr } }

// SetWebsiteGroupsStores uses a group slice and a table slice to set the groups associated
// to this website and the stores associated to this website. It returns an error if
// the data integrity is incorrect.
func SetWebsiteGroupsStores(tgs TableGroupSlice, tss TableStoreSlice) WebsiteOption {
	return func(w *Website) {
		groups := tgs.Filter(func(tg *TableGroup) bool {
			return tg.WebsiteID == w.Data.WebsiteID
		})

		w.Groups = make(GroupSlice, groups.Len(), groups.Len())
		for i, g := range groups {
			var err error
			w.Groups[i], err = NewGroup(g, SetGroupWebsite(w.Data), SetGroupConfig(w.cr), SetGroupStores(tss, nil))
			if err != nil {
				w.MultiErr = w.AppendErrors(errors.Wrapf(err, "[store] NewGroup. Group %#v Website Data: %#v", g, w.Data))
				return
			}
		}
		stores := tss.FilterByWebsiteID(w.Data.WebsiteID)
		w.Stores = make(StoreSlice, stores.Len(), stores.Len())
		for i, s := range stores {
			group, err := tgs.FindByGroupID(s.GroupID)
			if err != nil {
				w.MultiErr = w.AppendErrors(fmt.Errorf("Integrity error. A store %#v must be assigned to a group.\nGroupSlice: %#v\n\n", s, tgs))
				return
			}
			w.Stores[i], err = NewStore(s, w.Data, group, WithStoreConfig(w.cr))
			if err != nil {
				w.MultiErr = w.AppendErrors(errors.Wrapf(err, "[store] NewStore. Store %#v Website Data %#v Group %#v", s, w.Data, group))
				return
			}
		}
	}
}
