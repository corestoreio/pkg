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
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
	"golang.org/x/text/language"
)

const (
	// DefaultWebsiteID is always 0
	DefaultWebsiteID int64 = 0
)

type (
	// Website represents the overall parent structure of its children Group and Store.
	// A website defines the default group ID. A website can contain custom configuration
	// settings which overrides the default scope but get itself overridden by the Store scope.
	Website struct {
		cr config.Reader
		// Data raw website data
		Data *TableWebsite

		// Groups contains a slice to all groups associated to one website. This slice can be nil.
		Groups GroupSlice
		// Stores contains a slice to all stores associated to one website. This slice can be nil.
		Stores StoreSlice
	}
	// WebsiteSlice contains pointer to Website struct and some nifty method receivers.
	WebsiteSlice []*Website

	// WebsiteOption option func for NewWebsite()
	WebsiteOption func(*Website)
)

var (
	// ErrWebsiteNotFound when the website has not been found within a slice
	ErrWebsiteNotFound = errors.New("Website not found")
	// ErrWebsiteDefaultGroupNotFound the default group cannot be found
	ErrWebsiteDefaultGroupNotFound = errors.New("Website Default Group not found")
)

var _ config.ScopeIDer = (*Website)(nil)
var _ config.ScopeCoder = (*Website)(nil)

// SetWebsiteConfig sets the config.Reader to the Website.
// Default reader is config.DefaultManager
func SetWebsiteConfig(cr config.Reader) WebsiteOption {
	return func(w *Website) { w.cr = cr }
}

// NewWebsite returns a new pointer to a Website.
func NewWebsite(tw *TableWebsite, opts ...WebsiteOption) *Website {
	if tw == nil {
		panic(ErrStoreNewArgNil)
	}
	w := &Website{
		cr:   config.DefaultManager,
		Data: tw,
	}
	w.ApplyOptions(opts...)
	return w
}

// ApplyOptions sets the options on a Website
func (w *Website) ApplyOptions(opts ...WebsiteOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(w)
		}
	}
}

// ScopeID satisfies the interface ScopeIDer and mainly used in the StoreManager for selecting Website,Group ...
func (w *Website) ScopeID() int64 { return w.Data.WebsiteID }

// ScopeCode satisfies the interface ScopeCoder
func (w *Website) ScopeCode() string { return w.Data.Code.String }

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
	groups := tgs.FilterByWebsiteID(w.Data.WebsiteID)
	w.Groups = make(GroupSlice, groups.Len(), groups.Len())
	for i, g := range groups {
		w.Groups[i] = NewGroup(g, SetGroupWebsite(w.Data), SetGroupConfig(w.cr)).SetStores(tss, nil)
	}
	stores := tss.FilterByWebsiteID(w.Data.WebsiteID)
	w.Stores = make(StoreSlice, stores.Len(), stores.Len())
	for i, s := range stores {
		group, err := tgs.FindByID(s.GroupID)
		if err != nil {
			panic(fmt.Sprintf("Integrity error. A store %#v must be assigned to a group.\nGroupSlice: %#v\n\n", s, tgs))
		}
		w.Stores[i] = NewStore(s, w.Data, group, SetStoreConfig(w.cr))
	}
	return w
}

// ConfigString tries to get a value from the scopeStore if empty
// falls back to default global scope.
// If using etcd or consul maybe this can lead to round trip times because of network access.
func (w *Website) ConfigString(path ...string) string {
	val := w.cr.GetString(config.ScopeWebsite(w), config.Path(path...))
	if val == "" {
		val = w.cr.GetString(config.Path(path...))
	}
	return val
}

// @todo
func (w *Website) BaseCurrencyCode() (language.Currency, error) {
	var c string
	if w.ConfigString(PathPriceScope) == PriceScopeGlobal {
		c = w.cr.GetString(config.Path(directory.PathCurrencyBase))
	} else {
		c = w.ConfigString(directory.PathCurrencyBase)
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
	TableWebsite and TableWebsiteSlice method receivers
*/

// Load uses a dbr session to load all data from the core_website table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Website/Collection.php::Load()
func (s *TableWebsiteSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndexWebsite, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		return sb.OrderBy("main_table.sort_order ASC").OrderBy("main_table.name ASC")
	})...)
}

// FindByID returns a TableWebsite if found by id or an error
func (s TableWebsiteSlice) FindByID(id int64) (*TableWebsite, error) {
	for _, w := range s {
		if w != nil && w.WebsiteID == id {
			return w, nil
		}
	}
	return nil, ErrWebsiteNotFound
}

// FindByCode returns a TableWebsite if found by code or an error
func (s TableWebsiteSlice) FindByCode(code string) (*TableWebsite, error) {
	for _, w := range s {
		if w != nil && w.Code.Valid && w.Code.String == code {
			return w, nil
		}
	}
	return nil, ErrWebsiteNotFound
}

// Filter returns a new slice filtered by predicate f
func (s TableWebsiteSlice) Filter(f func(*TableWebsite) bool) TableWebsiteSlice {
	var tws TableWebsiteSlice
	for _, w := range s {
		if w != nil && f(w) {
			tws = append(tws, w)
		}
	}
	return tws
}

// @todo review Magento code because of column is_default
//func (s TableWebsite) IsDefault() bool {
//	return s.WebsiteID == DefaultWebsiteId
//}

/*
	@todo implement Magento\Store\Model\Website
*/
