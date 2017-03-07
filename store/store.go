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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// DefaultStoreID is always 0.
const DefaultStoreID int64 = 0

// Store represents the scope in which a shop runs. Everything is bound to a
// Store. A store knows its website ID, group ID and if its active. A store can
// have its own configuration settings which overrides the default scope and
// website scope.
type Store struct {
	// Config contains the scoped configuration which cannot be changed once the
	// object has been created.
	Config config.Scoped
	// Data underlying raw data
	Data *TableStore
	// Website points to the current website for this store. No integrity checks.
	// Can be nil.
	Website Website
	// Group points to the current store group for this store. No integrity
	// checks. Can be nil.
	Group Group
}

// NewStore creates a new Store. Returns an error if the first three arguments
// are nil. Returns an error if integrity checks fail. config.Getter will be
// also set to Group and Website.
func NewStore(cfg config.Getter, ts *TableStore, tw *TableWebsite, tg *TableGroup) (Store, error) {
	// ts cannot be nil but tw and tg can be nil
	s := Store{
		Data: ts,
	}
	if err := s.SetWebsiteGroup(cfg, tw, tg); err != nil {
		return Store{}, errors.Wrap(err, "[store] NewStore.SetWebsiteGroup")
	}
	return s, nil
}

// MustNewStore same as NewStore except that it panics on an error.
func MustNewStore(cfg config.Getter, ts *TableStore, tw *TableWebsite, tg *TableGroup) Store {
	s, err := NewStore(cfg, ts, tw, tg)
	if err != nil {
		panic(err)
	}
	return s
}

// SetWebsiteGroup uses a raw website and a table store slice to set the groups
// associated to this website and the stores associated to this website. It
// returns an error if the data integrity is incorrect.
func (s *Store) SetWebsiteGroup(cfg config.Getter, tw *TableWebsite, tg *TableGroup) error {

	s.Config = cfg.NewScoped(s.WebsiteID(), s.ID())

	if tw == nil && tg == nil {
		// avoid recursion and stack overflow
		return nil
	}

	var err error
	s.Website, err = NewWebsite(cfg, tw, nil, nil) // avoid recursion so store must be nil
	if err != nil {
		return errors.Wrapf(err, "[store] Store.SetWebsiteGroup.NewWebsite")
	}
	if s.Group, err = NewGroup(cfg, tg, nil, nil); err != nil { // avoid recursion so store must be nil
		return errors.Wrapf(err, "[store] TableGroup: %#v\nTableWebsite: %#v\n", tg, tw)
	}

	return s.Validate()
}

// Validate checks the internal integrity. May panic when the data has not been
// set.
func (s Store) Validate() (err error) {
	switch {
	case s.WebsiteID() != s.Website.ID():
		err = errors.NewNotValidf("[store] Store %d: WebsiteID %d != Website.ID %d", s.ID(), s.WebsiteID(), s.Website.ID())

	case s.Group.Website.Data != nil && s.Group.Website.ID() != s.WebsiteID():
		err = errors.NewNotValidf("[store] Store %d: Group.WebsiteID %d != Website.ID %d", s.ID(), s.Group.Website.ID(), s.WebsiteID())

	case s.GroupID() != s.Group.ID():
		err = errors.NewNotValidf("[store] Store %d: Store.GroupID %d != Group.ID %d", s.ID(), s.GroupID(), s.Group.ID())

	case s.Config.WebsiteID != s.WebsiteID():
		err = errors.NewNotValidf("[store] Store %d: Website ID %d != Config Website ID %d", s.ID(), s.WebsiteID(), s.Config.WebsiteID)

	case s.Config.StoreID != s.ID():
		err = errors.NewNotValidf("[store] Store %d: Store ID %d != Config Store ID %d", s.ID(), s.ID(), s.Config.StoreID)
	}
	return err
}

// ID returns the store id. If Data is nil returns -1.
func (s Store) ID() int64 {
	if s.Data == nil {
		return -1
	}
	return s.Data.StoreID
}

// Code returns the store code. Returns empty if Data is nil.
func (s Store) Code() string {
	if s.Data == nil {
		return ""
	}
	return s.Data.Code.String
}

// Name returns the store name. Returns empty if Data is nil.
func (s Store) Name() string {
	if s.Data == nil {
		return ""
	}
	return s.Data.Name
}

// GroupID returns the associated group ID. If data is nil returns -1.
func (s Store) GroupID() int64 {
	if s.Data == nil {
		return -1
	}
	return s.Data.GroupID
}

// WebsiteID returns the associated website ID. If data is nil returns -1.
func (s Store) WebsiteID() int64 {
	if s.Data == nil {
		return -1
	}
	return s.Data.WebsiteID
}

// IsActive returns true if data is not nil and store is active.
func (s Store) IsActive() bool {
	if s.Data == nil {
		return false
	}
	return s.Data.IsActive
}

// MarshalJSON satisfies interface for JSON marshalling. The TableStore
// struct will be encoded to JSON using Go's standard library.
func (s Store) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Data)
}

// MarshalLog implements the log.Marshaler interface
func (s Store) MarshalLog(kv log.KeyValuer) error {
	kv.AddInt64("store_id", s.ID())
	kv.AddString("store_code", s.Code())
	kv.AddString("store_name", s.Name())
	kv.AddInt64("store_website_id", s.WebsiteID())
	kv.AddInt64("store_group_id", s.GroupID())
	kv.AddInt64("store_website_website_id", s.Website.ID())
	kv.AddInt64("store_group_group_id", s.Group.ID())
	return nil
}

// String returns human readable information about a Store. The returned string
// may change in the future to provide better informations.
func (s Store) String() string {
	return fmt.Sprintf(`Store[ID:%d Code:%q Name:%q] Website[ID:%d Code:%q Name:%q] Group[ID:%d Name:%q] Valid: %t`,
		s.ID(),
		s.Name(),
		s.Code(),
		s.Website.ID(),
		s.Website.Code(),
		s.Website.Name(),
		s.Group.ID(),
		s.Group.Name(),
		s.Validate() == nil,
	)
}

// RootCategoryID returns the root category ID assigned to this store view.
func (s Store) RootCategoryID() int64 {
	return s.Group.Data.RootCategoryID
}
