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
	"strings"
	"time"

	"github.com/corestoreio/cspkg/store"
	"github.com/corestoreio/cspkg/store/scope"
)

// ProcessStoreCodeCookie can extract the store code from a cookie within an
// HTTP Request. Handles cookies to permanently set the store code under
// different conditions. This store code is then responsible for changing the
// runMode.
type ProcessStoreCodeCookie struct {
	// FieldName optional custom name, defaults to constant store.CodeFieldName.
	// Cannot be changed after the first call to FromRequest().
	FieldName string
	// URLFieldName optional custom name, defaults to constant
	// store.CodeURLFieldName. Cannot be changed after the first call to
	// FromRequest().
	URLFieldName string

	// CookieTemplate optional pre-configured cookie to set the store
	// code. Expiration time and value will get overwritten.
	CookieTemplate func(*http.Request) *http.Cookie
	// CookieExpiresSet defaults to one year expiration for the store code.
	CookieExpiresSet time.Time
	// CookieExpiresDelete defaults to minus ten years to delete the store code
	// cookie.
	CookieExpiresDelete time.Time

	// pre-calculated keys for strings.Contain
	keyFieldName    string
	keyURLFieldName string
}

func (e *ProcessStoreCodeCookie) keyURLFN() (string, string) {
	if e.URLFieldName == "" {
		e.URLFieldName = store.CodeURLFieldName
	}
	if e.keyURLFieldName == "" {
		e.keyURLFieldName = e.URLFieldName + "="
	}
	return e.URLFieldName, e.keyURLFieldName
}

// FromRequest returns from a GET request with a query string the value of the
// store code. If no code can be found in the query string, this function falls
// back to the cookie name defined in field FieldName. Valid has three values: 0
// not valid, 10 valid and code found in GET query string, 20 valid and code
// found in cookie. Implements interface store.CodeProcessor.
func (e *ProcessStoreCodeCookie) FromRequest(_ scope.TypeID, req *http.Request) string {
	fn, fnK := e.keyURLFN()
	if strings.Contains(req.URL.RawQuery, fnK) {
		code := req.URL.Query().Get(fn)
		if err := store.CodeIsValid(code); err == nil {
			return code
		}
	}
	return e.fromCookie(req)
}

func (e *ProcessStoreCodeCookie) keyFN() (string, string) {
	if e.FieldName == "" {
		e.FieldName = store.CodeFieldName
	}
	if e.keyFieldName == "" {
		e.keyFieldName = e.FieldName + "="
	}
	return e.FieldName, e.keyFieldName
}

// fromCookie extracts a store from a cookie using the field name FieldName as
// an identifier.
func (e *ProcessStoreCodeCookie) fromCookie(req *http.Request) string {
	fn, fnK := e.keyFN()
	if c := req.Header.Get("Cookie"); c != "" && strings.Contains(c, fnK) {
		// move cookie parsing after the check for the code in the cookie string
		if keks, err := req.Cookie(fn); err == nil {
			if err := store.CodeIsValid(keks.Value); err == nil {
				return keks.Value
			}
		}
	}
	return ""
}

func (a *ProcessStoreCodeCookie) newCookie(r *http.Request) *http.Cookie {
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
		Name:     store.CodeFieldName,
		Path:     "/", // we can sit behind a proxy, so path must be configurable
		Domain:   d,
		Secure:   isSecure,
		HttpOnly: true, // disable for JavaScript access
	}
}

func (a *ProcessStoreCodeCookie) setStoreCookie(storeCode string, w http.ResponseWriter, r *http.Request) {
	t := a.CookieExpiresSet
	if a.CookieExpiresSet.IsZero() {
		t = time.Now().AddDate(1, 0, 0) // one year valid
	}

	keks := a.newCookie(r)
	keks.Expires = t
	keks.Value = storeCode
	http.SetCookie(w, keks)
}

func (a *ProcessStoreCodeCookie) deleteStoreCookie(w http.ResponseWriter, r *http.Request) {
	t := a.CookieExpiresDelete
	if a.CookieExpiresDelete.IsZero() {
		t = time.Now().AddDate(-10, 0, 0) // -10 years
	}

	keks := a.newCookie(r)
	keks.Expires = t
	keks.Value = ""
	http.SetCookie(w, keks)
}

// ProcessDenied deletes the store code cookie, if a store cookie can be found.
// Implements interface store.CodeProcessor.
func (a *ProcessStoreCodeCookie) ProcessDenied(_ scope.TypeID, _, _ int64, w http.ResponseWriter, r *http.Request) {
	// if store code found in cookie and not valid anymore, delete the cookie.
	if c := a.fromCookie(r); c != "" {
		a.deleteStoreCookie(w, r)
	}
}

// ProcessAllowed deletes the store code cookie if found and stores are equal or
// sets a store code cookie if the stores differ. Implements interface
// store.CodeProcessor.
func (a *ProcessStoreCodeCookie) ProcessAllowed(_ scope.TypeID, oldStoreID, newStoreID int64, newStoreCode string, w http.ResponseWriter, r *http.Request) {
	c := a.fromCookie(r)

	if c != "" && oldStoreID == newStoreID {
		// cookie not needed anymore, so delete it.
		a.deleteStoreCookie(w, r)
		return
	}

	// no cookie found but the code changed, so set cookie once with the new code
	if c == "" && oldStoreID != newStoreID {
		a.setStoreCookie(newStoreCode, w, r)
	}
}

type nullCodeProcessor struct{}

func (nc nullCodeProcessor) FromRequest(_ scope.TypeID, _ *http.Request) string { return "" }
func (nc nullCodeProcessor) ProcessDenied(_ scope.TypeID, _, _ int64, _ http.ResponseWriter, _ *http.Request) {
}
func (nc nullCodeProcessor) ProcessAllowed(_ scope.TypeID, _, _ int64, _ string, _ http.ResponseWriter, _ *http.Request) {
}

var _ store.CodeProcessor = (*nullCodeProcessor)(nil)
