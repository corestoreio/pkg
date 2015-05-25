// Copyright 2015 CoreStore Authors
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
)

// ScopeOption function to be used as variadic argument in ScopeKey() and ScopeKeyValue()
type ScopeOption func(*arg)

// ScopeWebsite wrapper helper function. See Scope()
func ScopeWebsite(r ...Retriever) ScopeOption { return Scope(IDScopeWebsite, r...) }

// ScopeStore wrapper helper function. See Scope()
func ScopeStore(r ...Retriever) ScopeOption { return Scope(IDScopeStore, r...) }

// Scope sets the scope using the ScopeID and a variadic (0 or 1 arg) store.Retriever.
// A store.Retriever can contain an ID from a website or a store. Make sure the correct ScopeID has also been set.
// Retriever can only be left off when the ScopeID is default otherwise the scope will fallback to default scope.
func Scope(s ScopeID, r ...Retriever) ScopeOption {
	var ret Retriever
	hasR := len(r) == 1 && r[0] != nil
	if hasR {
		ret = r[0]
	}

	if s != IDScopeDefault && !hasR {
		s = IDScopeDefault
	}

	return func(a *arg) { a.s = s; a.r = ret }
}

// Path option function to specify the configuration path. If one argument has been
// provided then it must be a full valid path. If more than one argument has been provided
// then the arguments will be joined together. Panics if nil arguments will be provided.
func Path(paths ...string) ScopeOption {
	var p string
	lp := len(paths)
	if lp > 0 {
		p = paths[0]
	}

	if lp == 3 {
		p = p + "/" + paths[1] + "/" + paths[2]
	}
	return func(a *arg) { a.p = p }
}

// Value sets the value for a scope key.
func Value(v interface{}) ScopeOption { return func(a *arg) { a.v = v } }

// ValueReader sets the value for a scope key using the io.Reader interface.
func ValueReader(r io.Reader) ScopeOption {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		logger.WithField("Argument", "ValueReader").Error(err)
	}
	return func(a *arg) {
		a.v = data
	}
}

// ScopeKey generates the correct scope key e.g.: stores/2/system/currency/installed => scope/scope_id/path
// which is used by the underlaying configuration manager to fetch a value
func ScopeKey(opts ...ScopeOption) string {
	if len(opts) == 0 {
		return ""
	}
	return newArg(opts...).scopePath()
}

// ScopeKeyValue generates from the options the scope key and the value
func ScopeKeyValue(opts ...ScopeOption) (string, interface{}) {
	if len(opts) == 0 {
		return "", nil
	}
	a := newArg(opts...)
	return a.scopePath(), a.v
}

type arg struct {
	p string // p is the three level path e.g. a/b/c
	s ScopeID
	r Retriever
	v interface{} // value use for saving
}

// this "cache" should covers ~80% of all store setups
var int64Cache = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
}
var int64CacheLen = int64(len(int64Cache))

func newArg(opts ...ScopeOption) *arg {
	var a = new(arg)
	for _, opt := range opts {
		if opt != nil {
			opt(a)
		}
	}
	return a
}

func (a *arg) scopePath() string {
	// e.g.: stores/2/system/currency/installed => scope/scope_id/path
	if a.p == "" {
		return ""
	}
	return a.scopeData() + "/" + a.scopeID() + "/" + a.p
}

func (a *arg) scopeID() string {
	if a.r != nil {
		if a.r.ID() <= int64CacheLen {
			return int64Cache[a.r.ID()]
		}
		return strconv.FormatInt(a.r.ID(), 10)
	}
	return "0"
}

func (a *arg) scopeData() string {
	switch a.s {
	case IDScopeWebsite:
		return StringScopeWebsites
	case IDScopeStore:
		return StringScopeStores
	}
	return StringScopeDefault
}
