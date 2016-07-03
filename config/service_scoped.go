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
	"github.com/corestoreio/csfw/util/errors"
)

// ScopedGetter is equal to Getter but the underlying implementation takes care
// of providing the correct scope: default, website or store and bubbling up the
// scope chain from store -> website -> default if a value won't get found in
// the desired scope. The cfgpath.Route for each primitive type represents
// always a path like "section/group/element" without the scope string and scope
// ID.
//
// To restrict bubbling up you can provide a second argument scope.Scope. You
// can restrict a configuration path to be only used with the default, website
// or store scope. See the examples. This second argument will mainly be used by
// the cfgmodel package to use a defined scope in a config.Structure. If you
// access the ScopedGetter from a store.Store, store.Website type the second
// argument must already be internally pre-filled.
//
// Returned error has mostly the behaviour of not found.
type ScopedGetter interface {
	// Parent tells you the parent underlying scope and its ID. Store falls back to
	// website and website falls back to default.
	Parent() (scope.Scope, int64)
	// Scope tells you the current underlying scope and its ID.
	scope.Scoper
	// Byte traverses through the scopes store->website->default to find
	// a matching byte slice value.
	Byte(r cfgpath.Route, s ...scope.Scope) ([]byte, scope.Hash, error)
	// String see Byte()
	String(r cfgpath.Route, s ...scope.Scope) (string, scope.Hash, error)
	// Bool see Byte()
	Bool(r cfgpath.Route, s ...scope.Scope) (bool, scope.Hash, error)
	// Float64 see Byte()
	Float64(r cfgpath.Route, s ...scope.Scope) (float64, scope.Hash, error)
	// Int see Byte()
	Int(r cfgpath.Route, s ...scope.Scope) (int, scope.Hash, error)
	// Time see Byte()
	Time(r cfgpath.Route, s ...scope.Scope) (time.Time, scope.Hash, error)
}

// think about that segregation
//type ScopedStringer interface {
//  Parent() (scope.Scope, int64)
//	scope.Scoper
//	Bind(scope.Scope) ScopedGetter
//	String(r cfgpath.Route, s ...scope.Scope) (string, error)
//}
// and so on ...

type scopedService struct {
	root Getter
	// scp defines the scope bound to
	scp       scope.Scope
	websiteID int64
	storeID   int64
}

var _ ScopedGetter = (*scopedService)(nil)

// NewScopedService instantiates a ScopedGetter implementation. For internal use
// only. Exported because of the config/cfgmock package. Getter specifies the
// root Getter which does not know about any scope. WebsiteID and StoreID must
// be in a relation like enforced in the database tables via foreign keys. Empty
// storeID triggers the website scope. Empty websiteID and empty storeID are
// triggering the default scope.
func NewScopedService(r Getter, websiteID, storeID int64) ScopedGetter {
	ss := scopedService{
		root:      r,
		websiteID: websiteID,
		storeID:   storeID,
	}
	ss.scp, _ = ss.Scope()
	return ss
}

// Parent tells you the parent underlying scope and its ID. Store falls back to
// website and website falls back to default.
func (ss scopedService) Parent() (scope.Scope, int64) {
	if ss.storeID > 0 {
		return scope.Website, ss.websiteID
	}
	return scope.Default, 0
}

// Scope tells you the current underlying scope and its ID.
func (ss scopedService) Scope() (scope.Scope, int64) {
	switch {
	case ss.storeID > 0:
		return scope.Store, ss.storeID
	case ss.websiteID > 0:
		return scope.Website, ss.websiteID
	}
	return scope.Default, 0
}

func (ss scopedService) isAllowedStore(s ...scope.Scope) bool {
	scp := ss.scp
	if len(s) > 0 && s[0] > scope.Absent {
		scp = s[0]
	}
	return ss.storeID > 0 && scope.PermStoreReverse.Has(scp)
}

func (ss scopedService) isAllowedWebsite(s ...scope.Scope) bool {
	scp := ss.scp
	if len(s) > 0 && s[0] > scope.Absent {
		scp = s[0]
	}
	return ss.websiteID > 0 && scope.PermWebsiteReverse.Has(scp)
}

