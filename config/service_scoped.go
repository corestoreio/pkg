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

	"github.com/corestoreio/csfw/store/scope"
)

// ScopedGetter is equal to Getter but the underlying implementation takes
// care of providing the correct scope: default, website or store and bubbling
// up the scope chain from store -> website -> default.
//
// This interface is mainly implemented in the store package. The functions
// should be the same as in Getter but only the different is the paths
// argument. A path can be either one string containing a valid path like a/b/c
// or it can consists of 3 path parts like "a", "b", "c". All other arguments
// are invalid. Returned error is mostly of ErrKeyNotFound.
type ScopedGetter interface {
	scope.Scoper
	String(paths ...string) (string, error)
	Bool(paths ...string) (bool, error)
	Float64(paths ...string) (float64, error)
	Int(paths ...string) (int, error)
	DateTime(paths ...string) (time.Time, error)
}

type scopedService struct {
	root      Getter
	websiteID int64
	groupID   int64
	storeID   int64
}

var _ ScopedGetter = (*scopedService)(nil)

func newScopedService(r Getter, websiteID, groupID, storeID int64) scopedService {
	return scopedService{
		root:      r,
		websiteID: websiteID,
		groupID:   groupID,
		storeID:   storeID,
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
func (ss scopedService) String(paths ...string) (v string, err error) {
	// fallback to next parent scope if value does not exists
	switch {
	case ss.storeID > 0:
		v, err = ss.root.String(Scope(scope.StoreID, ss.storeID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.String(Scope(scope.GroupID, ss.groupID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.String(Scope(scope.WebsiteID, ss.websiteID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.String(Scope(scope.DefaultID, 0), Path(paths...))
	}
}

// Bool traverses through the scopes store->group->website->default to find
// a matching bool value.
func (ss scopedService) Bool(paths ...string) (v bool, err error) {
	// fallback to next parent scope if value does not exists
	switch {
	case ss.storeID > 0:
		v, err = ss.root.Bool(Scope(scope.StoreID, ss.storeID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.Bool(Scope(scope.GroupID, ss.groupID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.Bool(Scope(scope.WebsiteID, ss.websiteID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.Bool(Scope(scope.DefaultID, 0), Path(paths...))
	}
}

// Float64 traverses through the scopes store->group->website->default to find
// a matching float64 value.
func (ss scopedService) Float64(paths ...string) (v float64, err error) {
	// fallback to next parent scope if value does not exists
	switch {
	case ss.storeID > 0:
		v, err = ss.root.Float64(Scope(scope.StoreID, ss.storeID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.Float64(Scope(scope.GroupID, ss.groupID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.Float64(Scope(scope.WebsiteID, ss.websiteID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.Float64(Scope(scope.DefaultID, 0), Path(paths...))
	}

}

// Int traverses through the scopes store->group->website->default to find
// a matching int value.
func (ss scopedService) Int(paths ...string) (v int, err error) {
	// fallback to next parent scope if value does not exists
	switch {
	case ss.storeID > 0:
		v, err = ss.root.Int(Scope(scope.StoreID, ss.storeID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.Int(Scope(scope.GroupID, ss.groupID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.Int(Scope(scope.WebsiteID, ss.websiteID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.Int(Scope(scope.DefaultID, 0), Path(paths...))
	}

}

// DateTime traverses through the scopes store->group->website->default to find
// a matching time.Time value.
func (ss scopedService) DateTime(paths ...string) (v time.Time, err error) {
	// fallback to next parent scope if value does not exists
	switch {
	case ss.storeID > 0:
		v, err = ss.root.DateTime(Scope(scope.StoreID, ss.storeID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in store scope go to group scope
	case ss.groupID > 0:
		v, err = ss.root.DateTime(Scope(scope.GroupID, ss.groupID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in group scope go to website scope
	case ss.websiteID > 0:
		v, err = ss.root.DateTime(Scope(scope.WebsiteID, ss.websiteID), Path(paths...))
		if NotKeyNotFoundError(err) || err == nil {
			return // value found or err is not a KeyNotFound error
		}
		fallthrough // if not found in website scope go to default scope
	default:
		return ss.root.DateTime(Scope(scope.DefaultID, 0), Path(paths...))
	}
}
