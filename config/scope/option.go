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

	"github.com/corestoreio/csfw/utils"
)

// ErrUnsupportedScope whenever a not valid scope ID will be provided this
// error gets returned.
var ErrUnsupportedScope = errors.New("Unsupported Scope ID")

// OptionFunc for the Init() function. One of the three struct fields
// must be set or the function panics. Instead of the IDer interface you can
// also provide a *Coder interface. Order of scope calculation:
// Website -> Group -> Store. Be sure to set e.g. Website and Group to nil
// if you need initialization for store level.
type OptionFunc func(*option)

// option will be kept private to not confuse others ...
type option struct {
	Website  WebsiteIDer
	Group    GroupIDer
	Store    StoreIDer
	lastErrs []error
}

func NewOption(opts ...OptionFunc) (option, error) {
	o := option{}
	for _, opt := range opts {
		if nil != opt {
			opt(&o)
		}
	}
	if nil != o.lastErrs {
		return o, o
	}
	return o, nil
}

func ApplyWebsite(w WebsiteIDer) OptionFunc {
	if w == nil {
		return func(o *option) { o.lastErrs = append(o.lastErrs, errors.New("Website argument cannot be nil")) }
	}
	return func(o *option) {
		if o.Store != nil || o.Group != nil {
			o.lastErrs = append(o.lastErrs, errors.New("Store or Group already set"))
		} else {
			o.Website = w
		}
	}
}

func ApplyGroup(g GroupIDer) OptionFunc {
	if g == nil {
		return func(o *option) { o.lastErrs = append(o.lastErrs, errors.New("Group argument cannot be nil")) }
	}
	return func(o *option) {
		if o.Website != nil || o.Store != nil {
			o.lastErrs = append(o.lastErrs, errors.New("Website or Store already set"))
		} else {
			o.Group = g
		}
	}
}

func ApplyStore(s StoreIDer) OptionFunc {
	if s == nil {
		return func(o *option) { o.lastErrs = append(o.lastErrs, errors.New("Store argument cannot be nil")) }
	}
	return func(o *option) {
		if o.Website != nil || o.Group != nil {
			o.lastErrs = append(o.lastErrs, errors.New("Website or Group already set"))
		} else {
			o.Store = s
		}
	}
}

func ApplyCode(code string, scopeType Scope) OptionFunc {
	return func(o *option) {
		c := MockCode(code)
		// GroupID does not have a scope code
		switch scopeType {
		case WebsiteID:
			o.Website = c
		case StoreID:
			o.Store = c
		default:
			o.lastErrs = append(o.lastErrs, ErrUnsupportedScope)
		}
	}
}

func ApplyID(scopeID int64, scopeType Scope) OptionFunc {
	return func(o *option) {
		i := MockID(scopeID)
		switch scopeType {
		case WebsiteID:
			o.Website = i
		case GroupID:
			o.Group = i
		case StoreID:
			o.Store = i
		default:
			o.lastErrs = append(o.lastErrs, ErrUnsupportedScope)
		}
	}
}

func (o option) Scope() (s Scope) {
	s = AbsentID
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

func (o option) StoreCode() (code string) {
	if sc, ok := o.Store.(StoreCoder); ok {
		code = sc.StoreCode()
	}
	return
}

func (o option) WebsiteCode() (code string) {
	if wc, ok := o.Website.(WebsiteCoder); ok {
		code = wc.WebsiteCode()
	}
	return
}

var _ error = (*option)(nil)

// Error satisfy the error interface
func (o option) Error() string {
	return utils.Errors(o.lastErrs...)
}
