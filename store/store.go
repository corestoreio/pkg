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
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
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

// Validate checks the internal integrity. May panic when the data has not been
// set.
func (s Store) Validate() error {
	if s.WebsiteID() != s.Website.ID() {
		return errors.NewNotValidf("[store] NewStore: Store.WebsiteID (%d) != Website.ID (%d)", s.WebsiteID(), s.Website.ID())
	}
	if s.Group.Website.ID() != s.WebsiteID() {
		return errors.NewNotValidf("[store] NewStore: Group.WebsiteID (%d) != Website.ID (%d)", s.Group.Website.ID(), s.WebsiteID())
	}
	if s.GroupID() != s.Group.ID() {
		return errors.NewNotValidf("[store] NewStore: Store.GroupID (%d) != Group.ID (%d)", s.GroupID(), s.Group.ID())
	}
	if s.Config.WebsiteID != s.WebsiteID() {
		return errors.NewNotValidf("[store] Store.Validate: Config Website ID %d does not match Website ID %d", s.Config.WebsiteID, s.WebsiteID())
	}
	if s.Config.StoreID != s.ID() {
		return errors.NewNotValidf("[store] Store.Validate: Config Store ID %d does not match Store ID %d", s.Config.StoreID, s.ID())
	}
	return nil
}

// SetWebsiteGroup uses a raw website and a table store slice to set the groups
// associated to this website and the stores associated to this website. It
// returns an error if the data integrity is incorrect.
func (s *Store) SetWebsiteGroup(cfg config.Getter, tw *TableWebsite, tg *TableGroup) error {
	var err error
	s.Website, err = NewWebsite(cfg, tw, TableGroupSlice{tg}, TableStoreSlice{s.Data})
	if err != nil {
		return errors.Wrapf(err, "[store] Store.SetWebsiteGroup.NewWebsite")
	}
	if s.Group, err = NewGroup(cfg, tg, tw, TableStoreSlice{s.Data}); err != nil {
		return errors.Wrapf(err, "[store] TableGroup: %#v\nTableWebsite: %#v\n", tg, tw)
	}
	s.Config = cfg.NewScoped(tw.WebsiteID, s.ID())
	return s.Validate()
}

// ID returns the store id
func (s Store) ID() int64 {
	return s.Data.StoreID
}

// Code returns the store code.
func (s Store) Code() string {
	return s.Data.Code.String
}

// GroupID returns the associated group ID.
func (s Store) GroupID() int64 {
	return s.Data.GroupID
}

// WebsiteID returns the associated website ID.
func (s Store) WebsiteID() int64 {
	return s.Data.WebsiteID
}

// MarshalJSON satisfies interface for JSON marshalling. The TableStore
// struct will be encoded to JSON using Go's standard library.
func (s Store) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Data)
}

// MarshalLog implements the log.Marshaler interface
func (s Store) MarshalLog(kv log.KeyValuer) error {
	if s.Data != nil {
		kv.AddString("store_code", s.Data.Code.String)
		kv.AddInt64("store_id", s.Data.StoreID)
	}
	return nil
}

//// Path returns the sub path from the URL where CoreStore is installed
//func (s Store) Path() string {
//	url, err := s.BaseURL(config.URLTypeWeb, false)
//	if err != nil {
//		return "/"
//	}
//	return url.Path
//}

