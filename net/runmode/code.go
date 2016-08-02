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

	"strings"

	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
)

// StoreCodeProcessor gets used in the middleware WithRunMode() to extract a
// store code from a Request and modify the response; for example setting
// cookies to persists the selected store.
type StoreCodeProcessor interface {
	// FromRequest returns the valid non-empty store code. Returns an empty
	// store code on all other cases.
	FromRequest(req *http.Request) (code string)
	// ProcessDenied gets called in the middleware WithRunMode whenever a store
	// ID isn't allowed to proceed. The variable newStoreID reflects the denied
	// store ID. The ResponseWriter and Request variables can be used for
	// additional information writing and extracting. The error Handler  will
	// always be called.
	ProcessDenied(runMode scope.Hash, newID int64, w http.ResponseWriter, r *http.Request)
	// ProcessAllowed enables to adjust the ResponseWriter based on the new
	// store ID. The variable newStoreID contains the new ID, which can also be
	// 0. The code is guaranteed to be not empty, a valid store code, and always
	// points to an existing active store. The ResponseWriter and Request
	// variables can be used for additional information writing and extracting.
	// The next Handler in the chain will after this function be called.
	ProcessAllowed(runMode scope.Hash, newID int64, storeCode string, w http.ResponseWriter, r *http.Request)
}

// FieldName use in Cookies and JSON Web Tokens (JWT) to identify an active
// store besides from the default loaded store. This is the default value.
const FieldName = `store`

// URLFieldName name of the GET parameter to set a new store in a current
// website/group context.  This is the default value.
const URLFieldName = `___store`

// ProcessStoreCode can extract the store code from an HTTP Request. Handles
// cookies to permanently set the store code under different conditions. This
// code is then responsible for changing the runMode.
type ProcessStoreCode struct {
	// FieldName optional custom name, defaults to constant FieldName. Cannot be
	// changed after the first call to FromRequest().
	FieldName string
	// URLFieldName optional custom name, defaults to constant URLFieldName. Cannot
	// be changed after the first call to FromRequest().
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

func (e *ProcessStoreCode) keyURLFN() (string, string) {
	if e.URLFieldName == "" {
		e.URLFieldName = URLFieldName
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
// found in cookie. Implements interface StoreCodeExtracter.
func (e *ProcessStoreCode) FromRequest(req *http.Request) string {
	fn, fnK := e.keyURLFN()
	if strings.Contains(req.URL.RawQuery, fnK) {
		code := req.URL.Query().Get(fn)
		if err := store.CodeIsValid(code); err == nil {
			return code
		}
	}
	return e.fromCookie(req)
}

func (e *ProcessStoreCode) keyFN() (string, string) {
	if e.FieldName == "" {
		e.FieldName = FieldName
	}
	if e.keyFieldName == "" {
		e.keyFieldName = e.FieldName + "="
	}
	return e.FieldName, e.keyFieldName
}

// fromCookie extracts a store from a cookie using the field name FieldName as
// an identifier.
func (e *ProcessStoreCode) fromCookie(req *http.Request) string {
	fn, fnK := e.keyFN()
	if c := req.Header.Get("Cookie"); c != "" && strings.Contains(c, fnK) {
		if keks, err := req.Cookie(fn); err == nil {
			if err := store.CodeIsValid(keks.Value); err == nil {
				return keks.Value
			}
		}
	}
	return ""
}

func (a *ProcessStoreCode) newCookie(r *http.Request) *http.Cookie {
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

func (a *ProcessStoreCode) getCookieExpiresSet() time.Time {
	if a.CookieExpiresSet.IsZero() {
		return time.Now().AddDate(1, 0, 0) // one year valid
	}
	return a.CookieExpiresSet
}

func (a *ProcessStoreCode) getCookieExpiresDelete() time.Time {
	if a.CookieExpiresDelete.IsZero() {
		return time.Now().AddDate(-10, 0, 0) // -10 years
	}
	return a.CookieExpiresDelete
}

func (a *ProcessStoreCode) writeDeleteCookie(w http.ResponseWriter, r *http.Request) {
	if c := a.fromCookie(r); c != "" {
		keks := a.newCookie(r)
		keks.Expires = a.getCookieExpiresDelete()
		keks.Value = ""
		http.SetCookie(w, keks)
	}
}

// ProcessDenied deletes the store code cookie, if a store cookie can be found.
func (a *ProcessStoreCode) ProcessDenied(runMode scope.Hash, newStoreID int64, w http.ResponseWriter, r *http.Request) {
	// if store code found in cookie and not valid anymore, delete the cookie.
	a.writeDeleteCookie(w, r)
}

// ProcessAllowed deletes the store code cookie if found and stores are equal or
// sets a store code cookie if the stores differ.
func (a *ProcessStoreCode) ProcessAllowed(runMode scope.Hash, newStoreID int64, storeCode string, w http.ResponseWriter, r *http.Request) {

	if runMode.ID() == newStoreID {
		a.writeDeleteCookie(w, r)
		return
	}

	// no cookie found but the code changed

	// set cookie once with the new code
	keks := a.newCookie(r)
	keks.Expires = a.getCookieExpiresSet()
	keks.Value = storeCode
	http.SetCookie(w, keks)
}

type nullCodeProcessor struct{}

func (nc nullCodeProcessor) FromRequest(_ *http.Request) string { return "" }
func (nc nullCodeProcessor) ProcessDenied(_ scope.Hash, _ int64, _ http.ResponseWriter, _ *http.Request) {
}
func (nc nullCodeProcessor) ProcessAllowed(_ scope.Hash, _ int64, _ string, _ http.ResponseWriter, _ *http.Request) {
}

var _ StoreCodeProcessor = (*nullCodeProcessor)(nil)
