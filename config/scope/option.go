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

package scope

import (
	"errors"
)

// ErrUnsupportedScope whenever a not valid scope ID will be provided this
// error gets returned.
var ErrUnsupportedScope = errors.New("Unsupported Scope ID")

// Option for the Init() function. Instead of the IDer interface you can
// also provide a *Coder interface. Order of scope precedence:
// Website -> Group -> Store. Be sure to set e.g. Website and Group to nil
// if you need initialization for store level.
type Option struct {
	Website WebsiteIDer
	Group   GroupIDer
	Store   StoreIDer
}

func SetByCode(code string, scopeType Scope) (o Option, err error) {
	c := MockCode(code)
	// GroupID does not have a scope code
	switch scopeType {
	case WebsiteID:
		o.Website = c
	case StoreID:
		o.Store = c
	default:
		err = ErrUnsupportedScope
	}
	return
}

func SetByID(scopeID int64, scopeType Scope) (o Option, err error) {
	i := MockID(scopeID)
	// the order of the cases is important
	switch scopeType {
	case WebsiteID:
		o.Website = i
	case GroupID:
		o.Group = i
	case StoreID:
		o.Store = i
	default:
		err = ErrUnsupportedScope
	}
	return
}

func (o Option) Scope() (s Scope) {
	s = AbsentID
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

func (o Option) StoreCode() (code string) {
	if sc, ok := o.Store.(StoreCoder); ok {
		code = sc.StoreCode()
	}
	return
}

func (o Option) WebsiteCode() (code string) {
	if wc, ok := o.Website.(WebsiteCoder); ok {
		code = wc.WebsiteCode()
	}
	return
}
