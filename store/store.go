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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/corestoreio/csfw/backend"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
)

// DefaultStoreID is always 0.
const DefaultStoreID int64 = 0

// Store represents the scope in which a shop runs. Everything is bound to a
// Store. A store knows its website ID, group ID and if its active. A store can
// have its own configuration settings which overrides the default scope and
// website scope.
type Store struct {
	cr config.Getter // internal root config.Getter which can be overwritten
	// Config contains a config.Service which takes care of the scope based
	// configuration values.
	Config config.ScopedGetter
	// Website points to the current website for this store. No integrity checks.
	// Can be nil.
	Website *Website
	// Group points to the current store group for this store. No integrity
	// checks. Can be nil.
	Group *Group
	// Data underlying raw data
	Data *TableStore

	urlcache *struct {
		secure   *config.URLCache
		unsecure *config.URLCache
	}
}

// ErrStore* are general errors when handling with the Store type.
// They are self explanatory.
var (
	ErrStoreNotFound         = errors.New("Store not found")
	ErrStoreNotActive        = errors.New("Store not active")
	ErrArgumentCannotBeNil   = errors.New("An argument cannot be nil")
	ErrStoreIncorrectGroup   = errors.New("Incorrect group")
	ErrStoreIncorrectWebsite = errors.New("Incorrect website")
	ErrStoreCodeInvalid      = errors.New("The store code may contain only letters (a-z), numbers (0-9) or underscore(_). The first character must be a letter")
)

// NewStore creates a new Store. Returns an error if the first three arguments
// are nil. Returns an error if integrity checks fail. config.Getter will be
// also set to Group and Website.
func NewStore(ts *TableStore, tw *TableWebsite, tg *TableGroup, opts ...StoreOption) (s *Store, err error) {
	if ts == nil || tw == nil || tg == nil {
		return nil, ErrArgumentCannotBeNil
	}
	if ts.WebsiteID != tw.WebsiteID {
		return nil, ErrStoreIncorrectWebsite
	}
	if tg.WebsiteID != tw.WebsiteID {
		return nil, ErrStoreIncorrectWebsite
	}
	if ts.GroupID != tg.GroupID {
		return nil, ErrStoreIncorrectGroup
	}

	var nw *Website
	if nw, err = NewWebsite(tw); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("store.NewStore.NewWebsite", "err", err, "tw", tw)
		}
		return
	}

	var ng *Group
	if ng, err = NewGroup(tg, SetGroupWebsite(tw)); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("store.NewStore.NewGroup", "err", err, "tg", tg, "tw", tw)
		}
		return
	}

	s = &Store{
		Data:    ts,
		Website: nw,
		Group:   ng,
		urlcache: &struct {
			secure   *config.URLCache
			unsecure *config.URLCache
		}{
			secure:   config.NewURLCache(),
			unsecure: config.NewURLCache(),
		},
	}
	s.ApplyOptions(opts...)
	if _, err = s.Website.ApplyOptions(SetWebsiteConfig(s.cr)); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("store.Website.ApplyOptions", "err", err, "tg", tg, "tw", tw)
		}
		return
	}
	if _, err = s.Group.ApplyOptions(SetGroupConfig(s.cr)); err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("store.Group.ApplyOptions", "err", err, "tg", tg, "tw", tw)
		}
		return
	}
	return s, nil
}

// MustNewStore same as NewStore except that it panics on an error.
func MustNewStore(ts *TableStore, tw *TableWebsite, tg *TableGroup, opts ...StoreOption) *Store {
	s, err := NewStore(ts, tw, tg, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// ApplyOptions sets the options to the Store struct.
func (s *Store) ApplyOptions(opts ...StoreOption) *Store {
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	if nil != s.Website && nil != s.Group && s.cr != nil {
		s.Config = s.cr.NewScoped(s.Website.WebsiteID(), s.StoreID())
	}
	return s
}

/*
	TODO(cs) implement Magento\Store\Model\Store
*/

// StoreID satisfies the interface scope.StoreIDer and returns the store ID.
func (s *Store) StoreID() int64 {
	return s.Data.StoreID
}

// StoreCode satisfies the interface scope.StoreCoder and returns the store code.
func (s *Store) StoreCode() string {
	return s.Data.Code.String
}

// GroupID implements scope.GroupIDer interface
func (s *Store) GroupID() int64 {
	return s.Data.GroupID
}

// WebsiteID implements scope.WebsiteIDer interface
func (s *Store) WebsiteID() int64 {
	return s.Data.WebsiteID
}

// MarshalJSON satisfies interface for JSON marshalling. The TableStore
// struct will be encoded to JSON using Go's standard library.
func (s *Store) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Data)
}

