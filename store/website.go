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
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/utils"
	"golang.org/x/text/language"
)

const (
	// DefaultWebsiteID is always 0
	DefaultWebsiteID int64 = 0
)

// Website represents the overall parent structure of its children Group and Store.
// A website defines the default group ID. A website can contain custom configuration
// settings which overrides the default scope but get itself overridden by the Store scope.
type Website struct {
	cr config.Reader // internal root config.Reader which can be overriden
	// Config contains a config.Manager which takes care of the scope based
	// configuration values.
	Config config.ScopedReader
	// Data raw website data
	Data *TableWebsite

	// Groups contains a slice to all groups associated to one website. This slice can be nil.
	Groups GroupSlice
	// Stores contains a slice to all stores associated to one website. This slice can be nil.
	Stores StoreSlice
}

// WebsiteSlice contains pointer to Website struct and some nifty method receivers.
type WebsiteSlice []*Website

// WebsiteOption can be used as an argument in NewWebsite to configure a website.
type WebsiteOption func(*Website)

var (
	// ErrWebsiteNotFound when the website has not been found within a slice
	ErrWebsiteNotFound = errors.New("Website not found")
	// ErrWebsiteDefaultGroupNotFound the default group cannot be found
	ErrWebsiteDefaultGroupNotFound = errors.New("Website Default Group not found")
)

// SetWebsiteConfig sets the config.Reader to the Website.
// Default reader is config.DefaultManager
func SetWebsiteConfig(cr config.Reader) WebsiteOption { return func(w *Website) { w.cr = cr } }

// NewWebsite returns a new pointer to a Website.
func NewWebsite(tw *TableWebsite, opts ...WebsiteOption) *Website {
	if tw == nil {
		panic(ErrStoreNewArgNil)
	}
	w := &Website{
		cr:   config.DefaultManager,
		Data: tw,
	}
	return w.ApplyOptions(opts...)
}

// ApplyOptions sets the options on a Website
func (w *Website) ApplyOptions(opts ...WebsiteOption) *Website {
	for _, opt := range opts {
		if opt != nil {
			opt(w)
		}
	}
	w.Config = w.cr.NewScoped(w.WebsiteID(), 0, 0) // Scope Store and Group is not available
	return w
}

var _ scope.WebsiteIDer = (*Website)(nil)
var _ scope.WebsiteCoder = (*Website)(nil)

// WebsiteID satisfies the interface scope.WebsiteIDer and returns the website ID.
func (w *Website) WebsiteID() int64 { return w.Data.WebsiteID }

// WebsiteCode satisfies the interface scope.WebsiteCoder and returns the code.
func (w *Website) WebsiteCode() string { return w.Data.Code.String }

// MarshalJSON satisfies interface for JSON marshalling. The TableWebsite
// struct will be encoded to JSON.
func (w *Website) MarshalJSON() ([]byte, error) {
	// @todo while generating the TableStore structs we can generate the ffjson code ...
	return json.Marshal(w.Data)
}

// DefaultGroup returns the default Group or an error if not found
func (w *Website) DefaultGroup() (*Group, error) {
	for _, g := range w.Groups {
		if w.Data.DefaultGroupID == g.Data.GroupID {
			return g, nil
		}
	}
	return nil, ErrWebsiteDefaultGroupNotFound
}

// DefaultStore returns the default store which via the default group.
func (w *Website) DefaultStore() (*Store, error) {
	g, err := w.DefaultGroup()
	if err != nil {
		return nil, err
	}
	return g.DefaultStore()
}

// SetGroupsStores uses a group slice and a table slice to set the groups associated to this website
// and the stores associated to this website. It panics if the integrity is incorrect.
func (w *Website) SetGroupsStores(tgs TableGroupSlice, tss TableStoreSlice) *Website {
	groups := tgs.Filter(func(tg *TableGroup) bool {
		return tg.WebsiteID == w.Data.WebsiteID
	})
	w.Groups = make(GroupSlice, groups.Len(), groups.Len())
	for i, g := range groups {
		w.Groups[i] = NewGroup(g, SetGroupWebsite(w.Data), SetGroupConfig(w.cr)).SetStores(tss, nil)
	}
	stores := tss.FilterByWebsiteID(w.Data.WebsiteID)
	w.Stores = make(StoreSlice, stores.Len(), stores.Len())
	for i, s := range stores {
		group, err := tgs.FindByGroupID(s.GroupID)
		if err != nil {
			panic(fmt.Sprintf("Integrity error. A store %#v must be assigned to a group.\nGroupSlice: %#v\n\n", s, tgs))
		}
		w.Stores[i] = NewStore(s, w.Data, group, SetStoreConfig(w.cr))
	}
	return w
}

// @todo
func (w *Website) BaseCurrencyCode() (language.Currency, error) {
	var c string
	if w.Config.GetString(PathPriceScope) == PriceScopeGlobal {
		c, _ = w.cr.GetString(config.Path(directory.PathCurrencyBase)) // TODO check for error
	} else {
		c = w.Config.GetString(directory.PathCurrencyBase)
	}
	return language.ParseCurrency(c)
}

// @todo
func (w *Website) BaseCurrency() directory.Currency {
	return directory.Currency{}
}

/*
	WebsiteSlice method receivers
*/

// Sort convenience helper
func (ws *WebsiteSlice) Sort() *WebsiteSlice {
	sort.Sort(ws)
	return ws
}

func (ws WebsiteSlice) Len() int { return len(ws) }

func (ws *WebsiteSlice) Swap(i, j int) { (*ws)[i], (*ws)[j] = (*ws)[j], (*ws)[i] }

func (ws *WebsiteSlice) Less(i, j int) bool {
	return (*ws)[i].Data.SortOrder < (*ws)[j].Data.SortOrder
}

// Filter returns a new slice filtered by predicate f
func (ws WebsiteSlice) Filter(f func(*Website) bool) WebsiteSlice {
	var nws WebsiteSlice
	for _, v := range ws {
		if v != nil && f(v) {
			nws = append(nws, v)
		}
	}
	return nws
}

// Codes returns a StringSlice with all website codes
func (ws WebsiteSlice) Codes() utils.StringSlice {
	if len(ws) == 0 {
		return nil
	}
	var c utils.StringSlice
	for _, w := range ws {
		if w != nil {
			c.Append(w.Data.Code.String)
		}
	}
	return c
}

// IDs returns an Int64Slice with all website ids
func (ws WebsiteSlice) IDs() utils.Int64Slice {
	if len(ws) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, w := range ws {
		if w != nil {
			ids.Append(w.Data.WebsiteID)
		}
	}
	return ids
}

/*
	@todo implement Magento\Store\Model\Website
*/
