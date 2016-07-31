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

package runmode

import (
	"net"
	"net/http"
	"time"

	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
)

// StoreCodeProcesser knows how to extract a store code from
// a request.
type StoreCodeProcesser interface {
	// FromRequest returns the valid non-empty store code.
	FromRequest(req *http.Request) (code string)

	ProcessDenied(runMode scope.Hash, newStoreID int64) http.Handler
	ProcessAllowed(runMode scope.Hash, newStoreID int64, w http.ResponseWriter, r *http.Request)
}

// FieldName use in Cookies and JSON Web Tokens (JWT) to identify an active
// store besides from the default loaded store. This is the default value.
const FieldName = `store`

// URLFieldName name of the GET parameter to set a new store in a current
// website/group context.  This is the default value.
const URLFieldName = `___store`

// ExtractStoreCode can extract the store code from an HTTP Request. This code
// is then responsible for changing the runMode.
type ExtractStoreCode struct {
	// FieldName optional custom name, defaults to constant FieldName
	FieldName string
	// URLFieldName optional custom name, defaults to constant URLFieldName
	URLFieldName string

	// CookieTemplate optional pre-configured cookie to set the store
	// code. Expiration time and value will get overwritten.
	CookieTemplate func(*http.Request) *http.Cookie
	// CookieExpiresSet defaults to one year expiration for the store code.
	CookieExpiresSet time.Time
	// CookieExpiresDelete defaults to minus ten years to delete the store code
	// cookie.
	CookieExpiresDelete time.Time
}

// FromRequest returns from a GET request with a query string the value of the
// store code. If no code can be found in the query string, this function falls
// back to the cookie name defined in field FieldName. Valid has three values: 0
// not valid, 10 valid and code found in GET query string, 20 valid and code
// found in cookie. Implements interface StoreCodeExtracter.
func (e ExtractStoreCode) FromRequest(req *http.Request) (code string) {
	// todo find a better solution for the valid type
	hps := URLFieldName
	if e.URLFieldName != "" {
		hps = e.URLFieldName
	}
	code = req.URL.Query().Get(hps)
	if code == "" {
		return e.fromCookie(req)
	}
	if err := store.CodeIsValid(code); err == nil {
		return code
	}
	return ""
}

// fromCookie extracts a store from a cookie using the field name FieldName as
// an identifier.
func (e ExtractStoreCode) fromCookie(req *http.Request) (code string) {
	p := FieldName
	if e.FieldName != "" {
		p = e.FieldName
	}
	if keks, err := req.Cookie(p); err == nil {
		code = keks.Value
	}
	if err := store.CodeIsValid(code); err == nil {
		return code
	}
	return ""
}

func (a ExtractStoreCode) newCookie(r *http.Request) *http.Cookie {
	// sync.Pool for cookies, refactor Cookies.String() in stdlib
	if a.CookieTemplate != nil {
		return a.CookieTemplate(r)
	}
	d, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		d = r.Host // might be a bug ...
	}
	var isSecure bool
	if r.TLS != nil {
		isSecure = true
	}
	return &http.Cookie{
		Name:     FieldName,
		Path:     "/", // we can sit behind a proxy, so path must be configurable
		Domain:   d,
		Secure:   isSecure,
		HttpOnly: true, // disable for JavaScript access
	}
}

func (a ExtractStoreCode) getCookieExpiresSet() time.Time {
	if a.CookieExpiresSet.IsZero() {
		return time.Now().AddDate(1, 0, 0) // one year valid
	}
	return a.CookieExpiresSet
}

func (a ExtractStoreCode) getCookieExpiresDelete() time.Time {
	if a.CookieExpiresDelete.IsZero() {
		return time.Now().AddDate(-10, 0, 0) // -10 years
	}
	return a.CookieExpiresDelete
}

func (a ExtractStoreCode) ProcessDenied(runMode scope.Hash, newStoreID int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if store code found in cookie and not valid anymore, delete the cookie.
		if c := a.fromCookie(r); c != "" {
			keks := a.newCookie(r)
			keks.Expires = a.getCookieExpiresDelete()
			http.SetCookie(w, keks)
		}
	})
}

func (a ExtractStoreCode) ProcessAllowed(runMode scope.Hash, newStoreID int64, w http.ResponseWriter, r *http.Request) {
	if c := a.fromCookie(r); c == "" {

	}
	if reqStoreCodeValid < 20 { // no cookie found but the code changed
		// set cookie once with the new code
		keks := a.getCookie(r)
		keks.Expires = a.getCookieExpiresSet()
		http.SetCookie(w, keks)
	}
}

type nullCodeProcessor struct{}

func (nc nullCodeProcessor) FromRequest(_ *http.Request) string               { return "" }
func (nc nullCodeProcessor) ProcessDenied(_ scope.Hash, _ int64) http.Handler { return nil }
func (nc nullCodeProcessor) ProcessAllowed(runMode scope.Hash, newStoreID int64, w http.ResponseWriter, r *http.Request) {
}

var _ StoreCodeProcesser = (*nullCodeProcessor)(nil)
