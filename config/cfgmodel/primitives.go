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

package cfgmodel

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/juju/errors"
)

// Bool represents a path in config.Getter which handles bool values.
type Bool struct{ baseValue }

// NewBool creates a new Bool cfgmodel with a given path.
func NewBool(path string, opts ...Option) Bool {
	return Bool{
		baseValue: NewValue(path, opts...),
	}
}

// Get returns a bool value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (b Bool) Get(sg config.ScopedGetter) (bool, error) {
	// This code must be kept in sync with other Get() functions

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

// Byte represents a path in config.Getter which handles byte slices.
type Byte struct{ baseValue }

// NewByte creates a new Byte cfgmodel with a given path.
func NewByte(path string, opts ...Option) Byte {
	return Byte{baseValue: NewValue(path, opts...)}
}

// Get returns a byte slice from ScopedGetter, if empty the
// *element.Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *element.Field.Scopes is empty.
// The slice is owned by this function. You must copy it away for
// further modifications.
func (bt Byte) Get(sg config.ScopedGetter) ([]byte, error) {
	// This code must be kept in sync with other lookup*() functions

	var v []byte
	var scp = scope.DefaultID
	if bt.Field != nil {
		scp = bt.Field.Scopes.Top()
		var err error
		v, err = conv.ToByteE(bt.Field.Default)
		if err != nil {
			return nil, errors.Mask(err)
		}
	}

	val, err := sg.Byte(bt.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", bt.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a byte slice without validating it against the source.Slice.
func (str Byte) Write(w config.Writer, v []byte, s scope.Scope, scopeID int64) error {
	return str.baseValue.Write(w, v, s, scopeID)
}

// Str represents a path in config.Getter which handles string values.
// The name Str has been chosen to avoid conflict with the String() function
// in the Stringer interface.
type Str struct{ baseValue }

// NewStr creates a new Str cfgmodel with a given path.
func NewStr(path string, opts ...Option) Str {
	return Str{baseValue: NewValue(path, opts...)}
}

// Get returns a string value from ScopedGetter, if empty the
// *element.Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *element.Field.Scopes is empty.
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

// NewInt creates a new Int cfgmodel with a given path.
func NewInt(path string, opts ...Option) Int {
	return Int{baseValue: NewValue(path, opts...)}
}

// Get returns an int value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (i Int) Get(sg config.ScopedGetter) (int, error) {
	// This code must be kept in sync with other Get() functions

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

// NewFloat64 creates a new Float64 cfgmodel with a given path.
func NewFloat64(path string, opts ...Option) Float64 {
	return Float64{baseValue: NewValue(path, opts...)}
}

// Get returns a float64 value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (f Float64) Get(sg config.ScopedGetter) (float64, error) {
	// This code must be kept in sync with other Get() functions

	var v float64
	var scp = scope.DefaultID
	if f.Field != nil {
		scp = f.Field.Scopes.Top()
		if d := f.Field.Default; d != nil {
			var err error
			v, err = conv.ToFloat64E(d)
			if err != nil {
				return 0, errors.Mask(err)
			}
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
