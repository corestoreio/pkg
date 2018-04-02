// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/conv"
)

// Bool represents a path in config.Getter which handles bool values.
type Bool struct{ baseValue }

// NewBool creates a new Bool cfgmodel with a given path.
func NewBool(path string, opts ...Option) Bool {
	return Bool{
		baseValue: newBaseValue(path, opts...),
	}
}

// Get returns a bool value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (b Bool) Value(sg config.Scoped) (bool, error) {
	// This code must be kept in sync with other Value() functions

	if b.LastError != nil {
		return false, errors.Wrap(b.LastError, "[cfgmodel] Bool.Get.LastError")
	}

	var v bool
	var scp = b.initScope().Top()
	if b.HasField() {
		scp = b.Field.Scopes.Top()
		var err error
		v, err = conv.ToBoolE(b.Field.Default)
		if err != nil {
			return false, errors.NotValid.Newf("[cfgmodel] ToBoolE: %v", err)
		}
	}

	val, err := sg.Bool(b.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.NotFound.Match(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", b.route)
	default:
		// use default value v because sg found nothing
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset error
	}
	return v, err
}

// Write writes a bool value without validating it against the cfgsource.Slice.
func (b Bool) Write(w config.Writer, v bool, h scope.TypeID) error {
	return b.baseValue.Write(w, v, h)
}

// Byte represents a path in config.Getter which handles byte slices.
type Byte struct{ baseValue }

// NewByte creates a new Byte cfgmodel with a given path.
func NewByte(path string, opts ...Option) Byte {
	return Byte{baseValue: newBaseValue(path, opts...)}
}

// Get returns a byte slice from ScopedGetter, if empty the
// *element.Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *element.Field.Scopes is empty.
// The slice is owned by this function. You must copy it away for
// further modifications.
func (bt Byte) Value(sg config.Scoped) ([]byte, error) {
	// This code must be kept in sync with other Value() functions

	if bt.LastError != nil {
		return nil, errors.Wrap(bt.LastError, "[cfgmodel] Byte.Get.LastError")
	}

	var v []byte
	var scp = bt.initScope().Top()
	if bt.HasField() {
		scp = bt.Field.Scopes.Top()
		var err error
		v, err = conv.ToByteE(bt.Field.Default)
		if err != nil {
			return nil, errors.NotValid.Newf("[cfgmodel] ToByteE: %v", err)
		}
	}

	val, err := sg.Byte(bt.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.NotFound.Match(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", bt.route)
	default:
		// use default value v because sg found nothing
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset error
	}
	return v, err
}

// Write writes a byte slice without validating it against the cfgsource.Slice.
func (str Byte) Write(w config.Writer, v []byte, h scope.TypeID) error {
	return str.baseValue.Write(w, v, h)
}

// Str represents a path in config.Getter which handles string values.
// The name Str has been chosen to avoid conflict with the String() function
// in the fmt.Stringer interface.
type Str struct{ baseValue }

// NewStr creates a new Str cfgmodel with a given path.
func NewStr(path string, opts ...Option) Str {
	return Str{baseValue: newBaseValue(path, opts...)}
}

// Get returns a string value from ScopedGetter, if empty the
// *element.Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *element.Field.Scopes is empty.
func (str Str) Value(sg config.Scoped) (string, error) {
	// This code must be kept in sync with other Value() functions

	if str.LastError != nil {
		return "", errors.Wrap(str.LastError, "[cfgmodel] Str.Get.LastError")
	}

	var v string
	var scp = str.initScope().Top()
	if str.HasField() {
		scp = str.Field.Scopes.Top()
		var err error
		v, err = conv.ToStringE(str.Field.Default)
		if err != nil {
			return "", errors.NotValid.Newf("[cfgmodel] ToStringE: %v", err)
		}
	}

	val, err := sg.String(str.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.NotFound.Match(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", str.route)
	default:
		// use default value v because sg found nothing
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset error
	}
	return v, err
}

// Write writes a string value without validating it against the cfgsource.Slice.
func (str Str) Write(w config.Writer, v string, h scope.TypeID) error {
	return str.baseValue.Write(w, v, h)
}

// Int represents a path in config.Getter which handles int values.
type Int struct{ baseValue }

// NewInt creates a new Int cfgmodel with a given path.
func NewInt(path string, opts ...Option) Int {
	return Int{baseValue: newBaseValue(path, opts...)}
}

// Get returns an int value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (i Int) Value(sg config.Scoped) (int, error) {
	// This code must be kept in sync with other Value() functions

	if i.LastError != nil {
		return 0, errors.Wrap(i.LastError, "[cfgmodel] Int.Get.LastError")
	}

	var v int
	var scp = i.initScope().Top()
	if i.HasField() {
		scp = i.Field.Scopes.Top()
		var err error
		v, err = conv.ToIntE(i.Field.Default)
		if err != nil {
			return 0, errors.NotValid.Newf("[cfgmodel] ToIntE: %v", err)
		}
	}

	val, err := sg.Int(i.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.NotFound.Match(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", i.route)
	default:
		// use default value v because sg found nothing
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset error
	}
	return v, err
}

// Write writes an int value without validating it against the cfgsource.Slice.
func (i Int) Write(w config.Writer, v int, h scope.TypeID) error {
	return i.baseValue.Write(w, v, h)
}

// Float64 represents a path in config.Getter which handles float64 values.
type Float64 struct{ baseValue }

// NewFloat64 creates a new Float64 cfgmodel with a given path.
func NewFloat64(path string, opts ...Option) Float64 {
	return Float64{baseValue: newBaseValue(path, opts...)}
}

// Get returns a float64 value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
func (f Float64) Value(sg config.Scoped) (float64, error) {
	// This code must be kept in sync with other Value() functions

	if f.LastError != nil {
		return 0, errors.Wrap(f.LastError, "[cfgmodel] Float64.Get.LastError")
	}

	var v float64
	var scp = f.initScope().Top()
	if f.HasField() {
		scp = f.Field.Scopes.Top()
		if d := f.Field.Default; d != nil {
			var err error
			v, err = conv.ToFloat64E(d)
			if err != nil {
				return 0, errors.NotValid.Newf("[cfgmodel] ToFloat64E: %v", err)
			}
		}
	}

	val, err := sg.Float64(f.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.NotFound.Match(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", f.route)
	default:
		// use default value v because sg found nothing
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset error
	}
	return v, err
}

// Write writes a float64 value without validating it against the cfgsource.Slice.
func (f Float64) Write(w config.Writer, v float64, h scope.TypeID) error {
	return f.baseValue.Write(w, v, h)
}
