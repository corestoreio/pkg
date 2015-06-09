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
	"io"
	"io/ioutil"
	"strconv"

	"github.com/corestoreio/csfw/utils/log"
)

// PS path separator used in the database table core_config_data and in config.Manager
const PS = "/"

// ArgFunc Argument function to be used as variadic argument in ScopeKey() and ScopeKeyValue()
type ArgFunc func(*arg)

// ScopeWebsite wrapper helper function. See Scope()
func ScopeWebsite(r ScopeIDer) ArgFunc { return Scope(ScopeWebsiteID, r) }

// ScopeStore wrapper helper function. See Scope()
func ScopeStore(r ScopeIDer) ArgFunc { return Scope(ScopeStoreID, r) }

// Scope sets the scope using the ScopeGroup and a config.ScopeIDer.
// A config.ScopeIDer can contain an ID from a website or a store. Make sure
// the correct ScopeGroup has also been set. If config.ScopeIDer is nil
// the scope will fallback to default scope.
func Scope(s ScopeGroup, r ScopeIDer) ArgFunc {
	if s != ScopeDefaultID && r == nil {
		s = ScopeDefaultID
	}
	return func(a *arg) { a.s = s; a.r = r }
}

// Path option function to specify the configuration path. If one argument has been
// provided then it must be a full valid path. If more than one argument has been provided
// then the arguments will be joined together. Panics if nil arguments will be provided.
func Path(paths ...string) ArgFunc {
	var p string
	lp := len(paths)
	if lp > 0 {
		p = paths[0]
	}

	if lp == 3 {
		p = p + PS + paths[1] + PS + paths[2]
	}
	return func(a *arg) { a.p = p }
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
		log.Error("Argument=ValueReader::ReadAll", "err", err)
	}
	return func(a *arg) {
		a.v = data
	}
}

// arg responsible for the correct scope key e.g.: stores/2/system/currency/installed => scope/scope_id/path
// which is used by the underlying configuration manager to fetch or store a value
type arg struct {
	p  string // p is the three level path e.g. a/b/c
	s  ScopeGroup
	r  ScopeIDer
	nb bool        // noBubble, if false value search: (store|website) -> default
	v  interface{} // value use for saving
}

// this "cache" should covers ~80% of all store setups
var int64Cache = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
}
var int64CacheLen = int64(len(int64Cache))

// newArg creates an argument container which requires different options depending on the use case.
func newArg(opts ...ArgFunc) *arg {
	var a = new(arg)
	for _, opt := range opts {
		if opt != nil {
			opt(a)
		}
	}
	return a
}

func (a *arg) isDefault() bool { return a.s == ScopeDefaultID || a.s == ScopeAbsentID }

func (a *arg) isBubbling() bool { return !a.nb }

func (a *arg) scopePath() string {
	// first part of the path is called scope in Magento and in CoreStore ScopeRange
	// e.g.: stores/2/system/currency/installed => scope/scope_id/path
	// e.g.: websites/1/system/currency/installed => scope/scope_id/path
	if a.p == "" {
		return ""
	}
	return a.scopeRange() + PS + a.scopeID() + PS + a.p
}

func (a *arg) scopePathDefault() string {
	// e.g.: default/0/system/currency/installed => scope/scope_id/path
	if a.p == "" {
		return ""
	}
	return ScopeRangeDefault + PS + "0" + PS + a.p
}

func (a *arg) scopeID() string {
	if a.r != nil {
		if a.r.ScopeID() <= int64CacheLen {
			return int64Cache[a.r.ScopeID()]
		}
		return strconv.FormatInt(a.r.ScopeID(), 10)
	}
	return "0"
}

func (a *arg) scopeRange() string {
	switch a.s {
	case ScopeWebsiteID:
		return ScopeRangeWebsites
	case ScopeStoreID:
		return ScopeRangeStores
	}
	return ScopeRangeDefault
}