// BaseURL returns a parsed and maybe cached URL from config.ScopedReader.
// It returns a copy of url.URL or an error. Possible URLTypes are:
//     - config.URLTypeWeb
//     - config.URLTypeStatic
//     - config.URLTypeMedia
//func (s Store) BaseURL(ut config.URLType, isSecure bool) (url.URL, error) {
//
//	switch isSecure {
//	case true:
//		if pu := s.urlcache.secure.Get(ut); pu != nil {
//			return *pu, nil
//		}
//	case false:
//		if pu := s.urlcache.unsecure.Get(ut); pu != nil {
//			return *pu, nil
//		}
//	}
//
//	var p cfgmodel.BaseURL
//	switch ut {
//	case config.URLTypeWeb:
//		p = backend.Backend.WebUnsecureBaseURL
//		if isSecure {
//			p = backend.Backend.WebSecureBaseURL
//		}
//		break
//	case config.URLTypeStatic:
//		p = backend.Backend.WebUnsecureBaseStaticURL
//		if isSecure {
//			p = backend.Backend.WebSecureBaseStaticURL
//		}
//		break
//	case config.URLTypeMedia:
//		p = backend.Backend.WebUnsecureBaseMediaURL
//		if isSecure {
//			p = backend.Backend.WebSecureBaseMediaURL
//		}
//		break
//	case config.URLTypeAbsent: // hack to clear the cache :-( refactor that
//		_ = s.urlcache.unsecure.Clear()
//		return url.URL{}, s.urlcache.secure.Clear()
//	// TODO(cs) rethink that here and maybe add the other paths if needed.
//	default:
//		return url.URL{}, fmt.Errorf("Unsupported UrlType: %d", ut)
//	}
//
//	rawURL, _, err := p.Get(s.Config)
//	if err != nil {
//		return url.URL{}, err
//	}
//
//	if strings.Contains(rawURL, cfgmodel.PlaceholderBaseURL) {
//		// TODO(cs) replace placeholder with \Magento\Framework\App\Request\Http::getDistroBaseUrl()
//		// getDistroBaseUrl will be generated from the $_SERVER variable,
//		base, err := s.baseConfig.String(cfgpath.MustNewByParts(config.PathCSBaseURL))
//		if err != nil && !errors.IsNotFound(err) {
//			base = config.CSBaseURL
//		}
//		rawURL = strings.Replace(rawURL, cfgmodel.PlaceholderBaseURL, base, 1)
//	}
//	rawURL = strings.TrimRight(rawURL, "/") + "/"
//
//	if isSecure {
//		retURL, retErr := s.urlcache.secure.Set(ut, rawURL)
//		return *retURL, retErr
//	}
//	retURL, retErr := s.urlcache.unsecure.Set(ut, rawURL)
//	return *retURL, retErr
//}

// IsFrontURLSecure returns true from the config if the frontend must be secure.
//func (s Store) IsFrontURLSecure() bool {
//	return false // backend.Backend.WebSecureUseInFrontend.Get(s.Config)
//}

// IsCurrentlySecure checks if a request for a give store aka. scope is secure. Checks
// include if base URL has been set and if front URL is secure
// This function might gets executed on every request.
//func (s Store) IsCurrentlySecure(r *http.Request) bool {
//	return false
//if httputil.IsSecure(s.cr, r) {
//	return true
//}
//
//secureBaseURL, err := s.BaseURL(config.URLTypeWeb, true)
//if err != nil || false == s.IsFrontURLSecure() {
//	PkgLog.Debug("store.Store.IsCurrentlySecure.BaseURL", "err", err, "secureBaseURL", secureBaseURL)
//	return false
//}
//return secureBaseURL.Scheme == "https" && r.URL.Scheme == "https" // todo(cs) check for ports !? other schemes?
//}

// TODO move net related functions into the storenet package

// RootCategoryID returns the root category ID assigned to this store view.
func (s Store) RootCategoryID() int64 {
	return s.Group.Data.RootCategoryID
}

///*
//	Store Currency
//*/
//
//// CurrentCurrency TODO(cs)
//// @see app/code/Magento/Store/Model/Store.php::getCurrentCurrency
//func (s Store) CurrentCurrency() string {
//	/*
//		this returns just a string or string slice and no further
//		involvement of the directory package.
//
//		or those functions move directly into the directory package
//	*/
//	return ""
//}
//
//func (s Store) DefaultCurrency() string {
//	return ""
//}
//
//func (s Store) AvailableCurrencyCodes() []string {
//	return nil
//}
