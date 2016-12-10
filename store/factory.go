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
	"context"
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/util/errors"
	"golang.org/x/sync/errgroup"
)

// factory contains the raw slices from the database and can read from the
// database. It creates for each call to each of its method receivers new
// Stores, Groups or Websites.
type factory struct {
	// rootConfig parent config service. can only be set once.
	rootConfig config.Getter
	mu         sync.RWMutex
	websites   TableWebsiteSlice
	groups     TableGroupSlice
	stores     TableStoreSlice
}

// newFactory creates a new object which handles the raw data from the three
// database tables for website, group and store. You can either provide the raw
// data separately for each type or pass an option to load it from the database.
// To set the raw data either call the WithTable*() functions or use ReInit()
// and a DB connection.
func newFactory(cfg config.Getter, opts ...Option) (*factory, error) {
	s := &factory{
		rootConfig: cfg,
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(s); err != nil {
				return nil, errors.Wrap(err, "[store] NewStorage Option")
			}
		}
	}
	return s, nil
}

// website returns a TableWebsite by using the id.
func (f *factory) website(id int64) (*TableWebsite, bool) {
	return f.websites.FindByWebsiteID(id)
}

// Website creates a new Website  from an ID including all of its groups
// and all related stores. Returns a NotFound error behaviour.
func (f *factory) Website(id int64) (Website, error) {
	w, found := f.website(id)
	if !found {
		return Website{}, errors.NewNotFoundf("[store] WebsiteID %d", id)
	}
	return NewWebsite(f.rootConfig, w, f.groups, f.stores)
}

// Websites creates a slice containing all new pointers to Websites with its
// associated new groups and new store pointers. It returns an error if the
// integrity is incorrect or NotFound errors.
func (f *factory) Websites() (WebsiteSlice, error) {
	websites := make(WebsiteSlice, len(f.websites), len(f.websites))
	for i, w := range f.websites {
		var err error
		websites[i], err = NewWebsite(f.rootConfig, w, f.groups, f.stores)
		if err != nil {
			return nil, errors.Wrapf(err, "[store] Storage.Websites. WebsiteID: %d", w.WebsiteID)
		}
	}
	return websites, nil
}

// group returns a TableGroup by using a group id as argument.
func (f *factory) group(id int64) (*TableGroup, bool) {
	return f.groups.FindByGroupID(id)
}

// Group creates a new Group  for an ID which contains all related store-
// and its website-pointers.
func (f *factory) Group(id int64) (Group, error) {
	g, found := f.group(id)
	if !found {
		return Group{}, errors.NewNotFoundf("[store] Group %d", id)
	}

	w, found := f.website(g.WebsiteID)
	if !found {
		return Group{}, errors.NewNotFoundf("[store] Website. WebsiteID %d GroupID %v", g.WebsiteID, id)
	}
	return NewGroup(f.rootConfig, g, w, f.stores)
}

// Groups creates a slice containing all pointers to Groups with its associated
// new store- and new website-pointers. It returns an error if the integrity is
// incorrect or a NotFound error.
func (f *factory) Groups() (GroupSlice, error) {
	groups := make(GroupSlice, len(f.groups), len(f.groups))
	for i, g := range f.groups {
		w, found := f.website(g.WebsiteID)
		if !found {
			return nil, errors.NewNotFoundf("[store] WebsiteID %d not found in Group Slice with its Website IDs %v", g.WebsiteID, f.websites.Extract().WebsiteID())
		}
		var err error
		groups[i], err = NewGroup(f.rootConfig, g, w, f.stores)
		if err != nil {
			return nil, errors.Wrapf(err, "[store] GroupID %d WebsiteID %d", g.GroupID, g.WebsiteID)
		}
	}
	return groups, nil
}

// store returns a TableStore by an id.
func (f *factory) store(id int64) (*TableStore, bool) {
	return f.stores.FindByStoreID(id)
}

// Store creates a new Store  containing its group and its website.
// Returns an error if the integrity is incorrect. May return a NotFound error
// behaviour.
func (f *factory) Store(id int64) (Store, error) {
	var ns Store
	s, found := f.store(id)
	if !found {
		return ns, errors.NewNotFoundf("[store] Store: %d", id)
	}
	w, found := f.website(s.WebsiteID)
	if !found {
		return ns, errors.NewNotFoundf("[store] WebsiteID: %d", s.WebsiteID)
	}
	g, found := f.group(s.GroupID)
	if !found {
		return ns, errors.NewNotFoundf("[store] GroupID: %d", s.GroupID)
	}
	var err error
	ns, err = NewStore(f.rootConfig, s, w, g)
	if err != nil {
		return ns, errors.Wrapf(err, "[store] StoreID %d WebsiteID %d GroupID %d", s.StoreID, w.WebsiteID, g.GroupID)
	}
	if err := ns.Website.SetGroupsStores(f.groups, f.stores); err != nil {
		return ns, errors.Wrap(err, "")
	}
	if err := ns.Group.SetWebsiteStores(f.rootConfig, w, f.stores); err != nil {
		return ns, errors.Wrap(err, "[store] Storage.Store.Group.SetWebsiteStores")
	}
	return ns, nil
}

// Stores creates a new store slice with all of its new Group and new Website
// pointers. Can return an error when the website or the group cannot be found.
func (f *factory) Stores() (StoreSlice, error) {
	stores := make(StoreSlice, len(f.stores))
	for i, s := range f.stores {
		var err error
		if stores[i], err = f.Store(s.StoreID); err != nil {
			return nil, errors.Wrapf(err, "[store] StoreID %d", s.StoreID)
		}
	}
	return stores, nil
}

// DefaultStoreID traverses through the websites to find the default website
// and gets the default group which has the default store id assigned to. Only
// one website can be the default one.
func (f *factory) DefaultStoreID() (int64, error) {
	for _, w := range f.websites {
		if w.IsDefault.Bool && w.IsDefault.Valid {
			g, found := f.group(w.DefaultGroupID)
			if !found {
				return 0, errors.NewNotFoundf("[store] WebsiteID %d DefaultGroupID %d", w.WebsiteID, w.DefaultGroupID)
			}
			return g.DefaultStoreID, nil
		}
	}
	return 0, errors.NewNotFoundf(errStoreIDDefaultNotFound)
}

// LoadFromResource reloads all websites, groups and stores concurrently from a
// resource. On error  all internal slices will be reset to nil.
func (f *factory) LoadFromResource(twr TableWebsitesResourcer, tgr TableGroupsResourcer, tsr TableStoresResourcer) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// todo investigate how to use ctx to cancel the DB query and under which condition
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() (err error) {
		f.websites, err = twr.Select()
		return errors.Wrap(err, "[store] SQLSelect Websites")
	})
	eg.Go(func() (err error) {
		f.groups, err = tgr.Select()
		return errors.Wrap(err, "[store] SQLSelect Groups")
	})
	eg.Go(func() (err error) {
		f.stores, err = tsr.Select()
		return errors.Wrap(err, "[store] SQLSelect Stores")
	})

	return errors.Wrap(eg.Wait(), "[store] LoadFromDB.Wait")

}