// Path returns the sub path from the URL where CoreStore is installed
func (s *Store) Path() string {
	url, err := s.BaseURL(config.URLTypeWeb, false)
	if err != nil {
		return "/"
	}
	return url.Path
}

// BaseURL returns a parsed and maybe cached URL from config.ScopedReader.
// It returns a copy of url.URL or an error. Possible URLTypes are:
//     - config.URLTypeWeb
//     - config.URLTypeStatic
//     - config.URLTypeMedia
func (s *Store) BaseURL(ut config.URLType, isSecure bool) (url.URL, error) {

	switch isSecure {
	case true:
		if pu := s.urlcache.secure.Get(ut); pu != nil {
			return *pu, nil
		}
	case false:
		if pu := s.urlcache.unsecure.Get(ut); pu != nil {
			return *pu, nil
		}
	}

	var p cfgmodel.BaseURL
	switch ut {
	case config.URLTypeWeb:
		p = backend.Backend.WebUnsecureBaseURL
		if isSecure {
			p = backend.Backend.WebSecureBaseURL
		}
		break
	case config.URLTypeStatic:
		p = backend.Backend.WebUnsecureBaseStaticURL
		if isSecure {
			p = backend.Backend.WebSecureBaseStaticURL
		}
		break
	case config.URLTypeMedia:
		p = backend.Backend.WebUnsecureBaseMediaURL
		if isSecure {
			p = backend.Backend.WebSecureBaseMediaURL
		}
		break
	case config.URLTypeAbsent: // hack to clear the cache :-( refactor that
		_ = s.urlcache.unsecure.Clear()
		return url.URL{}, s.urlcache.secure.Clear()
	// TODO(cs) rethink that here and maybe add the other paths if needed.
	default:
		return url.URL{}, fmt.Errorf("Unsupported UrlType: %d", ut)
	}

	rawURL, err := p.Get(s.Config)
	if err != nil {
		return url.URL{}, err
	}

	if strings.Contains(rawURL, cfgmodel.PlaceholderBaseURL) {
		// TODO(cs) replace placeholder with \Magento\Framework\App\Request\Http::getDistroBaseUrl()
		// getDistroBaseUrl will be generated from the $_SERVER variable,
		base, err := s.cr.String(cfgpath.MustNewByParts(config.PathCSBaseURL))
		if config.NotKeyNotFoundError(err) {
			PkgLog.Debug("store.Store.BaseURL.String", "err", err, "path", config.PathCSBaseURL)
			base = config.CSBaseURL
		}
		rawURL = strings.Replace(rawURL, cfgmodel.PlaceholderBaseURL, base, 1)
	}
	rawURL = strings.TrimRight(rawURL, "/") + "/"

	if isSecure {
		retURL, retErr := s.urlcache.secure.Set(ut, rawURL)
		return *retURL, retErr
	}
	retURL, retErr := s.urlcache.unsecure.Set(ut, rawURL)
	return *retURL, retErr
}

// IsFrontURLSecure returns true from the config if the frontend must be secure.
func (s *Store) IsFrontURLSecure() bool {
	return false // backend.Backend.WebSecureUseInFrontend.Get(s.Config)
}

// IsCurrentlySecure checks if a request for a give store aka. scope is secure. Checks
// include if base URL has been set and if front URL is secure
// This function might gets executed on every request.
func (s *Store) IsCurrentlySecure(r *http.Request) bool {
	return false
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
}

// TOOD move net related functions into the storenet package

// RootCategoryID returns the root category ID assigned to this store view.
func (s *Store) RootCategoryID() int64 {
	return s.Group.Data.RootCategoryID
}

/*
	Store Currency
*/

// CurrentCurrency TODO(cs)
// @see app/code/Magento/Store/Model/Store.php::getCurrentCurrency
func (s *Store) CurrentCurrency() string {
	/*
		this returns just a string or string slice and no further
		involvement of the directory package.

		or those functions move directly into the directory package
	*/
	return ""
}

func (s *Store) DefaultCurrency() string {
	return ""
}

func (s *Store) AvailableCurrencyCodes() []string {
	return nil
}
