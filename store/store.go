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
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"fmt"

	"github.com/corestoreio/csfw/backend"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
)

const (
	// DefaultStoreID is always 0.
	DefaultStoreID int64 = 0
	// HTTPRequestParamStore name of the GET parameter to set a new store in a
	// current website/group context
	HTTPRequestParamStore = `___store`
	// ParamName important when the user selects a different store within the
	// current website/group context. This name will be used in a cookie or as
	// key value in a token to permanently save the new selected
	// store code.
	ParamName = `store`
	// PriceScopeGlobal prices are for all stores and websites the same.
	PriceScopeGlobal = 0
	// PriceScopeWebsite prices are in each website different.
	PriceScopeWebsite = 1
)

// Store represents the scope in which a shop runs. Everything is bound to a
// Store. A store knows its website ID, group ID and if its active. A store can
// have its own configuration settings which overrides the default scope and
// website scope.
type Store struct {
	cr config.Getter // internal root config.Reader which can be overwritten
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

// StoreSlice a collection of pointers to the Store structs.
// StoreSlice has some nifty method receivers.
type StoreSlice []*Store

// StoreOption can be used as an argument in NewStore to configure a store.
type StoreOption func(s *Store)

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

// WithStoreConfig sets the config.Reader to the Store. Default reader is
// config.DefaultService. You should call this function before calling other
// option functions otherwise your preferred config.Reader won't be inherited
// to a Website or a Group.
func WithStoreConfig(cr config.Getter) StoreOption { return func(s *Store) { s.cr = cr } }

// NewStore creates a new Store. Returns an error if the first three arguments
// are nil. Returns an error if integrity checks fail. config.Reader will be
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
		cr:      config.DefaultService,
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
	if nil != s.Website && nil != s.Group {
		s.Config = s.cr.NewScoped(s.Website.WebsiteID(), s.Group.GroupID(), s.StoreID())
	}
	return s
}

/*
	TODO(cs) implement Magento\Store\Model\Store
*/

var _ scope.StoreIDer = (*Store)(nil)
var _ scope.GroupIDer = (*Store)(nil)
var _ scope.WebsiteIDer = (*Store)(nil)
var _ scope.StoreCoder = (*Store)(nil)

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

	var p model.BaseURL
	switch ut {
	case config.URLTypeWeb:
		p = PathUnsecureBaseURL
		if isSecure {
			p = PathSecureBaseURL
		}
		break
	case config.URLTypeStatic:
		p = PathUnsecureBaseStaticURL
		if isSecure {
			p = PathSecureBaseStaticURL
		}
		break
	case config.URLTypeMedia:
		p = PathUnsecureBaseMediaURL
		if isSecure {
			p = PathSecureBaseMediaURL
		}
		break
	case config.URLTypeAbsent: // hack to clear the cache
		_ = s.urlcache.unsecure.Clear()
		return url.URL{}, s.urlcache.secure.Clear()
	// TODO(cs) rethink that here and maybe add the other paths if needed.
	default:
		return url.URL{}, fmt.Errorf("Unsupported UrlType: %d", ut)
	}

	rawURL := p.Get(PackageConfiguration, s.Config)

	if strings.Contains(rawURL, PlaceholderBaseURL) {
		// TODO(cs) replace placeholder with \Magento\Framework\App\Request\Http::getDistroBaseUrl()
		// getDistroBaseUrl will be generated from the $_SERVER variable,
		base, err := s.cr.String(config.Path(config.PathCSBaseURL))
		if config.NotKeyNotFoundError(err) {
			PkgLog.Debug("store.Store.BaseURL.String", "err", err, "path", config.PathCSBaseURL)
			base = config.CSBaseURL
		}
		rawURL = strings.Replace(rawURL, PlaceholderBaseURL, base, 1)
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
	return backend.PathWebSecureUseInFrontend.Get(PackageConfiguration, s.Config)
}

