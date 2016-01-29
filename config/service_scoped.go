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

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/store/scope"
)

// ScopedGetter is equal to Getter but the underlying implementation takes
// care of providing the correct scope: default, website or store and bubbling
// up the scope chain from store -> website -> default.
//
// This interface is mainly implemented in the store package. The functions
// should be the same as in Getter but only the different is the route
// argument. A route represents always "a/b/c".
// Returned error is mostly of ErrKeyNotFound.
type ScopedGetter interface {
	scope.Scoper
	String(r path.Route) (string, error)
	Bool(r path.Route) (bool, error)
	Float64(r path.Route) (float64, error)
	Int(r path.Route) (int, error)
	DateTime(r path.Route) (time.Time, error)
}

// think about that segregation
//type ScopedStringer interface {
//	scope.Scoper
//	String(r path.Route) (string, error)
//}
//
//type ScopedBooler interface {
//	scope.Scoper
//	Bool(r path.Route) (bool, error)
//}

type scopedService struct {
	root       Getter
	websiteID  int64
	websiteArg ArgFunc // 1x alloc, like a cache, can be nil
	groupID    int64
	groupArg   ArgFunc // 1x alloc, like a cache, can be nil
	storeID    int64
	storeArg   ArgFunc // 1x alloc, like a cache, can be nil
}

var _ ScopedGetter = (*scopedService)(nil)

func newScopedService(r Getter, websiteID, groupID, storeID int64) scopedService {

	var wa, ga, sa ArgFunc
	if websiteID > 0 {
		wa = ScopeWebsite(websiteID)
	}
	if groupID > 0 {
		ga = ScopeGroup(groupID)
	}
	if storeID > 0 {
		sa = ScopeStore(storeID)
	}

	return scopedService{
		root:       r,
		websiteID:  websiteID,
		websiteArg: wa,
		groupID:    groupID,
		groupArg:   ga,
		storeID:    storeID,
		storeArg:   sa,
	}
}

// Scope tells you the current underlying scope and its website, group or store ID
func (ss scopedService) Scope() (scope.Scope, int64) {
	switch {
	case ss.storeID > 0:
		return scope.StoreID, ss.storeID
	case ss.groupID > 0:
		return scope.GroupID, ss.groupID
	case ss.websiteID > 0:
		return scope.WebsiteID, ss.websiteID
	default:
		return scope.DefaultID, 0
	}
}

// String traverses through the scopes store->group->website->default to find
// a matching string value.
func (ss scopedService) String(r path.Route) (v string, err error) {
	// fallback to next parent scope if value does not exists
	rArg := Route(r)
	switch {
	case ss.storeID > 0:
		v, err = ss.root.String(ss.storeArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.String(ss.groupArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.String(ss.websiteArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.String(scopeDefaultArg, rArg)
	}
}

// Bool traverses through the scopes store->group->website->default to find
// a matching bool value.
func (ss scopedService) Bool(r path.Route) (v bool, err error) {
	// fallback to next parent scope if value does not exists
	rArg := Route(r)
	switch {
	case ss.storeID > 0:
		v, err = ss.root.Bool(rArg, ss.storeArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.Bool(rArg, ss.groupArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.Bool(rArg, ss.websiteArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.Bool(scopeDefaultArg, rArg)
	}
}

// Float64 traverses through the scopes store->group->website->default to find
// a matching float64 value.
func (ss scopedService) Float64(r path.Route) (v float64, err error) {
	// fallback to next parent scope if value does not exists
	rArg := Route(r)
	switch {
	case ss.storeID > 0:
		v, err = ss.root.Float64(ss.storeArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.Float64(ss.groupArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.Float64(ss.websiteArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.Float64(scopeDefaultArg, rArg)
	}

}

// Int traverses through the scopes store->group->website->default to find
// a matching int value.
func (ss scopedService) Int(r path.Route) (v int, err error) {
	// fallback to next parent scope if value does not exists
	rArg := Route(r)
	switch {
	case ss.storeID > 0:
		v, err = ss.root.Int(rArg, ss.storeArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.Int(rArg, ss.groupArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.Int(rArg, ss.websiteArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.Int(scopeDefaultArg, rArg)
	}

}

// DateTime traverses through the scopes store->group->website->default to find
// a matching time.Time value.
func (ss scopedService) DateTime(r path.Route) (v time.Time, err error) {
	// fallback to next parent scope if value does not exists
	rArg := Route(r)
	switch {
	case ss.storeID > 0:
		v, err = ss.root.DateTime(ss.storeArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.DateTime(ss.groupArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.DateTime(ss.websiteArg, rArg)
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.DateTime(scopeDefaultArg, rArg)
	}
}
