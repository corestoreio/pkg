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
	"fmt"

	"github.com/corestoreio/csfw/catalog/catconfig"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
)

// DefaultWebsiteID is always 0
const DefaultWebsiteID int64 = 0

// Website represents the overall parent structure of its children Group and Store.
// A website defines the default group ID. A website can contain custom configuration
// settings which overrides the default scope but get itself overridden by the Store scope.
type Website struct {
	cr config.Getter // internal root config.Reader which can be overridden

	// Config contains the scope based configuration reader.
	Config config.ScopedGetter
	// Data raw website data from DB table.
	Data *TableWebsite

	// Groups contains a slice to all groups associated to one website. This slice can be nil.
	Groups GroupSlice
	// Stores contains a slice to all stores associated to one website. This slice can be nil.
	Stores StoreSlice
	*cserr.MultiErr
}

// WebsiteOption can be used as an argument in NewWebsite to configure a website.
type WebsiteOption func(*Website)

// ErrWebsite* are general errors when handling with the Website type.
// They are self explanatory.
var (
	ErrWebsiteNotFound             = errors.New("Website not found")
	ErrWebsiteDefaultGroupNotFound = errors.New("Website Default Group not found")
)

// SetWebsiteConfig sets the config.Reader to the Website. Default reader is
// config.DefaultManager. You should call this function before calling other
// option functions otherwise your preferred config.Reader won't be inherited
// to a Group or Store.
func SetWebsiteConfig(cr config.Getter) WebsiteOption {
	return func(w *Website) {
		w.cr = cr
	}
}

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
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.SetWebsiteGroupsStores.NewGroup", "err", err, "g", g, "w", w.Data)
				}
				w.MultiErr = w.AppendErrors(errors.Mask(err))
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
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.SetWebsiteGroupsStores.NewStore", "err", err, "s", s, "w.Data", w.Data, "group", group)
				}
				w.MultiErr = w.AppendErrors(errors.Mask(err))
				return
			}
		}
	}
}

// NewWebsite creates a new website pointer with the config.DefaultManager.
func NewWebsite(tw *TableWebsite, opts ...WebsiteOption) (*Website, error) {
	if tw == nil {
		return nil, ErrArgumentCannotBeNil
	}
	w := &Website{
		cr:   config.DefaultService,
		Data: tw,
	}
	return w.ApplyOptions(opts...)
}

// MustNewWebsite same as NewWebsite but panics on error.
func MustNewWebsite(tw *TableWebsite, opts ...WebsiteOption) *Website {
	w, err := NewWebsite(tw, opts...)
	if err != nil {
		panic(err)
	}
	return w
}

// ApplyOptions sets the options on a Website
func (w *Website) ApplyOptions(opts ...WebsiteOption) (*Website, error) {
	for _, opt := range opts {
		if opt != nil {
			opt(w)
		}
	}
	if w.HasErrors() {
		return nil, w
	}
	w.Config = w.cr.NewScoped(w.WebsiteID(), 0) // Scope Store is not available
	return w, nil
}

// WebsiteID satisfies the interface scope.WebsiteIDer and returns the website ID.
func (w *Website) WebsiteID() int64 { return w.Data.WebsiteID }

// WebsiteCode satisfies the interface scope.WebsiteCoder and returns the code.
func (w *Website) WebsiteCode() string { return w.Data.Code.String }

// GroupID implements the GroupIDer interface and returns the default group ID.
func (w *Website) GroupID() int64 {
	return w.Data.DefaultGroupID
}

// StoreID implements the StoreIDer interface and returns the default store ID.
// It may return a scope.UnavailableStoreID when finding the DefaultGroup()
// returns an error. Error will be logged.
func (w *Website) StoreID() int64 {
	g, err := w.DefaultGroup()
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("store.Website.StoreID", "err", err, "Website", w)
		}
		return scope.UnavailableStoreID
	}
	return g.Data.DefaultStoreID
}

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

// BaseCurrency returns the base currency code of a website.
// 	1st argument should be a path to catalog/price/scope
// 	2nd argument should be a path to currency/options/base
func (w *Website) BaseCurrency(ps catconfig.PriceScope, cc directory.ConfigCurrency) (directory.Currency, error) {

	isGlobal, err := ps.IsGlobal(w.Config)
	if err != nil {
		return directory.Currency{}, errors.Mask(err)
	}
	if isGlobal {
		return cc.GetDefault(w.cr) // default scope
	}
	return cc.Get(w.Config) // website scope
}

/*
	@todo implement Magento\Store\Model\Website
*/
