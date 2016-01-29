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

package config

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
)

// ArgFunc Argument function to be used as variadic argument in ScopeKey() and ScopeKeyValue()
type ArgFunc func(*arg)

// scopeDefaultArg cached default scope argument function
var scopeDefaultArg = Scope(scope.DefaultID, 0)

// ScopeDefault wrapper helper function. See Scope(). Mainly used to show humans
// than a config value can only be set for a global scope.
func ScopeDefault() ArgFunc { return scopeDefaultArg }

// ScopeWebsite wrapper helper function. See Scope()
func ScopeWebsite(id int64) ArgFunc { return Scope(scope.WebsiteID, id) }

// ScopeGroup wrapper helper function. See Scope()
func ScopeGroup(id int64) ArgFunc { return Scope(scope.GroupID, id) }

// ScopeStore wrapper helper function. See Scope()
func ScopeStore(id int64) ArgFunc { return Scope(scope.StoreID, id) }

// Scope sets the scope using the scope.Group and a ID.
// The ID can contain an integer from a website or a store. Make sure
// the correct scope.Scope has also been set. If the ID is smaller
// than zero the scope will fallback to default scope.
func Scope(s scope.Scope, id int64) ArgFunc {
	if s != scope.DefaultID && id < 1 {
		id = 0
		s = scope.DefaultID
	}
	return func(a *arg) {
		a.Scope = s
		a.ID = id
		a.scopeSet = true
	}
}

// Path option function to set the configuration Path. If the Scope*()
// option functions have not been applied the scope and scope ID from the
// Path will be taken. You can then overwrite those two settings with the
// Scope*() functions. The other way round does not work. Once the Scope*()
// functions have been applied the Scope and ID from Path won't be applied.
func Path(p path.Path) ArgFunc {
	return func(a *arg) {
		if a.scopeSet { // only copy Route, do not overwrite Scope and ID
			a.Path.Route = p.Route
		} else {
			a.Path = p
		}
		if err := a.IsValid(); err != nil {
			a.lastErrors = append(a.lastErrors, err)
		}
	}
}

// PathScoped creates a new Path from a path string ps, a scope and the
// ID. This option function overrides everything.
func PathScoped(ps string, s scope.Scope, id int64) ArgFunc {
	return func(a *arg) {
		p, err := path.NewByParts(ps)
		if err != nil {
			a.lastErrors = append(a.lastErrors, err)
		}
		a.Path = p.Bind(s, id)
		if err := a.IsValid(); err != nil {
			a.lastErrors = append(a.lastErrors, err)
		}
	}
}

// Route option function to specify the configuration path without any scope
// or scope ID applied. You must call the Scope*() functions also. If not
// the default scope will be applied.
func Route(r path.Route) ArgFunc {
	return func(a *arg) {
		a.Route = r
		if err := a.IsValid(); err != nil {
			a.lastErrors = append(a.lastErrors, err)
		}
	}
}

// Value sets the value for a scope key.
func Value(v interface{}) ArgFunc { return func(a *arg) { a.v = v } }

// ValueReader sets the value for a scope key using the io.Reader interface.
// If asserting to a io.Closer is successful then Close() will be called.
func ValueReader(r io.Reader) ArgFunc {
	data, err := ioutil.ReadAll(r)
	if c, ok := r.(io.Closer); ok && c != nil {
		if err := c.Close(); err != nil {
			return func(a *arg) {
				a.lastErrors = append(a.lastErrors, fmt.Errorf("ValueReader.Close error %s", err))
			}
		}
	}
	if err != nil {
		return func(a *arg) {
			a.lastErrors = append(a.lastErrors, fmt.Errorf("ValueReader error %s", err))
		}
	}
	return func(a *arg) {
		a.v = data
	}
}

// arg responsible for the correct scope key e.g.: stores/2/system/currency/installed => scope/scope_id/path
// which is used by the underlying configuration Service to fetch or store a value
type arg struct {
	path.Path
	scopeSet   bool        // true if the Scope*() functions have been used
	v          interface{} // value use for saving
	lastErrors []error
}

// newArg creates an argument container which requires different options depending on the use case.
func newArg(opts ...ArgFunc) (arg, error) {
	a := arg{}
	return a.option(opts...)
}

// mustNewArg panics on error. useful for initialization process
func mustNewArg(opts ...ArgFunc) arg {
	a, err := newArg(opts...)
	if err != nil {
		panic(err)
	}
	return a
}

func (a arg) option(opts ...ArgFunc) (arg, error) {
	for _, opt := range opts {
		if opt != nil {
			opt(&a)
		}
	}
	if len(a.lastErrors) > 0 {
		a.v = nil // discard value on error
		return a, a
	}
	return a, nil
}

func (a arg) isDefault() bool { return a.Scope == scope.DefaultID || a.Scope == scope.AbsentID }

var _ error = (*arg)(nil)

func (a arg) Error() string {
	return util.Errors(a.lastErrors...)
}
