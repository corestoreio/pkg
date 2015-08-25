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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/directory"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
	"github.com/dgrijalva/jwt-go"
)

const (
	// DefaultStoreID is always 0.
	DefaultStoreID int64 = 0
	// HTTPRequestParamStore name of the GET parameter to set a new store in a current website/group context
	HTTPRequestParamStore = `___store`
	// CookieName important when the user selects a different store within the current website/group context.
	// This cookie permanently saves the new selected store code for one year.
	// The cookie must be removed when the default store of the current website if equal to the current store.
	CookieName = `store`

	// PriceScopeGlobal prices are for all stores and websites the same.
	PriceScopeGlobal = `0` // must be string
	// PriceScopeWebsite prices are in each website different.
	PriceScopeWebsite = `1` // must be string
)

type (
	// Store represents the scope in which a shop runs. Everything is bound to a Store. A store
	// knows its website ID, group ID and if its active. A store can have its own configuration settings
	// which overrides the default scope and website scope.
	Store struct {
		cr config.Reader
		// Website points to the current website for this store. No integrity checks. Can be nil.
		Website *Website
		// Group points to the current store group for this store. No integrity checks. Can be nil.
		Group *Group
		// Data underlying raw data
		Data *TableStore
	}
	// StoreSlice a collection of pointers to the Store structs. StoreSlice has some nifty method receviers.
	StoreSlice []*Store

	// StoreOption option func for NewStore()
	StoreOption func(s *Store)
)

var (
	ErrStoreNotFound         = errors.New("Store not found")
	ErrStoreNotActive        = errors.New("Store not active")
	ErrStoreNewArgNil        = errors.New("An argument cannot be nil")
	ErrStoreIncorrectGroup   = errors.New("Incorrect group")
	ErrStoreIncorrectWebsite = errors.New("Incorrect website")
	ErrStoreCodeInvalid      = errors.New("The store code may contain only letters (a-z), numbers (0-9) or underscore(_). The first character must be a letter")
)

var _ config.ScopeIDer = (*Store)(nil)
var _ config.ScopeCoder = (*Store)(nil)

// SetStoreConfig sets the config.Reader to the Store.
// Default reader is config.DefaultManager
func SetStoreConfig(cr config.Reader) StoreOption {
	return func(s *Store) { s.cr = cr }
}

// NewStore creates a new Store. Panics if TableGroup and TableWebsite have not been provided
// Panics if integrity checks fail. config.Reader will be set to Group and Website.
func NewStore(ts *TableStore, tw *TableWebsite, tg *TableGroup, opts ...StoreOption) *Store {
	if ts == nil || tw == nil || tg == nil { // group and website required so at least 2 args
		panic(ErrStoreNewArgNil)
	}
	if ts.WebsiteID != tw.WebsiteID {
		panic(ErrStoreIncorrectWebsite)
	}
	if ts.GroupID != tg.GroupID {
		panic(ErrStoreIncorrectGroup)
	}
	s := &Store{
		cr:      config.DefaultManager,
		Data:    ts,
		Website: NewWebsite(tw),
		Group:   NewGroup(tg),
	}
	s.ApplyOptions(opts...)
	s.Website.ApplyOptions(SetWebsiteConfig(s.cr))
	s.Group.ApplyOptions(SetGroupConfig(s.cr))
	return s
}

// ApplyOptions sets the options to the Store struct.
func (s *Store) ApplyOptions(opts ...StoreOption) *Store {
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	return s
}

/*
	@todo implement Magento\Store\Model\Store
*/

// ScopeID satisfies the interface ScopeIDer and mainly used in the StoreManager for selecting Website,Group ...
func (s *Store) ScopeID() int64 {
	return s.Data.StoreID
}

// ScopeCode satisfies the interface ScopeCoder
func (s *Store) ScopeCode() string {
	return s.Data.Code.String
}

// MarshalJSON satisfies interface for JSON marshalling. The TableStore
// struct will be encoded to JSON using Go's standard library.
func (s *Store) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Data)
}

