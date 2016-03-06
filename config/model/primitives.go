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
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/juju/errors"
)

// Bool represents a path in config.Getter which handles bool values.
type Bool struct{ baseValue }

// NewBool creates a new Bool model with a given path.
func NewBool(path string, opts ...Option) Bool {
	return Bool{
		baseValue: NewValue(path, opts...),
	}
}

// Get returns a bool value.
// Retrieves a bool value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (b Bool) Get(sg config.ScopedGetter) (bool, error) {
	// This code must be kept in sync with other lookup*() functions

	var v bool
	var scp = scope.DefaultID
	if b.Field != nil {
		scp = b.Field.Scopes.Top()
		var err error
		v, err = conv.ToBoolE(b.Field.Default)
		if err != nil {
			return false, errors.Mask(err)
		}
	}

	val, err := sg.Bool(b.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", b.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a bool value without validating it against the source.Slice.
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
// Retrieves a string value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (str Str) Get(sg config.ScopedGetter) (string, error) {
	// This code must be kept in sync with other lookup*() functions

	var v string
	var scp = scope.DefaultID
	if str.Field != nil {
		scp = str.Field.Scopes.Top()
		var err error
		v, err = conv.ToStringE(str.Field.Default)
		if err != nil {
			return "", errors.Mask(err)
		}
	}

	val, err := sg.String(str.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", str.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a string value without validating it against the source.Slice.
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
// Retrieves an int value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (i Int) Get(sg config.ScopedGetter) (int, error) {
	// This code must be kept in sync with other lookup*() functions

	var v int
	var scp = scope.DefaultID
	if i.Field != nil {
		scp = i.Field.Scopes.Top()
		var err error
		v, err = conv.ToIntE(i.Field.Default)
		if err != nil {
			return 0, errors.Mask(err)
		}
	}

	val, err := sg.Int(i.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", i.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes an int value without validating it against the source.Slice.
func (i Int) Write(w config.Writer, v int, s scope.Scope, scopeID int64) error {
	return i.baseValue.Write(w, v, s, scopeID)
}

// Float64 represents a path in config.Getter which handles float64 values.
type Float64 struct{ baseValue }

// NewFloat64 creates a new Float64 model with a given path.
func NewFloat64(path string, opts ...Option) Float64 {
	return Float64{baseValue: NewValue(path, opts...)}
}

// Get returns a float64 value.
// Retrieves a float64 value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (f Float64) Get(sg config.ScopedGetter) (float64, error) {
	// This code must be kept in sync with other lookup*() functions

	var v float64
	var scp = scope.DefaultID
	if f.Field != nil {
		scp = f.Field.Scopes.Top()
		var err error
		v, err = conv.ToFloat64E(f.Field.Default)
		if err != nil {
			return 0, errors.Mask(err)
		}
	}

	val, err := sg.Float64(f.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", f.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a float64 value without validating it against the source.Slice.
func (f Float64) Write(w config.Writer, v float64, s scope.Scope, scopeID int64) error {
	return f.baseValue.Write(w, v, s, scopeID)
}

// Time represents a path in config.Getter which handles time values.
type Time struct{ baseValue }

// NewTime creates a new Time model with a given path.
func NewTime(path string, opts ...Option) Time {
	return Time{baseValue: NewValue(path, opts...)}
}

// Get returns a time value.
// Retrieves a time value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
// Get is able to parse available time formats as defined in
// github.com/corestoreio/csfw/util/conv.StringToDate()
func (t Time) Get(sg config.ScopedGetter) (time.Time, error) {
	// This code must be kept in sync with other lookup*() functions

	var v time.Time
	var scp = scope.DefaultID
	if t.Field != nil {
		scp = t.Field.Scopes.Top()
		var err error
		v, err = conv.ToTimeE(t.Field.Default)
		if err != nil {
			return time.Time{}, errors.Mask(err)
		}
	}

	val, err := sg.Time(t.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", t.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a time value without validating it against the source.Slice.
func (t Time) Write(w config.Writer, v time.Time, s scope.Scope, scopeID int64) error {
	return t.baseValue.Write(w, v, s, scopeID)
}