// Byte traverses through the scopes store->website->default to find
// a matching byte slice value.
func (ss scopedService) Byte(r cfgpath.Route, s ...scope.Scope) ([]byte, scope.Hash, error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "[config] Byte. Route %q", r)
	}

	if ss.isAllowedStore(s...) {
		p = p.BindStore(ss.storeID)
		v, err := ss.root.Byte(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	}
	if ss.isAllowedWebsite(s...) {
		p = p.BindWebsite(ss.websiteID)
		v, err := ss.root.Byte(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	}
	p.ScopeHash = scope.DefaultHash
	v, err := ss.root.Byte(p)
	return v, scope.DefaultHash, err
}

// String traverses through the scopes store->website->default to find
// a matching string value.
func (ss scopedService) String(r cfgpath.Route, s ...scope.Scope) (string, scope.Hash, error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		return "", 0, errors.Wrapf(err, "[config] String. Route %q", r)
	}

	if ss.isAllowedStore(s...) {
		p = p.BindStore(ss.storeID)
		v, err := ss.root.String(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	}

	if ss.isAllowedWebsite(s...) {
		p = p.BindWebsite(ss.websiteID)
		v, err := ss.root.String(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	}
	p.ScopeHash = scope.DefaultHash
	v, err := ss.root.String(p)
	return v, scope.DefaultHash, err
}

// Bool traverses through the scopes store->website->default to find
// a matching bool value.
func (ss scopedService) Bool(r cfgpath.Route, s ...scope.Scope) (bool, scope.Hash, error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		return false, 0, errors.Wrapf(err, "[config] Bool. Route %q", r)
	}

	if ss.isAllowedStore(s...) {
		p = p.BindStore(ss.storeID)
		v, err := ss.root.Bool(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in store scope go to website scope

	if ss.isAllowedWebsite(s...) {
		p = p.BindWebsite(ss.websiteID)
		v, err := ss.root.Bool(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in website scope go to default scope
	p.ScopeHash = scope.DefaultHash
	v, err := ss.root.Bool(p)
	return v, scope.DefaultHash, err
}

// Float64 traverses through the scopes store->website->default to find
// a matching float64 value.
func (ss scopedService) Float64(r cfgpath.Route, s ...scope.Scope) (float64, scope.Hash, error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "[config] Float64. Route %q", r)
	}

	if ss.isAllowedStore(s...) {
		p = p.BindStore(ss.storeID)
		v, err := ss.root.Float64(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in store scope go to website scope

	if ss.isAllowedWebsite(s...) {
		p = p.BindWebsite(ss.websiteID)
		v, err := ss.root.Float64(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in website scope go to default scope
	p.ScopeHash = scope.DefaultHash
	v, err := ss.root.Float64(p)
	return v, scope.DefaultHash, err
}

// Int traverses through the scopes store->website->default to find
// a matching int value.
func (ss scopedService) Int(r cfgpath.Route, s ...scope.Scope) (int, scope.Hash, error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "[config] Int. Route %q", r)
	}

	if ss.isAllowedStore(s...) {
		p = p.BindStore(ss.storeID)
		v, err := ss.root.Int(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in store scope go to website scope

	if ss.isAllowedWebsite(s...) {
		p = p.BindWebsite(ss.websiteID)
		v, err := ss.root.Int(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in website scope go to default scope
	p.ScopeHash = scope.DefaultHash
	v, err := ss.root.Int(p)
	return v, scope.DefaultHash, err
}

// Time traverses through the scopes store->website->default to find
// a matching time.Time value.
func (ss scopedService) Time(r cfgpath.Route, s ...scope.Scope) (time.Time, scope.Hash, error) {
	// fallback to next parent scope if value does not exists
	p, err := cfgpath.New(r)
	if err != nil {
		return time.Time{}, 0, errors.Wrapf(err, "[config] Time. Route %q", r)
	}

	if ss.isAllowedStore(s...) {
		p = p.BindStore(ss.storeID)
		v, err := ss.root.Time(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in store scope go to website scope

	if ss.isAllowedWebsite(s...) {
		p = p.BindWebsite(ss.websiteID)
		v, err := ss.root.Time(p)
		if !errors.IsNotFound(err) || err == nil {
			// value found or err is not a NotFound error
			return v, p.ScopeHash, err
		}
	} // if not found in website scope go to default scope
	p.ScopeHash = scope.DefaultHash
	v, err := ss.root.Time(p)
	return v, scope.DefaultHash, err
}
