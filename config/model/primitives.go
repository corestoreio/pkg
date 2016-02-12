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

package model

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// Bool represents a path in config.Getter which handles bool values.
type Bool struct {
	baseValue
}

// NewBool creates a new Bool model with a given path.
func NewBool(path string, opts ...Option) Bool {
	return Bool{
		baseValue: NewValue(path, opts...),
	}
}

// Get returns a bool value.
// Errors gets appended to MultiErr
func (b Bool) Get(sg config.ScopedGetter) (v bool) {
	var err error
	if v, err = b.lookupBool(sg); err != nil {
		b.MultiErr = b.AppendErrors(err)
	}
	return
}

// Write writes a bool.
func (b Bool) Write(w config.Writer, v bool, s scope.Scope, scopeID int64) error {
	return b.baseValue.Write(w, v, s, scopeID)
}

// Str represents a path in config.Getter which handles string values.
// The name Str has been chosen to avoid conflict with the String() function
// in the Stringer interface.
type Str struct{ baseValue }

// NewStr creates a new Str model with a given path.
func NewStr(path string, opts ...Option) Str {
	return Str{baseValue: NewValue(path, opts...)}
}

// Get returns a string value
// Errors gets appended to MultiErr
func (str Str) Get(sg config.ScopedGetter) (v string) {
	var err error
	if v, err = str.lookupString(sg); err != nil {
		str.MultiErr = str.AppendErrors(err)
	}
	return
}

// Write writes a string value
func (str Str) Write(w config.Writer, v string, s scope.Scope, scopeID int64) error {
	return str.baseValue.Write(w, v, s, scopeID)
}

// Int represents a path in config.Getter which handles int values.
type Int struct{ baseValue }

// NewInt creates a new Int model with a given path.
func NewInt(path string, opts ...Option) Int {
	return Int{baseValue: NewValue(path, opts...)}
}

// Get returns an int value.
func (i Int) Get(sg config.ScopedGetter) (v int) {
	var err error
	if v, err = i.lookupInt(sg); err != nil {
		i.MultiErr = i.AppendErrors(err)
	}
	return
}

// Write writes an int value
func (i Int) Write(w config.Writer, v int, s scope.Scope, scopeID int64) error {
	return i.baseValue.Write(w, v, s, scopeID)
}

// Float64 represents a path in config.Getter which handles int values.
type Float64 struct{ baseValue }

// NewFloat64 creates a new Float64 model with a given path.
func NewFloat64(path string, opts ...Option) Float64 {
	return Float64{baseValue: NewValue(path, opts...)}
}

// Get returns a float64 value.
func (f Float64) Get(sg config.ScopedGetter) (v float64) {
	var err error
	if v, err = f.lookupFloat64(sg); err != nil {
		f.MultiErr = f.AppendErrors(err)
	}
	return
}

// Write writes a float64 value
func (f Float64) Write(w config.Writer, v float64, s scope.Scope, scopeID int64) error {
	return f.baseValue.Write(w, v, s, scopeID)
}