// Path returns the sub path from the URL where CoreStore is installed
func (s *Store) Path() string {
	url, err := url.ParseRequestURI(s.BaseURL(config.URLTypeWeb, false))
	if err != nil {
		return "/"
	}
	return url.Path
}

// BaseUrl returns the path from the URL or config where CoreStore is installed @todo
// @see https://github.com/magento/magento2/blob/0.74.0-beta7/app/code/Magento/Store/Model/Store.php#L539
func (s *Store) BaseURL(ut config.URLType, isSecure bool) string {
	var url string
	var p string
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
	// @todo rethink that here and maybe add the other paths if needed.
	default:
		panic("Unsupported UrlType")
	}

	url = s.ConfigString(p)

	if strings.Contains(url, PlaceholderBaseURL) {
		// @todo replace placeholder with \Magento\Framework\App\Request\Http::getDistroBaseUrl()
		// getDistroBaseUrl will be generated from the $_SERVER variable,
		url = strings.Replace(url, PlaceholderBaseURL, s.cr.GetString(config.Path(config.PathCSBaseURL)), 1)
	}
	url = strings.TrimRight(url, "/") + "/"

	return url
}

// ConfigString tries to get a value from the scopeStore if empty
// falls back to default global scope.
// If using etcd or consul maybe this can lead to round trip times because of network access.
func (s *Store) ConfigString(path ...string) string {
	val := s.cr.GetString(config.ScopeStore(s), config.Path(path...))
	if val == "" {
		val = s.cr.GetString(config.Path(path...))
	}
	return val
}

// NewCookie creates a new pre-configured cookie.
// @todo create cookie manager to stick to the limits of http://www.ietf.org/rfc/rfc2109.txt page 15
// @see http://browsercookielimits.squawky.net/
func (s *Store) NewCookie() *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
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

// AddClaim adds the store code to a JSON web token
func (s *Store) AddClaim(t *jwt.Token) {
	t.Claims[CookieName] = s.Data.Code.String
}

// RootCategoryId returns the root category ID assigned to this store view.
func (s *Store) RootCategoryId() int64 {
	return s.Group.Data.RootCategoryID
}

/*
	Store Currency
*/

// AllowedCurrencies returns all installed currencies from global scope.
func (s *Store) AllowedCurrencies() []string {
	return strings.Split(s.cr.GetString(config.Path(directory.PathSystemCurrencyInstalled)), ",")
}

// CurrentCurrency @todo
// @see app/code/Magento/Store/Model/Store.php::getCurrentCurrency
func (s *Store) CurrentCurrency() *directory.Currency {
	return nil
}

/*
	Global functions
*/
// GetClaim returns a valid store code from a JSON web token or nil
func GetCodeFromClaim(t *jwt.Token) config.ScopeIDer {
	if t == nil {
		return nil
	}
	c, ok := t.Claims[CookieName]
	if cs, okcs := c.(string); okcs && ok && nil == ValidateStoreCode(cs) {
		return config.ScopeCode(cs)
	}
	return nil
}

// GetCookie returns from a Request the value of the store cookie or nil.
func GetCodeFromCookie(req *http.Request) config.ScopeIDer {
	if req == nil {
		return nil
	}
	if keks, err := req.Cookie(CookieName); nil == err && nil == ValidateStoreCode(keks.Value) {
		return config.ScopeCode(keks.Value)
	}
	return nil
}

