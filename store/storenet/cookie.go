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

package storenet

import (
	"net/http"
	"time"

	"github.com/corestoreio/csfw/store"
)

// Cookie allows to set and delete the store cookie
type Cookie struct {
	Store *store.Store
}

// NewCookie creates a new pre-configured cookie.
// TODO(cs) create cookie manager to stick to the limits of http://www.ietf.org/rfc/rfc2109.txt page 15
// @see http://browsercookielimits.squawky.net/
func (c Cookie) New(path string) *http.Cookie {
	return &http.Cookie{
		Name:     ParamName,
		Value:    "",
		Path:     path,
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
	}
}

// Set adds a cookie which contains the store code and is valid for one year.
func (c Cookie) Set(res http.ResponseWriter) {
	if res != nil {
		keks := c.New()
		keks.Value = c.Store.Data.Code.String
		keks.Expires = time.Now().AddDate(1, 0, 0) // one year valid
		http.SetCookie(res, keks)
	}
}

// DeleteCookie deletes the store cookie
func (c Cookie) Delete(res http.ResponseWriter) {
	if res != nil {
		keks := c.New()
		keks.Expires = time.Now().AddDate(-10, 0, 0)
		http.SetCookie(res, keks)
	}
}
