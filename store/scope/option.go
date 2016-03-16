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

package scope

import (
	"errors"
)

// ErrUnsupportedScope gets returned when a string does not match
// StrDefault, StrWebsites or StrStores constants.
var ErrUnsupportedScope = errors.New("Unsupported StrScope type")

// ErrUnsupportedScopeID whenever a not valid scope ID will be provided.
// Neither a WebsiteID nor a GroupID nor a StoreID.
var ErrUnsupportedScopeID = errors.New("Unsupported Scope ID")

// Option takes care of the hierarchical level between Website, Group and Store.
// Option can be used as an argument in other functions.
// Instead of the [Website|Group|Store]IDer interface you can
// also provide a [Website|Group|Store]Coder interface.
// Order of scope precedence:
// Website -> Group -> Store. Be sure to set e.g. Website and Group to nil
// if you need initialization for store level.
type Option struct {
	Website WebsiteIDer
	Group   GroupIDer
	Store   StoreIDer
}

// SetByCode depending on the scopeType the code string gets converted into a
// StoreCoder or WebsiteCoder interface and the appropriate struct fields
// get assigned with the *Coder interface. scopeType can only be WebsiteID or
// StoreID because a Group code does not exists.
func SetByCode(scp Scope, code string) (o Option, err error) {
	c := MockCode(code)
	// GroupID does not have a scope code
	switch scp {
	case WebsiteID:
		o.Website = c
	case StoreID:
		o.Store = c
	default:
		err = ErrUnsupportedScopeID
	}
	return
}

// SetByID depending on the scopeType the scopeID int64 gets converted into a
// [Website|Group|Store]IDer.
func SetByID(scp Scope, id int64) (o Option, err error) {
	i := MockID(id)
	// the order of the cases is important
	switch scp {
	case WebsiteID:
		o.Website = i
	case GroupID:
		o.Group = i
	case StoreID:
		o.Store = i
	default:
		err = ErrUnsupportedScopeID
	}
	return
}

// Scope returns the underlying scope ID depending on which struct field is set.
// It maintains the hierarchical order: 1. Website, 2. Group, 3. Store.
// If no field has been set returns DefaultID.
func (o Option) Scope() (s Scope) {
	s = DefaultID
	// the order of the cases is important
	switch {
	case o.Website != nil:
		s = WebsiteID
	case o.Group != nil:
		s = GroupID
	case o.Store != nil:
		s = StoreID
	}
	return
}

// String is short hand for Option.Scope().String()
func (o Option) String() string {
	return o.Scope().String()
}

// StoreCode extracts the Store code. Checks if the interface StoreCoder
// is available.
func (o Option) StoreCode() (code string) {
	if sc, ok := o.Store.(StoreCoder); ok {
		code = sc.StoreCode()
	}
	return
}

// WebsiteCode extracts the Website code. Checks if the interface WebsiteCoder
// is available.
func (o Option) WebsiteCode() (code string) {
	if wc, ok := o.Website.(WebsiteCoder); ok {
		code = wc.WebsiteCode()
	}
	return
}