// ValidateStoreCode checks if a store code is valid. Returns an ErrStoreCodeInvalid if the
// first letter is not a-zA-Z and followed by a-zA-Z0-9_ or store code length is greater than 32 characters.
func ValidateStoreCode(c string) error {
	if c == "" || len(c) > 32 {
		return ErrStoreCodeInvalid
	}
	c1 := c[0]
	if false == ((c1 >= 'a' && c1 <= 'z') || (c1 >= 'A' && c1 <= 'Z')) {
		return ErrStoreCodeInvalid
	}
	if false == utils.StrIsAlNum(c) {
		return ErrStoreCodeInvalid
	}
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
func (s StoreSlice) Filter(f func(*Store) bool) StoreSlice {
	var stores StoreSlice
	for _, v := range s {
		if v != nil && f(v) {
			stores = append(stores, v)
		}
	}
	return stores
}

// Codes returns a StringSlice with all store codes
func (s StoreSlice) Codes() utils.StringSlice {
	if len(s) == 0 {
		return nil
	}
	var c utils.StringSlice
	for _, st := range s {
		if st != nil {
			c.Append(st.Data.Code.String)
		}
	}
	return c
}

// IDs returns an Int64Slice with all store ids
func (s StoreSlice) IDs() utils.Int64Slice {
	if len(s) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, st := range s {
		if st != nil {
			ids.Append(st.Data.StoreID)
		}
	}
	return ids
}

// LastItem returns the last item of this slice or nil
func (s StoreSlice) LastItem() *Store {
	if s.Len() > 0 {
		return s[s.Len()-1]
	}
	return nil
}

/*
	TableStore and TableStoreSlice method receivers
*/

// IsDefault returns true if the current store is the default store.
func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreID
}

// Load uses a dbr session to load all data from the core_store table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see https://github.com/magento/magento2/blob/0.74.0-beta7/app%2Fcode%2FMagento%2FStore%2FModel%2FResource%2FStore%2FCollection.php#L147
// regarding the sort order.
func (s *TableStoreSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return s.parentLoad(dbrSess, append(append([]csdb.DbrSelectCb{nil}, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		sb.OrderBy("CASE WHEN main_table.store_id = 0 THEN 0 ELSE 1 END ASC")
		sb.OrderBy("main_table.sort_order ASC")
		return sb.OrderBy("main_table.name ASC")
	}), cbs...)...)
}

// FindByID returns a TableStore if found by id or an error
func (s TableStoreSlice) FindByID(id int64) (*TableStore, error) {
	for _, st := range s {
		if st != nil && st.StoreID == id {
			return st, nil
		}
	}
	return nil, ErrStoreNotFound
}

// FindByCode returns a TableStore if found by id or an error
func (s TableStoreSlice) FindByCode(code string) (*TableStore, error) {
	for _, st := range s {
		if st != nil && st.Code.Valid && st.Code.String == code {
			return st, nil
		}
	}
	return nil, ErrStoreNotFound
}

// FilterByGroupID returns a new slice with all TableStores belonging to a group id
func (s TableStoreSlice) FilterByGroupID(id int64) TableStoreSlice {
	return s.Filter(func(ts *TableStore) bool {
		return ts.GroupID == id
	})
}

// FilterByWebsiteID returns a new slice with all TableStores belonging to a website id
func (s TableStoreSlice) FilterByWebsiteID(id int64) TableStoreSlice {
	return s.Filter(func(ts *TableStore) bool {
		return ts.WebsiteID == id
	})
}

// Filter returns a new slice containing TableStores filtered by predicate f
func (s TableStoreSlice) Filter(f func(*TableStore) bool) TableStoreSlice {
	if len(s) == 0 {
		return nil
	}
	var tss TableStoreSlice
	for _, v := range s {
		if v != nil && f(v) {
			tss = append(tss, v)
		}
	}
	return tss
}

// Codes returns a StringSlice with all store codes
func (s TableStoreSlice) Codes() utils.StringSlice {
	if len(s) == 0 {
		return nil
	}
	var c utils.StringSlice
	for _, store := range s {
		if store != nil {
			c.Append(store.Code.String)
		}
	}
	return c
}

// IDs returns an Int64Slice with all store ids
func (s TableStoreSlice) IDs() utils.Int64Slice {
	if len(s) == 0 {
		return nil
	}
	var ids utils.Int64Slice
	for _, store := range s {
		if store != nil {
			ids.Append(store.StoreID)
		}
	}
	return ids
}
