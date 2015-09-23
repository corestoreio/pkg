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

package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"

	"errors"
	"fmt"
	"strings"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/utils/log"
)

const hierarchyLevel int = 3 // a/b/c

// PS path separator used in the database table core_config_data and in config.Manager
const PS = "/"

// ErrPathEmpty when you provide an empty path in the function Path()
var ErrPathEmpty = errors.New("Path cannot be empty")

// ArgFunc Argument function to be used as variadic argument in ScopeKey() and ScopeKeyValue()
type ArgFunc func(*arg)

// ScopeWebsite wrapper helper function. See Scope()
func ScopeWebsite(id int64) ArgFunc { return Scope(scope.WebsiteID, id) }

// ScopeGroup wrapper helper function. See Scope()
func ScopeGroup(id int64) ArgFunc { return Scope(scope.GroupID, id) }

// ScopeStore wrapper helper function. See Scope()
func ScopeStore(id int64) ArgFunc { return Scope(scope.StoreID, id) }

// Scope sets the scope using the scope.Group and a scope.IDer.
// A scope.IDer can contain an ID from a website or a store. Make sure
// the correct scope.Group has also been set. If scope.IDer is nil
// the scope will fallback to default scope.
func Scope(s scope.Scope, id int64) ArgFunc {
	if s != scope.DefaultID && id < 1 {
		id = 0
		s = scope.DefaultID
	}
	return func(a *arg) { a.sg = s; a.si = id }
}

// Path option function to specify the configuration path. If one argument has been
// provided then it must be a full valid path. If more than one argument has been provided
// then the arguments will be joined together. Panics if nil arguments will be provided.
func Path(paths ...string) ArgFunc {
	// TODO(cs) validation of the path see typeConfigPath in app/code/Magento/Config/etc/system_file.xsd
	var p string
	lp := len(paths)
	if lp > 0 {
		p = paths[0]
	}
	if p == "" {
		return func(a *arg) {
			a.lastErrors = append(a.lastErrors, ErrPathEmpty)
		}
	}

	var paSlice []string
	if lp == hierarchyLevel {
		p = p + PS + paths[1] + PS + paths[2]
		paSlice = paths
	} else {
		paSlice = strings.Split(p, PS)
		if len(paSlice) != hierarchyLevel {
			return func(a *arg) {
				a.lastErrors = append(a.lastErrors, fmt.Errorf("Incorrect number of paths elements: want %d, have %d, Path: %v", hierarchyLevel, len(paSlice), paths))
			}
		}
	}
	return func(a *arg) {
		a.pa = p
		a.paSlice = paSlice
	}
}

// NoBubble disables the fallback to the default scope when a value in the current
// scope not exists.
func NoBubble() ArgFunc { return func(a *arg) { a.nb = true } }

// Value sets the value for a scope key.
func Value(v interface{}) ArgFunc { return func(a *arg) { a.v = v } }

// ValueReader sets the value for a scope key using the io.Reader interface.
// If asserting to a io.Closer is successful then Close() will be called.
func ValueReader(r io.Reader) ArgFunc {
	if c, ok := r.(io.Closer); ok {
		defer c.Close()
	}
	data, err := ioutil.ReadAll(r)
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
// which is used by the underlying configuration manager to fetch or store a value
type arg struct {
	pa         string   // p is the three level path e.g. a/b/c
	paSlice    []string // used for hierarchy for the pubSub system
	sg         scope.Scope
	si         int64       // scope ID
	nb         bool        // noBubble, if false value search: (store|website) -> default
	v          interface{} // value use for saving
	lastErrors []error
}

// this "cache" should covers ~80% of all store setups
var int64Cache = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
}
var int64CacheLen = int64(len(int64Cache))

// newArg creates an argument container which requires different options depending on the use case.
func newArg(opts ...ArgFunc) (arg, error) {
	var a = arg{}
	for _, opt := range opts {
		if opt != nil {
			opt(&a)
		}
	}
	if len(a.lastErrors) > 0 {
		return arg{}, a
	}
	return a, nil
}

// mustNewArg panics on error. useful for initialization process
func mustNewArg(opts ...ArgFunc) arg {
	a, err := newArg(opts...)
	if err != nil {
		log.Error("config.mustNewArg", "err", err)
		panic(err)
	}
	return a
}

func (a arg) isDefault() bool { return a.sg == scope.DefaultID || a.sg == scope.AbsentID }

func (a arg) isBubbling() bool { return !a.nb }

func (a arg) pathLevel1() string {
	return a.paSlice[0]
}

func (a arg) pathLevel2() string {
	return a.paSlice[0] + PS + a.paSlice[1]
}

func (a arg) pathLevel3() string {
	return a.paSlice[0] + PS + a.paSlice[1] + PS + a.paSlice[2]
}

func (a arg) scopePath() string {
	// first part of the path is called scope in Magento and in CoreStore ScopeRange
	// e.g.: stores/2/system/currency/installed => scope/scope_id/path
	// e.g.: websites/1/system/currency/installed => scope/scope_id/path
	if a.pa == "" {
		return ""
	}
	return a.scopeRange() + PS + a.scopeID() + PS + a.pa
}

func (a arg) scopePathDefault() string {
	// e.g.: default/0/system/currency/installed => scope/scope_id/path
	return scope.RangeDefault + PS + "0" + PS + a.pa
}

func (a arg) scopeID() string {
	if a.si > 0 {
		if a.si <= int64CacheLen {
			return int64Cache[a.si]
		}
		return strconv.FormatInt(a.si, 10)
	}
	return "0"
}

func (a arg) scopeRange() string {
	switch a.sg {
	case scope.WebsiteID:
		return scope.RangeWebsites
	case scope.StoreID:
		return scope.RangeStores
	}
	return scope.RangeDefault
}

var _ error = (*arg)(nil)

func (a arg) Error() string {
	var buf bytes.Buffer
	lle := len(a.lastErrors) - 1
	for i, e := range a.lastErrors {
		buf.WriteString(e.Error())
		if i < lle {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
