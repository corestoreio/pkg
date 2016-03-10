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
	"time"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errors"
)

// ScopedGetter is equal to Getter but the underlying implementation takes
// care of providing the correct scope: default, website or store and bubbling
// up the scope chain from store -> website -> default if a value won't get
// found in the desired scope.
//
// To restrict bubbling up you can provide a second argument scope.Scope.
// You can restrict a configuration path to be only used with the default,
// website or store scope. See the examples.
//
// This interface is mainly implemented in the store package. The functions
// should be the same as in Getter but only the different is the route
// argument. A route represents always "a/b/c".
// Returned error is mostly of ErrKeyNotFound.
type ScopedGetter interface {
	scope.Scoper
	String(r cfgpath.Route, s ...scope.Scope) (string, error)
	Bool(r cfgpath.Route, s ...scope.Scope) (bool, error)
	Float64(r cfgpath.Route, s ...scope.Scope) (float64, error)
	Int(r cfgpath.Route, s ...scope.Scope) (int, error)
	Time(r cfgpath.Route, s ...scope.Scope) (time.Time, error)
}

// think about that segregation
//type ScopedStringer interface {
//	scope.Scoper
//	Bind(scope.Scope) ScopedGetter
//	String(r cfgpath.Route, s ...scope.Scope) (string, error)
//}
// and so on ...

type scopedService struct {
	root      Getter
	websiteID int64
	storeID   int64
}

var _ ScopedGetter = (*scopedService)(nil)

// NewScopedService instantiates a ScopedGetter implementation.
// For internal use only. Exported because of the config/cfgmock package.
func NewScopedService(r Getter, websiteID, storeID int64) scopedService {
	return scopedService{
		root:      r,
		websiteID: websiteID,
		storeID:   storeID,
	}
}

// Scope tells you the current underlying scope and its website or store ID
func (ss scopedService) Scope() (scope.Scope, int64) {
	switch {
	case ss.storeID > 0:
		return scope.StoreID, ss.storeID
	case ss.websiteID > 0:
		return scope.WebsiteID, ss.websiteID
	default:
		return scope.DefaultID, 0
	}
}

// String traverses through the scopes store->website->default to find
// a matching string value.
func (ss scopedService) String(r cfgpath.Route, s ...scope.Scope) (v string, err error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		err = errors.Mask(err)
		return
	}

	if ss.storeID > 0 && scope.PermStoreReverse.Has(s...) {
		v, err = ss.root.String(p.Bind(scope.StoreID, ss.storeID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	}
	if ss.websiteID > 0 && scope.PermWebsiteReverse.Has(s...) {
		v, err = ss.root.String(p.Bind(scope.WebsiteID, ss.websiteID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	}
	return ss.root.String(p)
}

// Bool traverses through the scopes store->website->default to find
// a matching bool value.
func (ss scopedService) Bool(r cfgpath.Route, s ...scope.Scope) (v bool, err error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		err = errors.Mask(err)
		return
	}

	if ss.storeID > 0 && scope.PermStoreReverse.Has(s...) {
		v, err = ss.root.Bool(p.Bind(scope.StoreID, ss.storeID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in store scope go to website scope

	if ss.websiteID > 0 && scope.PermWebsiteReverse.Has(s...) {
		v, err = ss.root.Bool(p.Bind(scope.WebsiteID, ss.websiteID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in website scope go to default scope
	return ss.root.Bool(p)
}

// Float64 traverses through the scopes store->website->default to find
// a matching float64 value.
func (ss scopedService) Float64(r cfgpath.Route, s ...scope.Scope) (v float64, err error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		err = errors.Mask(err)
		return
	}

	if ss.storeID > 0 && scope.PermStoreReverse.Has(s...) {
		v, err = ss.root.Float64(p.Bind(scope.StoreID, ss.storeID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in store scope go to website scope

	if ss.websiteID > 0 && scope.PermWebsiteReverse.Has(s...) {
		v, err = ss.root.Float64(p.Bind(scope.WebsiteID, ss.websiteID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in website scope go to default scope
	return ss.root.Float64(p)
}

// Int traverses through the scopes store->website->default to find
// a matching int value.
func (ss scopedService) Int(r cfgpath.Route, s ...scope.Scope) (v int, err error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		err = errors.Mask(err)
		return
	}

	if ss.storeID > 0 && scope.PermStoreReverse.Has(s...) {
		v, err = ss.root.Int(p.Bind(scope.StoreID, ss.storeID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in store scope go to website scope

	if ss.websiteID > 0 && scope.PermWebsiteReverse.Has(s...) {
		v, err = ss.root.Int(p.Bind(scope.WebsiteID, ss.websiteID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in website scope go to default scope
	return ss.root.Int(p)
}

// Time traverses through the scopes store->website->default to find
// a matching time.Time value.
func (ss scopedService) Time(r cfgpath.Route, s ...scope.Scope) (v time.Time, err error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		err = errors.Mask(err)
		return
	}

	if ss.storeID > 0 && scope.PermStoreReverse.Has(s...) {
		v, err = ss.root.Time(p.Bind(scope.StoreID, ss.storeID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in store scope go to website scope

	if ss.websiteID > 0 && scope.PermWebsiteReverse.Has(s...) {
		v, err = ss.root.Time(p.Bind(scope.WebsiteID, ss.websiteID))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
	} // if not found in website scope go to default scope
	return ss.root.Time(p)
}