// IsCurrentlySecure checks if a request for a give store aka. scope is secure. Checks
// include if base URL has been set and if front URL is secure
// This function might gets executed on every request.
func (s *Store) IsCurrentlySecure(r *http.Request) bool {
	if httputil.IsSecure(s.cr, r) {
		return true
	}

	secureBaseURL, err := s.BaseURL(config.URLTypeWeb, true)
	if err != nil || false == s.IsFrontURLSecure() {
		PkgLog.Debug("store.Store.IsCurrentlySecure.BaseURL", "err", err, "secureBaseURL", secureBaseURL)
		return false
	}
	return secureBaseURL.Scheme == "https" && r.URL.Scheme == "https" // todo(cs) check for ports !? other schemes?
}

// NewCookie creates a new pre-configured cookie.
// TODO(cs) create cookie manager to stick to the limits of http://www.ietf.org/rfc/rfc2109.txt page 15
// @see http://browsercookielimits.squawky.net/
func (s *Store) NewCookie() *http.Cookie {
	return &http.Cookie{
		Name:     ParamName,
		Value:    "",
		Path:     s.Path(),
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
	}
}

// SetCookie adds a cookie which contains the store code and is valid for one year.
func (s *Store) SetCookie(res http.ResponseWriter) {
	if res != nil {
		keks := s.NewCookie()
		keks.Value = s.Data.Code.String
		keks.Expires = time.Now().AddDate(1, 0, 0) // one year valid
		http.SetCookie(res, keks)
	}
}

// DeleteCookie deletes the store cookie
func (s *Store) DeleteCookie(res http.ResponseWriter) {
	if res != nil {
		keks := s.NewCookie()
		keks.Expires = time.Now().AddDate(-10, 0, 0)
		http.SetCookie(res, keks)
	}
}

// AddClaim adds the store code to a JSON web token.
// tokenClaim may be *jwt.Token.Claim
func (s *Store) AddClaim(tokenClaim map[string]interface{}) {
	tokenClaim[ParamName] = s.Data.Code.String
}

// RootCategoryID returns the root category ID assigned to this store view.
func (s *Store) RootCategoryID() int64 {
	return s.Group.Data.RootCategoryID
}

/*
	Store Currency
*/

// CurrentCurrency TODO(cs)
// @see app/code/Magento/Store/Model/Store.php::getCurrentCurrency
func (s *Store) CurrentCurrency() *directory.Currency {
	return nil
}

/*
	StoreSlice method receivers
*/

// Sort convenience helper
func (ss *StoreSlice) Sort() *StoreSlice {
	sort.Sort(ss)
	return ss
}

func (ss StoreSlice) Len() int { return len(ss) }

func (ss *StoreSlice) Swap(i, j int) { (*ss)[i], (*ss)[j] = (*ss)[j], (*ss)[i] }

func (ss *StoreSlice) Less(i, j int) bool {
	return (*ss)[i].Data.SortOrder < (*ss)[j].Data.SortOrder
}

// Filter returns a new slice filtered by predicate f
func (ss StoreSlice) Filter(f func(*Store) bool) StoreSlice {
	var stores StoreSlice
	for _, v := range ss {
		if v != nil && f(v) {
			stores = append(stores, v)
		}
	}
	return stores
}

// Codes returns a StringSlice with all store codes
func (ss StoreSlice) Codes() util.StringSlice {
	if len(ss) == 0 {
		return nil
	}
	var c util.StringSlice
	for _, st := range ss {
		if st != nil {
			c.Append(st.Data.Code.String)
		}
	}
	return c
}

// IDs returns an Int64Slice with all store ids
func (ss StoreSlice) IDs() util.Int64Slice {
	if len(ss) == 0 {
		return nil
	}
	var ids util.Int64Slice
	for _, st := range ss {
		if st != nil {
			ids.Append(st.Data.StoreID)
		}
	}
	return ids
}

// LastItem returns the last item of this slice or nil
func (ss StoreSlice) LastItem() *Store {
	if ss.Len() > 0 {
		return ss[ss.Len()-1]
	}
	return nil
}
